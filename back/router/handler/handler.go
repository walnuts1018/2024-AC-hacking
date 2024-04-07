package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/2024-AC-hacking/config"
	"github.com/walnuts1018/2024-AC-hacking/psql"
)

type Handler struct {
	config       config.Config
	psqlClient   *psql.Client
	airconStatus AirconStatus
}

func NewHandler(config config.Config, psqlClient *psql.Client) Handler {
	return Handler{
		config:     config,
		psqlClient: psqlClient,
	}
}

func (h Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h Handler) ProxyLogin(c *gin.Context) {
	proxyPassword := c.PostForm("password")
	_ = proxyPassword
	// if proxyPassword != h.config.ProxyPassword {
	// 	return
	// }

	session := sessions.Default(c)
	session.Set("first_proxy", true)
	session.Save()

	c.Redirect(302, "/dashboard")
}

func (h Handler) GetAirconStatus(c *gin.Context) {
	c.JSON(200, h.airconStatus)
}
