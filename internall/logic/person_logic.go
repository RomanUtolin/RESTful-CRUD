package logic

import (
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
)

type PersonLogic struct {
	Rep entity.PersonRepository
}

func NewPersonLogic(rep entity.PersonRepository) entity.PersonLogic {
	return &PersonLogic{rep}
}

func (p *PersonLogic) GetAll(email, phone, firstName string, limit, offset int) ([]*entity.Person, int, error) {
	return p.Rep.GetAll(email, phone, firstName, limit, offset)
}

func (p *PersonLogic) GetByID(id int) (*entity.Person, error) {
	person, err := p.Rep.GetByID(id)
	if person == nil && err == nil {
		err = serverErr.ErrNotFound
	}
	return person, err
}

func (p *PersonLogic) GetByEmail(email string) (*entity.Person, error) {
	return p.Rep.GetByEmail(email)
}

func (p *PersonLogic) Create(req *entity.Person) (*entity.Person, error) {
	if err := p.findPerson(req.Email); err != nil {
		return nil, err
	}
	return p.Rep.Create(req)
}

func (p *PersonLogic) Update(id int, req *entity.Person) (*entity.Person, error) {
	if _, err := p.GetByID(id); err != nil {
		return nil, err
	}
	if err := p.findPerson(req.Email); err != nil {
		return nil, err
	}
	return p.Rep.Update(id, req)
}

func (p *PersonLogic) Delete(id int) error {
	return p.Rep.Delete(id)
}

func (p *PersonLogic) findPerson(email string) error {
	person, err := p.GetByEmail(email)
	if person != nil && err == nil {
		err = serverErr.ErrConflict
	}
	return err
}
