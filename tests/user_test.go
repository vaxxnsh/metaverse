package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func doRequest(t *testing.T, method, url string, body any, token string) (*http.Response, map[string]any) {
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

	var result map[string]any
	json.Unmarshal(respBody, &result)

	return resp, result
}

func TestSignupOnlyOnce(t *testing.T) {
	username := randomUsername()
	password := "123456"

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	resp2, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
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

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	resp, data := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]any{
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

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]any{
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

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	_, signinData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]any{
		"username": username,
		"password": password,
	}, "")

	token := signinData["token"].(string)

	_, avatarData := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/avatar", map[string]any{
		"imageUrl": "https://test.com/image.png",
		"name":     "Timmy",
	}, token)

	avatarId := avatarData["avatarId"].(string)

	resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/user/metadata", map[string]any{
		"avatarId": avatarId,
	}, token)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func TestCreateAndDeleteSpace(t *testing.T) {
	username := randomUsername()
	password := "123456"

	doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]any{
		"username": username,
		"password": password,
		"type":     "user",
	}, "")

	_, signinData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]any{
		"username": username,
		"password": password,
	}, "")

	token := signinData["token"].(string)

	_, spaceData := doRequest(t, "POST", BACKEND_URL+"/api/v1/space", map[string]any{
		"name":       "Test",
		"dimensions": "100x200",
	}, token)

	spaceId := spaceData["spaceId"].(string)

	resp, _ := doRequest(t, "DELETE", BACKEND_URL+"/api/v1/space/"+spaceId, nil, token)

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func TestArenaEndpoints(t *testing.T) {

	username := randomUsername()
	password := "123456"

	_, signupData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	adminId := signupData["userId"]

	_, signinData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username,
		"password": password,
	}, "")

	adminToken := signinData["token"].(string)

	_, userSignupData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username + "-user",
		"password": password,
		"type":     "user",
	}, "")

	userId := userSignupData["userId"]

	_, userSigninData := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username + "-user",
		"password": password,
	}, "")

	userToken := userSigninData["token"].(string)

	_, el1 := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
		"imageUrl": "https://test.com/a.png",
		"width":    1,
		"height":   1,
		"static":   true,
	}, adminToken)

	element1Id := el1["id"].(string)

	_, el2 := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
		"imageUrl": "https://test.com/b.png",
		"width":    1,
		"height":   1,
		"static":   true,
	}, adminToken)

	element2Id := el2["id"].(string)

	_, mapData := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/map", map[string]interface{}{
		"thumbnail":  "https://thumbnail.com/a.png",
		"dimensions": "100x200",
		"name":       "Default space",
		"defaultElements": []map[string]interface{}{
			{"elementId": element1Id, "x": 20, "y": 20},
			{"elementId": element1Id, "x": 18, "y": 20},
			{"elementId": element2Id, "x": 19, "y": 20},
		},
	}, adminToken)

	mapId := mapData["id"].(string)

	_, spaceData := doRequest(t, "POST", BACKEND_URL+"/api/v1/space", map[string]interface{}{
		"name":       "Test",
		"dimensions": "100x200",
		"mapId":      mapId,
	}, userToken)

	spaceId := spaceData["spaceId"].(string)

	t.Run("Incorrect spaceId returns 400", func(t *testing.T) {
		resp, _ := doRequest(t, "GET", BACKEND_URL+"/api/v1/space/123invalid", nil, userToken)

		if resp.StatusCode != 400 {
			t.Fatalf("expected 400 got %d", resp.StatusCode)
		}
	})

	t.Run("Correct spaceId returns all elements", func(t *testing.T) {
		_, data := doRequest(t, "GET", BACKEND_URL+"/api/v1/space/"+spaceId, nil, userToken)

		if data["dimensions"] != "100x200" {
			t.Fatal("dimensions mismatch")
		}

		elements := data["elements"].([]interface{})
		if len(elements) != 3 {
			t.Fatalf("expected 3 elements got %d", len(elements))
		}
	})

	t.Run("Delete endpoint deletes element", func(t *testing.T) {
		_, data := doRequest(t, "GET", BACKEND_URL+"/api/v1/space/"+spaceId, nil, userToken)
		elements := data["elements"].([]interface{})
		firstElement := elements[0].(map[string]interface{})
		elementId := firstElement["id"].(string)

		doRequest(t, "DELETE", BACKEND_URL+"/api/v1/space/element", map[string]interface{}{
			"id": elementId,
		}, userToken)

		_, newData := doRequest(t, "GET", BACKEND_URL+"/api/v1/space/"+spaceId, nil, userToken)
		newElements := newData["elements"].([]interface{})

		if len(newElements) != 2 {
			t.Fatalf("expected 2 elements got %d", len(newElements))
		}
	})

	t.Run("Adding element outside dimensions fails", func(t *testing.T) {
		resp, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/space/element", map[string]interface{}{
			"elementId": element1Id,
			"spaceId":   spaceId,
			"x":         10000,
			"y":         210000,
		}, userToken)

		if resp.StatusCode != 400 {
			t.Fatalf("expected 400 got %d", resp.StatusCode)
		}
	})

	t.Run("Adding element works", func(t *testing.T) {
		doRequest(t, "POST", BACKEND_URL+"/api/v1/space/element", map[string]interface{}{
			"elementId": element1Id,
			"spaceId":   spaceId,
			"x":         50,
			"y":         20,
		}, userToken)

		_, newData := doRequest(t, "GET", BACKEND_URL+"/api/v1/space/"+spaceId, nil, userToken)
		newElements := newData["elements"].([]interface{})

		if len(newElements) != 3 {
			t.Fatalf("expected 3 elements got %d", len(newElements))
		}
	})

	_ = adminId
	_ = userId
}
