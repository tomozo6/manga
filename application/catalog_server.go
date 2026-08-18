package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/tomozo6/manga/application/internal/catalog"
)

func openCatalogForServer(ctx context.Context) (*sql.DB, func(), error) {
	if path := os.Getenv("CATALOG_DB"); path != "" {
		db, err := catalog.OpenReadonly(path)
		return db, func() {}, err
	}
	source := os.Getenv("CATALOG_SOURCE_DIR")
	if source == "" {
		source = "catalog/mangas"
	}
	file, err := os.CreateTemp("", "manga-catalog-*.db")
	if err != nil {
		return nil, nil, err
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = os.Remove(path) }
	if err := catalog.Build(ctx, source, path); err != nil {
		cleanup()
		return nil, nil, err
	}
	db, err := catalog.OpenReadonly(path)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return db, cleanup, nil
}
