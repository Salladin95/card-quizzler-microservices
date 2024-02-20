package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	dbService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func mockTerm() *entities.CreateTermDto {
	return &entities.CreateTermDto{
		ID:          uuid.New().String(),
		Title:       fmt.Sprintf("sth-%v", time.Now()),
		Description: fmt.Sprintf("sth-%v", time.Now()),
	}
}

func mockTerms() []*entities.CreateTermDto {
	i := 1
	var terms []*entities.CreateTermDto
	for i < 10 {
		terms = append(terms, mockTerm())
		i++
	}
	return terms
}

func main() {
	server := echo.New()
	cfg, err := config.NewConfig()

	if err != nil {
		server.Logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.DbUrl)
	defer conn.Close(ctx)
	if err != nil {
		server.Logger.Fatal(err)
	}

	routes := server.Group("/v1/api")

	routes.GET("/folders/:userID", func(c echo.Context) error {
		userID := c.Param("userID")

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)
		folders, err := queries.GetFolderWithModulesForUser(secondCTX, userID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested folders", Data: folders},
		)
	})

	routes.GET("/folder/:folderID", func(c echo.Context) error {
		folderID := c.Param("folderID")
		parsedID, err := uuid.Parse(folderID)
		if err != nil {
			fmt.Println(err)
		}
		folderID = parsedID.String()

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)
		folder, err := queries.GetFolderByID(secondCTX, folderID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested folder", Data: folder},
		)
	})

	routes.POST("/folder", func(c echo.Context) error {
		var createFolderDTO entities.CreateFolderDto
		c.Bind(&createFolderDTO)

		createModuleParams, err := createFolderDTO.ToCreateFolderParams()

		if err != nil {
			fmt.Println(err)
		}

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)
		err = queries.CreateFolder(secondCTX, createModuleParams)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		createdFolder, err := queries.GetFolderByID(secondCTX, createModuleParams.ID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Created folder", Data: createdFolder},
		)
	})

	routes.POST("/module/:folderID", func(c echo.Context) error {
		folderID := c.Param("folderID")
		parsedID, err := uuid.Parse(folderID)
		if err != nil {
			fmt.Println(err)
		}
		folderID = parsedID.String()

		var createModuleDto entities.CreateModuleDto
		c.Bind(&createModuleDto)

		createModuleParams, err := createModuleDto.ToCreateModuleParams()

		if err != nil {
			fmt.Println(err)
		}

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)
		err = queries.CreateModule(secondCTX, createModuleParams)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = queries.AddModuleToFolder(ctx, dbService.AddModuleToFolderParams{
			FolderID: folderID,
			ModuleID: createModuleParams.ID,
		})

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		terms := mockTerms()
		for _, v := range terms {
			createTermParams, err := v.ToCreateTermDto(createModuleParams.ID)
			if err != nil {
				fmt.Println(err)
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = queries.CreateTerm(secondCTX, createTermParams)
			if err != nil {
				fmt.Println(err)
				return c.String(http.StatusBadRequest, err.Error())
			}
		}

		createdModule, err := queries.GetModuleByID(secondCTX, createModuleParams.ID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested module", Data: createdModule},
		)
	})

	routes.POST("/add-module-to-folder", func(c echo.Context) error {
		// Retrieve the folderID and moduleID query parameters
		folderID := c.QueryParam("folderID")
		moduleID := c.QueryParam("moduleID")

		secondCTX := c.Request().Context()

		// Check if both folderID and moduleID are provided
		if folderID == "" || moduleID == "" {
			// If any parameter is missing, return an error response
			return c.JSON(
				http.StatusBadRequest,
				entities.JsonResponse{Message: "Both folderID and moduleID are required", Data: nil},
			)
		}

		// You can use folderID and moduleID in your application logic here
		queries := dbService.New(conn)
		err = queries.AddModuleToFolder(secondCTX, dbService.AddModuleToFolderParams{
			FolderID: folderID,
			ModuleID: moduleID,
		})

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Return a success response
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested module added to folder", Data: nil})
	})

	routes.GET("/module/:id", func(c echo.Context) error {
		id := c.Param("id")

		parsedID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}
		id = parsedID.String()

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)

		mod, err := queries.GetModuleByID(secondCTX, id)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		var terms []models.Term

		err = json.Unmarshal(mod.Terms, &terms)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Construct the response data
		response := map[string]interface{}{
			"moduleID": mod.ModuleID,
			"title":    mod.ModuleTitle,
			"terms":    terms,
		}

		return c.JSON(
			http.StatusOK,
			response,
		)
	})

	routes.GET("/modules/:uid", func(c echo.Context) error {
		id := c.Param("uid")

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)

		modules, err := queries.GetAllModulesByUserID(secondCTX, id)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			modules,
		)
	})

	routes.POST("/module", func(c echo.Context) error {
		var createModuleDto entities.CreateModuleDto
		c.Bind(&createModuleDto)

		createModuleParams, err := createModuleDto.ToCreateModuleParams()

		if err != nil {
			fmt.Println(err)
		}

		secondCTX := c.Request().Context()

		queries := dbService.New(conn)
		err = queries.CreateModule(secondCTX, createModuleParams)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		terms := mockTerms()
		for _, v := range terms {
			createTermParams, err := v.ToCreateTermDto(createModuleParams.ID)
			if err != nil {
				fmt.Println(err)
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = queries.CreateTerm(secondCTX, createTermParams)
			if err != nil {
				fmt.Println(err)
				return c.String(http.StatusBadRequest, err.Error())
			}
		}

		createdModule, err := queries.GetModuleByID(secondCTX, createModuleParams.ID)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(
			http.StatusOK,
			entities.JsonResponse{Message: "Requested module", Data: createdModule},
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
