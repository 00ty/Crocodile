package bootstrap

import (
	"Crocodile6/internal/consts"
	"Crocodile6/internal/service"
	"Crocodile6/internal/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func InitRouter(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
	card := service.NewCardService()
	query := service.NewQueryService()

	e.Static("/", "public")
	e.GET("/ws", query.Query)
	// 登录
	e.POST("/user/login", card.Authorization)
	// 公告
	e.GET("/notice", func(c echo.Context) error {
		return c.JSON(consts.StatusOK, consts.JsonResult{
			Code: 200,
			Msg:  "获取成功",
			Data: consts.Conf.App.Notice,
		})
	})
	r := e.Group("/query")
	{
		r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     new(consts.JwtCustomClaims),
			SigningKey: utils.StringToByteSlice(consts.Conf.Server.JwtSign),
		}))
		r.GET("/get", card.Index)
	}
}
