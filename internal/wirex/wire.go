package wirex

import (
	"context"
	"gin-admin/internal/mods"
	"gin-admin/pkg/util"

	"github.com/google/wire"
)

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	wire.Build(
		InitCacher,
		InitDB,
		InitAuth,
		wire.NewSet(wire.Struct(new(util.Trans), "*")),
		wire.NewSet(wire.Struct(new(Injector), "*")),
		mods.Set,
	)
	return new(Injector), nil, nil
}
