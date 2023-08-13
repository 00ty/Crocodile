package service

import (
	"Crocodile6/internal/consts"
	"Crocodile6/internal/dao"
	"Crocodile6/internal/utils"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"log"
)

var (
	queryDao = dao.NewQueryDao()
	ws       = websocket.Upgrader{}
)

type QueryService struct{}

func NewQueryService() *QueryService {
	return &QueryService{}
}

// Client 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// ClientManager 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Unregister chan *Client
}

// Message 信息转JSON (包括：发送者、接收者、内容)
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Unregister: make(chan *Client),
}

// 数组分页
func paginateArray(data []string, pageSize int) [][]string {
	pageCount := (len(data) + pageSize - 1) / pageSize
	paginatedArray := make([][]string, pageCount)
	for i := 0; i < pageCount; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > len(data) {
			end = len(data)
		}
		paginatedArray[i] = data[start:end]
	}
	return paginatedArray
}

// 验证卡密
func verifyJWT(jwtString string) error {
	if jwtString == "" || len(jwtString) < 10 {
		return errors.New("卡密为空")
	}
	token, err := jwt.ParseWithClaims(jwtString, &consts.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.StringToByteSlice(consts.Conf.Server.JwtSign), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*consts.JwtCustomClaims); ok && token.Valid {
		return nil
	}
	return errors.New("解析token失败")
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c //关闭前，将所有信息都传入进去
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		queryInput := new(consts.WsQueryInput)

		err := c.Socket.ReadJSON(&queryInput) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			//log.Println("数据格式不正确", err)
			msg, _ := json.Marshal(consts.WebSocketReturn{
				Type: "error",
				Msg:  "数据格式不正确",
			})
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			// 直接关闭连接
			Manager.Unregister <- c
			_ = c.Socket.Close()
			continue
		}
		// 验证token
		log.Println("token->", queryInput.Token)
		err = verifyJWT(queryInput.Token)
		if err != nil {
			//log.Println("验证token失败", err)
			msg, _ := json.Marshal(consts.WebSocketReturn{
				Type: "error",
				Msg:  "验证token失败",
			})
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			continue
		}
		// 调用分页
		pageSize := 100
		paginatedArray := paginateArray(queryInput.List, pageSize)
		for _, page := range paginatedArray {
			// 调用服务
			// 真几把坑
			var Data []consts.QueryDataDto
			err = queryDao.QueryCenter(queryInput.Type, page, &Data)
			if err != nil {
				log.Println("查询错误", err)
				msg, _ := json.Marshal(consts.WebSocketReturn{
					Type: "error",
					Msg:  err.Error(),
				})
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				break
			}
			for _, data := range Data {
				msg, _ := json.Marshal(consts.WebSocketReturn{
					Type: "data",
					Msg:  "Success",
					Data: data,
				})
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

// Query 查询中心服务
func (q *QueryService) Query(c echo.Context) error {
	conn, err := ws.Upgrade(c.Response(), c.Request(), nil) // 升级成ws协议
	if err != nil {
		return c.JSON(200, map[string]string{
			"msg": "错误",
		})
	}
	// 创建一个用户实例
	client := &Client{
		ID:     "1231", //1发给2
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	go client.Read()
	return nil
	/*log.Println("测试")
	for {
		client.Socket.PongHandler()
		_, msg, err := client.Socket.ReadMessage() // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据格式不正确", err)
			Manager.Unregister <- client
			_ = client.Socket.Close()
			break
		}
		log.Println("消息来了", msg)
		err = client.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil*/
}
