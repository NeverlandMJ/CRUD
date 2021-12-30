package security

import (
	
	"context"
	
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//"golang.org/x/crypto/bcrypt"
)

var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInternal = errors.New("internal error")
var ErrExpired = errors.New("token is expired")
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

// Auth ....
type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

func (s *Service) AuthenticateCusomer(
	ctx context.Context, 
	token string,
) (id int64, err error){
	expiredTime := time.Now()
	nowTimeInSec := expiredTime.UnixNano()
	err = s.pool.QueryRow(ctx, `
		select customer_id from customers_tokens where token = $1
	`, token).Scan(&id)
	if err == pgx.ErrNoRows{
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, ErrInternal
	}
	if nowTimeInSec > expiredTime.UnixNano() {
		return -1, ErrExpired
	}
	return id, nil
}