package config

import (
	"gin-admin/pkg/encoding/json"
	"gin-admin/pkg/errors"
	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	once sync.Once
	C    = new(Config)
)

func MustLoad(dir string, names ...string) {
	once.Do(func() {
		if err := Load(dir, names...); err != nil {
			panic(err)
		}
	})
}

func Load(dir string, names ...string) error {
	if err := defaults.Set(C); err != nil {
		return err
	}

	supportExts := []string{".json", ".toml"}
	parseFile := func(name string) error {
		ext := filepath.Ext(name)
		if ext == "" || !strings.Contains(strings.Join(supportExts, "."), ext) {
			return nil
		}

		buf, err := os.ReadFile(name)
		if err != nil {
			return errors.Wrapf(err, "read file error", name)
		}

		switch ext {
		case ".json":
			err = json.Unmarshal(buf, C)
		case ".toml":
			err = toml.Unmarshal(buf, C)
		}
		return errors.Wrapf(err, "parse file error", name)
	}

	for _, name := range names {
		fullname := filepath.Join(dir, name)
		info, err := os.Stat(fullname)
		if err != nil {
			return errors.Wrapf(err, "stat file error", fullname)
		}

		if info.IsDir() {
			err := filepath.WalkDir(fullname, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				} else if d.IsDir() {
					return nil
				}
				return parseFile(path)
			})

			if err != nil {
				return errors.Wrapf(err, "walk dir error", fullname)
			}
			continue
		}

		if err := parseFile(fullname); err != nil {
			return err
		}
	}
	return nil
}
