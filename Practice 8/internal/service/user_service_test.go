package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"practice-8/internal/service/mocks"
	"practice-8/pkg/modules"

	"github.com/golang/mock/gomock"
)

func setupMockRepository(t testing.TB) *mocks.MockUserRepository {
	ctrl := gomock.NewController(t)
	return mocks.NewMockUserRepository(ctrl)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		req     modules.RegisterRequest
		wantErr bool
	}{
		{
			name: "successful registration",
			req: modules.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			req: modules.RegisterRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			req: modules.RegisterRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: modules.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "admin registration attempt",
			req: modules.RegisterRequest{
				Username: "adminuser",
				Email:    "admin@example.com",
				Password: "password123",
				Role:     "admin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupMockRepository(t)
			service := NewUserService(repo)

			_, err := service.Register(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "missing username",
			username: "",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "missing password",
			username: "testuser",
			password: "",
			wantErr:  true,
		},
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupMockRepository(t)
			service := NewUserService(repo)

			req := modules.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			_, _, err := service.Login(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		user    string
		email   string
		wantErr bool
	}{
		{
			name:    "successful update",
			id:      1,
			user:    "newuser",
			email:   "new@example.com",
			wantErr: false,
		},
		{
			name:    "invalid user id",
			id:      0,
			user:    "newuser",
			email:   "new@example.com",
			wantErr: true,
		},
		{
			name:    "missing username",
			id:      1,
			user:    "",
			email:   "new@example.com",
			wantErr: true,
		},
		{
			name:    "missing email",
			id:      1,
			user:    "newuser",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupMockRepository(t)
			service := NewUserService(repo)

			err := service.UpdateUser(tt.id, tt.user, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "successful delete",
			id:      1,
			wantErr: false,
		},
		{
			name:    "invalid user id",
			id:      0,
			wantErr: true,
		},
		{
			name:    "negative id",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupMockRepository(t)
			service := NewUserService(repo)

			err := service.DeleteUser(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRate(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		wantErr bool
	}{
		{
			name:    "missing from currency",
			from:    "",
			to:      "USD",
			wantErr: true,
		},
		{
			name:    "missing to currency",
			from:    "EUR",
			to:      "",
			wantErr: true,
		},
		{
			name:    "successful rate fetch",
			from:    "EUR",
			to:      "USD",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupMockRepository(t)
			service := NewUserService(repo)

			_, err := service.GetRate(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRateWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "latest") {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"base":"EUR","rates":{"USD":1.12,"GBP":0.86}}`)
	}))
	defer server.Close()

	repo := setupMockRepository(t)
	service := NewUserService(repo)

	_, err := service.GetRate("EUR", "USD")
	if err != nil {
		t.Errorf("GetRate() error = %v", err)
	}
}
