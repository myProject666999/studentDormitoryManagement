package routers

import (
	"student-dormitory-management/controllers"
	"student-dormitory-management/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user", controllers.GetCurrentUser)
			auth.PUT("/user/profile", controllers.UpdateProfile)
			auth.PUT("/user/password", controllers.ChangePassword)

			auth.POST("/upload/image", controllers.UploadImage)
			auth.POST("/upload/attachment", controllers.UploadAttachment)
			auth.GET("/download/:filename", controllers.DownloadAttachment)

			auth.GET("/notices", controllers.GetNotices)
			auth.GET("/notices/:id", controllers.GetNotice)

			admin := auth.Group("/admin")
			admin.Use(middleware.RoleAuth("admin"))
			{
				admin.POST("/notices", controllers.CreateNotice)
				admin.PUT("/notices/:id", controllers.UpdateNotice)
				admin.DELETE("/notices/:id", controllers.DeleteNotice)

				admin.GET("/admins", controllers.GetAdmins)
				admin.POST("/admins", controllers.CreateAdmin)
				admin.PUT("/admins/:id", controllers.UpdateAdmin)
				admin.DELETE("/admins/:id", controllers.DeleteAdmin)

				admin.GET("/students", controllers.GetStudents)
				admin.POST("/students", controllers.CreateStudent)
				admin.PUT("/students/:id", controllers.UpdateStudent)
				admin.DELETE("/students/:id", controllers.DeleteStudent)

				admin.GET("/workers", controllers.GetWorkers)
				admin.GET("/workers/all", controllers.GetAllWorkers)
				admin.POST("/workers", controllers.CreateWorker)
				admin.PUT("/workers/:id", controllers.UpdateWorker)
				admin.DELETE("/workers/:id", controllers.DeleteWorker)

				admin.GET("/dormitories", controllers.GetDormitories)
				admin.GET("/dormitories/available", controllers.GetAvailableDormitories)
				admin.GET("/dormitories/:id", controllers.GetDormitory)
				admin.POST("/dormitories", controllers.CreateDormitory)
				admin.PUT("/dormitories/:id", controllers.UpdateDormitory)
				admin.DELETE("/dormitories/:id", controllers.DeleteDormitory)

				admin.GET("/assignments", controllers.GetAssignments)
				admin.GET("/assignments/:id", controllers.GetAssignment)
				admin.POST("/assignments", controllers.CreateAssignment)
				admin.PUT("/assignments/:id", controllers.UpdateAssignment)
				admin.DELETE("/assignments/:id", controllers.DeleteAssignment)

				admin.GET("/repair-requests", controllers.GetRepairRequests)
				admin.GET("/repair-requests/:id", controllers.GetRepairRequest)
				admin.PUT("/repair-requests/:id", controllers.UpdateRepairRequest)
				admin.DELETE("/repair-requests/:id", controllers.DeleteRepairRequest)

				admin.GET("/repair-records", controllers.GetRepairRecords)
				admin.GET("/repair-records/:id", controllers.GetRepairRecord)
				admin.POST("/repair-records", controllers.CreateRepairRecord)
				admin.PUT("/repair-records/:id", controllers.UpdateRepairRecord)
				admin.DELETE("/repair-records/:id", controllers.DeleteRepairRecord)
			}

			student := auth.Group("/student")
			student.Use(middleware.RoleAuth("student"))
			{
				student.GET("/my-assignment", controllers.GetMyAssignment)
				student.GET("/my-repair-requests", controllers.GetMyRepairRequests)
				student.POST("/repair-requests", controllers.CreateRepairRequest)
				student.PUT("/repair-requests/:id", controllers.UpdateRepairRequest)
				student.GET("/my-repair-records", controllers.GetMyRepairRecords)
				student.GET("/workers/all", controllers.GetAllWorkers)
			}

			worker := auth.Group("/worker")
			worker.Use(middleware.RoleAuth("worker"))
			{
				worker.GET("/my-repair-requests", controllers.GetMyRepairRequests)
				worker.GET("/my-repair-records", controllers.GetMyRepairRecords)
				worker.POST("/repair-records", controllers.CreateRepairRecord)
				worker.PUT("/repair-records/:id", controllers.UpdateRepairRecord)
			}
		}
	}
}
