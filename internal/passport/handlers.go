package passport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/leeprovoost/go-rest-api-template/pkg/health"
	"github.com/leeprovoost/go-rest-api-template/pkg/status"
)

// respond writes a JSON response with the given status code.
func respond(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// --- Health & readiness ---

func (s *Server) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, health.Check{
		AppName: "go-rest-api-template",
		Version: s.version,
	})
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- Users ---

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := s.userStore.ListUsers(r.Context())
	if err != nil {
		s.logger.Error("failed to list users", "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "failed to list users",
		})
		return
	}

	// Pagination
	total := len(list)
	offset, limit := parsePagination(r)
	if offset > len(list) {
		list = []models.User{}
	} else {
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		list = list[offset:end]
	}

	respond(w, http.StatusOK, map[string]any{
		"users":  list,
		"count":  len(list),
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "invalid user id",
		})
		return
	}
	user, err := s.userStore.GetUser(r.Context(), uid)
	if err != nil {
		s.logger.Error("user not found", "id", uid, "error", err)
		respond(w, http.StatusNotFound, status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find user",
		})
		return
	}
	respond(w, http.StatusOK, user)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.logger.Error("malformed user object", "error", err)
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed user object",
		})
		return
	}
	if errs := validateUser(u); len(errs) > 0 {
		respond(w, http.StatusUnprocessableEntity, status.Response{
			Status:  strconv.Itoa(http.StatusUnprocessableEntity),
			Message: "validation failed",
			Errors:  errs,
		})
		return
	}
	u.ID = -1 // will be assigned by store
	user, _ := s.userStore.AddUser(r.Context(), u)
	respond(w, http.StatusCreated, user)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.logger.Error("malformed user object", "error", err)
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed user object",
		})
		return
	}
	if errs := validateUser(u); len(errs) > 0 {
		respond(w, http.StatusUnprocessableEntity, status.Response{
			Status:  strconv.Itoa(http.StatusUnprocessableEntity),
			Message: "validation failed",
			Errors:  errs,
		})
		return
	}
	user, err := s.userStore.UpdateUser(r.Context(), u)
	if err != nil {
		s.logger.Error("failed to update user", "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		})
		return
	}
	respond(w, http.StatusOK, user)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "invalid user id",
		})
		return
	}
	if err := s.userStore.DeleteUser(r.Context(), uid); err != nil {
		s.logger.Error("failed to delete user", "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Passports ---

func (s *Server) handleListUserPassports(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.PathValue("uid"))
	if err != nil {
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "invalid user id",
		})
		return
	}
	passports, err := s.passportStore.ListPassportsByUser(r.Context(), uid)
	if err != nil {
		s.logger.Error("failed to list passports", "userId", uid, "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "failed to list passports",
		})
		return
	}
	respond(w, http.StatusOK, map[string]any{
		"passports": passports,
		"count":     len(passports),
	})
}

func (s *Server) handleGetPassport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	passport, err := s.passportStore.GetPassport(r.Context(), id)
	if err != nil {
		s.logger.Error("passport not found", "id", id, "error", err)
		respond(w, http.StatusNotFound, status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find passport",
		})
		return
	}
	respond(w, http.StatusOK, passport)
}

func (s *Server) handleCreatePassport(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.PathValue("uid"))
	if err != nil {
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "invalid user id",
		})
		return
	}
	var p models.Passport
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.logger.Error("malformed passport object", "error", err)
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed passport object",
		})
		return
	}
	p.UserID = uid
	if errs := validatePassport(p); len(errs) > 0 {
		respond(w, http.StatusUnprocessableEntity, status.Response{
			Status:  strconv.Itoa(http.StatusUnprocessableEntity),
			Message: "validation failed",
			Errors:  errs,
		})
		return
	}
	passport, err := s.passportStore.AddPassport(r.Context(), p)
	if err != nil {
		s.logger.Error("failed to create passport", "error", err)
		respond(w, http.StatusConflict, status.Response{
			Status:  strconv.Itoa(http.StatusConflict),
			Message: err.Error(),
		})
		return
	}
	respond(w, http.StatusCreated, passport)
}

func (s *Server) handleUpdatePassport(w http.ResponseWriter, r *http.Request) {
	var p models.Passport
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.logger.Error("malformed passport object", "error", err)
		respond(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed passport object",
		})
		return
	}
	p.ID = r.PathValue("id")
	if errs := validatePassport(p); len(errs) > 0 {
		respond(w, http.StatusUnprocessableEntity, status.Response{
			Status:  strconv.Itoa(http.StatusUnprocessableEntity),
			Message: "validation failed",
			Errors:  errs,
		})
		return
	}
	passport, err := s.passportStore.UpdatePassport(r.Context(), p)
	if err != nil {
		s.logger.Error("failed to update passport", "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		})
		return
	}
	respond(w, http.StatusOK, passport)
}

func (s *Server) handleDeletePassport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.passportStore.DeletePassport(r.Context(), id); err != nil {
		s.logger.Error("failed to delete passport", "error", err)
		respond(w, http.StatusInternalServerError, status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Validation ---

func validateUser(u models.User) []string {
	var errs []string
	if u.FirstName == "" {
		errs = append(errs, "firstName is required")
	}
	if u.LastName == "" {
		errs = append(errs, "lastName is required")
	}
	if u.DateOfBirth.IsZero() {
		errs = append(errs, "dateOfBirth is required")
	}
	if u.LocationOfBirth == "" {
		errs = append(errs, "locationOfBirth is required")
	}
	return errs
}

func validatePassport(p models.Passport) []string {
	var errs []string
	if p.ID == "" {
		errs = append(errs, "id is required")
	}
	if p.DateOfIssue.IsZero() {
		errs = append(errs, "dateOfIssue is required")
	}
	if p.DateOfExpiry.IsZero() {
		errs = append(errs, "dateOfExpiry is required")
	}
	if p.Authority == "" {
		errs = append(errs, "authority is required")
	}
	return errs
}

// --- Helpers ---

func parsePagination(r *http.Request) (offset, limit int) {
	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return offset, limit
}
