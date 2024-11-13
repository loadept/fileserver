package handler

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/jsusmachaca/fileserver/internal/util"
	"github.com/jsusmachaca/fileserver/pkg/auth"
	"github.com/jsusmachaca/go-router/pkg/response"
	"github.com/mattn/go-sqlite3"
)

type Login struct {
	DB *sql.DB
}
type Register struct {
	DB *sql.DB
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body auth.UserModel
	authRep := &auth.AuthRepository{DB: h.DB}

	err := util.ValidateRequest(r.Body, &body)
	if err != nil {
		log.Printf("Error validating request: %v", err)

		if err == io.EOF {
			response.JsonErrorFromString(w, "Request body is required", http.StatusBadRequest)
			return
		}
		var fieldErr *auth.FieldRequiredError
		if errors.As(err, &fieldErr) {
			response.JsonErrorFromString(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, auth.ErrInvalidDataType) {
			response.JsonErrorFromString(w, "Incorrect request data", http.StatusBadRequest)
			return
		}
		response.JsonErrorFromString(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data, err := authRep.GetUser(body)
	if err != nil {
		log.Printf("Error querying by user: %v", err)

		if errors.Is(err, auth.ErrIncorrectCredentials) {
			response.JsonErrorFromString(w, "Incorrect credentials", http.StatusUnauthorized)
			return
		}
		response.JsonErrorFromString(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	token, err := util.CreateToken(data.ID, data.Username)
	if err != nil {
		response.JsonErrorFromString(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	message := util.ResponseWithToken{
		Message: "Session started successfully",
		Token:   token,
	}
	log.Printf("A user has logged in: %s", r.RemoteAddr)
	response.JsonResponse(w, message, http.StatusOK)
}

func (h *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body auth.UserModel
	authRep := &auth.AuthRepository{DB: h.DB}

	err := util.ValidateRegisterRequest(r.Body, &body)
	if err != nil {
		log.Printf("Error validating request: %v", err)

		if err == io.EOF {
			response.JsonErrorFromString(w, "Request body is required", http.StatusBadRequest)
			return
		}
		var fieldErr *auth.FieldRequiredError
		if errors.As(err, &fieldErr) {
			response.JsonErrorFromString(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, auth.ErrInvalidDataType) {
			response.JsonErrorFromString(w, "Incorrect request data", http.StatusBadRequest)
			return
		}
		response.JsonErrorFromString(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = authRep.RegisterUser(body)
	if err != nil {
		log.Printf("Error registering user: %v", err)

		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			response.JsonErrorFromString(w, "Registration is not possible at this time", http.StatusConflict)
			return
		}
		response.JsonErrorFromString(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	response.JsonResponse(w, map[string]string{
		"message": "Success",
	}, http.StatusOK)
}
