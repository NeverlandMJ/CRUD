package main 

import (
	"context"
	//"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/NeverlandMJ/CRUD/cmd/app"
	"github.com/NeverlandMJ/CRUD/pkg/customers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/dig"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host, port, dsn string) (err error){
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error){
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customers.NewService,
		
		func(server *app.Server) *http.Server{
			return &http.Server{
				Addr: net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps{
		err = container.Provide(dep)
		if err != nil{
			return err
		}
	}

	err = container.Invoke(func(server *app.Server){
		server.Init()
	})

	if err != nil{
		return err
	}
	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}

// func execute(host, port, dsn string) (err error){
// 	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
// 	pool, err := pgxpool.Connect(connectCtx, dsn)
// 	if err != nil{
// 		log.Println(err)
// 		return
// 	}
// 	defer pool.Close()

// 	mux := http.NewServeMux()
// 	customersSvc := customers.NewService(pool)
// 	server := app.NewServer(mux, customersSvc)
// 	server.Init()

// 	srv := &http.Server{
// 		Addr: net.JoinHostPort(host, port),
// 		Handler: server,
// 	}
// 	return srv.ListenAndServe()
// }



// func execute(host, port, dsn string) (err error) {
// 	db, err := sql.Open("pgx", dsn)
// 	if err != nil {
// 		return err
// 	}
// 	defer func ()  {
// 		if cerr := db.Close(); cerr != nil {
// 			if err == nil{
// 				err = cerr
// 				return
// 			}
// 			log.Println(err)
// 		}	
// 	}()
	
// 	mux := http.NewServeMux()
// 	customersSvc := customers.NewService(db)
// 	server := app.NewServer(mux, customersSvc)
// 	server.Init()
// 	srv := &http.Server{
// 		Addr: 	net.JoinHostPort(host, port),
// 		Handler: 	server,
// 	}

	
// 	return srv.ListenAndServe()
// }