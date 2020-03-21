package db

import (
	"context"
)

type UserCollection interface {
	Insert(ctx context.Context, user *User) error
	InsertMany(ctx context.Context, users []*User) error
}
