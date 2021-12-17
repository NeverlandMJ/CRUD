package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")


type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID 		int64 	`json: "id"`
	Name 		string 	`json: "name"`
	Phone 	string 	`json: "phone"`
	Active 	bool 		`json: "active"`
	Created 	time.Time 	`json: "created"`
}

func (s *Service) ById(ctx context.Context, id int64) (*Customer, error){
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows){
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

	rows, err := s.db.QueryContext(ctx, sqlStatement)
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

func (s *Service) GetAllActive(ctx context.Context) ([]Customer, error) {
	items := []Customer{}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active
	`)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for rows.Next() {
		item := Customer{}
		err = rows.Scan(
			&item.ID, 
			&item.Name, 
			&item.Phone, 
			&item.Active, 
			&item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, ErrNotFound
	}
	return items, nil
}

func (s *Service) Save(ctx context.Context, cutomer *Customer) (c *Customer, err error)  {
	item := &Customer{}

	if cutomer.ID == 0{
		err = s.db.QueryRowContext(ctx, `
			INSERT INTO customers(phone, name) VALUES ($1, $2) RETURNING *
		`).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	} else {
		err = s.db.QueryRowContext(ctx, `
			UPDATE customers SET name = &1, phone = &2 WHERE  id = &3 RETURNING *
		`).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

func (s *Service) RemoveById(ctx context.Context, id int64) (*Customer,  error) {
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
	  	DELET customers WHERE id = $1 RETURNING *
	`).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	
	if errors.Is(err, sql.ErrNoRows){
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

	err := s.db.QueryRowContext(ctx, `
		UPDATE customers SET active = false WHERE id = $1 RETURNING *
	`).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err == sql.ErrNoRows {
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

	err := s.db.QueryRowContext(ctx, `
		UPDATE customers SET active = true WHERE id = $1 RETURNING *
	`).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
