package cachex

import "github.com/dgraph-io/badger/v3"

type BadgerConfig struct {
	Path string
}

func NewBadgerCache(cfg BadgerConfig, opts ...Option) Cacher {
	defaultOpts := &options{
		Delimiter: defaultDelimiter,
	}

	for _, o := range opts {
		o(defaultOpts)
	}

	badgerOpts := badger.DefaultOptions(cfg.Path)
	badgerOpts = badgerOpts.WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(badgerOpts)
	if err != nil {
		panic(err)
	}

	return &badgerCache{
		opts: defaultOpts,
		db:   db,
	}
}

type badgerCache struct {
	opts *options
	db   *badger.DB
}
