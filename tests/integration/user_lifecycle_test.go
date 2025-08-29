package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"userService/internal/events"
	"userService/tests/testutil"
)

// doRequest helper remains unchanged
func doRequest(t *testing.T, method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	testApp.ServeHTTP(w, req)
	return w
}

func TestUserLifecycle(t *testing.T) {
	userID := uuid.New()
	logger := logging.GetLogger()
	logger.Infof("testing user lifecycle %s", userID)

	// --- Step 1: Produce "UserCreated" Kafka event ---
	payload := events.UserCreatedPayload{
		UserID:    userID,
		Firstname: "Sayan",
		Lastname:  "Seksenbayev",
		Email:     "sayan123serv@gmail.com",
	}

	err := container.Producer.Produce(events.UserCreated, payload)
	require.NoError(t, err, "failed to produce user.created event")

	// --- Step 2: Run consumer in goroutine ---
	go container.Consumer.Start()

	// Always close the consumer after the test ends
	t.Cleanup(func() {
		logger.Info("Shutting down consumer...")
		container.Consumer.Close()
	})

	fmt.Printf("Waiting for user to become active %d seconds\n", 5)
	time.Sleep(5 * time.Second)

	// --- Step 3: Fetch created user ---
	w := doRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/users/%s", userID.String()), nil, nil)
	require.Equal(t, http.StatusOK, w.Code)

	// --- Step 4: Update user (with JWT token) ---
	token := testutil.GenerateMockToken(userID.String())
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Mandatory fields
	_ = writer.WriteField("firstname", "integration_tester")
	_ = writer.WriteField("lastname", "integration_tester")
	_ = writer.WriteField("email", "sayan123serv@gmail.com")

	// Optional fields
	_ = writer.WriteField("about", "Just an integration test user")
	_ = writer.WriteField("dateOfBirth", "02.01.2006")
	_ = writer.WriteField("gender", "male")
	_ = writer.WriteField("location", "Almaty, Kazakhstan")
	_ = writer.WriteField("socials[]", "https://twitter.com/tester")
	_ = writer.WriteField("socials[]", "https://github.com/tester")

	fileWriter, _ := writer.CreateFormFile("avatar", "img.png")

	if _, err := fileWriter.Write([]byte("fake image bytes")); err != nil {
		t.Fatalf("failed to write fake image bytes: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  writer.FormDataContentType(),
	}

	w = doRequest(t, http.MethodPut, "/api/v1/users/"+userID.String(), body, headers)
	require.Equal(t, http.StatusOK, w.Code, "expected 200 on update")

	// --- Step 5: Delete user (with JWT token) ---
	headers = map[string]string{
		"Authorization": "Bearer " + token,
	}
	w = doRequest(t, http.MethodDelete, "/api/v1/users/"+userID.String(), nil, headers)
	require.Equal(t, http.StatusOK, w.Code, "expected 200 on delete")

	// --- Step 6: Confirm deletion ---
	w = doRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/users/%s", userID.String()), nil, nil)
	require.Equal(t, http.StatusOK, w.Code, "expected 200 after fetch")

	// Parse response JSON into a map
	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "failed to unmarshal response")

}

func TestUserLifecycle_NeedsCompletion(t *testing.T) {
	logger := logging.GetLogger()
	userID := uuid.New()
	payload := events.UserCreatedPayload{
		UserID: userID,
		Email:  "sayan123serv@gmail.com",
	}

	err := container.Producer.Produce(events.UserCreated, payload)
	require.NoError(t, err, "failed to produce user.created event")

	// --- Step 2: Run consumer in goroutine ---
	go container.Consumer.Start()

	// Always close the consumer after the test ends
	t.Cleanup(func() {
		logger.Info("Shutting down consumer...")
		container.Consumer.Close()
	})

	fmt.Printf("Waiting for user to become active %d seconds\n", 5)
	time.Sleep(5 * time.Second)

	w := doRequest(t, http.MethodGet, "/api/v1/users/"+userID.String(), nil, nil)
	require.Equal(t, http.StatusOK, w.Code, "expected 200 on delete")

	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "failed to unmarshal response")

	// Assert that needs_completion is true
	needsCompletion, ok := resp["needs_completion"].(bool)
	require.True(t, ok, "needs_completion field missing or not a bool")
	require.True(t, needsCompletion, "needs_completion should be true")
}
