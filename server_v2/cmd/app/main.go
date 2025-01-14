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
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	KEY := []byte(GetHasherKey())

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
	userDB := db.New[internal.User](database, internal.USER, internal.UserMapping)

	// Create services.
	userService := internal.NewUserService(userDB, passwordHasher, refresh)
	tokenService := internal.NewTokenService(authentication)

	// Create handlers.
	userHandler := internal.NewUserHandler(userService)
	tokenHandler := internal.NewTokenHandler(tokenService, authentication)

	// Create application server.
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.New().HandleError, DisableStartupMessage: true,
	})

	// Create middleware.
	authMiddleware := internal.NewHTTPAuthenticator(authentication)

	// Create routes.
	userGroup := app.Group("/users")
	userGroup.Post("/", userHandler.Create)
	userGroup.Post("/login", userHandler.Authenticate)
	userGroup.Get("/:id", authMiddleware.Authenticate(userHandler.Get))
	userGroup.Get("/:id", authMiddleware.Authenticate(userHandler.Delete))

	// Token groups.
	tokenGroup := app.Group("/tokens")
	tokenGroup.Post("/", tokenHandler.Create)

	return app.Listen(fmt.Sprintf(":%s", GetApplicationPort()))
}