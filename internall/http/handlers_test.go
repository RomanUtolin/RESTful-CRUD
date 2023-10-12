package http_test

import (
	"encoding/json"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	personHandler "github.com/RomanUtolin/RESTful-CRUD/internall/http"
	"github.com/RomanUtolin/RESTful-CRUD/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testPerson = &entity.Person{
	ID:        1,
	Email:     "test@test.ru",
	Phone:     "8999",
	FirstName: "test",
}
var invalidTestPerson = &entity.Person{
	Email:     "",
	Phone:     "8999",
	FirstName: "test",
}

func TestHandler_GetAllPerson(t *testing.T) {
	ListPerson := make([]*entity.Person, 0)
	ListPerson = append(ListPerson, testPerson)
	jsonListPerson, _ := json.Marshal(ListPerson)

	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		waitCode     int
		waitResponse string
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", mock.Anything).Return(ListPerson, nil)
			},
			waitCode:     http.StatusOK,
			waitResponse: string(jsonListPerson),
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", mock.Anything).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: `{"message":"internal Server Error"}`,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.GET, "", strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("person")

		handler := personHandler.Handler{Logic: mockUCase}
		err = handler.GetAllPerson(c)

		require.NoError(t, err)
		assert.Equal(t, test.waitCode, rec.Code)
		assert.Equal(t, test.waitResponse, strings.Trim(rec.Body.String(), "\n"))
		mockUCase.AssertExpectations(t)
	}
}

func TestHandler_GetPerson(t *testing.T) {
	PersonJson, _ := json.Marshal(testPerson)

	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		waitCode     int
		waitResponse string
		id           string
	}{
		{
			name: "valid",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(testPerson, nil)
			},
			waitCode:     http.StatusOK,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: `{"message":"internal Server Error"}`,
		},
		{
			name: "id valid, in db not found",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetByID", mock.Anything, testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitCode:     http.StatusNotFound,
			waitResponse: `{"message":"your requested item is not found"}`,
		},
		{
			name:         "id invalid",
			id:           "invalid",
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: `{"message":"your requested item is not found"}`,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.GET, "", strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("person/:id")
		c.SetParamNames("id")
		c.SetParamValues(test.id)
		handler := personHandler.Handler{Logic: mockUCase}
		err = handler.GetPerson(c)

		require.NoError(t, err)
		assert.Equal(t, test.waitCode, rec.Code)
		assert.Equal(t, test.waitResponse, strings.Trim(rec.Body.String(), "\n"))
		mockUCase.AssertExpectations(t)
	}
}

func TestHandler_CreatePerson(t *testing.T) {
	PersonJson, _ := json.Marshal(testPerson)
	invalidData, _ := json.Marshal(invalidTestPerson)

	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		waitCode     int
		waitResponse string
		data         []byte
	}{
		{
			name: "valid",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Create", mock.Anything, testPerson).Return(testPerson, nil)
			},
			waitCode:     http.StatusCreated,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Create", mock.Anything, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: `{"message":"internal Server Error"}`,
		},
		{
			name:         "given param is not valid",
			data:         invalidData,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusBadRequest,
			waitResponse: `{"message":"given param is not valid"}`,
		},
		{
			name: "Conflict Data in db",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Create", mock.Anything, testPerson).Return(nil, serverErr.ErrConflict)
			},
			waitCode:     http.StatusConflict,
			waitResponse: `{"message":"your email already exist, must be unique"}`,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.POST, "", strings.NewReader(string(test.data)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("person")
		handler := personHandler.Handler{Logic: mockUCase}
		err = handler.CreatePerson(c)

		require.NoError(t, err)
		assert.Equal(t, test.waitCode, rec.Code)
		assert.Equal(t, test.waitResponse, strings.Trim(rec.Body.String(), "\n"))
		mockUCase.AssertExpectations(t)
	}
}

func TestHandler_UpdatePerson(t *testing.T) {
	PersonJson, _ := json.Marshal(testPerson)
	invalidData, _ := json.Marshal(invalidTestPerson)

	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		waitCode     int
		waitResponse string
		data         []byte
		id           string
	}{
		{
			name: "valid",
			id:   "1",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Update", mock.Anything, testPerson.ID, testPerson).Return(testPerson, nil)
			},
			waitCode:     http.StatusCreated,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			id:   "1",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Update", mock.Anything, testPerson.ID, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: `{"message":"internal Server Error"}`,
		},
		{
			name:         "given param is not valid",
			id:           "1",
			data:         invalidData,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusBadRequest,
			waitResponse: `{"message":"given param is not valid"}`,
		},
		{
			name: "Conflict Data in db",
			id:   "1",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Update", mock.Anything, testPerson.ID, testPerson).Return(nil, serverErr.ErrConflict)
			},
			waitCode:     http.StatusConflict,
			waitResponse: `{"message":"your email already exist, must be unique"}`,
		},
		{
			name:         "id invalid",
			id:           "invalid",
			data:         PersonJson,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: `{"message":"your requested item is not found"}`,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.PUT, "", strings.NewReader(string(test.data)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("person/:id")
		c.SetParamNames("id")
		c.SetParamValues(test.id)
		handler := personHandler.Handler{Logic: mockUCase}
		err = handler.UpdatePerson(c)

		require.NoError(t, err)
		assert.Equal(t, test.waitCode, rec.Code)
		assert.Equal(t, test.waitResponse, strings.Trim(rec.Body.String(), "\n"))
		mockUCase.AssertExpectations(t)
	}
}
func TestHandler_DeletePerson(t *testing.T) {
	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		waitCode     int
		waitResponse string
		id           string
	}{
		{
			name: "valid",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(nil)
			},
			waitCode:     http.StatusNoContent,
			waitResponse: "",
		},
		{
			name: "store error",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: `{"message":"internal Server Error"}`,
		},
		{
			name: "id valid, in db not found",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Delete", mock.Anything, testPerson.ID).Return(serverErr.ErrNotFound)
			},
			waitCode:     http.StatusNotFound,
			waitResponse: `{"message":"your requested item is not found"}`,
		},
		{
			name:         "id invalid",
			id:           "invalid",
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: `{"message":"your requested item is not found"}`,
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.DELETE, "", strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("person/:id")
		c.SetParamNames("id")
		c.SetParamValues(test.id)
		handler := personHandler.Handler{Logic: mockUCase}
		err = handler.DeletePerson(c)

		require.NoError(t, err)
		assert.Equal(t, test.waitCode, rec.Code)
		assert.Equal(t, test.waitResponse, strings.Trim(rec.Body.String(), "\n"))
		mockUCase.AssertExpectations(t)
	}
}
