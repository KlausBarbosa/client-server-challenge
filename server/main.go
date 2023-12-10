package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
