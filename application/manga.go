package main

import (
	"context"
	"database/sql"
	"errors"
)

type manga struct {
	ID     string
	Title  string
	Author string
}

type volume struct {
	ID            string
	Number        int
	Title         string
	PageCount     int
	PageExtension string
}

func (a app) findManga(ctx context.Context, id string) (manga, bool, error) {
	var item manga
	err := a.db.QueryRowContext(ctx, `SELECT id, title, author_name FROM mangas WHERE id = ?`, id).Scan(&item.ID, &item.Title, &item.Author)
	if errors.Is(err, sql.ErrNoRows) {
		return manga{}, false, nil
	}
	if err != nil {
		return manga{}, false, err
	}
	return item, true, nil
}

func (a app) findVolume(ctx context.Context, mangaID, volumeID string) (manga, volume, bool, error) {
	var item manga
	var selectedVolume volume
	err := a.db.QueryRowContext(ctx, `
SELECT m.id, m.title, m.author_name, v.id, v.number, v.title, v.page_count, v.page_extension
FROM mangas m JOIN volumes v ON v.manga_id = m.id
WHERE m.id = ? AND v.id = ?`, mangaID, volumeID).Scan(
		&item.ID, &item.Title, &item.Author, &selectedVolume.ID, &selectedVolume.Number, &selectedVolume.Title, &selectedVolume.PageCount, &selectedVolume.PageExtension)
	if errors.Is(err, sql.ErrNoRows) {
		return manga{}, volume{}, false, nil
	}
	if err != nil {
		return manga{}, volume{}, false, err
	}
	return item, selectedVolume, true, nil
}
