package main

import (
	"Crocodile6/internal/bootstrap"
	"Crocodile6/internal/consts"
	"Crocodile6/internal/utils"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"time"
)

var (
	Name     = "GoQuery802"
	flagConf string
	err      error
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "配置文件地址, 比如: -conf config.yaml")
}

func main() {
	flag.Parse()
	fmt.Println(Name)
	fmt.Println(flagConf)
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	defer c.Close()

	if err = c.Load(); err != nil {
		panic(err)
	}

	if err = c.Scan(&consts.Conf); err != nil {
		panic(err)
	}
	// 监听配置
	if err := c.Watch("app.notice", func(key string, value config.Value) {
		notice, err := value.String()
		if err != nil {
			return
		}
		log.Println("公告修改为", notice)
		consts.Conf.App.Notice = notice
	}); err != nil {
		log.Println("配置监听失败", err)
	}

	err = utils.InitDB(consts.Conf.Mysql)
	if err != nil {
		log.Panic("连接数据库失败，请检查参数：", err)
	}
	log.Println("数据库连接成功")

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Panic("无法加载上海时区:", err)
		return
	}
	consts.CurrentTime = time.Now().In(location)
	//d := dao.NewQueryDao()
	//var r []consts.QueryDataDto
	/*err = d.QueryCenter("GetQQUidForMobile", []string{
		"15944850489",
	}, &r)*/ // GetMobileUidForQQ

	/*err = d.QueryCenter("GetMobileUidForQQ", []string{
		"259213601",
		"1561576257",
		"3092093560",
		"2012725357",
		"2646174195",
		"690420902",
		"2784344025",
		"1073943080",
		"3278203310",
		"1310989912",
		"2993626980",
		"643091445",
		"2155726893",
		"669630546",
		"453312321",
		"2557497600",
		"1207419059",
		"1709451622",
		"362364470",
	}, &r)
	if err != nil {
		log.Println("查询错误->", err)
		return
	}
	log.Println("查询结果为->", r)*/

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, ip=${remote_ip} time=${time_rfc3339}\n",
	}))
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true
	bootstrap.InitRouter(e)
	e.Logger.Fatal(e.Start(
		fmt.Sprintf(":%d", consts.Conf.Server.Port),
	))
}
