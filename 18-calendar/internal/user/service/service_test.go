package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/mocks"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	user := dto.RegisterUser{
		Username: "testuser",
		Password: "qwerty123",
	}

	// Проверяем, что CreateUser вызовется с захешированным паролем
	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, u domain.User) error {
			assert.Equal(t, "testuser", u.Username)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("qwerty123")))
			return nil
		})

	err := svc.RegisterUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	user := dto.RegisterUser{
		Username: "exists",
		Password: "qwerty123",
	}

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(repo.ErrUserExists)

	err := svc.RegisterUser(context.Background(), user)
	assert.ErrorIs(t, err, service.ErrUserExists)
}

func TestRegisterUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	user := dto.RegisterUser{
		Username: "testuser",
		Password: "qwerty123",
	}

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	err := svc.RegisterUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	password := "qwerty123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.EXPECT().
		User(gomock.Any(), "testuser").
		Return(domain.User{ID: 1, Username: "testuser", PasswordHash: string(hash)}, nil)

	mockManager.EXPECT().
		NewToken(1, time.Hour).
		Return("mocktoken123", nil)

	token, err := svc.Login(context.Background(), dto.LoginUser{
		Username: "testuser",
		Password: password,
	})

	assert.NoError(t, err)
	assert.Equal(t, "mocktoken123", token)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	mockRepo.EXPECT().
		User(gomock.Any(), "notfound").
		Return(domain.User{}, errors.New("user not found"))

	token, err := svc.Login(context.Background(), dto.LoginUser{
		Username: "notfound",
		Password: "whatever",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Empty(t, token)
}

func TestLogin_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	mockRepo.EXPECT().
		User(gomock.Any(), "testuser").
		Return(domain.User{ID: 1, Username: "testuser", PasswordHash: string(hash)}, nil)

	token, err := svc.Login(context.Background(), dto.LoginUser{
		Username: "testuser",
		Password: "wrong-password",
	})

	assert.ErrorIs(t, errors.Unwrap(err), service.ErrInvalidCredentials)
	assert.Empty(t, token)
}

func TestLogin_TokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockManager := mocks.NewMockTokenManager(ctrl)

	svc := service.NewUser(mockRepo, mockManager, time.Hour)

	password := "qwerty123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.EXPECT().
		User(gomock.Any(), "testuser").
		Return(domain.User{ID: 1, Username: "testuser", PasswordHash: string(hash)}, nil)

	mockManager.EXPECT().
		NewToken(1, time.Hour).
		Return("", errors.New("token generation failed"))

	token, err := svc.Login(context.Background(), dto.LoginUser{
		Username: "testuser",
		Password: password,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token generation failed")
	assert.Empty(t, token)
}
