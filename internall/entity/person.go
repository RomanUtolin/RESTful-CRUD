package entity

type Person struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
}

type PersonRepository interface {
	GetAll(email, phone, firstName string, page, limit int) ([]*Person, error)
	GetByID(id int) (*Person, error)
	GetByEmail(email string) (*Person, error)
	Create(req *Person) (*Person, error)
	Update(id int, req *Person) (*Person, error)
	Delete(id int) error
	ParseData(data []byte) (*Person, error)
	Count() (int, error)
}

type PersonLogic interface {
	GetAll(email, phone, firstName string, page, limit int) ([]*Person, error)
	GetByID(id int) (*Person, error)
	GetByEmail(email string) (*Person, error)
	Create(req *Person) (*Person, error)
	Update(id int, req *Person) (*Person, error)
	Delete(id int) error
	Count() (int, error)
}
