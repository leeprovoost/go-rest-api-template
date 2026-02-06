package passport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandler() http.Handler {
	srv := NewTestServer()
	return srv.middleware(srv.routes())
}

// --- Health & readiness ---

func TestHealthcheckHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GNU Terry Pratchett", w.Header().Get("X-Clacks-Overhead"))
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "go-rest-api-template", body["appName"])
	assert.Equal(t, "0.0.0", body["version"])
}

func TestReadyHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body["status"])
}

// --- Users ---

func TestListUsersHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GNU Terry Pratchett", w.Header().Get("X-Clacks-Overhead"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(2), body["count"])
	assert.Equal(t, float64(2), body["total"])
	assert.Equal(t, float64(0), body["offset"])
	assert.Equal(t, float64(25), body["limit"])
	assert.NotNil(t, body["users"])
}

func TestListUsersHandlerPagination(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users?offset=0&limit=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(1), body["count"])
	assert.Equal(t, float64(2), body["total"])
	assert.Equal(t, float64(0), body["offset"])
	assert.Equal(t, float64(1), body["limit"])
}

func TestGetUserHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users/0", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "John", body["firstName"])
	assert.Equal(t, "Doe", body["lastName"])
}

func TestGetUserNotFound(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateUserHandler(t *testing.T) {
	handler := newTestHandler()
	body := `{"firstName":"Apple","lastName":"Jack","dateOfBirth":"1972-03-07T00:00:00Z","locationOfBirth":"Cambridge"}`
	r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	var user map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &user))
	assert.Equal(t, float64(2), user["id"])
	assert.Equal(t, "Apple", user["firstName"])
}

func TestCreateUserValidationFails(t *testing.T) {
	handler := newTestHandler()
	body := `{"firstName":"","lastName":"","dateOfBirth":"0001-01-01T00:00:00Z","locationOfBirth":""}`
	r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "validation failed", resp["message"])
	errors := resp["errors"].([]any)
	assert.Len(t, errors, 4)
}

func TestGetUserInvalidID(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid user id", resp["message"])
}

func TestCreateUserMalformedJSON(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{bad json`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "malformed user object", resp["message"])
}

func TestUpdateUserHandler(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":0,"firstName":"John","lastName":"Updated","dateOfBirth":"1985-12-31T00:00:00Z","locationOfBirth":"Manchester"}`
	r := httptest.NewRequest(http.MethodPut, "/users/0", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var user map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &user))
	assert.Equal(t, "Updated", user["lastName"])
	assert.Equal(t, "Manchester", user["locationOfBirth"])
}

func TestUpdateUserMalformedJSON(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodPut, "/users/0", strings.NewReader(`not json`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUserValidationFails(t *testing.T) {
	handler := newTestHandler()
	body := `{"firstName":"","lastName":"","dateOfBirth":"0001-01-01T00:00:00Z","locationOfBirth":""}`
	r := httptest.NewRequest(http.MethodPut, "/users/0", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "validation failed", resp["message"])
	assert.Len(t, resp["errors"].([]any), 4)
}

func TestUpdateUserStoreError(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":999,"firstName":"Ghost","lastName":"User","dateOfBirth":"1990-01-01T00:00:00Z","locationOfBirth":"Nowhere"}`
	r := httptest.NewRequest(http.MethodPut, "/users/999", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListUsersNegativeOffset(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users?offset=-5&limit=200", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(0), body["offset"])
	assert.Equal(t, float64(25), body["limit"]) // 200 > 100, so clamped to default 25
}

func TestDeleteUserHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserInvalidID(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserNotFound(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListUsersOffsetBeyondTotal(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users?offset=100", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(0), body["count"])
	assert.Equal(t, float64(2), body["total"])
}

// --- Passports ---

func TestListUserPassportsHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users/0/passports", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(1), body["count"])
}

func TestGetPassportHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/passports/012345678", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "012345678", body["id"])
	assert.Equal(t, "HMPO", body["authority"])
}

func TestGetPassportNotFound(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/passports/000000000", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreatePassportHandler(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":"111222333","dateOfIssue":"2024-01-01T00:00:00Z","dateOfExpiry":"2034-01-01T00:00:00Z","authority":"HMPO"}`
	r := httptest.NewRequest(http.MethodPost, "/users/0/passports", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	var passport map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &passport))
	assert.Equal(t, "111222333", passport["id"])
	assert.Equal(t, float64(0), passport["userId"])
}

func TestListUserPassportsInvalidID(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodGet, "/users/abc/passports", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePassportInvalidUserID(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":"999888777","dateOfIssue":"2024-01-01T00:00:00Z","dateOfExpiry":"2034-01-01T00:00:00Z","authority":"HMPO"}`
	r := httptest.NewRequest(http.MethodPost, "/users/abc/passports", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePassportMalformedJSON(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodPost, "/users/0/passports", strings.NewReader(`{bad`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePassportValidationFails(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":"","dateOfIssue":"0001-01-01T00:00:00Z","dateOfExpiry":"0001-01-01T00:00:00Z","authority":""}`
	r := httptest.NewRequest(http.MethodPost, "/users/0/passports", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "validation failed", resp["message"])
	assert.Len(t, resp["errors"].([]any), 4)
}

func TestCreatePassportDuplicate(t *testing.T) {
	handler := newTestHandler()
	body := `{"id":"012345678","dateOfIssue":"2024-01-01T00:00:00Z","dateOfExpiry":"2034-01-01T00:00:00Z","authority":"HMPO"}`
	r := httptest.NewRequest(http.MethodPost, "/users/0/passports", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdatePassportHandler(t *testing.T) {
	handler := newTestHandler()
	body := `{"dateOfIssue":"2021-06-01T00:00:00Z","dateOfExpiry":"2031-06-01T00:00:00Z","authority":"IPS"}`
	r := httptest.NewRequest(http.MethodPut, "/passports/012345678", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var passport map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &passport))
	assert.Equal(t, "012345678", passport["id"])
	assert.Equal(t, "IPS", passport["authority"])
}

func TestUpdatePassportMalformedJSON(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodPut, "/passports/012345678", strings.NewReader(`{bad`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassportValidationFails(t *testing.T) {
	handler := newTestHandler()
	body := `{"dateOfIssue":"0001-01-01T00:00:00Z","dateOfExpiry":"0001-01-01T00:00:00Z","authority":""}`
	r := httptest.NewRequest(http.MethodPut, "/passports/012345678", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUpdatePassportNotFound(t *testing.T) {
	handler := newTestHandler()
	body := `{"dateOfIssue":"2021-06-01T00:00:00Z","dateOfExpiry":"2031-06-01T00:00:00Z","authority":"IPS"}`
	r := httptest.NewRequest(http.MethodPut, "/passports/000000000", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeletePassportHandler(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodDelete, "/passports/012345678", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePassportNotFound(t *testing.T) {
	handler := newTestHandler()
	r := httptest.NewRequest(http.MethodDelete, "/passports/000000000", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
