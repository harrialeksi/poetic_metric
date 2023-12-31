package user

import (
	"github.com/th0th/poeticmetric/backend/pkg/depot"
	"github.com/th0th/poeticmetric/backend/pkg/model"
)

func Read(dp *depot.Depot, id uint64) (*User, error) {
	user := &User{}

	err := dp.Postgres().
		Model(&model.User{}).
		Where("id = ?", id).
		First(user).
		Error

	return user, err
}
