// internal/domain/location/entity.go: 都道府県情報を追加。
package location

import (
	"errors"
	"time"
)

type Location struct {
	ID         int       `json:"id"`
	PlaceID    string    `json:"place_id"`
	Title      string    `json:"title"`
	Address    string    `json:"address"`
	Prefecture string    `json:"prefecture"` // 追加: 都道府県 (例: 東京都)
	Category   string    `json:"category"`
	Comment    string    `json:"comment"`
	Color      string    `json:"color"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	CreatedAt  time.Time `json:"created_at"`
}

var (
	ErrNotFound = errors.New("指定された地点は見つかりません")
	ErrInvalidCoordinate = errors.New("座標が不正です")
)
