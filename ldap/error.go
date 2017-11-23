package ldap

import "fmt"

var (
	UserAlreadyExists  = fmt.Errorf("User already exists")
	GroupAlreadyExists = fmt.Errorf("User already exists")

	UserNotFound  = fmt.Errorf("User not found")
	GroupNotFound = fmt.Errorf("Group not found")

	UserAuthFailure = fmt.Errorf("User authentication failed")
)
