package main

import (
	"context"

	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	pgxlog := logrusadapter.NewLogger(log)
	conn, err := pgx.Connect(context.Background(), "postgresql://sqlserverip:5432/gosql?user=postgres&password=")
	if err != nil {
		pgxlog.Log(pgx.LogLevelError, "unable to connect", nil)
	}
	defer conn.Close(context.Background())
	var name string
	var annoyance int64
	err = conn.QueryRow(context.Background(), "select name, ability_to_annoy_people from gosql_test1 where name=$1", "annoying_persons_name").Scan(&name, &annoyance)
	if err != nil {
		pgxlog.Log(pgx.LogLevelError, "error executing query", nil)
	}
	log.Info(name, annoyance)
	m := make(map[string]interface{})
	m[name] = annoyance
	pgxlog.Log(pgx.LogLevelInfo, "got values", m)
}
