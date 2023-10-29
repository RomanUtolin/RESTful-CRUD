package logic_test

import (
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/logic"
	"github.com/RomanUtolin/RESTful-CRUD/mocks"
	"github.com/stretchr/testify/assert"
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
				mockUCase.On("GetAll", "", "", "", 10, 0).Return(ListPerson, 1, nil)
			},
			waitErr:    nil,
			waitResult: ListPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetAll", "", "", "", 10, 0).Return(nil, 1, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		persons, data, err := personLogic.GetAll("", "", "", 10, 0)
		_ = data
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
				mockUCase.On("GetByID", testPerson.ID).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", testPerson.ID).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "not found",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitErr:    serverErr.ErrNotFound,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.GetByID(testPerson.ID)
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
				mockUCase.On("GetByEmail", testPerson.Email).Return(nil, nil)
				mockUCase.On("Create", testPerson).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByEmail", testPerson.Email).Return(nil, nil)
				mockUCase.On("Create", testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "Conflict Data in db",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByEmail", testPerson.Email).Return(testPerson, nil)
			},
			waitErr:    serverErr.ErrConflict,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.Create(testPerson)
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
				mockUCase.On("GetByID", testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", testPerson.Email).Return(nil, nil)
				mockUCase.On("Update", testPerson.ID, testPerson).Return(testPerson, nil)
			},
			waitErr:    nil,
			waitResult: testPerson,
		},
		{
			name: "invalid id",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitErr:    serverErr.ErrNotFound,
			waitResult: nil,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", testPerson.Email).Return(nil, nil)
				mockUCase.On("Update", testPerson.ID, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitErr:    serverErr.ErrInternalServer,
			waitResult: nil,
		},
		{
			name: "Conflict Data in db",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("GetByID", testPerson.ID).Return(testPerson, nil)
				mockUCase.On("GetByEmail", testPerson.Email).Return(testPerson, nil)
			},
			waitErr:    serverErr.ErrConflict,
			waitResult: nil,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		person, err := personLogic.Update(testPerson.ID, testPerson)
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
				mockUCase.On("Delete", testPerson.ID).Return(nil)
			},
			waitErr: nil,
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("Delete", testPerson.ID).Return(serverErr.ErrInternalServer)
			},
			waitErr: serverErr.ErrInternalServer,
		},
		{
			name: "not found",
			mockFunc: func(mockUCase *mocks.PersonRepository) {
				mockUCase.On("Delete", testPerson.ID).Return(serverErr.ErrNotFound)
			},
			waitErr: serverErr.ErrNotFound,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonRepository)
		test.mockFunc(mockUCase)
		personLogic := logic.NewPersonLogic(mockUCase)
		err := personLogic.Delete(testPerson.ID)
		assert.Equal(t, test.waitErr, err)

		mockUCase.AssertExpectations(t)
	}
}
