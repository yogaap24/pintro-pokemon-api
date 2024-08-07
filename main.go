package main

import (
	"context"
	"errors"
	"flag"
	"mda/pokemon"
	"mda/users"
	"mda/userspokemon"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "c", "config.yml", "Config file name")

	flag.Parse()

	cfg := defaultConfig()
	cfg.loadFromEnv()

	if len(configFileName) > 0 {
		err := loadConfigFromFile(configFileName, &cfg)
		if err != nil {
			log.Warn().Str("file", configFileName).Err(err).Msg("cannot load config file, use defaults")
		}
	}

	log.Debug().Any("config", cfg).Msg("config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBConfig.ConnStr())

	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	users.SetPool(pool)
	userspokemon.SetPool(pool)

	adminUsername := "admin"
	adminPassword := "secret"

	adminUser, err := users.CreateAdminUser(ctx, adminUsername, adminPassword)
	if err != nil {
		if errors.Is(err, errors.New("admin user already exists")) {
			log.Debug().Msg("Admin user already exists")
		} else {
			log.Error().Err(err).Msg("Failed to create admin user")
		}
	} else {
		log.Printf("Admin user created: %v", adminUser)
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Mount("/pokemon", pokemon.Router())
	r.Mount("/users", users.Router())
	r.Mount("/users-pokemon", userspokemon.Router())

	log.Info().Msg("Starting up server...")

	if err := http.ListenAndServe(cfg.Listen.Addr(), r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return
	}

	log.Info().Msg("Server Stopped")
}
