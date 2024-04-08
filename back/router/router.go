package router

import (
	"fmt"
	"log/slog"

	"github.com/abrander/ginproxy"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/2024-AC-hacking/config"
	"github.com/walnuts1018/2024-AC-hacking/psql"
	"github.com/walnuts1018/2024-AC-hacking/router/handler"
)

func NewRouter(config config.Config, psqlClient *psql.Client) (*gin.Engine, error) {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(sloggin.New(slog.Default()))

	if config.LogLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// session
	store, err := postgres.NewStore(psqlClient.DB().DB, []byte("secret"))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create session store: %v", err))
	}

	r.Use(sessions.Sessions("session", store))
	handler, err := handler.NewHandler(config, psqlClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %v", err)
	}

	// Nextjs側にそのままProxy
	ginProxy, err := ginproxy.NewGinProxy(config.FrontURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy: %v", err)
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/aircon")
	})

	r.GET("/ping", handler.Ping)
	r.GET("/proxy-login", ginProxy.Handler)
	r.POST("/proxy-login", handler.ProxyLogin)
	r.Any("/_next/*all", ginProxy.Handler)

	aircon := r.Group("/aircon")
	aircon.Use(sessionMiddleware())
	{
		aircon.GET("", ginProxy.Handler)
		aircon.GET("/status", ginProxy.Handler).Use(loginSessionMiddleware())
		aircon.GET("/status-json", handler.GetAirconStatus).Use(loginSessionMiddleware())
		aircon.GET("/operate", ginProxy.Handler).Use(loginSessionMiddleware()).Use(adminSessionMiddleware())
		aircon.POST("/operate", handler.OperationAircon).Use(loginSessionMiddleware()).Use(adminSessionMiddleware())
	}

	r.GET("/login", ginProxy.Handler).Use(sessionMiddleware())
	r.POST("/login", handler.Login).Use(sessionMiddleware())
	r.GET("/logout", ginProxy.Handler).Use(sessionMiddleware())
	r.POST("/logout", handler.Logout).Use(sessionMiddleware())
	r.GET("/register", ginProxy.Handler).Use(sessionMiddleware())
	r.POST("/register", handler.Register).Use(sessionMiddleware())
	r.GET("/check-login", handler.CheckLogin).Use(sessionMiddleware())
	r.GET("/user", ginProxy.Handler).Use(sessionMiddleware()).Use(loginSessionMiddleware())
	r.GET("/userapi", handler.GetUser).Use(sessionMiddleware()).Use(loginSessionMiddleware())
	r.PUT("/userapi", handler.UpdateUser).Use(sessionMiddleware()).Use(loginSessionMiddleware())

	return r, nil
}

func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		firstProxy := session.Get("first_proxy")
		if firstProxy == nil {
			slog.Info("Proxy login required")
			c.Redirect(302, "/proxy-login")
			c.Abort()
			return
		}

		c.Set("first_proxy", firstProxy)
		c.Next()
	}
}

func loginSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			slog.Info("Login required")
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func adminSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			slog.Info("Login required")
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		if userID != "admin" {
			slog.Info("Admin required")
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
