package main

import (
	"danger-dodgers/internal"
	"danger-dodgers/pkg/auth"
	"danger-dodgers/pkg/db"
	"danger-dodgers/pkg/errors"
	"danger-dodgers/pkg/passwords"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/o1egl/paseto"
	"googlemaps.github.io/maps"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	KEY := []byte(GetHasherKey())

	gMaps, err := maps.NewClient(maps.WithAPIKey(GetMapsAPIKey()))
	if err != nil {
		return err
	}

	database, err := badger.Open(badger.DefaultOptions(GetDatabasePath()))
	if err != nil {
		return err
	}
	defer database.Close()

	passwordHasher := passwords.New(func(bh *passwords.BCryptHasher) {
		bh.Cost = 5
	})

	paseto := paseto.NewV2()
	refresh, err := auth.NewPasetoAuthenticator(KEY, paseto, func(pa *auth.PasetoAuthenticator) {
		pa.Type = auth.Refresh
	})
	if err != nil {
		return err
	}

	authentication, err := auth.NewPasetoAuthenticator(KEY, paseto, func(pa *auth.PasetoAuthenticator) {
		pa.Type = auth.Authentication
	})
	if err != nil {
		return err
	}

	// Create databases.
	userDB := db.New[internal.User](database, db.Kind(internal.USER), internal.UserMapping)
	reportDB := db.New[internal.Report](database, db.Kind(internal.REPORT), internal.ReportMapping)
	activityDB := db.New[internal.Activity](database, db.Kind(internal.ACTIVITY), internal.ActivityMapping)
	positionDB := db.New[internal.ActivityPosition](database, db.Kind(internal.ACTIVITY_POSITION), internal.ActivityPositionMapping)

	// Create services.
	userService := internal.NewUserService(userDB, passwordHasher, refresh)
	tokenService := internal.NewTokenService(authentication)
	reportService := internal.NewReportService(reportDB)
	activityService := internal.NewActivityService(activityDB)
	positionService := internal.NewPositionService(positionDB, activityDB)
	directionsService := internal.NewDirectionsService(gMaps)

	// Create handlers.
	userHandler := internal.NewUserHandler(userService)
	tokenHandler := internal.NewTokenHandler(tokenService, authentication)
	reportHandler := internal.NewReportHandler(reportService)
	activityHandler := internal.NewActivityHandler(activityService)
	directionsHandler := internal.NewDirectionsHandler(directionsService)
	positionHandler := internal.NewPositionHandler(positionService)

	// Create application server.
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.New().HandleError, DisableStartupMessage: true,
	})

	// Create middleware.
	authMiddleware := internal.NewHTTPAuthenticator(authentication)

	// User group.
	userGroup := app.Group("/users")
	userGroup.Post("/", userHandler.Create)
	userGroup.Post("/login", userHandler.Authenticate)
	userGroup.Get("/:id", authMiddleware.Authenticate(userHandler.Get))
	userGroup.Delete("/:id", authMiddleware.Authenticate(userHandler.Delete))

	// Token group.
	tokenGroup := app.Group("/tokens")
	tokenGroup.Post("/", tokenHandler.Create)

	// Report group.
	reportGroup := app.Group("/reports")
	reportGroup.Post("/", authMiddleware.Provide(reportHandler.Create))
	reportGroup.Get("/", authMiddleware.Provide(reportHandler.List))
	reportGroup.Get("/:id", authMiddleware.Provide(reportHandler.Get))
	reportGroup.Delete("/:id", authMiddleware.Provide(reportHandler.Delete))

	// Activity group.
	activityGroup := app.Group("/activities")
	activityGroup.Post("/", authMiddleware.Provide(activityHandler.Create))
	activityGroup.Get("/", authMiddleware.Provide(activityHandler.List))
	activityGroup.Delete("/:id", authMiddleware.Provide(activityHandler.Delete))

	// Directions group.
	directionsGroup := app.Group("/directions")
	directionsGroup.Post("/", authMiddleware.Provide(directionsHandler.Route))

	positionGroup := app.Group("/positions")
	positionGroup.Post("/", authMiddleware.Provide(positionHandler.Create))
	positionGroup.Get("/:id", authMiddleware.Provide(positionHandler.List))

	return app.Listen(fmt.Sprintf(":%s", GetApplicationPort()))
}
