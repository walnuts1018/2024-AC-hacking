package handler

import (
	"fmt"
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/2024-AC-hacking/config"
	"github.com/walnuts1018/2024-AC-hacking/psql"
)

type Handler struct {
	config     config.Config
	psqlClient *psql.Client

	airconStatus AirconStatus
}

func NewHandler(config config.Config, psqlClient *psql.Client) (Handler, error) {

	return Handler{
		config:     config,
		psqlClient: psqlClient,
	}, nil
}

func (h Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type req struct {
	Password string `json:"password"`
}

func (h Handler) ProxyLogin(c *gin.Context) {
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}

	// if r.Password != h.config.ProxyPassword {
	// 	return
	// }

	session := sessions.Default(c)
	session.Set("first_proxy", true)
	session.Save()

	c.Redirect(302, "/aircon")
}

func (h Handler) GetAirconStatus(c *gin.Context) {
	c.JSON(200, h.airconStatus)
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h Handler) Login(c *gin.Context) {
	var r loginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}
	user, err := h.psqlClient.Login(r.Username, r.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "login failed",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.Username)
	session.Save()

	c.JSON(200, gin.H{
		"message": "login success",
	})
}

func (h Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(200, gin.H{
		"message": "logout success",
	})
}

func (h Handler) Register(c *gin.Context) {
	var r loginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}

	if err := h.psqlClient.Register(r.Username, r.Password); err != nil {
		c.JSON(400, gin.H{
			"message": "register failed",
		})
		slog.Error(fmt.Sprintf("Failed to register: %v", err))
		return
	}

	c.JSON(200, gin.H{
		"message": "register success",
	})
}

func (h Handler) CheckLogin(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(200, gin.H{
			"loggined": false,
			"username": "",
		})
		return
	}

	c.JSON(200, gin.H{
		"loggined": true,
		"username": userID,
	})
}

func (h Handler) GetUser(c *gin.Context) {
	session := sessions.Default(c)

	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(400, gin.H{
			"message": "not loggined",
		})
		return
	}

	userName, ok := userID.(string)
	if !ok {
		c.JSON(400, gin.H{
			"message": "invalid user id",
		})
		return
	}

	users, err := h.psqlClient.GetUser(userName)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "failed to get user",
		})
		return
	}

	c.JSON(200, users)
}

type updateUserReq struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

func (h Handler) UpdateUser(c *gin.Context) {
	session := sessions.Default(c)

	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(400, gin.H{
			"message": "not loggined",
		})
		return
	}

	userName, ok := userID.(string)
	if !ok {
		c.JSON(400, gin.H{
			"message": "invalid user id",
		})
		return
	}

	var r updateUserReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}

	if err := h.psqlClient.ChangePassword(userName, r.OldPassword, r.Password); err != nil {
		c.JSON(400, gin.H{
			"message": "failed to change password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "password changed",
	})
}
