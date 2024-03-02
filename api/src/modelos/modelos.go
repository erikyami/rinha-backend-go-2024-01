package modelos

import "time"

// type Cliente struct {
// 	ID     uint64 `json:"id"`
// 	Nome   string `json:"nome"`
// 	Limite uint64 `json:"limite"`
// 	Saldo  int64  `json:"saldo"`
// }

type Transacao struct {
	Valor       uint64    `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Extrato struct {
	Saldo             Saldo       `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int    `json:"total"`
	DataExtrato string `json:"data_extrato"`
	Limite      int    `json:"limite"`
}

type Req_transacao struct {
	Valor     uint64 `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type Resp_transacao struct {
	Limite uint64 `json:"limite"`
	Saldo  int64  `json:"saldo"`
}
