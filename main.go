package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uwaifo/video_server_api/infrastructure/auth"
	"github.com/uwaifo/video_server_api/infrastructure/persistence"
	"github.com/uwaifo/video_server_api/interfaces"
	"github.com/uwaifo/video_server_api/interfaces/fileupload"
	"github.com/uwaifo/video_server_api/interfaces/middleware"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no enviroment variable found")
	}

}

func main() {

	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Latter add some redis stuff
	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbdriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	// Run migration if none
	services.AutoMigrate()

	//redis connection
	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fd := fileupload.NewFileUpload()

	users := interfaces.NewUsers(services.User, redisService.Auth, tk)
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
	r.PUT("/food/:food_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), foods.UpdateFood)
	r.GET("/food/:food_id", foods.GetFoodAndCreator)
	r.DELETE("/food/:food_id", middleware.AuthMiddleware(), foods.DeleteFood)
	r.GET("/food", foods.GetAllFood)

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
