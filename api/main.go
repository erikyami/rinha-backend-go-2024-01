package main

import (
	"api/src/banco"
	"api/src/config"
	"api/src/modelos"
	"api/src/respostas"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.Carregar()

	db, err := banco.Conectar(context.Background())
	if err != nil {
		fmt.Println("Erro ao conectar no banco")
		log.Fatal(err)
	}
	defer db.Close()

	rota := mux.NewRouter()

	fmt.Printf("Escutando na porta %d\n", config.API_PORT)
	rota.HandleFunc("/clientes/{id}/transacoes", insereTransacao(db)).Methods("POST")
	rota.HandleFunc("/clientes/{id}/extrato", getExtrato(db)).Methods("GET")
	rota.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.API_PORT), rota))

}

func getExtrato(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, erro := strconv.ParseInt(vars["id"], 10, 64)
		if erro != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}

		if ID < 1 || ID > 5 {
			erro := errors.New("cliente não encontrado")
			respostas.Erro(w, http.StatusNotFound, erro)
			return
		}

		var saldoTotal, limite int
		var extrato modelos.Extrato
		dataExtrato := time.Now().Format("2006-01-02T15:04:05.999999Z")

		erro = db.QueryRow(context.Background(), "SELECT saldo, limite FROM clientes WHERE id = $1", ID).Scan(&saldoTotal, &limite)
		if erro != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}

		// Consulta para obter as últimas transações do cliente com ID especifico
		transacoesRows, erro := db.Query(context.Background(), "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10", ID)
		if erro != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}
		defer transacoesRows.Close()

		// Variável para armazenar as últimas transações do cliente
		var ultimasTransacoes []modelos.Transacao

		for transacoesRows.Next() {
			var transacao modelos.Transacao
			erro := transacoesRows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
			if erro != nil {
				respostas.Erro(w, http.StatusUnprocessableEntity, erro)
				return
			}
			ultimasTransacoes = append(ultimasTransacoes, transacao)
		}
		// Preencher o extrato
		extrato = modelos.Extrato{
			Saldo: modelos.Saldo{
				Total:       saldoTotal,
				DataExtrato: dataExtrato,
				Limite:      limite,
			},
			UltimasTransacoes: ultimasTransacoes,
		}

		respostas.JSON(w, http.StatusOK, extrato)

	}
}

func insereTransacao(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		ID, erro := strconv.ParseInt(vars["id"], 10, 64)
		if erro != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}

		if ID < 1 || ID > 5 {
			erro := errors.New("cliente não encontrado")
			respostas.Erro(w, http.StatusNotFound, erro)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var req_transacao modelos.Req_transacao

		err := decoder.Decode(&req_transacao)
		if err != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, err)
			return
		}

		if req_transacao.Tipo != "c" && req_transacao.Tipo != "d" || len(req_transacao.Descricao) < 1 || len(req_transacao.Descricao) > 10 {
			erro := errors.New("requisição Inválida")
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}

		var resp_transacao modelos.Resp_transacao

		erro = db.QueryRow(context.Background(),
			"SELECT v_saldo, v_limite FROM create_transaction_func($1, $2, $3, $4)",
			ID, req_transacao.Valor, req_transacao.Tipo, req_transacao.Descricao).Scan(
			&resp_transacao.Saldo,
			&resp_transacao.Limite)
		if erro != nil {
			respostas.Erro(w, http.StatusUnprocessableEntity, erro)
			return
		}

		respostas.JSON(w, http.StatusOK, resp_transacao)

	}
}
