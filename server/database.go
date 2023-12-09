package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func CreateTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS exchanges(" +
		"id TEXT PRIMARY KEY," +
		"code TEXT," +
		"code_in TEXT," +
		"name TEXT," +
		"high TEXT," +
		"low TEXT," +
		"var_bid TEXT," +
		"pct_change TEXT," +
		"bid TEXT," +
		"ask TEXT," +
		"timestamp TEXT," +
		"create_date TEXT)")

	if err != nil {
		log.Println("Erro ao criar tabela...")
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("Erro ao criar tabela...")
	}
}

func NewExchangeRate(db *sql.DB, ex *Usdbrl) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	stmt, err := db.Prepare("insert into exchanges(" +
		"id, code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) values (" +
		"? ,? , ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx,
		uuid.New(), ex.Code, ex.Codein, ex.Name, ex.High, ex.Low, ex.VarBid,
		ex.PctChange, ex.Bid, ex.Ask, ex.Timestamp, ex.CreateDate)
	if err != nil {
		log.Println("Erro - Tempo excedido na chamada do BD")
		return err
	}
	return nil
}
