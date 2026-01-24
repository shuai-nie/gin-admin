package rbac

import (
	"bytes"
	"context"
	"fmt"
	"gin-admin/internal/config"
	"gin-admin/internal/mods/rbac/dal"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/logging"
	"gin-admin/pkg/util"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type Casbinx struct {
	enforcer        *atomic.Value
	ticker          *time.Ticker
	Cache           cachex.Cacher
	MenuDAL         *dal.Menu
	MenuResourceDAL *dal.MenuResource
	RoleDAL         *dal.Role
}

func (a *Casbinx) GetEnforcer() *casbin.Enforcer {
	if v := a.enforcer.Load(); v != nil {
		return v.(*casbin.Enforcer)
	}
	return nil
}

type policyQueItem struct {
	RoleID    string
	Resources schema.MenuResources
}

func (a *Casbinx) Load(ctx context.Context) error {
	if config.C.Middleware.Casbin.Disable {
		return nil
	}

	a.enforcer = new(atomic.Value)
	if err := a.Load(ctx); err != nil {
		return err
	}

	go a.autoLoad(ctx)
	return nil
}

func (a *Casbinx) load(ctx context.Context) error {
	start := time.Now()
	roleResult, err := a.RoleDAL.Query(ctx, schema.RoleQueryParam{
		Status: schema.RoleStatusEnabled,
	}, schema.RoleQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return nil
	}

	var resCount int32
	queue := make(chan *policyQueItem, len(roleResult.Data))
	threadNum := config.C.Middleware.Casbin.LoadThread
	lock := new(sync.Mutex)
	buf := new(bytes.Buffer)

	wg := new(sync.WaitGroup)
	wg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go func() {
			defer wg.Done()

			ibuf := new(bytes.Buffer)
			for item := range queue {
				for _, res := range item.Resources {
					_, _ = ibuf.WriteString(fmt.Sprintf("%s,%s,%s\n", item.RoleID, res.Path, res.Method))
				}
			}
			lock.Lock()

			_, _ = buf.Write(ibuf.Bytes())
			lock.Unlock()
		}()
	}

	for _, item := range roleResult.Data {
		resources, err := a.queryRoleResources(ctx, item.ID)
		if err != nil {
			logging.Context(ctx).Error("query role resources error", zap.Error(err))
			continue
		}
		atomic.AddInt32(&resCount, int32(len(resources)))
		queue <- &policyQueItem{
			RoleID:    item.ID,
			Resources: resources,
		}
	}

	close(queue)
	wg.Wait()

	if buf.Len() > 0 {
		policyFile := filepath.Join(config.C.General.WorkDir, config.C.Middleware.Casbin.GenPolicyFile)
		_ = os.Rename(policyFile, policyFile+".bak")
		_ = os.WriteFile(policyFile, buf.Bytes(), 0755)
		if err := os.WriteFile(policyFile, buf.Bytes(), 066); err != nil {
			logging.Context(ctx).Error("Failed to write policy file", zap.Error(err))
			return err
		}

		_ = os.Chmod(policyFile, 0444)

		modelFile := filepath.Join(config.C.General.WorkDir, config.C.Middleware.Casbin.ModelFile)
		e, err := casbin.NewEnforcer(modelFile, policyFile)
		if err != nil {
			logging.Context(ctx).Error("Failed to create casbin enforcer", zap.Error(err))
			return err
		}

		e.EnableLog(config.C.IsDebug())
		a.enforcer.Store(e)
	}

	logging.Context(ctx).Info("Load casbin policy success",
		zap.Duration("cost", time.Since(start)),
		zap.Int("roles", len(roleResult.Data)),
		zap.Int32("resources", resCount),
		zap.Int("bytes", buf.Len()),
	)
	return nil
}

func (a *Casbinx) queryRoleResources(ctx context.Context, roleID string) (schema.MenuResources, error) {
	menuResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
		RoleID: roleID,
		Status: schema.MenuStatusEnabled,
	}, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id", "parent_id", "parent_path"},
		},
	})

	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, nil
	}

	menuIDs := make([]string, 0, len(menuResult.Data))
	menuIDMapper := make(map[string]struct{})
	for _, item := range menuResult.Data {
		if _, ok := menuIDMapper[item.ID]; ok {
			continue
		}

		menuIDs = append(menuIDs, item.ID)
		menuIDMapper[item.ID] = struct{}{}

		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, util.TreePathDelimiter) {
				if pid == "" {
					continue
				}
				if _, ok := menuIDMapper[pid]; ok {
					continue
				}

				menuIDs = append(menuIDs, pid)
				menuIDMapper[pid] = struct{}{}
			}
		}
	}

	menuResourceResult, err := a.MenuResourceDAL.Query(ctx, schema.MenuResourceQueryParam{
		MenuIDs: menuIDs,
	})
	if err != nil {
		return nil, err
	}
	return menuResourceResult.Data, nil
}

func (a *Casbinx) autoLoad(ctx context.Context) {
	var lastUpdated int64
	a.ticker = time.NewTicker(time.Duration(config.C.Middleware.Casbin.AutoLoadInterval) * time.Second)
	for range a.ticker.C {
		val, ok, err := a.Cache.Get(ctx, config.CacheNSForRole, config.CacheKeyForSyncToCasbin)
		if err != nil {
			logging.Context(ctx).Error("get cache error", zap.Error(err), zap.String("key", config.CacheKeyForSyncToCasbin))
			continue
		} else if !ok {
			continue
		}

		updated, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			logging.Context(ctx).Error("parse cache value error", zap.Error(err), zap.String("key", config.CacheKeyForSyncToCasbin))
			continue
		}

		if lastUpdated < updated {
			if err := a.load(ctx); err != nil {
				logging.Context(ctx).Error("load casbin policy error", zap.Error(err))
			} else {
				lastUpdated = updated
			}
		}
	}
}

func (a *Casbinx) Release(ctx context.Context) error {
	if a.ticker != nil {
		a.ticker.Stop()
	}
	return nil
}
