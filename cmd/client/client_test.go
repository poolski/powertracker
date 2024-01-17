package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gotest.tools/v3/assert"
)

func TestClient_Connect_SuccessfulConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the websocket upgrade request
		upgrader := websocket.Upgrader{}
		conn, _ := upgrader.Upgrade(w, r, nil)

		// Read the initial message
		initMsg := map[string]interface{}{
			"type": "init",
		}
		assert.NilError(t, conn.WriteJSON(initMsg), "write initial message failed")

		// Read the authentication message
		var authMsg map[string]interface{}
		assert.NilError(t, conn.ReadJSON(&authMsg), "read auth message failed")

		// Check the authentication message
		assert.Equal(t, authMsg["type"], "auth", "unexpected auth message type")
		assert.Equal(t, authMsg["access_token"], "test_token", "unexpected access token")

		// Send the authentication response
		authResp := map[string]interface{}{
			"type": "auth_ok",
		}
		assert.NilError(t, conn.WriteJSON(authResp), "write auth response failed")
	}))
	// Set up the client
	client := &Client{
		Config: Config{
			Insecure: true,
		},
	}

	// Set up the test environment
	viper.Set("url", s.URL)
	viper.Set("api_key", "test_token")

	// Call the Connect method
	err := client.Connect()

	// Check the error condition
	assert.NilError(t, err)
}

func TestClient_Connect_ErrorStates(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		apiKey   string
		expected string
	}{
		{
			name:     "Empty URL",
			url:      "",
			apiKey:   "test_token",
			expected: "url is required",
		},
		{
			name:     "Invalid URL",
			url:      "http://192.168.0.%31/",
			apiKey:   "test_token",
			expected: "parse \"http://192.168.0.%31/\": invalid URL escape \"%31\"",
		},
		{
			name:     "Malformed URL",
			url:      "htp:\\example.com",
			apiKey:   "test_token",
			expected: "dial: malformed ws or wss URL",
		},
		{
			name:     "Bad Handshake",
			url:      "http://example.com",
			apiKey:   "test_token",
			expected: "dial: websocket: bad handshake",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up the client
			client := &Client{
				Config: Config{
					Insecure: true,
				},
			}

			// Set up the test environment
			viper.Set("url", test.url)
			viper.Set("api_key", test.apiKey)

			// Call the Connect method
			err := client.Connect()

			// Check the error condition
			assert.ErrorContains(t, err, test.expected)
		})
	}
}
