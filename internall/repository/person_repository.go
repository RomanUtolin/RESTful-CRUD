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

func (r *PersonRepository) getPersons(ctx context.Context, query string, args ...interface{}) ([]*entity.Person, error) {
	persons := make([]*entity.Person, 0)
	rows, err := r.db.Query(ctx, query, args...)
	defer rows.Close()
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

func (r *PersonRepository) getOnePerson(ctx context.Context, query string, args ...interface{}) (*entity.Person, error) {
	person := new(entity.Person)
	err := r.db.QueryRow(ctx, query, args...).
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

func (r *PersonRepository) count(ctx context.Context, query, value string) (int, error) {
	var count int
	var err error
	if value == "" {
		err = r.db.QueryRow(ctx, query).Scan(&count)
	} else {
		err = r.db.QueryRow(ctx, query, value).Scan(&count)
	}
	return count, err
}

func (r *PersonRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			ORDER BY id
			LIMIT $1
			OFFSET $2;`
	return r.getPersons(ctx, sql, limit, offset)
}

func (r *PersonRepository) GetAllByEmail(ctx context.Context, email string, limit, offset int) ([]*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			WHERE email = $1
			ORDER BY id
			LIMIT $2
			OFFSET $3;`
	return r.getPersons(ctx, sql, email, limit, offset)
}

func (r *PersonRepository) GetAllByPhone(ctx context.Context, phone string, limit, offset int) ([]*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			WHERE phone = $1
			ORDER BY id
			LIMIT $2
			OFFSET $3;`
	return r.getPersons(ctx, sql, phone, limit, offset)
}

func (r *PersonRepository) GetAllByName(ctx context.Context, firstName string, limit, offset int) ([]*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			WHERE first_name = $1
			ORDER BY id
			LIMIT $2 
			OFFSET $3;`
	return r.getPersons(ctx, sql, firstName, limit, offset)
}

func (r *PersonRepository) GetByID(ctx context.Context, id int) (*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			WHERE id = $1;`
	return r.getOnePerson(ctx, sql, id)
}

func (r *PersonRepository) GetByEmail(ctx context.Context, email string) (*entity.Person, error) {
	sql := `SELECT id, email, phone, first_name
			FROM persons
			WHERE email = $1;`
	return r.getOnePerson(ctx, sql, email)
}

func (r *PersonRepository) Create(ctx context.Context, req *entity.Person) (*entity.Person, error) {
	sql := `INSERT INTO persons (email, phone, first_name, created_at)
			VALUES ($1,$2,$3,$4)
			RETURNING id;`
	err := r.db.QueryRow(ctx, sql, req.Email, req.Phone, req.FirstName, time.Now().Format(time.DateTime)).Scan(&req.ID)
	return req, err
}

func (r *PersonRepository) Update(ctx context.Context, id int, req *entity.Person) (*entity.Person, error) {
	sql := `UPDATE persons
			SET email = $1, phone = $2, first_name = $3, updated_at = $4
            WHERE id = $5
            RETURNING id;`
	err := r.db.QueryRow(ctx, sql, req.Email, req.Phone, req.FirstName, time.Now().Format(time.DateTime), id).Scan(&req.ID)
	return req, err
}

func (r *PersonRepository) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM persons
       		WHERE id = $1`

	result, err := r.db.Exec(ctx, sql, id)
	if result.RowsAffected() != 1 {
		err = serverErr.ErrNotFound
	}
	return err
}

func (r *PersonRepository) CountAll(ctx context.Context) (int, error) {
	sql := `SELECT COUNT(id) FROM persons;`
	return r.count(ctx, sql, "")
}

func (r *PersonRepository) CountAllByEmail(ctx context.Context, email string) (int, error) {
	sql := `SELECT COUNT(id) FROM persons WHERE email = $1;`
	return r.count(ctx, sql, email)
}

func (r *PersonRepository) CountAllByPhone(ctx context.Context, phone string) (int, error) {
	sql := `SELECT COUNT(id) FROM persons WHERE phone = $1;`
	return r.count(ctx, sql, phone)
}

func (r *PersonRepository) CountAllByName(ctx context.Context, name string) (int, error) {
	sql := `SELECT COUNT(id) FROM persons WHERE first_name = $1;`
	return r.count(ctx, sql, name)
}

func (r *PersonRepository) ParseData(data []byte) (*entity.Person, error) {
	person := new(entity.Person)
	err := json.Unmarshal(data, &person)
	return person, err
}
