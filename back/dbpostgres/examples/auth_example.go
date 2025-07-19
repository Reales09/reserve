package main

import (
	"fmt"
	"log"

	"dbpostgres/db/models"
	"dbpostgres/pkg/auth"
	"dbpostgres/pkg/services"

	"gorm.io/gorm"
)

// Ejemplo de uso del sistema de autenticación y autorización
func main() {
	// Asumiendo que ya tienes una conexión a la base de datos
	var db *gorm.DB
	// db = tu_conexion_a_db

	// Inicializar servicios
	authService := auth.NewAuthService(db)
	userService := services.NewUserService(db)
	restaurantService := services.NewRestaurantService(db)
	roleService := services.NewRoleService(db)

	// Ejemplo 1: Crear un super administrador
	createSuperAdmin(db, userService, roleService)

	// Ejemplo 2: Crear un restaurante con marca blanca
	createRestaurantWithBranding(db, restaurantService)

	// Ejemplo 3: Crear usuarios del restaurante
	createRestaurantUsers(db, userService, roleService)

	// Ejemplo 4: Autenticación y autorización
	authenticationExample(db, authService, userService)
}

// createSuperAdmin crea un super administrador del sistema
func createSuperAdmin(db *gorm.DB, userService *services.UserService, roleService *services.RoleService) {
	fmt.Println("=== Creando Super Administrador ===")

	// Crear usuario super admin
	superAdmin := &models.User{
		Name:     "Super Administrador",
		Email:    "admin@platform.com",
		Password: "hashed_password_here", // TODO: Usar bcrypt
		Phone:    "+1234567890",
		IsActive: true,
	}

	// Asignar rol de super admin (sin restaurantes específicos)
	err := userService.CreateUser(superAdmin, []string{models.ROLE_SUPER_ADMIN}, []uint{})
	if err != nil {
		log.Printf("Error creando super admin: %v", err)
		return
	}

	fmt.Printf("Super administrador creado: %s (%s)\n", superAdmin.Name, superAdmin.Email)
}

// createRestaurantWithBranding crea un restaurante con configuración de marca blanca
func createRestaurantWithBranding(db *gorm.DB, restaurantService *services.RestaurantService) {
	fmt.Println("\n=== Creando Restaurante con Marca Blanca ===")

	restaurant := &models.Restaurant{
		Name:        "Restaurante El Buen Sabor",
		Code:        "el-buen-sabor",
		Timezone:    "America/Bogota",
		Address:     "Calle 123 #45-67, Bogotá",
		Description: "El mejor restaurante de comida colombiana",

		// Configuración de marca blanca
		LogoURL:        "https://cdn.example.com/logos/el-buen-sabor.png",
		PrimaryColor:   "#1f2937",
		SecondaryColor: "#3b82f6",
		CustomDomain:   "elbuensabor.com",
		IsActive:       true,

		// Funcionalidades habilitadas
		EnableDelivery:     true,
		EnablePickup:       true,
		EnableReservations: true,
	}

	err := restaurantService.CreateRestaurant(restaurant)
	if err != nil {
		log.Printf("Error creando restaurante: %v", err)
		return
	}

	fmt.Printf("Restaurante creado: %s (ID: %d)\n", restaurant.Name, restaurant.ID)
	fmt.Printf("URL personalizada: https://%s\n", restaurant.CustomDomain)
	fmt.Printf("URL por defecto: https://platform.com/%s\n", restaurant.Code)
}

