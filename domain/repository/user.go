//go:generate mockgen -source=$GOFILE -destination=../../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package repository

import (
	"context"
	"share-basket-auth-service/domain/model"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (model.User, error)
}
