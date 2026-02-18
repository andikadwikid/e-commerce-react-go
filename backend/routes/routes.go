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

			// Route Baru: Get All (WAJIB diletakkan SEBELUM route /:id)
			admin.GET("/permissions/all", middlewares.Permission("permissions-index"), adminController.GetAllPermissions)

			// Route Baru: Create Permission
			admin.POST("/permissions", middlewares.Permission("permissions-create"), adminController.CreatePermission)

			// :id adalah parameter dinamis
			admin.GET("/permissions/:id", middlewares.Permission("permissions-show"), adminController.GetPermissionDetail)

			// Route Baru: Update Permission
			admin.PUT("/permissions/:id", middlewares.Permission("permissions-update"), adminController.UpdatePermission)

			admin.DELETE("/permissions/:id", middlewares.Permission("permissions-delete"), adminController.DeletePermission)
		}
	}

	return router
}
