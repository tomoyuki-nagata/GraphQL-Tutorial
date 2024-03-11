package graph

import (
	"context"
	"errors"
	"graphql-tutorial/internal"
	"graphql-tutorial/middlewares/auth"

	"github.com/99designs/gqlgen/graphql"
)

var Directive internal.DirectiveRoot = internal.DirectiveRoot{
	IsAuthenticated: IsAuthenticated,
}

func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if _, ok := auth.GetUserName(ctx); !ok {
		return nil, errors.New("not authenticated")
	}
	return next(ctx)
}
