// internal/routes/routes.go
package routes

import (
	"pusat-layanan-kemasan/backend-go/internal/controllers"
	"pusat-layanan-kemasan/backend-go/internal/middleware"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Grup untuk rute Pengguna
		users := api.Group("/users")
		{
			users.POST("/register", controllers.Register)
			users.POST("/login", controllers.Login)
			users.GET("/profile", middleware.Protect(), controllers.GetProfile)
			users.GET("/", middleware.Protect(), middleware.IsAdmin(), controllers.GetUsers)
			reports := api.Group("/reports")
			reports.Use(middleware.Protect())
			{
				reports.GET("/sales-summary", middleware.HasAnyRole(models.RoleManajer, models.RoleAdmin), controllers.GetSalesSummary)
				reports.GET("/sales-summary/export", middleware.HasAnyRole(models.RoleManajer, models.RoleAdmin), controllers.ExportSalesSummaryToExcel)
				reports.GET("/sales-chart", middleware.HasAnyRole(models.RoleManajer, models.RoleAdmin), controllers.GetSalesChartData)
			}
		}

		// --- TAMBAHKAN GRUP BARU UNTUK RUTE PESANAN ---
		orders := api.Group("/orders")
		orders.Use(middleware.Protect())
		{
			orders.POST("", middleware.HasAnyRole(models.RoleKasir, models.RoleAdmin), controllers.CreateOrder)

			orders.GET("/:id", controllers.GetOrderByID) // Dapat diakses semua peran yang login

			// Rute untuk Designer
			orders.GET("/queue/designer", middleware.HasAnyRole(models.RoleDesigner, models.RoleAdmin), controllers.GetDesignerQueue)

			// --- TAMBAHKAN RUTE BARU UNTUK OPERATOR ---
			orders.GET("/queue/operator", middleware.HasAnyRole(models.RoleOperator, models.RoleAdmin), controllers.GetOperatorQueue)

			orders.GET("/queue/kasir", middleware.HasAnyRole(models.RoleKasir, models.RoleAdmin), controllers.GetKasirQueue)

			// Rute ini sekarang bisa diakses oleh berbagai peran yang diizinkan
			orders.PATCH("/:id/status", middleware.HasAnyRole(models.RoleDesigner, models.RoleOperator, models.RoleKasir, models.RoleAdmin), controllers.UpdateOrderStatus)
			customers := api.Group("/customers")
			// Hanya Kasir dan Admin yang bisa mengelola pelanggan
			customers.Use(middleware.Protect(), middleware.HasAnyRole(models.RoleKasir, models.RoleAdmin))
			{
				customers.POST("", controllers.CreateCustomer)
				customers.GET("", controllers.GetCustomers)
			}
			orders.GET("/monitoring", middleware.HasAnyRole(models.RoleManajer, models.RoleAdmin), controllers.GetMonitoringOrders)
		}
		notifications := api.Group("/notifications")
		notifications.Use(middleware.Protect()) // Semua rute notifikasi butuh login
		{
			notifications.GET("", controllers.GetNotifications)
			notifications.PATCH("/:id/read", controllers.MarkAsRead)
		}
	}
}
