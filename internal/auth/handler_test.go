package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginAndMe(t *testing.T) {
	service, err := NewService("password123")
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	secret := []byte("test-secret-must-have-at-least-32-chars")
	handler := NewHandler(service, secret)
	router := newTestRouter(handler, secret)

	loginRequest := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"worker@fieldtask.com","password":"password123"}`))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginRequest)

	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", loginRecorder.Code, loginRecorder.Body.String())
	}

	var login loginResponse
	if err := json.Unmarshal(loginRecorder.Body.Bytes(), &login); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if login.AccessToken == "" || login.User.ID != "user-1" {
		t.Fatalf("unexpected login response: %+v", login)
	}

	meRequest := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	meRequest.Header.Set("Authorization", "Bearer "+login.AccessToken)
	meRecorder := httptest.NewRecorder()
	router.ServeHTTP(meRecorder, meRequest)

	if meRecorder.Code != http.StatusOK {
		t.Fatalf("me status = %d, body = %s", meRecorder.Code, meRecorder.Body.String())
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	service, err := NewService("password123")
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	secret := []byte("test-secret-must-have-at-least-32-chars")
	router := newTestRouter(NewHandler(service, secret), secret)
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"worker@fieldtask.com","password":"wrong-password"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("login status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
