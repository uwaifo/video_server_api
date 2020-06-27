package persistence

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/uwaifo/video_server_api/domian/entity"
	"github.com/uwaifo/video_server_api/domian/repository"
)

// PhotoWorkRepo  . ..
type PhotoWorkRepo struct {
	db *gorm.DB
}

// NewPhotoWorkRepository ...
func NewPhotoWorkRepository(db *gorm.DB) *PhotoWorkRepo {
	return &PhotoWorkRepo{db}
}

//PhotoWorkRepo implements the repository.PhotoWorkRepo interface
var _ repository.PhotoWorkRepository = &PhotoWorkRepo{}

// SavePhotoWork . . .
func (r *PhotoWorkRepo) SavePhotoWork(photoWork *entity.PhotoWork) (*entity.PhotoWork, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	//photoWork.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage

	err := r.db.Debug().Create(&photoWork).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "work title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return photoWork, nil
}

// GetPhotoWork . . .
func (r *PhotoWorkRepo) GetPhotoWork(id uint64) (*entity.PhotoWork, error) {
	var photoWork entity.PhotoWork
	err := r.db.Debug().Where("id = ?", id).Take(&photoWork).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("photo work not found")
	}
	return &photoWork, nil
}

// GetAllPhotoWork . . . .
func (r *PhotoWorkRepo) GetAllPhotoWork() ([]entity.PhotoWork, error) {
	var photoWorks []entity.PhotoWork
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&photoWorks).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return photoWorks, nil
}

// UpdatePhotoWork . . . .
func (r *PhotoWorkRepo) UpdatePhotoWork(photoWork *entity.PhotoWork) (*entity.PhotoWork, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&photoWork).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return photoWork, nil
}

/*
func (r *FoodRepo) DeleteFood(id uint64) error {
	var food entity.Food
	err := r.db.Debug().Where("id = ?", id).Delete(&food).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
*/
