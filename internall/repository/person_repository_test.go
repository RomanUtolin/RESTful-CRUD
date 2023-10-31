package repository_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testPerson1         = &entity.Person{ID: 1, Email: "test@test.ru", Phone: "1234", FirstName: "test"}
	testPerson2         = &entity.Person{ID: 2, Email: "test2@test.ru", Phone: "5678", FirstName: "test2"}
	testPersonForUpdate = &entity.Person{Email: "updated@test.ru", Phone: "9999", FirstName: "updated"}
	testPerson3         = &entity.Person{ID: 3, Email: "test3@test.ru", Phone: "1234", FirstName: "test"}
)

func GetTestDb() *pgxpool.Pool {
	dbHost := "localhost"
	dbPort := "5433"
	dbUser := "postgres"
	dbPass := "password"
	dbName := "devtest"
	sslMode := "disable"
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		logrus.Warning(err)
	}
	return dbPool
}

func truncate(ctx context.Context, db *pgxpool.Pool) {
	sql := `TRUNCATE persons CASCADE;`
	db.Exec(ctx, sql)
}

func TestPersonRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson1)
	ListPerson = append(ListPerson, testPerson2)
	result, err := rep.GetAll(ctx, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ListPerson, result)
}

func TestPersonRepository_GetAllByEmail(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson1)
	result, err := rep.GetAllByEmail(ctx, testPerson1.Email, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ListPerson, result)
}
func TestPersonRepository_GetAllByName(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson1)
	ListPerson = append(ListPerson, testPerson3)
	result, err := rep.GetAllByName(ctx, testPerson1.FirstName, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ListPerson, result)
}

func TestPersonRepository_GetAllByPhone(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson1)
	ListPerson = append(ListPerson, testPerson3)
	result, err := rep.GetAllByPhone(ctx, testPerson1.Phone, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ListPerson, result)
}

func TestPersonRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	result, err := rep.GetByID(ctx, testPerson1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson1, result)
}

func TestPersonRepository_GetByEmail(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	result, err := rep.GetByEmail(ctx, testPerson1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson1, result)
}

func TestPersonRepository_Create(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	result, err := rep.Create(ctx, testPerson1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPerson1, result)
}
func TestPersonRepository_Update(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	person, _ := rep.Create(ctx, testPerson1)
	result, err := rep.Update(ctx, person.ID, testPersonForUpdate)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPersonForUpdate, result)
}
func TestPersonRepository_Delete(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	person, _ := rep.Create(ctx, testPerson1)
	err := rep.Delete(ctx, person.ID)
	assert.NoError(t, err)
	err = rep.Delete(ctx, person.ID)
	assert.Equal(t, err, serverErr.ErrNotFound)
}

func TestPersonRepository_CountAll(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	result, err := rep.CountAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result)
}

func TestPersonRepository_CountAllByEmail(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	result, err := rep.CountAllByEmail(ctx, testPerson1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result)
}

func TestPersonRepository_CountAllByName(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	result, err := rep.CountAllByName(ctx, testPerson1.FirstName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result)
}

func TestPersonRepository_CountAllByPhone(t *testing.T) {
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	rep.Create(ctx, testPerson1)
	rep.Create(ctx, testPerson2)
	rep.Create(ctx, testPerson3)
	result, err := rep.CountAllByPhone(ctx, testPerson1.Phone)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result)
}

func TestPersonRepository_ParseData(t *testing.T) {
	data, _ := json.Marshal(testPerson1)
	ctx := context.Background()
	dbPoll := GetTestDb()
	defer func() {
		truncate(ctx, dbPoll)
		dbPoll.Close()
	}()
	rep := repository.NewPersonRepository(dbPoll)
	result, err := rep.ParseData(data)
	assert.NoError(t, err)
	assert.Equal(t, testPerson1, result)
}
