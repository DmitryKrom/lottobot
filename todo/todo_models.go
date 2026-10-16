package todo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ToDo struct {
	Id        int64
	Text      string
	CreatedAt string
	Deadline  *string
	Done      bool
	UserId    int
}

type TgUser struct {
	Id   int
	Name string
}

type ToDos struct {
	Todos []ToDo
}

func New(text string, now string, userId int) *ToDo {
	return &ToDo{
		Text:      text,
		CreatedAt: now,
		Done:      false,
		UserId:    userId,
	}
}

func (t *ToDo) AddDone() {
	t.Done = true
}

func (t *ToDo) AddDeadline(deadline string) {
	t.Deadline = &deadline
}

func (t *ToDo) AddUserId(id int) {
	t.UserId = id
}

func (t *ToDo) IdUser() int {
	return t.UserId
}

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase() *Database {
	fmt.Println("----- NewDatabase -----")
	return &Database{}
}

func (d *Database) Connect(conn string) {
	fmt.Println("----- Connect -----")
	var err error
	d.pool, err = pgxpool.New(context.Background(), conn)
	if err != nil {
		fmt.Println("[ERROR] pgxpool.New() : ", err)
	}
	// defer d.pool.Close()
	err = d.pool.Ping(context.Background())
	if err != nil {
		fmt.Println("[ERROR] Ping() : ", err)
		return
	}
	fmt.Println("Ping successful !!!")
}
