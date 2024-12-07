package main

import (
	"app/shared"
	"encoding/json"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	_, api := humatest.New(t)
	setupRoutes(api)

	t.Run("無効なトークンのテスト", func(t *testing.T) {
		res := api.Get("/api/v1/todos?token=invalid")

		assert.Equal(t, 401, res.Code)
	})

	t.Run("Todo一覧取得のテスト", func(t *testing.T) {
		res := api.Get("/api/v1/todos?token=" + shared.Token)

		assert.Equal(t, 200, res.Code)
	})

	t.Run("Todo詳細取得のテスト", func(t *testing.T) {
		res := api.Get("/api/v1/todos/first?token=" + shared.Token)

		assert.Equal(t, 200, res.Code)
	})

	t.Run("Todo作成のテスト", func(t *testing.T) {

		data := map[string]any{
			"title":     "test",
			"completed": false,
		}
		res := api.Post("/api/v1/todos?token="+shared.Token, data)

		assert.Equal(t, 201, res.Code)
		var body struct {
			Todo TodoBody `json:"todo"`
		}
		// レスポンスの body をパースして Todo を取得
		json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, data["title"], body.Todo.Title)
		assert.Equal(t, data["completed"], body.Todo.Completed)
	})
}
