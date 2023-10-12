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
)

var testPerson = &entity.Person{
	ID:        1,
	Email:     "test@test.ru",
	Phone:     "8999",
	FirstName: "test",
}

func TestPersonLogic_GetAll(t *testing.T) {
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson)
	tests := []struct {
		name       string
		mockFunc   func(mockUCase *mocks.PersonRepository)
		waitErr    error
		waitResult []*entity.Person
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAll", mock.Anything).Return(ListPerson, nil)
			},
			waitErr:    nil,
			waitResult: ListPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAll", mock.Anything).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		persons, err := personLogic.GetAll(context.Background())
		assert.Equal(t, test.waitErr, err)
		assert.Equal(t, test.waitResult, persons)

		mockUCase.AssertExpectations(t)
	}
}
func TestPersonLogic_GetByID(t *testing.T) {
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
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.GetByID(context.Background(), testPerson.ID)
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
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.Create(context.Background(), testPerson)
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
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.Update(context.Background(), testPerson.ID, testPerson)
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
		personLogic := logic.NewPersonLogic(mockUCase)
		err := personLogic.Delete(context.Background(), testPerson.ID)
		assert.Equal(t, test.waitErr, err)

		mockUCase.AssertExpectations(t)
	}
}
