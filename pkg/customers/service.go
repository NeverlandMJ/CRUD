package customers

import (
	"context"
	//"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")


type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `select * from customers where id=$1`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrNotFound
		}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

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

func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error)  {
	item := &Customer{}

	if customer.ID == 0{
		err = s.pool.QueryRow(ctx, `
			INSERT INTO customers(id, name, phone) VALUES ($1, $2, $3) RETURNING *
		`, customer.ID, customer.Name, customer.Phone).Scan(
			&item.ID, 
			&item.Name, 
			&item.Phone, 
			&item.Active, 
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `
			UPDATE customers SET name = $1, phone = $2 WHERE  id = $3 RETURNING *
		`, customer.Name, customer.Phone, customer.ID).Scan(
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

func (s *Service) RemoveById(ctx context.Context, id int64) (*Customer,  error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
	  	DELETE customers WHERE id = $1 RETURNING *
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	
	if errors.Is(err, pgx.ErrNoRows){
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) BlockById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = false WHERE id = $1 RETURNING *
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) UnblockById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = true WHERE id = $1 RETURNING *
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
