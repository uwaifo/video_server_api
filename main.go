package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uwaifo/video_server_api/infrastructure/auth"
	"github.com/uwaifo/video_server_api/infrastructure/persistence"
	"github.com/uwaifo/video_server_api/interfaces"
	"github.com/uwaifo/video_server_api/interfaces/fileupload"
	"github.com/uwaifo/video_server_api/interfaces/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no enviroment variable found")
	}

}

func main() {

	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Latter add some redis stuff
	//redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbDriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	// Run migration if none
	services.AutoMigrate()

	//redis connection
	redisService, err := auth.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fd := fileupload.NewFileUpload()

	users := interfaces.NewUsers(services.User, redisService.Auth, tk)
	photoWorks := interfaces.NewPhotoWork(services.PhotoWork, services.User, redisService.Auth, tk)
	foods := interfaces.NewFood(services.Food, services.User, fd, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)

	r := gin.Default()
	//r.Use(middleware.CORS(middleware())
	// Start the application

	//setup interfaces.
	//users := in

	//User Routes
	r.POST("/users", users.SaveUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:user_id", users.GetUser)

	//post routes
	r.POST("/food", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), foods.SaveFood)
	r.POST("/addfood", foods.SaveFood)
	r.PUT("/food/:food_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), foods.UpdateFood)
	r.GET("/food/:food_id", foods.GetFoodAndCreator)
	r.DELETE("/food/:food_id", middleware.AuthMiddleware(), foods.DeleteFood)
	r.GET("/food", foods.GetAllFood)

	//photo works routes
	r.POST("/work", middleware.AuthMiddleware(), photoWorks.SavePhotoWork)
	r.GET("/work/:photo_work_id", middleware.AuthMiddleware(), photoWorks.GetUserPhotoWork)
	r.GET("/work", photoWorks.GetAllPhotoWork)
	r.PUT("/work/:photo_work_id", middleware.AuthMiddleware(), photoWorks.UpdatePhotoWork)

	//authentication routes
	r.POST("/login", authenticate.Login)
	r.POST("/logout", authenticate.Logout)
	r.POST("/refresh", authenticate.Refresh)

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8888"

	}
	log.Fatal(r.Run(":" + appPort))

}
