package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL tidak di-set di .env")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Gagal koneksi dalam database:", err)
	}

	//tes koneksi
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	DB = pool
	fmt.Println("Berhasil koneksi ke PostgreSQL")
}
