package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	migrations "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/db"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	server := echo.New()
	cfg, err := config.NewConfig()

	if err != nil {
		server.Logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	routes := server.Group("/v1/api")

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	migrations.Migrate(db)

	repository := repo{db}

	// creates a module
	routes.POST("/module", func(c echo.Context) error {
		var createModuleDto entities.CreateModuleDto
		err := lib.BindBodyAndVerify(c, &createModuleDto)

		module, err := repository.CreateModule(createModuleDto)

		if err != nil {
			fmt.Printf("*** err - %v ***\n", err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested module", Data: module},
		)
	})

	// gets module by moduleID
	routes.GET("/module/:id", func(c echo.Context) error {
		id := c.Param("id")

		moduleID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		module, err := repository.GetModuleByID(moduleID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			module,
		)
	})

	// get modules by userID
	routes.GET("/modules/:uid", func(c echo.Context) error {
		uid := c.Param("uid")

		modules, err := repository.GetModulesByUID(uid)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			modules,
		)
	})

	// gets module by moduleID
	routes.DELETE("/module/:id", func(c echo.Context) error {
		id := c.Param("id")

		moduleID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.DeleteModule(moduleID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusNoContent,
			"Module is deleted",
		)
	})

	// creates folder
	routes.POST("/folder", func(c echo.Context) error {
		var createFolderDTO entities.CreateFolderDto
		err := lib.BindBodyAndVerify(c, &createFolderDTO)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		folder, err := repository.CreateFolder(createFolderDTO)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Created folder", Data: folder},
		)
	})

	// add module to the folder
	routes.POST("/add-module-to-folder", func(c echo.Context) error {
		// Retrieve the folderID and moduleID query parameters
		fID := c.QueryParam("folderID")
		mID := c.QueryParam("moduleID")

		folderID, err := uuid.Parse(fID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		moduleID, err := uuid.Parse(mID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.AddModuleToFolder(folderID, moduleID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested module added to folder", Data: nil})
	})

	// get folders by userID
	routes.GET("/folders/:userID", func(c echo.Context) error {
		userID := c.Param("userID")

		folders, err := repository.GetFoldersByUID(userID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested folders", Data: folders},
		)
	})

	// gets folder by folderID
	routes.GET("/folder/:folderID", func(c echo.Context) error {
		id := c.Param("folderID")
		folderID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}

		folder, err := repository.GetFolderByID(folderID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested folder", Data: folder},
		)
	})

	// creates module and adds fo folder with passed folderID
	routes.POST("/module/:folderID", func(c echo.Context) error {
		id := c.Param("folderID")
		folderID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}

		var createModuleDto entities.CreateModuleDto
		err = lib.BindBodyAndVerify(c, &createModuleDto)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		module, err := repository.CreateModule(createModuleDto)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.AddModuleToFolder(folderID, module.ID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested module", Data: module},
		)
	})

	// gets folder by folderID
	routes.DELETE("/folder/:folderID", func(c echo.Context) error {
		id := c.Param("folderID")
		folderID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}

		err = repository.DeleteFolder(folderID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusNoContent,
			entities.JsonResponse{Message: "Folder is deleted", Data: nil},
		)
	})

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", cfg.GrpcPort)
		if err := server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Panic("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shut down the Echo server.
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) CreateFolder(dto entities.CreateFolderDto) (models.Folder, error) {
	folder, err := dto.ToModel()

	if err != nil {
		return folder, err
	}

	createdModule := r.db.Create(&folder)
	if createdModule.Error != nil {
		return folder, createdModule.Error
	}

	return folder, nil
}

func (r *repo) GetFoldersByUID(uid string) ([]models.Folder, error) {
	var folders []models.Folder
	res := r.db.Preload("Modules.Terms").Find(&folders).Where("userID = ?", uid)
	return folders, res.Error
}

func (r *repo) GetFolderByID(id uuid.UUID) (models.Folder, error) {
	var folder models.Folder
	res := r.db.Preload("Modules.Terms").First(&folder).Where("id", id)
	return folder, res.Error
}

func (r *repo) CreateModule(dto entities.CreateModuleDto) (models.Module, error) {
	module, terms, err := dto.ToModels()

	if err != nil {
		return module, err
	}

	createdModule := r.db.Create(&module)
	if createdModule.Error != nil {
		return module, createdModule.Error
	}

	createdTerms := r.db.Create(&terms)
	if createdTerms.Error != nil {
		return module, createdTerms.Error
	}
	return module, nil
}

func (r *repo) GetModulesByUID(uid string) ([]models.Module, error) {
	var userModules []models.Module
	res := r.db.Preload("Terms").Preload("Folders").Find(&userModules).Where("userID = ?", uid)
	return userModules, res.Error
}

func (r *repo) GetModuleByID(id uuid.UUID) (models.Module, error) {
	var module models.Module
	res := r.db.Preload("Terms").Preload("Folders").First(&module).Where("id", id)
	return module, res.Error
}

func (r *repo) AddModuleToFolder(folderID uuid.UUID, moduleID uuid.UUID) error {
	var folder models.Folder
	if err := r.db.First(&folder, folderID).Error; err != nil {
		return err
	}

	// Create the association between the module and the folder
	res := r.db.Model(&folder).Association("Modules").Append(&models.Module{ID: moduleID})

	return res
}

func (r *repo) DeleteModule(id uuid.UUID) error {
	module, err := r.GetModuleByID(id)
	if err != nil {
		return err
	}
	// Delete all of a user's has one, has many, and many2many associations & all terms
	err = r.db.Select("Terms", clause.Associations).Delete(&module).Error
	return err
}

func (r *repo) DeleteFolder(id uuid.UUID) error {
	folder, err := r.GetFolderByID(id)
	if err != nil {
		return err
	}
	// Delete all of a user's has one, has many, and many2many associations
	err = r.db.Select(clause.Associations).Delete(&folder).Error
	return err
}
