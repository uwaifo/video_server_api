package persistence

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/uwaifo/video_server_api/domian/entity"
	"github.com/uwaifo/video_server_api/domian/repository"
)

type Repositories struct {
	User      repository.UserRepository
	Food      repository.FoodRepository
	PhotoWork repository.PhotoWorkRepository

	db *gorm.DB
}

func NewRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User:      NewUserRepository(db),
		Food:      NewFoodRepository(db),
		PhotoWork: NewPhotoWorkRepository(db),
		db:        db,
	}, nil
}

//closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

//This migrate all tables
func (s *Repositories) AutoMigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.PhotoWork{}, &entity.Food{}).Error
}
