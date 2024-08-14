package demo

type Repository interface {
	Create(i int) (int, error)
}

type service struct{}

func (r *service) Create(i int) (int, error) {
	return 100, nil
}

func NewRepository(_repository Repository) Repository {
	return &service{}
}
