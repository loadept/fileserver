package util

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/jsusmachaca/fileserver/pkg/auth"
)

func ValidateRequest(r io.Reader, body *auth.UserModel) error {
	err := json.NewDecoder(r).Decode(body)
	if err != nil {
		var errInvalid *json.UnmarshalTypeError
		if errors.As(err, &errInvalid) {
			return auth.ErrInvalidDataType
		}
		return err
	}
	if body.Username == "" && body.Email == "" {
		return auth.ErrFieldRequired("username or email")
	}
	if body.Password == "" {
		return auth.ErrFieldRequired("password")
	}

	return nil
}

func ValidateRegisterRequest(r io.Reader, body *auth.UserModel) error {
	err := json.NewDecoder(r).Decode(body)
	if err != nil {
		var errInvalid *json.UnmarshalTypeError
		if errors.As(err, &errInvalid) {
			return auth.ErrInvalidDataType
		}
		return err
	}

	if body.Username == "" {
		return auth.ErrFieldRequired("username")
	}
	if body.FirstName == "" {
		return auth.ErrFieldRequired("first_name")
	}
	if body.Email == "" {
		return auth.ErrFieldRequired("email")
	}
	if body.Password == "" {
		return auth.ErrFieldRequired("password")
	}
	return nil
}
