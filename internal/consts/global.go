package consts

import (
	"Crocodile6/internal/conf"
	"gorm.io/gorm"
	"time"
)

var (
	DB          *gorm.DB
	Conf        conf.AutoGenerated
	CurrentTime time.Time
)