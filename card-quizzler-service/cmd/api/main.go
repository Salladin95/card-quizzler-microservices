package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	//migrations.Migrate(db)

	repository := repositories.NewRepo(db)

	// creates a user
	routes.POST("/user/:uid", func(c echo.Context) error {
		uid := c.Param("uid")
		if uid == "" {
			fmt.Printf("*** err - %v ***\n", err)
			return c.String(http.StatusBadRequest, "User id is required")
		}

		err := repository.CreateUser(uid)

		if err != nil {
			fmt.Printf("*** err - %v ***\n", err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "User is created"},
		)
	})

	// adds module to user
	routes.POST("/add-module-to-user", func(c echo.Context) error {
		// Retrieve the termID and moduleID query parameters
		uid := c.QueryParam("userID")
		mID := c.QueryParam("moduleID")

		if uid == "" {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "User ID is required")
		}

		moduleID, err := uuid.Parse(mID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.AddModuleToUser(uid, moduleID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "The module is added to user", Data: nil})
	})

	// adds term to the module
	routes.POST("/folder-to-user", func(c echo.Context) error {
		// Retrieve the termID and moduleID query parameters
		tID := c.QueryParam("termID")
		mID := c.QueryParam("moduleID")

		termID, err := uuid.Parse(tID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		moduleID, err := uuid.Parse(mID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.AddTermToModule(termID, moduleID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "The term is added to module", Data: nil})
	})

	// creates a module
	routes.POST("/module", func(c echo.Context) error {
		fmt.Println("processing create module request")
		var createModuleDto entities.CreateModuleDto
		err := lib.BindBodyAndVerify(c, &createModuleDto)

		module, err := repository.CreateModule(createModuleDto)

		if err != nil {
			fmt.Printf("*** err - %v ***\n", err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Created module", Data: module},
		)
	})

	// updates module
	routes.PATCH("/module/:id", func(c echo.Context) error {
		id := c.Param("id")

		moduleID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		var updateModuleDto entities.UpdateModuleDto
		err = lib.BindBodyAndVerify(c, &updateModuleDto)

		module, err := repository.UpdateModule(moduleID, updateModuleDto)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			module,
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

	// deletes module by moduleID
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

	// adds term to the module
	routes.POST("/add-term-to-module", func(c echo.Context) error {
		// Retrieve the termID and moduleID query parameters
		tID := c.QueryParam("termID")
		mID := c.QueryParam("moduleID")

		termID, err := uuid.Parse(tID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		moduleID, err := uuid.Parse(mID)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.AddTermToModule(termID, moduleID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "The term is added to module", Data: nil})
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

	// add module to the folder
	routes.DELETE("/delete-module-from-folder", func(c echo.Context) error {
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

		err = repository.DeleteModuleFromFolder(folderID, moduleID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested module deleted from folder", Data: nil})
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