// createRestaurantUsers crea usuarios para el restaurante
func createRestaurantUsers(db *gorm.DB, userService *services.UserService, roleService *services.RoleService) {
	fmt.Println("\n=== Creando Usuarios del Restaurante ===")

	// Obtener el restaurante
	restaurant, err := services.NewRestaurantService(db).GetRestaurantByCode("el-buen-sabor")
	if err != nil {
		log.Printf("Error obteniendo restaurante: %v", err)
		return
	}

	// Crear propietario del restaurante
	owner := &models.User{
		Name:     "Juan Pérez",
		Email:    "juan@elbuensabor.com",
		Password: "hashed_password_here",
		Phone:    "+573001234567",
		IsActive: true,
	}

	err = userService.CreateUser(owner, []string{models.ROLE_RESTAURANT_OWNER}, []uint{restaurant.ID})
	if err != nil {
		log.Printf("Error creando propietario: %v", err)
		return
	}

	// Crear gerente del restaurante
	manager := &models.User{
		Name:     "María García",
		Email:    "maria@elbuensabor.com",
		Password: "hashed_password_here",
		Phone:    "+573001234568",
		IsActive: true,
	}

	err = userService.CreateUser(manager, []string{models.ROLE_RESTAURANT_MANAGER}, []uint{restaurant.ID})
	if err != nil {
		log.Printf("Error creando gerente: %v", err)
		return
	}

	// Crear mesero
	waiter := &models.User{
		Name:     "Carlos López",
		Email:    "carlos@elbuensabor.com",
		Password: "hashed_password_here",
		Phone:    "+573001234569",
		IsActive: true,
	}

	err = userService.CreateUser(waiter, []string{models.ROLE_WAITER}, []uint{restaurant.ID})
	if err != nil {
		log.Printf("Error creando mesero: %v", err)
		return
	}

	fmt.Printf("Usuarios del restaurante creados:\n")
	fmt.Printf("- Propietario: %s (%s)\n", owner.Name, owner.Email)
	fmt.Printf("- Gerente: %s (%s)\n", manager.Name, manager.Email)
	fmt.Printf("- Mesero: %s (%s)\n", waiter.Name, waiter.Email)
}

// authenticationExample muestra ejemplos de autenticación y autorización
func authenticationExample(db *gorm.DB, authService *auth.AuthService, userService *services.UserService) {
	fmt.Println("\n=== Ejemplos de Autenticación y Autorización ===")

	// Ejemplo 1: Autenticar super admin
	fmt.Println("\n1. Autenticando Super Administrador:")
	superAdminCtx, err := authService.AuthenticateUser("admin@platform.com", "password")
	if err != nil {
		log.Printf("Error autenticando super admin: %v", err)
		return
	}

	fmt.Printf("Usuario autenticado: %s\n", superAdminCtx.Email)
	fmt.Printf("Roles: %v\n", superAdminCtx.Roles)
	fmt.Printf("Es super usuario: %t\n", authService.IsSuperUser(superAdminCtx))
	fmt.Printf("Puede gestionar restaurantes: %t\n", authService.HasPermission(superAdminCtx, models.RESOURCE_RESTAURANTS, models.ACTION_MANAGE))

	// Ejemplo 2: Autenticar propietario de restaurante
	fmt.Println("\n2. Autenticando Propietario de Restaurante:")
	ownerCtx, err := authService.AuthenticateUser("juan@elbuensabor.com", "password")
	if err != nil {
		log.Printf("Error autenticando propietario: %v", err)
		return
	}

	fmt.Printf("Usuario autenticado: %s\n", ownerCtx.Email)
	fmt.Printf("Restaurantes: %v\n", ownerCtx.RestaurantIDs)
	fmt.Printf("Roles: %v\n", ownerCtx.Roles)
	fmt.Printf("Es super usuario: %t\n", authService.IsSuperUser(ownerCtx))
	fmt.Printf("Puede gestionar reservas: %t\n", authService.HasPermission(ownerCtx, models.RESOURCE_RESERVATIONS, models.ACTION_MANAGE))
	fmt.Printf("Puede configurar restaurante: %t\n", authService.HasPermission(ownerCtx, models.RESOURCE_RESTAURANTS, models.ACTION_UPDATE))

	// Ejemplo 3: Autenticar mesero
	fmt.Println("\n3. Autenticando Mesero:")
	waiterCtx, err := authService.AuthenticateUser("carlos@elbuensabor.com", "password")
	if err != nil {
		log.Printf("Error autenticando mesero: %v", err)
		return
	}

	fmt.Printf("Usuario autenticado: %s\n", waiterCtx.Email)
	fmt.Printf("Restaurantes: %v\n", waiterCtx.RestaurantIDs)
	fmt.Printf("Roles: %v\n", waiterCtx.Roles)
	fmt.Printf("Puede ver reservas: %t\n", authService.HasPermission(waiterCtx, models.RESOURCE_RESERVATIONS, models.ACTION_READ))
	fmt.Printf("Puede gestionar reservas: %t\n", authService.HasPermission(waiterCtx, models.RESOURCE_RESERVATIONS, models.ACTION_MANAGE))

	// Ejemplo 4: Verificar acceso a restaurante
	fmt.Println("\n4. Verificando Acceso a Restaurante:")
	restaurantID := uint(1) // ID del restaurante creado

	fmt.Printf("Super admin puede acceder al restaurante %d: %t\n",
		restaurantID, authService.CanAccessRestaurant(superAdminCtx, restaurantID))

	fmt.Printf("Propietario puede acceder al restaurante %d: %t\n",
		restaurantID, authService.CanAccessRestaurant(ownerCtx, restaurantID))

	fmt.Printf("Mesero puede acceder al restaurante %d: %t\n",
		restaurantID, authService.CanAccessRestaurant(waiterCtx, restaurantID))

	// Ejemplo 5: Middleware de autorización
	fmt.Println("\n5. Ejemplos de Middleware de Autorización:")

	// Middleware que requiere permiso específico
	requireManageReservations := authService.RequirePermission(models.RESOURCE_RESERVATIONS, models.ACTION_MANAGE)

	if err := requireManageReservations(ownerCtx); err != nil {
		fmt.Printf("Error de autorización (propietario): %v\n", err)
	} else {
		fmt.Printf("Propietario autorizado para gestionar reservas\n")
	}

	if err := requireManageReservations(waiterCtx); err != nil {
		fmt.Printf("Error de autorización (mesero): %v\n", err)
	} else {
		fmt.Printf("Mesero autorizado para gestionar reservas\n")
	}

	// Ejemplo 6: Usuario con múltiples restaurantes
	createMultiRestaurantUser(db, userService)
}

