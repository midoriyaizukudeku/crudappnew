package mocks

import (
	"dream.website/internal/model"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@Example.com":
		return model.DuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "zee@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, model.InvalidCredientials
}

func (m *UserModel) Exist(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
