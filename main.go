package main

import (
	"app/repository"
	"app/shared"
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-gonic/gin"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// TodoスキーマのBodyデータ定義
type TodoBody struct {
	ID        string `json:"id" example:"1" doc:"TodoのID"`
	Title     string `json:"title" example:"XXXに連絡する" doc:"Todoのタイトル"`
	Completed bool   `json:"completed" example:"false" doc:"Todoの完了状態"`
}

// Todo詳細取得のレスポンスデータ定義
type TodoOutput struct {
	Body struct {
		Todo TodoBody `json:"todo" doc:"Todoの詳細"`
	}
}

// Todo一覧取得のレスポンスデータ定義
type TodosOutput struct {
	Body struct {
		Todos []TodoBody `json:"todos" doc:"Todoの一覧"`
	}
}

// Todo作成のリクエストデータ定義
type CreateTodoInput struct {
	Body struct {
		Title     string `json:"title" minLength:"1" maxLength:"100" example:"XXXに連絡する" doc:"Todoのタイトル"`
		Completed bool   `json:"completed" example:"false" doc:"Todoの完了状態"`
	}
}

// APIのルーティング設定
func setupRoutes(api huma.API) {

	api.UseMiddleware(
		func(ctx huma.Context, next func(ctx huma.Context)) {

			token := ctx.Query("token")

			// トークンが一致しない場合は 401 を返す
			if shared.IsInvalidToken(token) {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("Invalid token"))
				return
			}

			next(ctx)
		},
	)

	// 一覧取得
	huma.Register(api, huma.Operation{
		OperationID: "getTodos",
		Method:      http.MethodGet,
		Path:        "/api/v1/todos",
		Summary:     "Todo一覧を取得",
		Tags:        []string{"todos"},
		Security: []map[string][]string{
			{"queryToken": {}},
		},
	}, func(ctx context.Context, input *struct{}) (*TodosOutput, error) {
		res := &TodosOutput{}

		todos := repository.GetTodos()

		for _, t := range todos {
			res.Body.Todos = append(res.Body.Todos, TodoBody{
				ID:        t.ID,
				Title:     t.Title,
				Completed: t.Completed,
			})
		}

		return res, nil
	})

	// 詳細取得
	huma.Register(api, huma.Operation{
		OperationID: "getTodo",
		Method:      http.MethodGet,
		Path:        "/api/v1/todos/{id}",
		Summary:     "Todo詳細を取得",
		Tags:        []string{"todos"},
		Security: []map[string][]string{
			{"queryToken": {}},
		},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" required:"true" doc:"TodoのID"`
	}) (*TodoOutput, error) {
		res := &TodoOutput{}

		t, ok := repository.GetTodoById(input.ID)
		if !ok {
			return nil, huma.Error404NotFound("Todo not found")
		}

		res.Body.Todo = TodoBody{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
		}

		return res, nil
	})

	// 作成
	huma.Register(api, huma.Operation{
		OperationID:   "createTodo",
		Method:        http.MethodPost,
		Path:          "/api/v1/todos",
		Summary:       "Todoを作成",
		Tags:          []string{"todos"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"queryToken": {}},
		},
	}, func(ctx context.Context, input *CreateTodoInput) (*TodoOutput, error) {
		res := &TodoOutput{}

		todo := repository.CreateTodo(repository.TodoForCreate{
			Title:     input.Body.Title,
			Completed: input.Body.Completed,
		})

		res.Body.Todo = TodoBody{
			ID:        todo.ID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}

		return res, nil
	})
}

// Options はコマンドライン引数を格納するための構造体
type Options struct {
	Port     int    `help:"Port to listen on" short:"p" default:"8888"`
	Hostname string `help:"Hostname to listen on" short:"n" default:"localhost"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		// Humaginアダプターを利用することで、HumaでGinを利用したAPIを作成できる
		engine := gin.Default()

		config := huma.DefaultConfig("Todo API", "1.0.0")

		// セキュリティスキーム（クエリパラメータでトークンを検証）を定義
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"queryToken": {
				Type: "apiKey",
				In:   "query",
				Name: "token",
			},
		}

		api := humagin.New(engine, config)

		setupRoutes(api)

		// そのままGinを利用してルートを追加することもできるが、
		// このルートはHumaのAPIドキュメントには反映されない
		engine.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// サーバー起動時の処理をフックに登録
		hooks.OnStart(func() {
			engine.Run(fmt.Sprintf("%s:%d", options.Hostname, options.Port))
		})
	})

	cli.Run()
}
