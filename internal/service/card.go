package service

import (
	"Crocodile6/internal/consts"
	"Crocodile6/internal/dao"
	"fmt"
	"github.com/labstack/echo"
)

type CardService struct{}

func NewCardService() *CardService {
	return &CardService{}
}

func (ctr *CardService) Index(c echo.Context) error {

	return c.String(consts.StatusOK, "这是首页")
}

type AuthorizationDto struct {
	SecretKey string `json:"secret_key" validate:"required"`
	Uid       string `json:"uid" validate:"required"`
}

// Authorization 验证卡密
func (ctr *CardService) Authorization(c echo.Context) (err error) {
	dto := new(AuthorizationDto)
	if err = c.Bind(dto); err != nil {
		return c.JSON(consts.StatusOK, consts.JsonResult{
			Code: 500,
			Msg:  "请求失败，请注意参数1",
		})
	}
	if err = c.Validate(dto); err != nil {
		return c.JSON(consts.StatusOK, consts.JsonResult{
			Code: 500,
			Msg:  "请求失败，请注意参数2",
		})
	}
	fmt.Println(dto)
	// 调用服务
	cardService := dao.NewCardDao()
	k, err := cardService.Verify(dto.SecretKey, dto.Uid)
	if err != nil {
		return c.JSON(consts.StatusOK, consts.JsonResult{
			Code: 500,
			Msg:  err.Error(),
		})
	}
	return c.JSON(consts.StatusOK, consts.JsonResult{
		Code: 200,
		Msg:  "请求成功",
		Data: k,
	})
}
