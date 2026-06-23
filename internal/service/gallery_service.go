package service

import (
	"fmt"
	"time"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

func ListGallery() ([]model.GalleryItem, error) {
	rows, err := db.DB.Query(
		`SELECT id, url, type, caption, sort_order, created_at FROM gallery_items ORDER BY sort_order, created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.GalleryItem
	for rows.Next() {
		var g model.GalleryItem
		if err := rows.Scan(&g.ID, &g.URL, &g.Type, &g.Caption, &g.SortOrder, &g.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, g)
	}
	if items == nil {
		items = []model.GalleryItem{}
	}
	return items, nil
}

func CreateGalleryItem(g model.GalleryItem) (model.GalleryItem, error) {
	if g.ID == "" {
		g.ID = fmt.Sprintf("gal_%s", time.Now().Format("20060102150405999"))
	}
	if g.Type == "" {
		g.Type = "image"
	}
	_, err := db.DB.Exec(
		`INSERT INTO gallery_items (id, url, type, caption, sort_order) VALUES ($1,$2,$3,$4,$5)`,
		g.ID, g.URL, g.Type, g.Caption, g.SortOrder,
	)
	if err != nil {
		return g, err
	}
	row := db.DB.QueryRow(`SELECT id, url, type, caption, sort_order, created_at FROM gallery_items WHERE id=$1`, g.ID)
	row.Scan(&g.ID, &g.URL, &g.Type, &g.Caption, &g.SortOrder, &g.CreatedAt)
	return g, nil
}

func UpdateGalleryItem(id string, g model.GalleryItem) error {
	if g.Type == "" {
		g.Type = "image"
	}
	_, err := db.DB.Exec(
		`UPDATE gallery_items SET url=$1, type=$2, caption=$3, sort_order=$4 WHERE id=$5`,
		g.URL, g.Type, g.Caption, g.SortOrder, id,
	)
	return err
}

func DeleteGalleryItem(id string) error {
	_, err := db.DB.Exec(`DELETE FROM gallery_items WHERE id=$1`, id)
	return err
}
