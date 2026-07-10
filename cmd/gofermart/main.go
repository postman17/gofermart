package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	clients "github.com/postman17/gofermart/internal/client"
	handlers "github.com/postman17/gofermart/internal/handlers"
	middlewares "github.com/postman17/gofermart/internal/middlewares"
	repo "github.com/postman17/gofermart/internal/repository"
)

func newDBRepository(ctx context.Context, config Config) (repo.DBRepository, *sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, nil, err
	}

	return repo.NewDBStorage(ctx, db), db, nil
}

func main() {
	config := parseFlags()

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	client := clients.NewAccrualSystemClient(appCtx, config.AccrualSystemAddress)
	repository, db, err := newDBRepository(appCtx, config)
	if err != nil {
		panic(fmt.Sprintf("failed to init repository: %v", err))
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Post("/api/user/register", handlers.RegisterUser(appCtx, repository))
	r.Post("/api/user/login", handlers.LoginUser(appCtx, repository))
	r.Post("/api/user/orders", middlewares.AuthMiddleware(appCtx, repository)(handlers.AddOrder(appCtx, repository, client)).ServeHTTP)
	r.Get("/api/user/orders", middlewares.AuthMiddleware(appCtx, repository)(handlers.GetOrders(appCtx, repository)).ServeHTTP)
	r.Get("/api/user/balance", middlewares.AuthMiddleware(appCtx, repository)(handlers.GetBalance(appCtx, repository)).ServeHTTP)
	r.Post("/api/user/balance/withdraw", middlewares.AuthMiddleware(appCtx, repository)(handlers.Withdraw(appCtx, repository)).ServeHTTP)
	r.Get("/api/user/withdrawals", middlewares.AuthMiddleware(appCtx, repository)(handlers.UserWithdrawals(appCtx, repository)).ServeHTTP)

	log.Printf("starting server on %s", config.RunAddr)
	if err := http.ListenAndServe(config.RunAddr, r); err != nil {
		panic(fmt.Sprintf("server failed: %v", err))
	}
}
