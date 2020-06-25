package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uwaifo/video_server_api/infrastructure/persistence"
	"log"
	"os"
)

func init() {

}

func main() {
	
	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	
	// Latter add some redis stuff
	
	services, err := persistence.NewRepositories(dbdriver, user, password, port,, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	// Run migration if none
	services.AutoMigrate()

	//redis connection



	r := gin.Default()
	//r.Use(middleware.CORS(middleware())
	// Start the application



	//setup interfaces.
	//users := in

	//User Routes
	r.POST("/users", user.SaveUser)




	app_port := os.Getenv("PORT")
	if app_port == "" {
		//ie loacalhost
		app_port = "8888"

	}
	log.Fatal()



}
