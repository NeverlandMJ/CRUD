package security

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)


type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}
type Manager struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Salary      int64     `json:"salary"`
	Plan        int64     `json:"plan"`
	BossID      int64     `json:"boss_id"`
	Departament string    `json:"departament"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func (s *Service) Auth(login, password string) (ok bool) {
	ctx := context.Background()
	temPass := ""
	err := s.pool.QueryRow(ctx, `
	select password from managers where login = $1
	`, login).Scan(&temPass)
	if errors.Is(err, pgx.ErrNoRows) {
		return false
	}

	if err != nil {
		log.Print(err)
		return false
	}
	if temPass == password {
		return true 
	}

	return false
}