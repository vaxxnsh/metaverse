package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func doRequest(t *testing.T, method, url string, body interface{}, token string) (*http.Response, map[string]interface{}) {
	var reqBody io.Reader

	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return resp, result
}

func TestSignupOnlyOnce(t *testing.T) {
	username := randomUsername()
	password := "123456"

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	resp2, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	if resp2.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", resp2.StatusCode)
	}
}

func TestSigninSuccess(t *testing.T) {
	username := randomUsername()
	password := "123456"

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	resp, data := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username,
		"password": password,
	}, "")

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	if data["token"] == nil {
		t.Fatal("expected token in response")
	}
}

func TestSigninFailure(t *testing.T) {
	username := randomUsername()
	password := "123456"

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": "wrongUser",
		"password": password,
	}, "")

	if resp.StatusCode != 403 {
		t.Fatalf("expected 403 got %d", resp.StatusCode)
	}
}

func TestUserMetadataUpdate(t *testing.T) {
	username := randomUsername()
	password := "123456"

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	_, signinData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username,
		"password": password,
	}, "")

	token := signinData["token"].(string)

	_, avatarData := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/avatar", map[string]interface{}{
		"imageUrl": "https://test.com/image.png",
		"name":     "Timmy",
	}, token)

	avatarId := avatarData["avatarId"].(string)

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/user/metadata", map[string]interface{}{
		"avatarId": avatarId,
	}, token)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func TestCreateAndDeleteSpace(t *testing.T) {
	username := randomUsername()
	password := "123456"

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "user",
	}, "")

	_, signinData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username,
		"password": password,
	}, "")

	token := signinData["token"].(string)

	_, spaceData := doRequest(t, "POST", BACKEND_URL+"/api/v1/space", map[string]interface{}{
		"name":       "Test",
		"dimensions": "100x200",
	}, token)

	spaceId := spaceData["spaceId"].(string)

	resp, _ := doRequest(t, "DELETE", BACKEND_URL+"/api/v1/space/"+spaceId, nil, token)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}
