package security

import (
	"time"

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
	Phone       string    `json:"phone"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func (s *Service) Auth(login, password string) (ok bool) {

	return true
}