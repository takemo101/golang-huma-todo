package repository

import (
	"math/rand/v2"
)

// Todo モデル
type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// 作成用のTodoモデル
type TodoForCreate struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// In-memory で Todo を管理するためのスライス
var todos []Todo = []Todo{
	{ID: "first", Title: "test1", Completed: false},
	{ID: "second", Title: "test2", Completed: true},
}

// Todo の一覧を返す
func GetTodos() []Todo {
	return todos
}

// Todo を取得
func GetTodoById(id string) (Todo, bool) {
	for _, t := range todos {
		if t.ID == id {
			return t, true
		}
	}

	return Todo{}, false
}

// Todo を作成
func CreateTodo(todo TodoForCreate) Todo {
	created := Todo{
		ID:        createRandomID(),
		Title:     todo.Title,
		Completed: todo.Completed,
	}

	todos = append(todos, created)

	return created
}

// ランダムな文字列IDを生成
func createRandomID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}
