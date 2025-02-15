//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"merch-shop/internal/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_E2E(t *testing.T) {
	httpHost := "http://localhost:8081"

	client := &http.Client{}

	authReq := models.AuthRequest{
		Username: "bob",
		Password: "password",
	}

	authBody, err := json.Marshal(authReq)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", httpHost+"/api/auth", bytes.NewReader(authBody))
	req.Header.Add("Content-type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	authResponse := models.AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	assert.NoError(t, err)

	assert.NotEmpty(t, authResponse.Token)
}
