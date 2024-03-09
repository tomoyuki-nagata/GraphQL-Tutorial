package services

import (
	"context"
	"graphql-tutorial/graph/model"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Services interface {
	UserService
	// issueテーブルを扱うIssueServiceなど、他のサービスインターフェースができたらそれらを追加していく
}

type UserService interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

type services struct {
	*userService
	// issueテーブルを扱うissueServiceなど、他のサービス構造体ができたらフィールドを追加していく
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		userService: &userService{exec: exec},
	}
}
