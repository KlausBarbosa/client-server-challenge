package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"time"
)

type Cotacao struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}
type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
type UsdResp struct {
	Dolar string `json:"dolar"`
}

func main() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	CreateTable(db)

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := UsdBrlPrice()
	if err != nil {
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
	w.Header().Set("Content-Type", "application.json")
	w.WriteHeader(http.StatusOK)
	cotacaoResp := UsdResp{Dolar: cotacao.Usdbrl.Bid}
	err = json.NewEncoder(w).Encode(cotacaoResp)

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var exg Usdbrl = cotacao.Usdbrl
	NewExchangeRate(db, &exg)
	if err != nil {
		log.Println("Erro - Tempo excedido ao salvar no BD")
	}
}

func UsdBrlPrice() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro - Tempo excedido na chamada da API")
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

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
