package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/rest"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/mocks"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/service"
)

// TODO: посмотреть тесты

func TestSignUp_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupApp(mockUser, mockValidator)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
}

func TestSignUp_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
	mockValidator.EXPECT().FormatValidationErrors(gomock.Any()).Return(map[string]string{"username": "required"})

	body := `{"username": "", "password": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
}

func TestSignUp_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(service.ErrUserExists)

	body := `{"username": "test", "password": "qwerty123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
}

func TestSignUp_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	body := `{"username": "test", "password": "qwerty123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
}

func TestSignUp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"username": "newuser", "password": "qwerty123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func setupApp(u rest.User, v rest.Validator) *fiber.App {
	app := fiber.New()
	h := rest.NewUserHandler(slogdiscard.NewDiscardLogger(), u, v)
	app.Post("/api/auth/sign-up", h.SignUp)
	return app
}

func TestSignIn_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupSignInApp(mockUser, mockValidator)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
	assert.NotEmpty(t, r.Data)
}

func TestSignIn_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupSignInApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
	mockValidator.EXPECT().FormatValidationErrors(gomock.Any()).Return(map[string]string{"username": "required"})

	body := `{"username": "", "password": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
	assert.Equal(t, map[string]interface{}{"username": "required"}, r.Data)
}

func TestSignIn_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupSignInApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().Login(gomock.Any(), gomock.Any()).Return("", service.ErrInvalidCredentials)

	body := `{"username": "wronguser", "password": "pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
	assert.Equal(t, "invalid credentials", r.Data)
}

func TestSignIn_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupSignInApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().Login(gomock.Any(), gomock.Any()).Return("", errors.New("db error"))

	body := `{"username": "user", "password": "pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var r response.Resp
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "error", r.Status)
	assert.Equal(t, "something went wrong, try again later", r.Data)
}

func TestSignIn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := setupSignInApp(mockUser, mockValidator)

	mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
	mockUser.EXPECT().Login(gomock.Any(), gomock.Any()).Return("mocktoken123", nil)

	body := `{"username": "user", "password": "pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var r map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&r)
	assert.Equal(t, "OK", r["status"])
	assert.Equal(t, "mocktoken123", r["token"])
}

func setupSignInApp(u rest.User, v rest.Validator) *fiber.App {
	app := fiber.New()
	h := rest.NewUserHandler(slogdiscard.NewDiscardLogger(), u, v)
	app.Post("/api/auth/sign-in", h.SignIn)
	return app
}
