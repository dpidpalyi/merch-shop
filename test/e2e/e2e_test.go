//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"merch-shop/internal/config"
	"merch-shop/internal/dbinit"
	"merch-shop/internal/models"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AuthUser(t *testing.T, username, password string) (int, string) {
	httpHost := "http://localhost:8081"
	client := &http.Client{}

	authReq := models.AuthRequest{
		Username: username,
		Password: password,
	}

	authBody, err := json.Marshal(authReq)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", httpHost+"/api/auth", bytes.NewReader(authBody))
	assert.NoError(t, err)
	req.Header.Add("Content-type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, ""
	}

	authResponse := models.AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	assert.NoError(t, err)

	assert.NotEmpty(t, authResponse.Token)
	return resp.StatusCode, authResponse.Token
}

func Test_Auth_E2E(t *testing.T) {
	// Users not exists
	tests := []struct {
		name           string
		username       string
		password       string
		wantStatusCode int
		wantToken      bool
	}{
		{
			name:           "valid request",
			username:       "bob",
			password:       "password",
			wantStatusCode: http.StatusOK,
			wantToken:      true,
		},
		{
			name:           "valid request, we need second user",
			username:       "alice",
			password:       "password",
			wantStatusCode: http.StatusOK,
			wantToken:      true,
		},
		{
			name:           "invalid request, empty username",
			username:       "",
			password:       "password",
			wantStatusCode: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name:           "invalid request, empty password",
			username:       "fake",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name:           "invalid request, too long password",
			username:       "fake",
			password:       string(make([]byte, 76)),
			wantStatusCode: http.StatusBadRequest,
			wantToken:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatusCode, token := AuthUser(t, tt.username, tt.password)

			assert.Equal(t, tt.wantStatusCode, gotStatusCode)

			if tt.wantToken {
				assert.NotEmpty(t, token)
			} else {
				assert.Empty(t, token)
			}
		})
	}

	// Second test, after users add to DB
	// [0]bob fails, [1]alice success, others still fail
	tests[0].password = "wrong password"
	tests[0].wantStatusCode = http.StatusUnauthorized
	tests[0].wantToken = false

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatusCode, token := AuthUser(t, tt.username, tt.password)

			assert.Equal(t, tt.wantStatusCode, gotStatusCode)

			if tt.wantToken {
				assert.NotEmpty(t, token)
			} else {
				assert.Empty(t, token)
			}
		})
	}
}

func Test_BuyItem(t *testing.T) {
	// For tests default coins in DB for users set to 100

	httpHost := "http://localhost:8081"
	client := &http.Client{}

	os.Chdir("../..")
	cfg, err := config.New(".")
	assert.NoError(t, err)
	cfg.DB.Port = "5433"
	cfg.DB.Name = "shop_test"

	db, err := dbinit.OpenDB(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		username       string
		password       string
		item           string
		wantStatusCode int
		wantErr        bool
		errResponse    string
		wantBalance    int
	}{
		{
			name:           "valid request, success buy",
			username:       "carl",
			password:       "password",
			item:           "book", // 50
			wantStatusCode: http.StatusOK,
			wantErr:        false,
			wantBalance:    50,
		},
		{
			name:           "valid request, not enough coins",
			username:       "sarah",
			password:       "password",
			item:           "pink-hoody", // 500
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
			errResponse:    "not enough coins",
			wantBalance:    100,
		},
		{
			name:           "invalid request, item not exists",
			username:       "kenny",
			password:       "password",
			item:           "beer",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
			errResponse:    "item: record not found",
			wantBalance:    100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, token := AuthUser(t, tt.username, tt.password)

			reqURL := fmt.Sprintf("%s/api/buy/%s", httpHost, tt.item)
			req, err := http.NewRequest("GET", reqURL, nil)
			assert.NoError(t, err)

			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)

			if tt.wantErr {
				errorResponse := &models.ErrorResponse{}
				err := json.NewDecoder(resp.Body).Decode(errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.errResponse, errorResponse.Errors)
			}

			query := `
			    SELECT balance
			    FROM coins
			    JOIN active_users ON coins.user_id = active_users.id
			    WHERE active_users.username = $1`

			var balance int
			err = db.QueryRow(query, tt.username).Scan(&balance)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantBalance, balance)
		})
	}
}
