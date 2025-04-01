package global

import (
	"uploader/model/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ServerConfig *config.Server
	Logger       *zap.Logger
	DB           *gorm.DB
)
