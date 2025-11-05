package data

type UserRepository interface {
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
