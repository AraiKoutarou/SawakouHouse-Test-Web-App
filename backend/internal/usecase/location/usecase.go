// internal/usecase/location/usecase.go: パラメータ更新。
package location

import (
	"github.com/arakou0812/backend/internal/domain/location"
)

type Usecase struct {
	repo location.Repository
}

func NewUsecase(repo location.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetAll() ([]location.Location, error) {
	return u.repo.GetAll()
}

func (u *Usecase) Create(placeID, title, address, prefecture, category, comment, color string, lat, lng float64) (location.Location, error) {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return location.Location{}, location.ErrInvalidCoordinate
	}
	return u.repo.Create(placeID, title, address, prefecture, category, comment, color, lat, lng)
}
