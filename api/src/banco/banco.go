package banco

import (
	"api/src/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Conectar(ctx context.Context) (*pgxpool.Pool, error) {

	StringDeConexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DB_HOST,
		config.DB_USER,
		config.DB_SENHA,
		config.DB_NOME,
		config.DB_PORT,
	)

	config, err := pgxpool.ParseConfig(StringDeConexao)
	if err != nil {
		return nil, err
	}
	maxConns, minConns := 30, 5
	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)

	return pgxpool.NewWithConfig(ctx, config)

}
