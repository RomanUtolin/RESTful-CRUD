package repository_test

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	"github.com/RomanUtolin/RESTful-CRUD/internall/repository"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var testPerson = &entity.Person{
	ID:        1,
	Email:     "test@test.ru",
	Phone:     "8999",
	FirstName: "test",
}

func TestPersonRepository_GetAll(t *testing.T) {
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson)
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		db.Close()
	}()

	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.PostgreSQL,
	}
	sess := conn.NewSession(nil)

	query := "SELECT (.+) FROM persons"

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)

	mock.ExpectQuery(query).WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.DbSession, sess)
	personRepository := repository.NewPersonRepository()
	result, err := personRepository.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ListPerson, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPersonRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		db.Close()
	}()

	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.PostgreSQL,
	}
	sess := conn.NewSession(nil)
	query := "SELECT (.+) FROM persons WHERE (.+)"

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)

	mock.ExpectQuery(query).WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.DbSession, sess)
	personRepository := repository.NewPersonRepository()
	result, err := personRepository.GetByID(ctx, testPerson.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPersonRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		db.Close()
	}()

	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.PostgreSQL,
	}
	sess := conn.NewSession(nil)
	query := "SELECT (.+) FROM persons WHERE (.+)"

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)

	mock.ExpectQuery(query).WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.DbSession, sess)
	personRepository := repository.NewPersonRepository()
	result, err := personRepository.GetByEmail(ctx, testPerson.Email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPersonRepository_ParseData(t *testing.T) {
	personRepository := repository.NewPersonRepository()
	PersonJson, _ := json.Marshal(testPerson)
	result, err := personRepository.ParseData(PersonJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson, result)
}
