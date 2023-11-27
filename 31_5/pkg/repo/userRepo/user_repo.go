package userRepo

import (
	"31_5/pkg/user"
	"fmt"
)

type UserStorage map[string]*user.User

type Service struct {
	Store map[string]*user.User
}

type Storage struct {
	Users []*user.User
}

func NewUserStorage() UserStorage {
	return make(map[string]*user.User)
}

func (us UserStorage) Put(u *user.User) {
	us[u.Name] = u
}

func (us UserStorage) Get(userName string) (*user.User, error) {
	u, ok := us[userName]
	if !ok {
		return nil, fmt.Errorf("no such user")
	}

	return u, nil
}
