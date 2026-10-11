package users

import (
	"encoding/json"
	"errors"
	"go-users-auth-template/internal/httpx"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Define router which contain all user-logic endpoint
func (uh *UserHandlers) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", uh.getAllUsers)
	r.Post("/", uh.registerUser)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", uh.getUserProfile)
		r.Put("/", uh.updateUserProfile)
		r.Delete("/", uh.deleteUser)
	})

	return r
}

type UserServiceAPI interface {
	RegisterUser(nickname, email, password string) (*User, error)
	GetAllUserProfiles() ([]*User, error)
	GetUserProfile(id int) (*User, error)
	UpdateUserProfile(user User) error
	DeleteUserProfile(id int) error
}

type UserHandlers struct {
	us *UserService
}

func NewUserHandlers(us *UserService) *UserHandlers {
	return &UserHandlers{us}
}

func (uh *UserHandlers) registerUser(w http.ResponseWriter, r *http.Request) {
	var reg RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		handleError(w, err)
		return
	}

	user, err := uh.us.RegisterUser(reg.Nickname, reg.Email, reg.Password)

	if err != nil {
		handleError(w, err)
		return
	}

	if err := httpx.WriteJson(w, http.StatusOK, user); err != nil {
		httpx.LogError("failer to create user while register", err, r)
	}
}

func (uh *UserHandlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	if users, err := uh.us.GetAllUserProfiles(); err != nil {
		handleError(w, err)
	} else {
		if err := httpx.WriteJson(w, http.StatusOK, users); err != nil {
			httpx.LogError("failed to write HTTP response", err, r)
		}
	}
}

func (uh *UserHandlers) getUserProfile(w http.ResponseWriter, r *http.Request) {
	id, err := handleUserID(r)
	if err != nil {
		handleError(w, ErrUserIDIdentifier)
		return
	}

	user, err := uh.us.GetUserProfile(id)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := httpx.WriteJson(w, http.StatusOK, user); err != nil {
		httpx.LogError("failed to write HTTP response", err, r)
	}
}

func (uh *UserHandlers) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	id, err := handleUserID(r)
	if err != nil {
		handleError(w, ErrUserIDIdentifier)
		return
	}

	var userUpdatingData UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdatingData); err != nil {
		handleError(w, ErrInvalidRequestFormat)
		return
	}

	if err := uh.us.UpdateUserProfile(id, userUpdatingData); err != nil {
		handleError(w, err)
		return
	}

	if err := httpx.WriteJson(w, http.StatusOK, userUpdatingData); err != nil {
		httpx.LogError("failed to write HTTP response", err, r)
	}
}

func (uh *UserHandlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := handleUserID(r)
	if err != nil {
		handleError(w, ErrUserIDIdentifier)
		return
	}

	if err := uh.us.DeleteUserProfile(id); err != nil {
		handleError(w, err)
		return
	}

	if err := httpx.WriteJson(w, http.StatusOK, "Profile deleted successfully"); err != nil {
		httpx.LogError("failed to write HTTP response", err, r)
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "user note found")
	case errors.Is(err, ErrInvalidRequestFormat):
		httpx.WriteError(w, http.StatusBadRequest, "invalid request")
	case errors.Is(err, ErrInvalidEmail):
		httpx.WriteError(w, http.StatusBadRequest, "Invalid user email")
	case errors.Is(err, ErrEmailTaken):
		httpx.WriteError(w, http.StatusBadRequest, "Email already taken")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func handleUserID(r *http.Request) (int, error) {
	strId := chi.URLParam(r, "id")
	return strconv.Atoi(strId)
}
