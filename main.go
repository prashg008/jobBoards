package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prashg008/jobBoards/database"
	docs "github.com/prashg008/jobBoards/docs"
	"github.com/prashg008/jobBoards/handlers"
	"github.com/prashg008/jobBoards/middleware"
	"github.com/prashg008/jobBoards/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Job Boards API
// @version 1.0
// @description API for managing job boards
// @host localhost:8080
// @BasePath /
// @SecurityDefinitions jwt
// @SecurityScheme jwt apiKey Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// Initialize database connection
	db := database.Init()
	// Migrate database
	err = utils.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	// Serve Swagger documentation
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", handlers.RegisterUser)
		userRoutes.POST("/login", handlers.LoginUser)
		userRoutes.GET("/me", middleware.Authenticate(), handlers.GetUserDetails)
		userRoutes.GET("/:id/applications", middleware.Authenticate(), handlers.ListApplicationsByUser)
	}
	jobGroup := router.Group("/jobs")
	{
		// Endpoint to create a job posting (accessible only to companies)
		jobGroup.POST("", middleware.Authenticate(), middleware.Authorize("company"), handlers.CreateJob)

		// Endpoint to list job postings (accessible to any logged-in user)
		jobGroup.GET("", middleware.Authenticate(), handlers.ListJobs)

		// Endpoint to create a job posting (accessible only to companies)
		jobGroup.DELETE("/:id", middleware.Authenticate(), middleware.Authorize("company"), handlers.DeleteJob)
	}
	applicationGroup := router.Group("/applications")
	{
		// Endpoint to apply for a job (accessible to any logged-in user)
		applicationGroup.POST("/", middleware.Authenticate(), handlers.ApplyForJob)
		// Endpoint to update application status by company (accessible only to companies)
		applicationGroup.PUT("/:id", middleware.Authenticate(), middleware.Authorize("company"), handlers.UpdateApplicationStatus)
	}

	router.Run(":8080")

}
