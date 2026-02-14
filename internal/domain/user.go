package domain

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
}

type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user *User) error
}

type UserUsecase interface {
	GetAllUsers() ([]User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
}
