package entity

import (
	"context"
)

type Person struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
}

type PersonRepository interface {
	GetAll(ctx context.Context) ([]*Person, error)
	GetByID(ctx context.Context, id int) (*Person, error)
	GetByEmail(ctx context.Context, email string) (*Person, error)
	Create(ctx context.Context, req *Person) (*Person, error)
	Update(ctx context.Context, id int, req *Person) (*Person, error)
	Delete(ctx context.Context, id int) error
	ParseData(data []byte) (*Person, error)
}

type PersonLogic interface {
	GetAll(ctx context.Context) ([]*Person, error)
	GetByID(ctx context.Context, id int) (*Person, error)
	GetByEmail(ctx context.Context, email string) (*Person, error)
	Create(ctx context.Context, req *Person) (*Person, error)
	Update(ctx context.Context, id int, req *Person) (*Person, error)
	Delete(ctx context.Context, id int) error
}
