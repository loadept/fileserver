package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jsusmachaca/fileserver/api/handler"
	"github.com/jsusmachaca/fileserver/internal/database"
	"github.com/jsusmachaca/fileserver/internal/util"
	"github.com/jsusmachaca/fileserver/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandle(t *testing.T) {
	os.Setenv("DB_NAME", "db.test")
	defer os.Unsetenv("DB_NAME")
	defer func() {
		err := os.Remove("db.test")
		if err != nil {
			t.Errorf("Error to delete test database: %v", err)
		}
	}()

	db, _ := database.GetConnection()
	database.Migrate(db)

	query := `INSERT INTO user(
		id, first_name,
		username,
		email,
		password
	) VALUES (?, ?, ?, ?, ?)
	ON CONFLICT (username)
	DO NOTHING`

	stmt, err := db.Prepare(query)
	if err != nil {
		t.Errorf("Error at insert test user: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("superstrongp@ssword"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Error at create hash: %v", err)
	}
	_, err = stmt.Exec(uuid.NewString(), "John Doe", "example", "example@example.com", string(hashedPassword))
	if err != nil {
		t.Errorf("Error at insert test user: %v", err)
	}
	defer db.Close()

	h := handler.Login{DB: db}

	user := auth.UserModel{
		Username: "example",
		Password: "superstrongp@ssword",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Error at trying create body: %v\n", err)
	}

	t.Run("Valid Login", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v\n", status, http.StatusOK)
		}

		var respone util.ResponseWithToken
		err = json.Unmarshal(rr.Body.Bytes(), &respone)
		if err != nil {
			t.Errorf("handler returned unexpected body: got %s want %s\n", rr.Body.String(), err)
		}
		if respone.Message != "Session started successfully" || respone.Token == "" {
			t.Errorf("Unexpected response: got %+v", respone)
		}
	})
}
