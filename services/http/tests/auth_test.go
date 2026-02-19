package tests

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"
)

const BACKEND_URL = "http://localhost:3000"

func randomUsername() string {
	return "vaxxnsh-" + string(rune(rand.Intn(1000000)))
}

func postRequest(t *testing.T, url string, body interface{}) (*http.Response, map[string]interface{}) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	var data map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&data)

	return resp, data
}

func TestAuthentication(t *testing.T) {

	t.Run("User is able to sign up only once", func(t *testing.T) {
		username := randomUsername()
		password := "123456"

		resp, _ := postRequest(t, BACKEND_URL+"/api/v1/signup", map[string]interface{}{
			"username": username,
			"password": password,
			"type":     "admin",
		})

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		resp2, _ := postRequest(t, BACKEND_URL+"/api/v1/signup", map[string]interface{}{
			"username": username,
			"password": password,
			"type":     "admin",
		})

		if resp2.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp2.StatusCode)
		}
	})

	t.Run("Signup fails if username is empty", func(t *testing.T) {
		password := "123456"

		resp, _ := postRequest(t, BACKEND_URL+"/api/v1/signup", map[string]interface{}{
			"password": password,
		})

		if resp.StatusCode != 400 {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Signin succeeds with correct credentials", func(t *testing.T) {
		username := randomUsername()
		password := "123456"

		postRequest(t, BACKEND_URL+"/api/v1/signup", map[string]interface{}{
			"username": username,
			"password": password,
			"type":     "admin",
		})

		resp, data := postRequest(t, BACKEND_URL+"/api/v1/signin", map[string]interface{}{
			"username": username,
			"password": password,
		})

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		if data["token"] == nil {
			t.Errorf("expected token to be defined")
		}
	})

	t.Run("Signin fails with incorrect credentials", func(t *testing.T) {
		username := randomUsername()
		password := "123456"

		postRequest(t, BACKEND_URL+"/api/v1/signup", map[string]interface{}{
			"username": username,
			"password": password,
			"role":     "admin",
		})

		resp, _ := postRequest(t, BACKEND_URL+"/api/v1/signin", map[string]interface{}{
			"username": "WrongUsername",
			"password": password,
		})

		if resp.StatusCode != 403 {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
	})
}
