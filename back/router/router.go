package router

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/2024-AC-hacking/config"
	"github.com/walnuts1018/2024-AC-hacking/psql"
	"github.com/walnuts1018/2024-AC-hacking/router/handler"
)

func NewRouter(config config.Config, psqlClient *psql.Client) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(sloggin.New(slog.Default()))

	if config.LogLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// session
	store, err := postgres.NewStore(psqlClient.DB().DB, []byte("secret"))
	if err != nil {
		slog.Error("Failed to create session store")
	}

	r.Use(sessions.Sessions("session", store))

	handler := handler.NewHandler(config, psqlClient)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", handler.Ping)
		v1.POST("/proxy-login", handler.ProxyLogin)

		aircon := v1.Group("/aircon")
		aircon.Use(sessionMiddleware())
		{
			aircon.GET("/status", handler.GetAirconStatus)
		}
	}

	return r
}

func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		firstProxy := session.Get("first_proxy")
		if firstProxy == nil {
			slog.Info("Proxy login required")
			c.Redirect(302, "/api/v1/proxy-login")
			c.Abort()
			return
		}

		c.Set("first_proxy", firstProxy)
		c.Next()
	}
}
