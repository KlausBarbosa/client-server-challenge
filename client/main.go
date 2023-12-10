package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoRes struct {
	CotacaoDolar string `json:"dolar"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		//panic(err)
		log.Println("Erro requestctxt")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro - Tempo excedido na chamada da api")
	}
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)

	}
	var data CotacaoRes
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println("Erro - Tempo excedido na chamada da api")
		panic(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", data.CotacaoDolar))
	if err != nil {
		panic(err)
	}
}
