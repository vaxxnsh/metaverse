package tests

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func waitForMessage(t *testing.T, conn *websocket.Conn) map[string]interface{} {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatal("failed to read websocket message:", err)
	}

	var data map[string]interface{}
	json.Unmarshal(msg, &data)
	return data
}

func TestWebsocketFlow(t *testing.T) {

	////////////////////////////////////////////////////
	//////////////////// HTTP SETUP ////////////////////
	////////////////////////////////////////////////////

	username := randomUsername()
	password := "123456"

	// Admin signup
	_, adminSignup := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username,
		"password": password,
		"type":     "admin",
	}, "")

	adminUserId := adminSignup["userId"].(string)

	_, adminSignin := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username,
		"password": password,
	}, "")

	adminToken := adminSignin["token"].(string)

	// User signup
	_, userSignup := doRequest(t, "POST", BACKEND_URL+"/api/v1/signup", map[string]interface{}{
		"username": username + "-user",
		"password": password,
		"type":     "user",
	}, "")

	userId := userSignup["userId"].(string)

	_, userSignin := doRequest(t, "POST", BACKEND_URL+"/api/v1/signin", map[string]interface{}{
		"username": username + "-user",
		"password": password,
	}, "")

	userToken := userSignin["token"].(string)

	// Create element
	_, el := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/element", map[string]interface{}{
		"imageUrl": "https://test.com/a.png",
		"width":    1,
		"height":   1,
		"static":   true,
	}, adminToken)

	elementId := el["id"].(string)

	// Create map
	_, mapResp := doRequest(t, "POST", BACKEND_URL+"/api/v1/admin/map", map[string]interface{}{
		"thumbnail":  "thumb.png",
		"dimensions": "100x200",
		"name":       "Default space",
		"defaultElements": []map[string]interface{}{
			{"elementId": elementId, "x": 20, "y": 20},
		},
	}, adminToken)

	mapId := mapResp["id"].(string)

	// Create space
	_, spaceResp := doRequest(t, "POST", BACKEND_URL+"/api/v1/space", map[string]interface{}{
		"name":       "Test",
		"dimensions": "100x200",
		"mapId":      mapId,
	}, userToken)

	spaceId := spaceResp["spaceId"].(string)

	////////////////////////////////////////////////////
	//////////////////// WS SETUP //////////////////////
	////////////////////////////////////////////////////

	u := url.URL{Scheme: "ws", Host: "localhost:3001", Path: "/"}

	ws1, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("ws1 connection failed:", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("ws2 connection failed:", err)
	}
	defer ws2.Close()

	////////////////////////////////////////////////////
	//////////////// JOIN SPACE ////////////////////////
	////////////////////////////////////////////////////

	ws1.WriteJSON(map[string]interface{}{
		"type": "join",
		"payload": map[string]interface{}{
			"spaceId": spaceId,
			"token":   adminToken,
		},
	})

	msg1 := waitForMessage(t, ws1)

	ws2.WriteJSON(map[string]interface{}{
		"type": "join",
		"payload": map[string]interface{}{
			"spaceId": spaceId,
			"token":   userToken,
		},
	})

	msg2 := waitForMessage(t, ws2)
	msg3 := waitForMessage(t, ws1)

	if msg1["type"] != "space-joined" {
		t.Fatal("admin should receive space-joined")
	}

	if msg2["type"] != "space-joined" {
		t.Fatal("user should receive space-joined")
	}

	if msg3["type"] != "user-joined" {
		t.Fatal("admin should receive user-joined event")
	}

	////////////////////////////////////////////////////
	//////////////// INVALID MOVE //////////////////////
	////////////////////////////////////////////////////

	ws1.WriteJSON(map[string]interface{}{
		"type": "move",
		"payload": map[string]interface{}{
			"x": 100000,
			"y": 100000,
		},
	})

	rejectMsg := waitForMessage(t, ws1)

	if rejectMsg["type"] != "movement-rejected" {
		t.Fatal("movement outside boundary should be rejected")
	}

	////////////////////////////////////////////////////
	//////////////// VALID MOVE ////////////////////////
	////////////////////////////////////////////////////

	ws1.WriteJSON(map[string]interface{}{
		"type": "move",
		"payload": map[string]interface{}{
			"x": 21,
			"y": 20,
		},
	})

	moveMsg := waitForMessage(t, ws2)

	if moveMsg["type"] != "movement" {
		t.Fatal("valid movement should broadcast")
	}

	////////////////////////////////////////////////////
	//////////////// USER LEAVE ////////////////////////
	////////////////////////////////////////////////////

	ws1.Close()

	leaveMsg := waitForMessage(t, ws2)

	if leaveMsg["type"] != "user-left" {
		t.Fatal("user-left event expected")
	}

	if leaveMsg["payload"].(map[string]interface{})["userId"] != adminUserId {
		t.Fatal("wrong user left event")
	}

	_ = userId
}
