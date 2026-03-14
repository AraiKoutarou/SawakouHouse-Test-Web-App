// internal/domain/location/repository.go: インターフェースの更新。
package location

type Repository interface {
	GetAll() ([]Location, error)
	Create(placeID, title, address, prefecture, category, comment, color string, lat, lng float64) (Location, error)
}
