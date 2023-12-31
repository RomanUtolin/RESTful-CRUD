package logic_test

import (
	"context"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/logic"
	"github.com/RomanUtolin/RESTful-CRUD/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var testPerson = &entity.Person{
	ID:        1,
	Email:     "test@test.ru",
	Phone:     "8999",
	FirstName: "test",
}

func TestPersonLogic_GetPersons(t *testing.T) {
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson)
	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonRepository)
		waitErr      error
		waitResult   []*entity.Person
		waitCount    int
		waitPage     int
		waitLastPage int
		email        string
		phone        string
		firstName    string
	}{
		{
			name: "GetAllValid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAll", mock.Anything, 10, 0).Return(ListPerson, nil)
				mockUCase.On("CountAll", mock.Anything).Return(1, nil)
			},
			waitErr:      nil,
			waitResult:   ListPerson,
			waitCount:    1,
			waitPage:     1,
			waitLastPage: 1,
			email:        "",
			phone:        "",
			firstName:    "",
		},
		{
			name: "GetAllByEmailValid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAllByEmail", mock.Anything, testPerson.Email, 10, 0).Return(ListPerson, nil)
				mockUCase.On("CountAllByEmail", mock.Anything, testPerson.Email).Return(1, nil)
			},
			waitErr:      nil,
			waitResult:   ListPerson,
			waitCount:    1,
			waitPage:     1,
			waitLastPage: 1,
			email:        testPerson.Email,
			phone:        "",
			firstName:    "",
		},
		{
			name: "GetAllByPhoneValid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAllByPhone", mock.Anything, testPerson.Phone, 10, 0).Return(ListPerson, nil)
				mockUCase.On("CountAllByPhone", mock.Anything, testPerson.Phone).Return(1, nil)
			},
			waitErr:      nil,
			waitResult:   ListPerson,
			waitCount:    1,
			waitPage:     1,
			waitLastPage: 1,
			email:        "",
			phone:        testPerson.Phone,
			firstName:    "",
		},
		{
			name: "GetAllByNameValid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAllByName", mock.Anything, testPerson.FirstName, 10, 0).Return(ListPerson, nil)
				mockUCase.On("CountAllByName", mock.Anything, testPerson.FirstName).Return(1, nil)
			},
			waitErr:      nil,
			waitResult:   ListPerson,
			waitCount:    1,
			waitPage:     1,
			waitLastPage: 1,
			email:        "",
			phone:        "",
			firstName:    testPerson.FirstName,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAll", mock.Anything, 10, 0).Return(nil, serverErr.ErrInternalServer)
				mockUCase.On("CountAll", mock.Anything).Return(0, serverErr.ErrInternalServer)
			},
			waitErr:      serverErr.ErrInternalServer,
			waitResult:   nil,
			waitCount:    0,
			waitPage:     1,
			waitLastPage: 0,
			email:        "",
			phone:        "",
			firstName:    "",
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase, time.Second*2)
		persons, count, page, lastPage, err := personLogic.GetPersons(context.TODO(), test.email, test.phone, test.firstName, 0, 0)

		assert.Equal(t, test.waitErr, err)
		assert.Equal(t, test.waitResult, persons)
		assert.Equal(t, test.waitCount, count)
		assert.Equal(t, test.waitPage, page)
		assert.Equal(t, test.waitLastPage, lastPage)

		mockUCase.AssertExpectations(t)
	}
}

func TestPersonLogic_GetOnePerson(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(mockUCase *mocks.PersonRepository)
		waitErr    error
		waitResult *entity.Person
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "not found",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitErr:    serverErr.ErrNotFound,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase, time.Second*2)
		person, err := personLogic.GetOnePerson(context.TODO(), testPerson.ID)
		assert.Equal(t, test.waitErr, err)
		assert.Equal(t, test.waitResult, person)

		mockUCase.AssertExpectations(t)
	}
}

func TestPersonLogic_Create(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(mockUCase *mocks.PersonRepository)
		waitErr    error
		waitResult *entity.Person
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(nil, nil)
				mockUCase.On("Create", mock.Anything, testPerson).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(nil, nil)
				mockUCase.On("Create", mock.Anything, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "Conflict Data in db",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(testPerson, nil)
			},
			waitErr:    serverErr.ErrConflict,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase, time.Second*2)
		person, err := personLogic.Create(context.TODO(), testPerson)
		assert.Equal(t, test.waitErr, err)
		assert.Equal(t, test.waitResult, person)

		mockUCase.AssertExpectations(t)
	}
}

func TestPersonLogic_Update(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(mockUCase *mocks.PersonRepository)
		waitErr    error
		waitResult *entity.Person
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(nil, nil)
				mockUCase.On("Update", mock.Anything, testPerson.ID, testPerson).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "invalid id",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitErr:    serverErr.ErrNotFound,
			waitResult: nil,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(nil, nil)
				mockUCase.On("Update", mock.Anything, testPerson.ID, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "Conflict Data in db",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", mock.Anything, testPerson.Email).Return(testPerson, nil)
			},
			waitErr:    serverErr.ErrConflict,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase, time.Second*2)
		person, err := personLogic.Update(context.TODO(), testPerson.ID, testPerson)
		assert.Equal(t, test.waitErr, err)
		assert.Equal(t, test.waitResult, person)

		mockUCase.AssertExpectations(t)
	}
}

func TestPersonLogic_Delete(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(mockUCase *mocks.PersonRepository)
		waitErr  error
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(nil)
			},
			waitErr: nil,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(serverErr.ErrInternalServer)
			},
			waitErr: serverErr.ErrInternalServer,
		},
		{
			name: "not found",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(serverErr.ErrNotFound)
			},
			waitErr: serverErr.ErrNotFound,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase, time.Second*2)
		err := personLogic.Delete(context.TODO(), testPerson.ID)
		assert.Equal(t, test.waitErr, err)

		mockUCase.AssertExpectations(t)
	}
}
