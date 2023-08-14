package consts

import "github.com/golang-jwt/jwt"

const StatusOK = 200

// JwtCustomClaims 验证
type JwtCustomClaims struct {
	Email string `json:"email"`
	*jwt.StandardClaims
}

type JsonResult struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type QueryDataDto struct {
	QQ      string `json:"qq"`
	Mobile  string `json:"mobile"`
	Uid     string `json:"uid"`
	Address string `json:"address"`
}

type WsQueryInput struct {
	Token string `json:"token"`
	/*
		查询类型：
		QQ查询手机号和微博ID   mobile_uid_by_qq

	*/
	Type string   `json:"type"`
	List []string `json:"list"`
}

// QueryDto 真坑
type QueryDto struct {
	Address  string `json:"address"`
	UserName string `gorm:"column:username" json:"qq"`
	Mobile   string `gorm:"column:mobile" json:"mobile"`
	Uid      string `gorm:"column:uid" json:"uid"`
}

//type WeiboDto struct {
//	Uid    string `gorm:"column:uid" json:"uid"`
//	Mobile string `gorm:"column:mobile" json:"mobile"`
//}

// WebSocketReturn 定义一个WebSocket返回消息规范
type WebSocketReturn struct {
	Type string      `json:"type"` // data, error
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
