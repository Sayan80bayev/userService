package integration

import (
	"bytes"
	"context"
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
	// Use background ctx only for producing the event; the shared consumer is started in TestMain.
	ctx := context.Background()

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

	err := container.Producer.Produce(ctx, events.UserCreated, payload)
	require.NoError(t, err, "failed to produce user.created event")

	// Give the system some time to process the event (polled checks are used later for important asserts)
	fmt.Printf("Waiting briefly for consumer to pick up event\n")
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

	// --- Step 6: Confirm fetch still works (or check a status) ---
	w = doRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/users/%s", userID.String()), nil, nil)
	require.Equal(t, http.StatusOK, w.Code, "expected 200 after fetch")

	// Parse response JSON into a map
	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "failed to unmarshal response")
}

func TestUserLifecycle_NeedsCompletion(t *testing.T) {
	// Produce event using a cancellable context (so produce call can be cancelled if needed).
	ctx := context.Background()

	userID := uuid.New()
	payload := events.UserCreatedPayload{
		UserID: userID,
		Email:  "sayan123serv@gmail.com",
	}

	// Produce event
	err := container.Producer.Produce(ctx, events.UserCreated, payload)
	require.NoError(t, err)

	// Wait until user is processed — poll for eventual consistency
	waitForUserNeedsCompletion(t, userID, 10*time.Second)
}

// waitForUserNeedsCompletion polls the API until the needs_completion field becomes true or timeout expires.
func waitForUserNeedsCompletion(t *testing.T, userID uuid.UUID, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		w := doRequest(t, http.MethodGet, "/api/v1/users/"+userID.String(), nil, nil)
		if w.Code == http.StatusOK {
			var resp map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err == nil {
				if needsCompletion, ok := resp["needs_completion"].(bool); ok && needsCompletion {
					return // ✅ condition met
				}
			}
		}
		time.Sleep(200 * time.Millisecond) // retry
	}
	t.Fatalf("user %s did not reach needs_completion=true within %s", userID, timeout)
}
