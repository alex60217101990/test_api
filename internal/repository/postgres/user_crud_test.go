package postgres

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/alex60217101990/test_api/internal/repository/mock"
	"github.com/google/uuid"
)

var (
	repo     repository.Repository
	m        *mock.Repository
	confFile = "../../../deploy/configs/application.yaml"
)

func init() {
	InitRepoTestEnvironment()
}

func InitRepoTestEnvironment() {
	// Load configs file
	err := configs.ReadConfigFile(confFile)
	if err != nil {
		logger.CmdError.Println(err)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	// Init loggers
	logger.InitLoggerSettings()

	switch configs.Conf.DB.RepoType {
	case configs.RepoPostgres:
		repo = NewPostgresRepository()
	default:
		logger.AppLogger.Fatal(fmt.Errorf("config error: invalid repository type %v", configs.Conf.DB.RepoType))
	}
	repo.Connect(context.Background())

	m = &mock.Repository{}
	m.Connect(context.Background())

	configs.Conf.Keys.PubKeyRepo = ".." + string(os.PathSeparator) + configs.Conf.Keys.PubKeyRepo
	configs.Conf.Keys.PrvKeyRepo = ".." + string(os.PathSeparator) + configs.Conf.Keys.PrvKeyRepo
}

func TestGetUserByCreeds(t *testing.T) {
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	ctx := context.Background()
	defer func() {
		repo.Close()
		//cancel()
	}()

	testings := repository.GenerateTestCredentials()
	for _, creed := range testings {
		user, err := repo.GetUserByCreeds(ctx, creed, struct{}{})
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
}

func TestInsertUser(t *testing.T) {
	fmt.Println(configs.Conf.Keys.PubKeyRepo)
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	testings, err := m.GenerateUser(ctx)
	if err != nil {
		t.Error(err)
	}

	err = repo.InsertUser(ctx, testings)
	if err != nil {
		t.Error(err)
	}

	// for _, creed := range testings {
	// 	user, err := repo.GetUserByCreeds(ctx, creed)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	t.Log(user)
	// }
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	uuid, err := uuid.Parse("645d7be2-9d70-43eb-9d81-37ff121a118b")
	if err != nil {
		t.Error(err)
	}

	user := &models.User{
		Base: models.Base{
			PublicID: uuid,
		},
		Username: "Alex",
		Email:    "some@gmail.com",
		Password: hex.EncodeToString([]byte("some_password")),
	}

	err = repo.UpdateUser(ctx, user)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSoftUser(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	err := repo.DeleteSoft(ctx, "feee81d7-4bfd-464c-9547-62ca3829ca83")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteHardUser(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	err := repo.DeleteHard(ctx, "feee81d7-4bfd-464c-9547-62ca3829ca83")
	if err != nil {
		t.Error(err)
	}
}
