package repository

import (
	"github.com/uwaifo/video_server_api/domian/entity"
)

// PhotoWorkRepositry  . . .
type PhotoWorkRepositry interface {
	SavePhotoWork(*entity.PhotoWork) (*entity.PhotoWork, map[string]string)
}
