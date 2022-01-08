package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dsn string
	}
}

type AppStatus struct {
	Status string `json:"status"`
	Environment string `json:"environment"`
	Version string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	// flagでconfigのプロパティを初期化する
	// 引数は変数のポインタ(メモリのアドレス値)、フラグの名前、デフォルト値、使い方の説明
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://localhost/go_movies?sslmode=disable", "Postgres connection string")
	// flag.StringVar(&cfg.db.dsn, "dsn", "postgres://tcs@localhost/go_movies?sslmode=disable", "Postgres connection string")
	// flag.StringVar(&cfg.jwt.secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	// 引数のフラグを解析しcfgにバインドする
	flag.Parse()

	// Loggerオブジェクトを生成して出力フォーマットを設定する
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// DBと接続する
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// main()がreturnするときにDBとの接続を閉じる
	defer db.Close()

	// アプリケーションの設定をする(参照渡しをおこなうためにポインタを使用)
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	// サーバー設定をカスタマイズする
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: 10 * time.Minute,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	// サーバをlistenする
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// DBへアクセスする(接続はまだ確立されない)
	db, err := sql.Open("postgres", cfg.db.dsn)
	// エラー処理
	if err != nil {
		return nil, err
	}

	// 5sでタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	// openDB()がreturnするまで実行されない
	defer cancel()

	// DBとの接続を検証する
	err = db.PingContext(ctx)
	//エラー処理
	if err != nil {
		return nil, err
	}

	return db, nil
}