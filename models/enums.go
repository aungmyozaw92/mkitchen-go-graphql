package models

import (
	"errors"
	"io"
	"strconv"
)

type UserRole string

const (
	UserRoleAdmin  UserRole = "A"
	UserRoleOwner  UserRole = "O"
	UserRoleCustom UserRole = "C"
)

func (p UserRole) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(p))))
}

func (p *UserRole) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("user role must be string")
	}

	userRole := map[string]UserRole{
		"A": UserRoleAdmin,
		"O": UserRoleOwner,
		"C": UserRoleCustom,
	}

	*p, ok = userRole[str]
	if !ok {
		return errors.New("invalid user role")
	}
	return nil
}