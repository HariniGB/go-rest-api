package common

import "fmt"

type User struct {
	Username  string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u User) ValidateLogin() error {
	if u.Username == "" {
		fmt.Errorf("id can not be empty")
	}

	if u.Password == "" {
		fmt.Errorf("password can not be empty")
	}
	return nil
}

func (u User) ValidateSignup() error {
	err := u.ValidateLogin()
	if err != nil {
		return err
	}

	if u.FirstName == "" {
		fmt.Errorf("first name can not be empty")
	}

	if u.LastName == "" {
		fmt.Errorf("last name can not be empty")
	}

	if u.Email == "" {
		fmt.Errorf("email can not be empty")
	}

	return nil
}
