package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PersonRepository struct {
	db *pgxpool.Pool
}

func NewPersonRepository(db *pgxpool.Pool) entity.PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) GetAll(email, phone, firstName string, page, limit int) ([]*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name 
			  FROM persons
			  WHERE email LIKE '%'||$1||'%' AND phone LIKE '%'||$2||'%' AND first_name LIKE '%'||$3||'%'
			  ORDER BY id LIMIT $4 OFFSET $5;`
	persons := make([]*entity.Person, 0)
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	rows, err := conn.Query(context.Background(), sql, email, phone, firstName, limit, (page-1)*limit)
	for rows.Next() {
		p := new(entity.Person)
		err = rows.Scan(
			&p.ID,
			&p.Email,
			&p.Phone,
			&p.FirstName,
		)
		persons = append(persons, p)
	}
	return persons, err
}

func (r *PersonRepository) GetByID(id int) (*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name FROM persons WHERE id = $1;`
	person := new(entity.Person)
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	err = conn.QueryRow(context.Background(), sql, id).
		Scan(&person.ID,
			&person.Email,
			&person.Phone,
			&person.FirstName,
		)
	if errors.Is(err, pgx.ErrNoRows) {
		err, person = nil, nil
	}
	return person, err
}

func (r *PersonRepository) GetByEmail(email string) (*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name FROM persons WHERE email = $1;`
	person := new(entity.Person)
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	err = conn.QueryRow(context.Background(), sql, email).
		Scan(&person.ID,
			&person.Email,
			&person.Phone,
			&person.FirstName,
		)
	if errors.Is(err, pgx.ErrNoRows) {
		err, person = nil, nil
	}
	return person, err
}

func (r *PersonRepository) Create(req *entity.Person) (*entity.Person, error) {
	sql := `INSERT INTO persons (email, phone, first_name, created_at) VALUES ($1,$2,$3,$4) RETURNING id;`
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	err = conn.QueryRow(context.Background(), sql, req.Email, req.Phone, req.FirstName, getTime()).
		Scan(&req.ID)

	return req, err
}

func (r *PersonRepository) Update(id int, req *entity.Person) (*entity.Person, error) {
	sql := `UPDATE persons SET email = $1, phone = $2, first_name = $3, updated_at = $4 WHERE id = $5 RETURNING id;`
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	err = conn.QueryRow(context.Background(), sql, req.Email, req.Phone, req.FirstName, getTime(), id).
		Scan(&req.ID)
	return req, err
}

func (r *PersonRepository) Delete(id int) error {
	sql := `DELETE FROM persons WHERE id = $1`
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	result, err := conn.Exec(context.Background(), sql, id)
	if result.RowsAffected() != 1 {
		err = serverErr.ErrNotFound
	}
	return err
}

func (r *PersonRepository) ParseData(data []byte) (*entity.Person, error) {
	person := new(entity.Person)
	err := json.Unmarshal(data, &person)
	return person, err
}

func (r *PersonRepository) Count() (int, error) {
	result := 0
	sql := `SELECT COUNT(id) FROM persons;`
	conn, err := r.db.Acquire(context.Background())
	defer conn.Release()
	err = conn.QueryRow(context.Background(), sql).Scan(&result)
	return result, err
}

func getTime() string {
	return time.Now().Format(time.DateTime)
}
