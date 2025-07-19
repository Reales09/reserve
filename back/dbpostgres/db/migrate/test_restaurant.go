package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/log"
)

// InitializeTestRestaurant crea el restaurante de pruebas por defecto
func InitializeTestRestaurant(systemRepo repository.SystemRepository, logger log.ILogger) error {
	logger.Info().Msg("Inicializando restaurante de pruebas...")

	// Verificar si ya existe el restaurante de pruebas
	exists, err := systemRepo.ExistsByCode("test-restaurant")
	if err != nil {
		return err
	}

	if exists {
		logger.Info().Str("restaurant_code", "test-restaurant").Msg("✅ Restaurante de pruebas ya existe, saltando migración de restaurante de pruebas")
		return nil
	}

	// Crear el restaurante de pruebas
	testRestaurant := models.Restaurant{
		Name:               "Restaurante de Pruebas",
		Code:               "test-restaurant",
		Timezone:           "America/Bogota",
		Address:            "Calle de Pruebas #123, Bogotá",
		Description:        "Restaurante de pruebas para desarrollo y testing",
		LogoURL:            "",
		PrimaryColor:       "#1f2937",
		SecondaryColor:     "#3b82f6",
		CustomDomain:       "",
		IsActive:           true,
		EnableDelivery:     true,
		EnablePickup:       true,
		EnableReservations: true,
	}

	if err := systemRepo.CreateRestaurant(&testRestaurant); err != nil {
		logger.Error().Err(err).Msg("Error al crear restaurante de pruebas")
		return err
	}

	// Crear algunas mesas de prueba
	if err := createTestTables(systemRepo, logger, testRestaurant.ID); err != nil {
		return err
	}

	logger.Info().Str("restaurant_name", testRestaurant.Name).Str("restaurant_code", testRestaurant.Code).Msg("✅ Restaurante de pruebas creado exitosamente")
	return nil
}

// createTestTables crea algunas mesas de prueba para el restaurante
func createTestTables(systemRepo repository.SystemRepository, logger log.ILogger, restaurantID uint) error {
	logger.Info().Msg("Creando mesas de prueba...")

	testTables := []models.Table{
		{RestaurantID: restaurantID, Number: 1, Capacity: 2},
		{RestaurantID: restaurantID, Number: 2, Capacity: 4},
		{RestaurantID: restaurantID, Number: 3, Capacity: 6},
		{RestaurantID: restaurantID, Number: 4, Capacity: 8},
		{RestaurantID: restaurantID, Number: 5, Capacity: 2},
		{RestaurantID: restaurantID, Number: 6, Capacity: 4},
		{RestaurantID: restaurantID, Number: 7, Capacity: 6},
		{RestaurantID: restaurantID, Number: 8, Capacity: 10},
	}

	// Verificar si todas las mesas ya existen
	allExist := true
	for _, table := range testTables {
		exists, err := systemRepo.ExistsByRestaurantAndNumber(table.RestaurantID, table.Number)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		logger.Info().Int("tables_count", len(testTables)).Msg("✅ Mesas de prueba ya existen, saltando migración de mesas")
		return nil
	}

	// Inicializar mesas usando el repositorio
	if err := systemRepo.InitializeTables(testTables); err != nil {
		return err
	}

	logger.Info().Int("tables_count", len(testTables)).Msg("✅ Mesas de prueba creadas correctamente")
	return nil
}
