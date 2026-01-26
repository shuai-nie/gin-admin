package bootstrap

import (
	"context"
	"fmt"
	"gin-admin/internal/config"
	"gin-admin/internal/utility/prom"
	"gin-admin/internal/wirex"
	"gin-admin/pkg/logging"
	"gin-admin/pkg/util"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"go.uber.org/zap"
)

type RunConfig struct {
	WorkDir   string
	Configs   string
	StaticDir string
}

func Run(ctx context.Context, runCfg RunConfig) error {
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf("failed tp sync zap logger: %s \n", err.Error())
		}
	}()

	workDir := runCfg.WorkDir
	staticDir := runCfg.StaticDir
	config.MustLoad(workDir, strings.Split(runCfg.Configs, ",")...)
	config.C.General.WorkDir = workDir
	config.C.Middleware.Static.Dir = staticDir
	config.C.Print()
	config.C.PreLoad()

	cleanLoggerFn, err := logging.InitWithConfig(ctx, &config.C.Logger, initLoggerHook)
	if err != nil {
		return err
	}

	ctx = logging.NewTag(ctx, logging.TagKeyMain)

	logging.Context(ctx).Info("starting service ...",
		zap.String("version", config.C.General.Version),
		zap.Int("pid", os.Getpid()),
		zap.String("workdir", workDir),
		zap.String("config", runCfg.Configs),
		zap.String("static", staticDir),
	)

	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("starting pprof server ...", zap.String("addr", addr))
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("pprof server error", zap.Error(err))
			}
		}()
	}

	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return err
	}

	if err := injector.M.Init(ctx); err != nil {
		return err
	}

	prom.Init()

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
		cleanHTTPServerFn, err := startHTTPServer(ctx, injector)
		if err != nil {
			return cleanInjectorFn, err
		}

		return func() {
			if err := injector.M.Release(ctx); err != nil {
				logging.Context(ctx).Error("failed to release mods", zap.Error(err))
			}

			if cleanHTTPServerFn != nil {
				cleanHTTPServerFn()
			}

			if cleanInjectorFn != nil {
				cleanInjectorFn()
			}

			if cleanLoggerFn != nil {
				cleanLoggerFn()
			}
		}, nil

	})
}
