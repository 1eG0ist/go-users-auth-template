package users

import "fmt"

type UserRepository interface {
	CreateUser(nickname, email, password string) (*User, error)
	GetAllUsers() ([]*User, error)
	GetUserById(id int) (*User, error)
	UpdateUser(id int, nickname, email string) error
	DeleteUser(id int) error
}

type UserService struct {
	ur UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{ur}
}

// method will be called only by auth microsercie after email confirmation
func (us *UserService) RegisterUser(nickname, email, password string) (*User, error) {
	if user, err := us.ur.CreateUser(nickname, email, password); err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	} else {
		return user, nil
	}
}

func (us *UserService) GetUserProfile(id int) (*User, error) {
	return us.ur.GetUserById(id)
}

func (us *UserService) GetAllUserProfiles() ([]*User, error) {
	if users, err := us.ur.GetAllUsers(); err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	} else {
		return users, nil
	}
}

func (us *UserService) UpdateUserProfile(id int, user UserUpdateRequest) error {
	if err := us.ur.UpdateUser(id, user.Nickname, user.Email); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func (us *UserService) DeleteUserProfile(id int) error {
	if err := us.ur.DeleteUser(id); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}
