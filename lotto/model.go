package lotto

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBOperation interface {
	Insert(date string, pool *pgxpool.Pool) error
	Select(date string, pool *pgxpool.Pool) error
	Update(date string, pool *pgxpool.Pool) error
	Delete(date string, pool *pgxpool.Pool) error
}

type EurojackpotZahlen struct {
	firstNum  string
	secondNum string
	thirdNum  string
	fourthNum string
	fifthNum  string
	supOne    string
	supTwo    string
}

func NewJP() DBOperation {
	return &EurojackpotZahlen{}
}
func (eu *EurojackpotZahlen) Insert(date string, pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "INSERT INTO eurojackpot (one, two, three, four, five, sup1, sup2, datum) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT DO NOTHING", &eu.firstNum, &eu.secondNum, &eu.thirdNum, &eu.fourthNum, &eu.fifthNum, &eu.supOne, &eu.supTwo, date)
	if err != nil {
		fmt.Println("[ERROR] pool.Exec() : ", err)
		return err
	}
	fmt.Printf("%s : Zahlen geschrieben\n", date)
	return nil
}
func (eu *EurojackpotZahlen) Select(date string, pool *pgxpool.Pool) error {
	err := pool.QueryRow(context.Background(), "SELECT one, two, three, four, five, sup1, sup2 FROM eurojackpot WHERE datum=$1", date).Scan(&eu.firstNum, &eu.secondNum, &eu.thirdNum, &eu.fourthNum, &eu.fifthNum, &eu.supOne, &eu.supTwo)
	if err != nil {
		fmt.Println("[ERROR] pool.QueryRow().Scan() : ", err)
		return err
	}
	fmt.Println("Zahlen aus der DB geholt!!! ")
	return nil
}
func (eu *EurojackpotZahlen) Update(date string, pool *pgxpool.Pool) error { return nil }
func (eu *EurojackpotZahlen) Delete(date string, pool *pgxpool.Pool) error { return nil }

type Lotto6aus49Zahlen struct {
	firstNum  string
	secondNum string
	thirdNum  string
	fourthNum string
	fifthNum  string
	sixthNum  string
	super     string
}

func New6aus49() DBOperation {
	return &Lotto6aus49Zahlen{}
}
func (l *Lotto6aus49Zahlen) Insert(date string, pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "INSERT INTO lotto6aus49 (one, two, three, four, five, six, sup1, datum) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT DO NOTHING", &l.firstNum, &l.secondNum, &l.thirdNum, &l.fourthNum, &l.fifthNum, &l.sixthNum, &l.super, date)
	if err != nil {
		fmt.Println("[ERROR] pool.Exec() : ", err)
		return err
	}
	fmt.Printf("%s : Zahlen geschrieben\n", date)
	return nil
}
func (l *Lotto6aus49Zahlen) Select(date string, pool *pgxpool.Pool) error {
	err := pool.QueryRow(context.Background(), "SELECT one, two, three, four, five, six, sup1 FROM lotto6aus49 WHERE datum=$1", date).Scan(&l.firstNum, &l.secondNum, &l.thirdNum, &l.fourthNum, &l.fifthNum, &l.sixthNum, &l.super)
	if err != nil {
		fmt.Println("[ERROR] pool.QueryRow().Scan() : ", err)
		return err
	}
	return nil
}
func (l *Lotto6aus49Zahlen) Update(date string, pool *pgxpool.Pool) error { return nil }
func (l *Lotto6aus49Zahlen) Delete(date string, pool *pgxpool.Pool) error { return nil }

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase() *Database {
	fmt.Println("----- NewDatabase -----")
	return &Database{}
}

func (d *Database) Connect(conn string) {
	fmt.Println("----- Connect -----")
	var err error
	d.Pool, err = pgxpool.New(context.Background(), conn)
	if err != nil {
		fmt.Println("[ERROR] pgxpool.New() : ", err)
	}
	// defer d.pool.Close()
	err = d.Pool.Ping(context.Background())
	if err != nil {
		fmt.Println("[ERROR] Ping() : ", err)
		return
	}
	fmt.Println("Ping successful !!!")
}
