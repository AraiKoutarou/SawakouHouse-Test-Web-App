// internal/infrastructure/persistence/location/repository.go: 都道府県カラムを追加。
package location

import (
	"database/sql"
	"log"

	"github.com/arakou0812/backend/internal/domain/location"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS locations (
		id         SERIAL PRIMARY KEY,
		place_id   VARCHAR(255),
		title      VARCHAR(255) NOT NULL,
		address    TEXT,
		prefecture VARCHAR(50),              -- 都道府県名
		category   VARCHAR(100),
		comment    TEXT,
		color      VARCHAR(50) DEFAULT '#3b82f6',
		latitude   DOUBLE PRECISION NOT NULL,
		longitude  DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	if _, err := r.db.Exec(query); err != nil {
		log.Fatalf("locationsテーブル作成失敗: %v", err)
	}
}

func (r *PostgresRepository) GetAll() ([]location.Location, error) {
	rows, err := r.db.Query(`SELECT id, place_id, title, address, prefecture, category, comment, color, latitude, longitude, created_at FROM locations ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []location.Location
	for rows.Next() {
		var l location.Location
		if err := rows.Scan(&l.ID, &l.PlaceID, &l.Title, &l.Address, &l.Prefecture, &l.Category, &l.Comment, &l.Color, &l.Latitude, &l.Longitude, &l.CreatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	if locations == nil {
		locations = []location.Location{}
	}
	return locations, nil
}

func (r *PostgresRepository) Create(placeID, title, address, prefecture, category, comment, color string, lat, lng float64) (location.Location, error) {
	var l location.Location
	err := r.db.QueryRow(
		`INSERT INTO locations (place_id, title, address, prefecture, category, comment, color, latitude, longitude) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		 RETURNING id, place_id, title, address, prefecture, category, comment, color, latitude, longitude, created_at`,
		placeID, title, address, prefecture, category, comment, color, lat, lng,
	).Scan(&l.ID, &l.PlaceID, &l.Title, &l.Address, &l.Prefecture, &l.Category, &l.Comment, &l.Color, &l.Latitude, &l.Longitude, &l.CreatedAt)
	return l, err
}
