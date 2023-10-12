package logic

import (
	"context"
	"github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	serverErr "github.com/RomanUtolin/RESTful-CRUD/internall/errors"
)

type PersonLogic struct {
	Rep entity.PersonRepository
}

func NewPersonLogic(rep entity.PersonRepository) entity.PersonLogic {
	return &PersonLogic{rep}
}

func (p *PersonLogic) GetAll(ctx context.Context) ([]*entity.Person, error) {
	return p.Rep.GetAll(ctx)
}

func (p *PersonLogic) GetByID(ctx context.Context, id int) (*entity.Person, error) {
	person, err := p.Rep.GetByID(ctx, id)
	if person == nil && err == nil {
		err = serverErr.ErrNotFound
	}
	return person, err
}

func (p *PersonLogic) GetByEmail(ctx context.Context, email string) (*entity.Person, error) {
	return p.Rep.GetByEmail(ctx, email)
}

func (p *PersonLogic) Create(ctx context.Context, req *entity.Person) (*entity.Person, error) {
	if err := p.findPerson(ctx, req.Email); err != nil {
		return nil, err
	}
	return p.Rep.Create(ctx, req)
}

func (p *PersonLogic) Update(ctx context.Context, id int, req *entity.Person) (*entity.Person, error) {
	if _, err := p.GetByID(ctx, id); err != nil {
		return nil, err
	}
	if err := p.findPerson(ctx, req.Email); err != nil {
		return nil, err
	}
	return p.Rep.Update(ctx, id, req)
}

func (p *PersonLogic) Delete(ctx context.Context, id int) error {
	return p.Rep.Delete(ctx, id)
}

func (p *PersonLogic) findPerson(ctx context.Context, email string) error {
	person, err := p.GetByEmail(ctx, email)
	if person != nil && err == nil {
		err = serverErr.ErrConflict
	}
	return err
}
