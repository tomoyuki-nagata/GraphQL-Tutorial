package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"graphql-tutorial/graph"
	"graphql-tutorial/graph/services"
	"graphql-tutorial/internal"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"

	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultPort = "8080"
	dbFile      = "./mygraphql.db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbFile))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQLBoilerによって発行されるSQLクエリをログ出力させるデバッグオプション
	// boil.DebugMode = true

	service := services.New(db)

	srv := handler.NewDefaultServer(internal.NewExecutableSchema(internal.Config{
		Resolvers: &graph.Resolver{
			Srv:     service,
			Loaders: graph.NewLoaders(service),
		},
		Complexity: graph.ComplexityConfig(),
	}))
	srv.Use(extension.FixedComplexityLimit(30))

	// クライアントからリクエストを受け取ったときに最初に呼ばれる
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		log.Println("before OperationHandler")
		ctx = context.WithValue(ctx, "traceId", "random traceId")
		oc := graphql.GetOperationContext(ctx)
		log.Println(fmt.Printf("%s: %s", ctx.Value("traceId"), oc.RawQuery))
		res := next(ctx)
		defer log.Println("after OperationHandler")
		return res
	})
	// クライアントに返すレスポンスを作成するという段階の前後処理を担う
	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		log.Println("before ResponseHandler")
		res := next(ctx)
		defer log.Println("after ResponseHandler")
		return res
	})
	// レスポンスデータ全体を作成するルートリゾルバの実行前後に処理を挿入するミドルウェア
	srv.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
		log.Println("before RootResolver")
		res := next(ctx)
		defer func() {
			var b bytes.Buffer
			res.MarshalGQL(&b)
			log.Println("after RootResolver", b.String())
		}()
		return res
	})
	// レスポンスに含めるjsonフィールドを1つ作る処理の前後にロジックを組み込むためのミドルウェア
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		log.Println("before Resolver")
		res, err = next(ctx)
		defer log.Println("after Resolver", res)
		return
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
