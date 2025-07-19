package main

import (
	"context"
	"dbpostgres/db"
	"dbpostgres/db/migrate"
	"dbpostgres/pkg/env"
	"dbpostgres/pkg/log"
	"time"
)

func main() {
	// 1. Inicializar el logger
	logger := log.New()
	logger.Info().Msg("Iniciando el proceso de migración")

	// 2. Cargar configuración desde variables de entorno
	config, err := env.New(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error al cargar la configuración")
	}

	// 3. Inicializar y conectar a la base de datos
	database := db.New(logger, config)

	// Usamos un contexto con timeout para la conexión inicial
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := database.Connect(ctx); err != nil {
		logger.Fatal().Err(err).Msg("No se pudo conectar a la base de datos")
	}

	// Defer para cerrar la conexión al final de la ejecución de main
	defer func() {
		logger.Info().Msg("Cerrando la conexión de la base de datos.")
		if err := database.Close(); err != nil {
			logger.Error().Err(err).Msg("Error al cerrar la conexión de la base de datos")
		}
	}()

	// 4. Ejecutar migraciones e inicializar datos del sistema
	logger.Info().Msg("Ejecutando migraciones e inicializando datos del sistema...")
	if err := migrate.MigrateDB(database.Conn(context.Background()), logger, config); err != nil {
		logger.Fatal().Err(err).Msg("Falló la migración e inicialización de datos del sistema")
	}

	logger.Info().Msg("Migración e inicialización completadas exitosamente.")
}
