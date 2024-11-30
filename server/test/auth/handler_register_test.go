package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jsusmachaca/fileserver/api/handler"
	"github.com/jsusmachaca/fileserver/internal/database"
	"github.com/jsusmachaca/fileserver/pkg/auth"
)

func TestRegisterHandle(t *testing.T) {
	os.Setenv("DB_NAME", "db.test")
	defer os.Unsetenv("DB_NAME")
	defer func() {
		err := os.Remove("db.test")
		if err != nil {
			t.Errorf("Error to delete test database: %v", err)
		}
	}()

	db, _ := database.GetConnection()
	defer db.Close()
	database.Migrate(db)

	h := handler.Register{DB: db}

	t.Run("Valid Register", func(t *testing.T) {
		user := auth.UserModel{
			FirstName: "John Doe",
			Username:  "example",
			Email:     "example@example.com",
			Password:  "superstrongp@ssword",
		}
		body, err := json.Marshal(user)
		if err != nil {
			t.Errorf("Error at trying create body: %v\n", err)
		}

		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v\n", status, http.StatusOK)
		}

		expected := "{\"message\":\"Success\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})

	t.Run("Existend User", func(t *testing.T) {
		user := auth.UserModel{
			FirstName: "John Doe",
			Username:  "example",
			Email:     "example@example.com",
			Password:  "superstrongp@ssword",
		}
		body, err := json.Marshal(user)
		if err != nil {
			t.Errorf("Error at trying create body: %v\n", err)
		}

		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v\n", status, http.StatusConflict)
		}

		expected := "{\"error\":\"Registration is not possible at this time\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})

	t.Run("Incomplete Field", func(t *testing.T) {
		user := auth.UserModel{
			FirstName: "John Doe",
			Email:     "example@example.com",
			Password:  "superstrongp@ssword",
		}
		body, err := json.Marshal(user)
		if err != nil {
			t.Errorf("Error at trying create body: %v\n", err)
		}

		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v\n", status, http.StatusBadRequest)
		}

		expected := "{\"error\":\"This field is required: username\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})

	t.Run("Invalid Data Type", func(t *testing.T) {
		user := map[string]any{
			"first_name": 1,
			"username":   "example",
			"email":      "example@example.com",
			"password":   "superstrongp@ssword",
		}
		body, err := json.Marshal(user)
		if err != nil {
			t.Errorf("Error at trying create body: %v\n", err)
		}

		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v\n", status, http.StatusBadRequest)
		}

		expected := "{\"error\":\"Incorrect request data\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})
}
