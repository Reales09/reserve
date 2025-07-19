package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/env"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// MigrateDB ejecuta todas las migraciones
func MigrateDB(db *gorm.DB, logger log.ILogger, config env.IConfig) error {
	logger.Info().Msg("Iniciando migraciones...")

	// Migrar modelos
	err := db.AutoMigrate(
		&models.Restaurant{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RestaurantStaff{},
		&models.Client{},
		&models.Table{},
		&models.Reservation{},
		&models.ReservationStatus{},
		&models.ReservationStatusHistory{},
	)
	if err != nil {
		return err
	}

	// Crear repositorio del sistema
	systemRepo := repository.NewSystemRepository(db, logger)

	// Inicializar datos del sistema usando el repositorio
	if err := InitializeSystemData(systemRepo, logger, config); err != nil {
		return err
	}

	logger.Info().Msg("Migraciones completadas exitosamente")
	return nil
}
