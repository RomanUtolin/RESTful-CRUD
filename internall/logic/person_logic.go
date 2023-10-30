package logic

import (
	"context"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"math"
)

type PersonLogic struct {
	Rep entity.PersonRepository
}

func NewPersonLogic(rep entity.PersonRepository) entity.PersonLogic {
	return &PersonLogic{rep}
}

func (p *PersonLogic) GetPersons(ctx context.Context, email, phone, firstName string, page, limit int) ([]*entity.Person, int, int, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var persons []*entity.Person
	var count int
	var err error
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	switch {
	case email != "":
		persons, err = p.Rep.GetAllByEmail(ctx, email, limit, offset)
		count, err = p.Rep.CountAllByEmail(ctx, email)
	case phone != "":
		persons, err = p.Rep.GetAllByPhone(ctx, phone, limit, offset)
		count, err = p.Rep.CountAllByPhone(ctx, phone)
	case firstName != "":
		persons, err = p.Rep.GetAllByName(ctx, firstName, limit, offset)
		count, err = p.Rep.CountAllByFirstName(ctx, firstName)
	default:
		persons, err = p.Rep.GetAll(ctx, limit, offset)
		count, err = p.Rep.CountAll(ctx)
	}
	lastPage := int(math.Ceil(float64(count) / float64(limit)))
	return persons, count, page, lastPage, err
}

func (p *PersonLogic) GetOnePerson(ctx context.Context, id int) (*entity.Person, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if id == 0 {
		return nil, serverErr.ErrNotFound
	}
	person, err := p.Rep.GetByID(ctx, id)
	if person == nil && err == nil {
		err = serverErr.ErrNotFound
	}
	return person, err
}

func (p *PersonLogic) Create(ctx context.Context, req *entity.Person) (*entity.Person, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := isRequestValid(req)
	if err != nil {
		return nil, err
	}
	err = p.findPerson(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return p.Rep.Create(ctx, req)
}

func (p *PersonLogic) Update(ctx context.Context, id int, req *entity.Person) (*entity.Person, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if id == 0 {
		return nil, serverErr.ErrNotFound
	}
	err := isRequestValid(req)
	if err != nil {
		return nil, err
	}
	_, err = p.GetOnePerson(ctx, id)
	if err != nil {
		return nil, err
	}
	err = p.findPerson(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return p.Rep.Update(ctx, id, req)
}

func (p *PersonLogic) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if id == 0 {
		return serverErr.ErrNotFound
	}
	return p.Rep.Delete(ctx, id)
}

func (p *PersonLogic) findPerson(ctx context.Context, email string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	person, err := p.Rep.GetByEmail(ctx, email)
	if person != nil && err == nil {
		err = serverErr.ErrConflict
	}
	return err
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
