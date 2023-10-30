package entity

import "context"

type Person struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
}

type PersonRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*Person, error)
	GetAllByEmail(ctx context.Context, email string, limit, offset int) ([]*Person, error)
	GetAllByPhone(ctx context.Context, phone string, limit, offset int) ([]*Person, error)
	GetAllByName(ctx context.Context, firstName string, limit, offset int) ([]*Person, error)
	GetByID(ctx context.Context, id int) (*Person, error)
	GetByEmail(ctx context.Context, email string) (*Person, error)
	Create(ctx context.Context, req *Person) (*Person, error)
	Update(ctx context.Context, id int, req *Person) (*Person, error)
	Delete(ctx context.Context, id int) error
	CountAll(ctx context.Context) (int, error)
	CountAllByEmail(ctx context.Context, email string) (int, error)
	CountAllByPhone(ctx context.Context, phone string) (int, error)
	CountAllByFirstName(ctx context.Context, name string) (int, error)
	ParseData(data []byte) (*Person, error)
}

type PersonLogic interface {
	GetPersons(ctx context.Context, email, phone, firstName string, page, limit int) ([]*Person, int, int, int, error)
	GetOnePerson(ctx context.Context, id int) (*Person, error)
	Create(ctx context.Context, req *Person) (*Person, error)
	Update(ctx context.Context, id int, req *Person) (*Person, error)
	Delete(ctx context.Context, id int) error
}
