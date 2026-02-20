package routes

import (
	"github.com/gin-gonic/gin"

	controllers "backend-commerce/controllers"
	adminController "backend-commerce/controllers/admin"
	middlewares "backend-commerce/middlewares"
)

func SetupRouter() *gin.Engine {

	// Initialize gin
	router := gin.Default()

	// auth routes (no auth required)
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		// Admin Routes
		// Kita buat group baru "/admin" di dalam "/api" -> URL menjadi /api/admin/...
		admin := api.Group("/admin")
		admin.Use(middlewares.AuthMiddleware())
		{
			// Route Permissions
			admin.GET("/permissions", middlewares.Permission("permissions-index"), adminController.FindPermissions)
			admin.POST("/permissions", middlewares.Permission("permissions-create"), adminController.CreatePermission)
			admin.GET("/permissions/all", middlewares.Permission("permissions-index"), adminController.GetAllPermissions)
			admin.GET("/permissions/:id", middlewares.Permission("permissions-show"), adminController.GetPermissionDetail)
			admin.PUT("/permissions/:id", middlewares.Permission("permissions-update"), adminController.UpdatePermission)
			admin.DELETE("/permissions/:id", middlewares.Permission("permissions-delete"), adminController.DeletePermission)

			// Route Roles
			admin.GET("/roles", middlewares.Permission("roles-index"), adminController.FindRoles)
			admin.POST("/roles", middlewares.Permission("roles-create"), adminController.CreateRole)
			admin.GET("/roles/all", middlewares.Permission("roles-index"), adminController.GetAllRoles)
			admin.GET("/roles/:id", middlewares.Permission("roles-show"), adminController.GetRoleDetail)
			admin.PUT("/roles/:id", middlewares.Permission("roles-update"), adminController.UpdateRole)
			admin.DELETE("/roles/:id", middlewares.Permission("roles-delete"), adminController.DeleteRole)

			// Route User
			admin.GET("/users", middlewares.Permission("users-index"), adminController.FindUsers)
			admin.POST("/users", middlewares.Permission("users-create"), adminController.CreateUser)
			admin.PUT("/users/:id", middlewares.Permission("users-update"), adminController.UpdateUser)
			admin.GET("/users/:id", middlewares.Permission("users-show"), adminController.GetUserDetail)
			admin.DELETE("/users/:id", middlewares.Permission("users-delete"), adminController.DeleteUser)
		}
	}

	return router
}
