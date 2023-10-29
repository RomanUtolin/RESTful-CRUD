package http

import (
	"errors"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"math"
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

func (h *Handler) GetAllPerson(c echo.Context) error {
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
	firstName := c.QueryParam("first_name")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	offset := (page - 1) * limit
	persons, count, err := h.Logic.GetAll(email, phone, firstName, limit, offset)
	if err != nil {
		return getError(c, err)
	}
	lastPage := int(math.Ceil(float64(count) / float64(limit)))
	data := &ResponseData{
		Data:     persons,
		Total:    count,
		Page:     page,
		LastPage: lastPage,
	}
	logrus.Info("Get all person Successful")
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) GetPerson(c echo.Context) error {
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	person, err := h.Logic.GetByID(id)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Get person id = %v Successful", id)
	return c.JSON(http.StatusOK, person)
}

func (h *Handler) CreatePerson(c echo.Context) error {
	req := &entity.Person{}
	if err := newPerson(c, req); err != nil {
		return getError(c, err)
	}
	person, err := h.Logic.Create(req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Create person id = %v Successful", person.ID)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) UpdatePerson(c echo.Context) error {
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	req := &entity.Person{}
	if err = newPerson(c, req); err != nil {
		return getError(c, err)
	}
	person, err := h.Logic.Update(id, req)
	if err != nil {
		return getError(c, err)
	}
	logrus.Infof("Update person id = %v Successful", id)
	return c.JSON(http.StatusCreated, person)
}

func (h *Handler) DeletePerson(c echo.Context) error {
	id, err := parseId(c)
	if err != nil {
		return getError(c, err)
	}
	err = h.Logic.Delete(id)
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
