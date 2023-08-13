package dao

import (
	"Crocodile6/internal/consts"
	"Crocodile6/internal/models"
	"Crocodile6/internal/utils"
	"errors"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"log"
	"time"
)

type CardDao struct{}

func NewCardDao() *CardDao {
	return &CardDao{}
}

func isExpired(expirationTime time.Time) bool {
	currentTime := time.Now()
	return currentTime.After(expirationTime)
}

// Verify 验证卡密
func (s *CardDao) Verify(secretKey string, Uid string) (models.Kami, error) {
	var err error
	var kami = models.Kami{}
	// 新增需求 修改时间：2023年8月13日12:32:46
	err = consts.DB.
		Where("text = ? AND state = ?",
			secretKey, 0).
		First(&kami).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kami, errors.New("找不到此卡密信息")
		}
		return kami, errors.New("查询失败")
	}
	// 判断
	if kami.Username != Uid && kami.Username != "" {
		return kami, errors.New("卡密已被绑定")
	}
	if isExpired(kami.EndDate) {
		log.Println(kami.EndDate, "已过期")
		return kami, errors.New("卡密已经过期")
	}
	// 生成用户
	err = consts.DB.Model(&kami).
		Updates(map[string]interface{}{
			"username": Uid,
			"use_date": time.Now(),
		}).Error

	if err != nil {
		log.Println("更新失败", err)
		return kami, errors.New("更新失败")
	}
	claims := &consts.JwtCustomClaims{
		Email: Uid,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(utils.StringToByteSlice(consts.Conf.Server.JwtSign))
	if err != nil {
		return kami, errors.New("生成Token失败")
	}
	kami.Ext = t
	kami.EndStr = kami.EndDate.Format("20060102")
	return kami, nil
}
