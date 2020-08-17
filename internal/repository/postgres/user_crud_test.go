package postgres

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/helpers"
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

// InitRepoTestEnvironment
func InitRepoTestEnvironment() {
	// Load configs file
	err := configs.ReadConfigFile(confFile)
	if err != nil {
		logger.CmdError.Println(err)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	// Init loggers
	logger.InitLoggerSettings()

	configs.Conf.DB.DbName = "pgx_test"

	switch configs.Conf.DB.RepoType {
	case configs.RepoPostgres:
		repo = &Repository{}
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
	fmt.Println(configs.Conf.DB.DbName)
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

	err := repo.DeleteSoftUser(ctx, "feee81d7-4bfd-464c-9547-62ca3829ca83")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteHardUser(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	err := repo.DeleteHardUser(ctx, "feee81d7-4bfd-464c-9547-62ca3829ca83")
	if err != nil {
		t.Error(err)
	}
}

func TestGetCategories(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	cat, err := repo.GetCategories(ctx, &models.Pagination{
		// From: "4db46ca6-513a-4e0c-bf19-fec107013793",
		// To:   10,
	}, &models.SortedBy{
		FieldName: "name",
		Desc:      true,
	}, struct{}{})
	if err != nil {
		t.Error(err)
	}
	t.Log(cat)
}

func TestGetCategoryByNameOrID(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	cat, err := repo.GetCategoryByNameOrID(ctx,
		"1966ee86-df0d-11ea-a0f7-acde48001122",
		struct{}{})
	if err != nil {
		t.Error(err)
	}
	t.Log(cat)
}

func TestInsertCategory(t *testing.T) {
	ctx := helpers.AddToContext(context.Background(), repository.UserSessionKey,
		&models.User{
			Base: models.Base{
				ID: 1,
			},
		})
	defer func() {
		repo.Close()
	}()

	cat, err := m.GenerateCategory(ctx)
	if err != nil {
		t.Error(err)
	}

	id1, err := uuid.Parse("cdf9bcc8-deec-11ea-a297-acde48001122")
	if err != nil {
		t.Error(err)
		return
	}
	id2, err := uuid.Parse("9b36ce7d-f017-4b25-9e67-54bb09530930")
	if err != nil {
		t.Error(err)
		return
	}
	cat.Products = []*models.Product{
		&models.Product{
			Base: models.Base{PublicID: id1},
		},
		&models.Product{
			Base: models.Base{PublicID: id2},
		},
	}

	err = repo.InsertCategory(ctx, cat)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateCategory(t *testing.T) {
	ctx := helpers.AddToContext(context.Background(), repository.UserSessionKey,
		&models.User{
			Base: models.Base{
				ID: 4,
			},
		})
	defer func() {
		repo.Close()
	}()

	uuid, err := uuid.Parse("b08702d5-3e59-4b1a-878f-e2f608a18d33")
	if err != nil {
		t.Error(err)
	}

	cat := &models.Category{
		Base: models.Base{
			PublicID: uuid,
		},
		Name: "NewName",
	}

	err = repo.UpdateCategory(ctx, cat)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSoftCategory(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	err := repo.DeleteSoftCategory(ctx, "4db46ca6-513a-4e0c-bf19-fec107013793")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteHardCategory(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	err := repo.DeleteHardCategory(ctx, "4db46ca6-513a-4e0c-bf19-fec107013793")
	if err != nil {
		t.Error(err)
	}
}

func TestGetProducts(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	cat, err := repo.GetProducts(ctx, &models.Pagination{
		From: "809fe121-2d4d-44b1-8e1f-29e7031f8ad6",
		To:   10,
	}, &models.SortedBy{
		FieldName: "name",
		Desc:      true,
	}, struct{}{})
	if err != nil {
		t.Error(err)
	}
	t.Log(cat)
}

func TestGetProductByNameOrID(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	cat, err := repo.GetProductByNameOrID(ctx,
		"voluptas",
		struct{}{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("--->", cat.Categories[0])
	t.Log(cat)
}

func TestInsertProduct(t *testing.T) {
	ctx := helpers.AddToContext(context.Background(), repository.UserSessionKey,
		&models.User{
			Base: models.Base{
				ID: 1,
			},
		})
	// ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	prod, err := m.GenerateProduct(ctx)
	if err != nil {
		t.Error(err)
	}

	id1, err := uuid.Parse("934b0839-858f-4368-b764-f4958a2bdbbf")
	if err != nil {
		t.Error(err)
		return
	}
	id2, err := uuid.Parse("9b36ce7d-f017-4b25-9e67-54bb09530930")
	if err != nil {
		t.Error(err)
		return
	}
	prod.Categories = []*models.Category{
		&models.Category{
			Base: models.Base{PublicID: id1},
		},
		&models.Category{
			Base: models.Base{PublicID: id2},
		},
	}

	st, _ := json.MarshalIndent(prod, "", "\t")
	fmt.Println(string(st))

	// err = repo.InsertProduct(ctx, prod)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// s := prod.PublicID.String()
	// fmt.Println(s)
	// err = repo.AddRelationCategory(ctx, prod.PublicID.String(), "9b36ce7d-f017-4b25-9e67-54bb09530930")
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestUpdateProdyct(t *testing.T) {
	ctx := helpers.AddToContext(context.Background(), repository.UserSessionKey,
		&models.User{
			Base: models.Base{
				ID: 1,
			},
		})
	defer func() {
		repo.Close()
	}()

	uuid, err := uuid.Parse("3a64ddce-dfd9-11ea-bdb8-acde48001122")
	if err != nil {
		t.Error(err)
	}

	prod := &models.Product{
		Base: models.Base{
			PublicID: uuid,
		},
		Name: "temporibus2",
	}

	err = repo.UpdateProduct(ctx, prod)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSoftProduct(t *testing.T) {
	ctx := context.Background()
	defer func() {
		repo.Close()
	}()

	// err := repo.AddRelationCategory(ctx, "6cc14f6a-deed-11ea-8ebd-acde48001122", "9b36ce7d-f017-4b25-9e67-54bb09530930")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	err := repo.DelRelationCategory(ctx, "6cc14f6a-deed-11ea-8ebd-acde48001122", "9b36ce7d-f017-4b25-9e67-54bb09530930")
	if err != nil {
		t.Error(err)
		return
	}

	// err := repo.DeleteSoftProduct(ctx, "6cc14f6a-deed-11ea-8ebd-acde48001122")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
}
