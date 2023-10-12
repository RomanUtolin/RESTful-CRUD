package repository

import (
	"context"
	"encoding/json"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/gocraft/dbr/v2"
	"time"
)

type PersonRepository struct {
}

func NewPersonRepository() entity.PersonRepository {
	return &PersonRepository{}
}

func (r *PersonRepository) GetAll(ctx context.Context) ([]*entity.Person, error) {
	persons := make([]*entity.Person, 0)
	dbSess := getDb(ctx)
	_, err := dbSess.
		Select("id", "email", "phone", "first_name").
		From("persons").
		OrderBy("id").
		Load(&persons)
	return persons, err
}

func (r *PersonRepository) GetByID(ctx context.Context, id int) (*entity.Person, error) {
	var person *entity.Person
	dbSess := getDb(ctx)
	_, err := dbSess.
		Select("id", "email", "phone", "first_name").
		From("persons").
		Where("id = ?", id).
		Load(&person)
	return person, err
}

func (r *PersonRepository) GetByEmail(ctx context.Context, email string) (*entity.Person, error) {
	var person *entity.Person
	dbSess := getDb(ctx)
	_, err := dbSess.
		Select("id", "email", "phone", "first_name").
		From("persons").
		Where("email = ?", email).
		Load(&person)
	return person, err
}

func (r *PersonRepository) Create(ctx context.Context, req *entity.Person) (*entity.Person, error) {
	dbSess := getDb(ctx)
	err := dbSess.
		InsertInto("persons").
		Columns("email", "phone", "first_name", "created_at").
		Values(req.Email, req.Phone, req.FirstName, time.Now().Format(time.DateTime)).
		Returning("id").
		Load(&req)
	return req, err
}

func (r *PersonRepository) Update(ctx context.Context, id int, req *entity.Person) (*entity.Person, error) {
	dbSess := getDb(ctx)
	err := dbSess.
		Update("persons").
		Set("email", req.Email).
		Set("phone", req.Phone).
		Set("first_name", req.FirstName).
		Set("updated_at", time.Now().Format(time.DateTime)).
		Where("id = ?", id).
		Returning("id").
		Load(&req)
	return req, err
}

func (r *PersonRepository) Delete(ctx context.Context, id int) error {
	dbSess := getDb(ctx)
	del, err := dbSess.
		DeleteFrom("persons").
		Where("id = ?", id).
		Exec()
	if affected, _ := del.RowsAffected(); affected == 0 {
		err = serverErr.ErrNotFound
	}
	return err
}

func (r *PersonRepository) ParseData(data []byte) (*entity.Person, error) {
	var person *entity.Person
	err := json.Unmarshal(data, &person)
	return person, err
}

func getDb(ctx context.Context) *dbr.Session {
	dbSess := ctx.Value(constants.DbSession).(*dbr.Session)
	return dbSess
}
