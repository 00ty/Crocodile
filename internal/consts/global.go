package consts

import (
	"Crocodile6/internal/conf"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Conf conf.AutoGenerated
)
