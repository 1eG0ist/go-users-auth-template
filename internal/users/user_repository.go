package users

import "sync"

type UserInMemoryRepository struct {
	// here will be db connections pool, now map
	mu           sync.Mutex
	users        map[int]*User
	currentMaxId int
}

func NewUserInMemoryRepository() *UserInMemoryRepository {
	ur := &UserInMemoryRepository{
		users: make(map[int]*User),
	}
	ur.currentMaxId = 0
	return ur
}

func (ur *UserInMemoryRepository) CreateUser(nickname, email, password string) (*User, error) {
	userId := ur.currentMaxId
	ur.mu.Lock()
	ur.currentMaxId++
	ur.mu.Unlock()
	ur.users[userId] = &User{
		ID:       userId,
		Nickname: nickname,
		Email:    email,
		Password: password,
	}

	return ur.users[userId], nil
}

func (ur *UserInMemoryRepository) GetUserById(id int) (*User, error) {
	if user, ok := ur.users[id]; ok {
		return user, nil
	}

	return nil, ErrNotFound
}

func (ur *UserInMemoryRepository) GetAllUsers() ([]*User, error) {
	var users = make([]*User, 0, len(ur.users))
	for _, user := range ur.users {
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserInMemoryRepository) UpdateUser(id int, nickname, email string) error {
	oldUser, ok := ur.users[id]
	if !ok {
		return ErrNotFound
	}
	ur.users[id] = &User{
		ID:       id,
		Nickname: nickname,
		Email:    email,
		Password: oldUser.Password,
	}
	return nil
}

func (ur *UserInMemoryRepository) DeleteUser(id int) error {
	if _, ok := ur.users[id]; !ok {
		return ErrNotFound
	}
	delete(ur.users, id)
	return nil
}
