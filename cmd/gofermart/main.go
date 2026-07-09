package gofermart

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	handlers "github.com/postman17/gofermart/internal/handlers"
	repo "github.com/postman17/gofermart/internal/repository"
)

func newDBRepository(ctx context.Context, config Config) (repo.DBRepository, *sql.DB, error) {
	db, err := dbconfig.Open(ctx, config.DatabaseURI)
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

	client := NewAccrualSystemClient(appCtx, config.AccrualSystemAddress)
	repository := newDBRepository(appCtx, config)

	r := chi.NewRouter()
	r.Post("/api/user/register", handlers.RegisterUser(appCtx, repository))
	r.Post("/api/user/login", handlers.LoginUser(appCtx, repository))
	r.Post("/api/user/orders", AuthMiddleware(appCtx, repository)(handlers.AddOrder(appCtx, repository, client)))
	r.Get("/api/user/orders", AuthMiddleware(appCtx, repository)(handlers.GetOrders(appCtx, repository)))
	r.Get("/api/user/balance", AuthMiddleware(appCtx, repository)(handlers.GetBalance(appCtx, repository)))
	r.Post("/api/user/balance/withdraw", AuthMiddleware(appCtx, repository)(handlers.Withdraw(appCtx, repository)))
	r.Get("/api/user/withdrawals", AuthMiddleware(appCtx, repository)(handlers.UserWithdrawals(appCtx, repository)))
}