// createMultiRestaurantUser crea un usuario que trabaja en múltiples restaurantes
func createMultiRestaurantUser(db *gorm.DB, userService *services.UserService) {
	fmt.Println("\n=== Creando Usuario con Múltiples Restaurantes ===")

	// Crear segundo restaurante
	restaurant2 := &models.Restaurant{
		Name:        "Restaurante La Parrilla",
		Code:        "la-parrilla",
		Timezone:    "America/Bogota",
		Address:     "Calle 456 #78-90, Bogotá",
		Description: "El mejor restaurante de parrilla",
		IsActive:    true,
	}

	err := services.NewRestaurantService(db).CreateRestaurant(restaurant2)
	if err != nil {
		log.Printf("Error creando segundo restaurante: %v", err)
		return
	}

	// Obtener ambos restaurantes
	restaurant1, err := services.NewRestaurantService(db).GetRestaurantByCode("el-buen-sabor")
	if err != nil {
		log.Printf("Error obteniendo primer restaurante: %v", err)
		return
	}

	// Crear usuario que trabaja en ambos restaurantes
	multiUser := &models.User{
		Name:     "Ana Rodríguez",
		Email:    "ana@multirestaurante.com",
		Password: "hashed_password_here",
		Phone:    "+573001234570",
		IsActive: true,
	}

	// Asignar rol de gerente y ambos restaurantes
	err = userService.CreateUser(multiUser, []string{models.ROLE_RESTAURANT_MANAGER}, []uint{restaurant1.ID, restaurant2.ID})
	if err != nil {
		log.Printf("Error creando usuario multi-restaurante: %v", err)
		return
	}

	fmt.Printf("Usuario multi-restaurante creado: %s (%s)\n", multiUser.Name, multiUser.Email)
	fmt.Printf("Trabaja en restaurantes: %d y %d\n", restaurant1.ID, restaurant2.ID)

	// Ejemplo de autenticación del usuario multi-restaurante
	multiUserCtx, err := auth.NewAuthService(db).AuthenticateUser("ana@multirestaurante.com", "password")
	if err != nil {
		log.Printf("Error autenticando usuario multi-restaurante: %v", err)
		return
	}

	fmt.Printf("Usuario autenticado: %s\n", multiUserCtx.Email)
	fmt.Printf("Restaurantes: %v\n", multiUserCtx.RestaurantIDs)
	fmt.Printf("Puede acceder al restaurante 1: %t\n", auth.NewAuthService(db).CanAccessRestaurant(multiUserCtx, restaurant1.ID))
	fmt.Printf("Puede acceder al restaurante 2: %t\n", auth.NewAuthService(db).CanAccessRestaurant(multiUserCtx, restaurant2.ID))
	fmt.Printf("Puede acceder al restaurante 3: %t\n", auth.NewAuthService(db).CanAccessRestaurant(multiUserCtx, 999)) // ID inexistente
}
