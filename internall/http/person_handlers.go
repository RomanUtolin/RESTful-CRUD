package http

import (
	"context"
	"errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

var lock = sync.Mutex{}

type Handler struct {
	Logic entity.PersonLogic
}

type ResponseError struct {
	Message string `json:"message"`
}

func (h *Handler) GetAllPerson(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	ctx := newContext(c)
	persons, err := h.Logic.GetAll(ctx)
	if err != nil {
		return getError(c, err)
	}
	logrus.Info("Get all person Successful")
	return c.JSON(http.StatusOK, persons)
}

func (h *Handler) GetPerson(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	ctx := newContext(c)
	person, err := h.Logic.GetByID(ctx, id)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Get person id = %v Successful", id)
	return c.JSON(http.StatusOK, person)
}

func (h *Handler) CreatePerson(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	req := &entity.Person{}
	if err := newPerson(c, req); err != nil {
		return getError(c, err)
	}
	ctx := newContext(c)
	person, err := h.Logic.Create(ctx, req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Create person id = %v Successful", person.ID)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) UpdatePerson(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	req := &entity.Person{}
	if err := newPerson(c, req); err != nil {
		return getError(c, err)
	}
	ctx := newContext(c)
	person, err := h.Logic.Update(ctx, id, req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Update person id = %v Successful", id)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) DeletePerson(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	ctx := newContext(c)
	err = h.Logic.Delete(ctx, id)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Delete person id = %v Successful", id)
	return c.NoContent(http.StatusNoContent)
}

func NewHandler(e *echo.Echo, logic entity.PersonLogic) {
	handler := &Handler{Logic: logic}
	e.GET("/person", handler.GetAllPerson)
	e.GET("/person/:id", handler.GetPerson)
	e.POST("/person", handler.CreatePerson)
	e.PUT("/person/:id", handler.UpdatePerson)
	e.DELETE("/person/:id", handler.DeletePerson)
}

func isRequestValid(req *entity.Person) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err,
			"Email":      req.Email,
			"Phone":      req.Phone,
			"First_name": req.FirstName,
		}).Error("validate err")
		err = serverErr.ErrBadParamInput
	}
	return err
}

func newPerson(c echo.Context, req *entity.Person) error {
	err := c.Bind(req)
	if err != nil {
		logrus.Error(err)
	} else {
		err = isRequestValid(req)
	}
	return err
}

func parseId(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = serverErr.ErrNotFound
	}
	return id, err
}

func getError(c echo.Context, err error) error {
	var code int
	logrus.Error(err)
	switch {
	case errors.Is(err, serverErr.ErrBadParamInput):
		code = http.StatusBadRequest
	case errors.Is(err, serverErr.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, serverErr.ErrConflict):
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}
	return c.JSON(code, ResponseError{Message: err.Error()})
}

func newContext(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), constants.DbSession, c.Get(constants.DbSession))
}
