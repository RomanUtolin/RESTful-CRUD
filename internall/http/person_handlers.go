package http

import (
	"errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	Logic entity.PersonLogic
}

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Data     []*entity.Person `json:"data"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	LastPage int              `json:"last_page"`
}

func (h *Handler) GetPersons(c echo.Context) error {
	ctx := c.Request().Context()
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
	firstName := c.QueryParam("first_name")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	persons, count, page, lastPage, err := h.Logic.GetPersons(ctx, email, phone, firstName, page, limit)
	if err != nil {
		return getError(c, err)
	}
	data := &ResponseData{
		Data:     persons,
		Total:    count,
		Page:     page,
		LastPage: lastPage,
	}
	logrus.Info("Get Persons Successful")
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) GetPerson(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	person, err := h.Logic.GetOnePerson(ctx, id)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Get person id = %v Successful", id)
	return c.JSON(http.StatusOK, person)
}

func (h *Handler) CreatePerson(c echo.Context) error {
	ctx := c.Request().Context()
	req := &entity.Person{}
	err := c.Bind(req)
	if err != nil {
		return getError(c, err)
	}
	person, err := h.Logic.Create(ctx, req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Create person id = %v Successful", person.ID)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) UpdatePerson(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	req := &entity.Person{}
	err := c.Bind(req)
	if err != nil {
		return getError(c, err)
	}
	person, err := h.Logic.Update(ctx, id, req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Update person id = %v Successful", id)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) DeletePerson(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Logic.Delete(ctx, id)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Delete person id = %v Successful", id)
	return c.NoContent(http.StatusNoContent)
}

func NewHandler(e *echo.Echo, logic entity.PersonLogic) {
	handler := &Handler{Logic: logic}
	e.GET("/person", handler.GetPersons)
	e.GET("/person/:id", handler.GetPerson)
	e.POST("/person", handler.CreatePerson)
	e.PUT("/person/:id", handler.UpdatePerson)
	e.DELETE("/person/:id", handler.DeletePerson)
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
