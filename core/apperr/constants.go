package apperr

import "errors"

var (
	ErrInvalidData          = errors.New("入力値が正しくありません")
	ErrDataNotFound         = errors.New("データが見つかりません")
	ErrDuplicatedKey        = errors.New("重複したキーが存在します")
	ErrInvalidCredentials   = errors.New("ユーザー名またはパスワードが正しくありません")
	ErrInvalidToken         = errors.New("無効なトークンです")
	ErrTokenExpired         = errors.New("トークンの有効期限が切れています")
	ErrUnauthenticated      = errors.New("認証に失敗しました")
	ErrExpiredCodeException = errors.New("検証コードの有効期限が切れています")
)
