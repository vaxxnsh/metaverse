package tests

import "testing"

func TestAdminEndpoints(t *testing.T) {

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

	t.Run("User cannot access admin endpoints", func(t *testing.T) {

		resp1, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
			"imageUrl": "https://test.com/a.png",
			"width":    1,
			"height":   1,
			"static":   true,
		}, userToken)

		resp2, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/map", map[string]interface{}{
			"thumbnail":       "thumb.png",
			"dimensions":      "100x200",
			"name":            "test space",
			"defaultElements": []interface{}{},
		}, userToken)

		resp3, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/avatar", map[string]interface{}{
			"imageUrl": "https://test.com/avatar.png",
			"name":     "Timmy",
		}, userToken)

		resp4, _ := doRequest(t, "PUT", BACKEND_URL+"/api/v1/admin/element/123", map[string]interface{}{
			"imageUrl": "https://test.com/new.png",
		}, userToken)

		if resp1.StatusCode != 403 ||
			resp2.StatusCode != 403 ||
			resp3.StatusCode != 403 ||
			resp4.StatusCode != 403 {
			t.Fatal("user should not access admin endpoints")
		}
	})

	t.Run("Admin can access admin endpoints", func(t *testing.T) {

		resp1, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
			"imageUrl": "https://test.com/a.png",
			"width":    1,
			"height":   1,
			"static":   true,
		}, adminToken)

		resp2, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/map", map[string]interface{}{
			"thumbnail":       "thumb.png",
			"name":            "Space",
			"dimensions":      "100x200",
			"defaultElements": []interface{}{},
		}, adminToken)

		resp3, _ := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/avatar", map[string]interface{}{
			"imageUrl": "https://test.com/avatar.png",
			"name":     "Timmy",
		}, adminToken)

		if resp1.StatusCode != 200 ||
			resp2.StatusCode != 200 ||
			resp3.StatusCode != 200 {
			t.Fatal("admin should access admin endpoints")
		}
	})

	t.Run("Admin can update element imageUrl", func(t *testing.T) {

		_, elementData := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
			"imageUrl": "https://test.com/original.png",
			"width":    1,
			"height":   1,
			"static":   true,
		}, adminToken)

		elementId := elementData["id"].(string)

		resp, _ := doRequest(t, "PUT", BACKEND_URL+"/api/v1/admin/element/"+elementId, map[string]interface{}{
			"imageUrl": "https://test.com/updated.png",
		}, adminToken)

		if resp.StatusCode != 200 {
			t.Fatal("admin should be able to update element")
		}
	})

	_ = adminId
	_ = userId
}
