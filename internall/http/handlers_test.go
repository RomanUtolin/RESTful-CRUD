package http_test

import (
	"encoding/json"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	personHandler "github.com/RomanUtolin/RESTful-CRUD/internall/http"
	"github.com/RomanUtolin/RESTful-CRUD/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testPerson = &entity.Person{
	ID:        1,
	Email:     "test@test.ru",
	Phone:     "1234",
	FirstName: "test",
}

var testPerson2 = &entity.Person{
	ID:        2,
	Email:     "test2@test.ru",
	Phone:     "5678",
	FirstName: "test2",
}

var invalidTestPerson = &entity.Person{
	Email:     "",
	Phone:     "8999",
	FirstName: "test",
}

var (
	jsonErrServer, _   = json.Marshal(personHandler.ResponseError{Message: serverErr.ErrInternalServer.Error()})
	jsonErrConflict, _ = json.Marshal(personHandler.ResponseError{Message: serverErr.ErrConflict.Error()})
	jsonErrNotFound, _ = json.Marshal(personHandler.ResponseError{Message: serverErr.ErrNotFound.Error()})
	jsonErrBadParam, _ = json.Marshal(personHandler.ResponseError{Message: serverErr.ErrBadParamInput.Error()})
)

func TestHandler_GetAllPerson(t *testing.T) {
	ListPerson := make([]*entity.Person, 0)
	ListOnePerson := append(ListPerson, testPerson)
	ListTwoPerson := append(ListOnePerson, testPerson2)
	dataOnePerson := &personHandler.ResponseData{
		Data:     ListOnePerson,
		Total:    1,
		Page:     1,
		LastPage: 1,
	}
	dataTwoPerson := &personHandler.ResponseData{
		Data:     ListTwoPerson,
		Total:    2,
		Page:     1,
		LastPage: 1,
	}
	dataTwoPersonLimit := &personHandler.ResponseData{
		Data:     ListOnePerson,
		Total:    2,
		Page:     1,
		LastPage: 2,
	}
	jsonOnePerson, _ := json.Marshal(*dataOnePerson)
	jsonTwoPerson, _ := json.Marshal(*dataTwoPerson)
	jsonTwoPersonLimit, _ := json.Marshal(*dataTwoPersonLimit)

	tests := []struct {
		name         string
		mockFunc     func(mockUCase *mocks.PersonLogic)
		path         string
		waitCode     int
		waitResponse string
	}{
		{
			name: "valid",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "", "", "", 10, 0).Return(ListTwoPerson, 2, nil)
			},
			path:         "person",
			waitCode:     http.StatusOK,
			waitResponse: string(jsonTwoPerson),
		},
		{
			name: "valid with param email",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "test@test.ru", "", "", 10, 0).Return(ListOnePerson, 1, nil)
			},
			path:         "person?email=test@test.ru",
			waitCode:     http.StatusOK,
			waitResponse: string(jsonOnePerson),
		},
		{
			name: "valid with param phone",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "", "1234", "", 10, 0).Return(ListOnePerson, 1, nil)
			},
			path:         "person?phone=1234",
			waitCode:     http.StatusOK,
			waitResponse: string(jsonOnePerson),
		},
		{
			name: "valid with param name",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "", "", "test", 10, 0).Return(ListOnePerson, 1, nil)
			},
			path:         "person?first_name=test",
			waitCode:     http.StatusOK,
			waitResponse: string(jsonOnePerson),
		},
		{
			name: "valid page=1&limit=1 all page 2",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "", "", "", 1, 0).Return(ListOnePerson, 2, nil)
			},
			path:         "person?page=1&limit=1",
			waitCode:     http.StatusOK,
			waitResponse: string(jsonTwoPersonLimit),
		},
		{
			name: "store error",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetAll", "", "", "", 10, 0).Return(nil, 0, serverErr.ErrInternalServer)
			},
			path:         "person",
			waitCode:     http.StatusInternalServerError,
			waitResponse: string(jsonErrServer),
		},
	}
	for _, test := range tests {
		mockUCase := new(mocks.PersonLogic)
		test.mockFunc(mockUCase)

		e := echo.New()

		req, err := http.NewRequest(echo.GET, test.path, strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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
				mockUCase.On("GetByID", testPerson.ID).Return(testPerson, nil)
			},
			waitCode:     http.StatusOK,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetByID", testPerson.ID).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: string(jsonErrServer),
		},
		{
			name: "id valid, in db not found",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("GetByID", testPerson.ID).Return(nil, serverErr.ErrNotFound)
			},
			waitCode:     http.StatusNotFound,
			waitResponse: string(jsonErrNotFound),
		},
		{
			name:         "id invalid",
			id:           "invalid",
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: string(jsonErrNotFound),
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
				mockUCase.On("Create", testPerson).Return(testPerson, nil)
			},
			waitCode:     http.StatusCreated,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Create", testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: string(jsonErrServer),
		},
		{
			name:         "given param is not valid",
			data:         invalidData,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusBadRequest,
			waitResponse: string(jsonErrBadParam),
		},
		{
			name: "Conflict Data in db",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Create", testPerson).Return(nil, serverErr.ErrConflict)
			},
			waitCode:     http.StatusConflict,
			waitResponse: string(jsonErrConflict),
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
				mockUCase.On("Update", testPerson.ID, testPerson).Return(testPerson, nil)
			},
			waitCode:     http.StatusCreated,
			waitResponse: string(PersonJson),
		},
		{
			name: "store error",
			id:   "1",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Update", testPerson.ID, testPerson).Return(nil, serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: string(jsonErrServer),
		},
		{
			name:         "given param is not valid",
			id:           "1",
			data:         invalidData,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusBadRequest,
			waitResponse: string(jsonErrBadParam),
		},
		{
			name: "Conflict Data in db",
			id:   "1",
			data: PersonJson,
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Update", testPerson.ID, testPerson).Return(nil, serverErr.ErrConflict)
			},
			waitCode:     http.StatusConflict,
			waitResponse: string(jsonErrConflict),
		},
		{
			name:         "id invalid",
			id:           "invalid",
			data:         PersonJson,
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: string(jsonErrNotFound),
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
				mockUCase.On("Delete", testPerson.ID).Return(nil)
			},
			waitCode:     http.StatusNoContent,
			waitResponse: "",
		},
		{
			name: "store error",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Delete", testPerson.ID).Return(serverErr.ErrInternalServer)
			},
			waitCode:     http.StatusInternalServerError,
			waitResponse: string(jsonErrServer),
		},
		{
			name: "id valid, in db not found",
			id:   "1",
			mockFunc: func(mockUCase *mocks.PersonLogic) {
				mockUCase.On("Delete", testPerson.ID).Return(serverErr.ErrNotFound)
			},
			waitCode:     http.StatusNotFound,
			waitResponse: string(jsonErrNotFound),
		},
		{
			name:         "id invalid",
			id:           "invalid",
			mockFunc:     func(mockUCase *mocks.PersonLogic) {},
			waitCode:     http.StatusNotFound,
			waitResponse: string(jsonErrNotFound),
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
