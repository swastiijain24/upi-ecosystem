package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(connString string) *pgxpool.Pool{
	pool, err := pgxpool.New(context.Background(), connString) //pool of connections for every req we dont create a new conn instead we have a pool of conn we take one connect run and return the conn, so connections are reused instead of recreated 
	if err!= nil{
		log.Fatal("unable to connect to database:", err)
	}
	err = pool.Ping(context.Background())
	if err!=nil{
		log.Fatal("database ping failed", err)
	}
	log.Println("postgres connected")

	return pool 

}