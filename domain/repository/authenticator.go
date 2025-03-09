//go:generate mockgen -source=$GOFILE -destination=../../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package repository

import "context"

type Authenticator interface {
	Login(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, email, password string) (string, error)
	SignUpConfirm(ctx context.Context, email, confirmationCode string) error
	VerifyToken(ctx context.Context, token string) (string, error)
}
