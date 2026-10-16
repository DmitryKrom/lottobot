package todo

import "lottobot/lotto"

func GetTodo() {

	db := lotto.NewDatabase()

	db.Connect("postgres://postgres:postgres@localhost:5432/todo")
	defer db.Pool.Close()
}
