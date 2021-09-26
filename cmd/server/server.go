package server

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	v1 "go_grpc/api/server/v1"
	"go_grpc/server"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	GRPCPORT     string
	DATADBHOST   string
	DATADBUSER   string
	DATADBPWD    string
	DATADBSCHEMA string
}

func RunServer() error {
	ctx := context.Background()
	var cfg Config
	flag.StringVar(&cfg.GRPCPORT, "grpc-port", ":9000", "grpc port bind")
	flag.StringVar(&cfg.DATADBHOST, "host", "127.0.0.1", "grpc host")
	flag.StringVar(&cfg.DATADBUSER, "user", "root", "root")
	flag.StringVar(&cfg.DATADBPWD, "password", "sds3229339", "password")
	flag.StringVar(&cfg.DATADBSCHEMA, "schema", "gin_grpc", "go_grpc")

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&%s", cfg.DATADBUSER, cfg.DATADBPWD, cfg.DATADBHOST, cfg.DATADBSCHEMA, param)
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	v1api := v1.NewToDoServer(db)
	return server.RunServer(ctx, v1api, cfg.GRPCPORT)
}
