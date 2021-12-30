package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	//"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


var ErrNotFound = errors.New("item not found")
var ErrNoSuchUser = errors.New("no such user")
var ErrInvalididPassword = errors.New("invalid password")
var ErrInternal = errors.New("internal error")
type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

//Customer ...
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Password string    `json:"password"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) TokenForCustomer(
	ctx context.Context,
	phone string,
	password string,
) (token string, err error){
	var hash string
	var id int64
	err = s.pool.QueryRow(ctx, `
		select id, password from customers where phone = $1
	`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}

	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalididPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)

	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `
		insert into customer_tokens(token, customer_id) values ($1, $2)
	`, token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}

// SaveCustomer save customer with pass from JSON
func (s *Service) SaveCustomer(ctx context.Context, item *Customer) (*Customer, error) {

	if item.ID == 0 {

		hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		log.Print(hex.EncodeToString(hash))

		res := &Customer{}
		err = s.pool.QueryRow(ctx, `
				INSERT INTO customers(name, phone, password) VALUES ($1, $2, $3)
				RETURNING id, name, phone, password, active, created
			`, item.Name, item.Phone, hash).Scan(&res.ID, &res.Name, &res.Phone, &res.Password, &res.Active, &res.Created)

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		return res, nil
	}

	
	return nil, ErrInternal
}

//All ....
func (s *Service) All(ctx context.Context) (cs []*Customer, err error) {

	sqlStatement := `select * from customers`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//AllActive ....
func (s *Service) AllActive(ctx context.Context) (cs []*Customer, err error) {

	sqlStatement := `select * from customers where active=true`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `select * from customers where id=$1`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//ChangeActive ...
func (s *Service) ChangeActive(ctx context.Context, id int64, active bool) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `update customers set active=$2 where id=$1 returning *`
	err := s.pool.QueryRow(ctx, sqlStatement, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Delete ...
func (s *Service) Delete(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `delete from customers  where id=$1 returning *`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Save ...
func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
		sqlStatement := `insert into customers(name, phone) values($1, $2) returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	} else {
		sqlStatement := `update customers set name=$1, phone=$2 where id=$3 returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

