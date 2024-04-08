package handler

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

func NewHandler(config config.Config, psqlClient *psql.Client) (*Handler, error) {

	return &Handler{
		config:     config,
		psqlClient: psqlClient,
	}, nil
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type req struct {
	Password string `json:"password"`
}

func (h *Handler) ProxyLogin(c *gin.Context) {
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

type statusResp struct {
	Power    string `json:"power"`
	Mode     string `json:"mode"`
	Temp     string `json:"temp"`
	Fan      string `json:"fan"`
	Swing    string `json:"swing"`
	Eco      string `json:"eco"`
	OnTimer  string `json:"ontimer"`
	OffTimer string `json:"offtimer"`
}

func (h *Handler) GetAirconStatus(c *gin.Context) {
	power := "0"
	if h.airconStatus.PowerOn {
		power = "1"
	}
	swing := "0"
	if h.airconStatus.Swing {
		swing = "1"
	}
	eco := "0"
	if h.airconStatus.Eco {
		eco = "1"
	}

	status := statusResp{
		Power:    power,
		Mode:     fmt.Sprintf("%d", h.airconStatus.Mode),
		Temp:     fmt.Sprintf("%d", h.airconStatus.Temp),
		Fan:      fmt.Sprintf("%d", h.airconStatus.Fan),
		Swing:    swing,
		Eco:      eco,
		OnTimer:  fmt.Sprintf("%d", h.airconStatus.OnTimerHour),
		OffTimer: fmt.Sprintf("%d", h.airconStatus.OffTimerHour),
	}

	c.JSON(200, status)
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
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

func (h *Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(200, gin.H{
		"message": "logout success",
	})
}

func (h *Handler) Register(c *gin.Context) {
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

	session := sessions.Default(c)
	session.Set("user_id", r.Username)
	session.Save()

	c.JSON(200, gin.H{
		"message": "register success",
	})
}

func (h *Handler) CheckLogin(c *gin.Context) {
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

func (h *Handler) GetUser(c *gin.Context) {
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

func (h *Handler) UpdateUser(c *gin.Context) {
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

type OperationReq struct {
	Power    string `json:"power"`
	Mode     string `json:"mode"`
	Temp     string `json:"temp"`
	Fan      string `json:"fan"`
	Swing    string `json:"swing"`
	Eco      string `json:"eco"`
	OnTimer  string `json:"ontimer"`
	OffTimer string `json:"offtimer"`
}

func (h *Handler) OperationAircon(c *gin.Context) {

	var r OperationReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
		slog.Error(fmt.Sprintf("Failed to bind json: %v", err), slog.Any("request", r))
		return
	}

	power, err := strconv.Atoi(r.Power)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid power",
		})
		return
	}

	if power != 0 && power != 1 {
		c.JSON(400, gin.H{
			"message": "invalid power",
		})
		return
	}

	mode, err := strconv.Atoi(r.Mode)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid mode",
		})
		return
	}
	if mode < 0 || mode > 4 {
		c.JSON(400, gin.H{
			"message": "invalid mode",
		})
		return
	}
	fan, err := strconv.Atoi(r.Fan)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid mode",
		})
		return
	}
	if fan < 0 || fan > 2 {
		c.JSON(400, gin.H{
			"message": "invalid fan",
		})
		return
	}

	temp, err := strconv.Atoi(r.Temp)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid temp",
		})
		return
	}

	swing, err := strconv.Atoi(r.Swing)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid swing",
		})
		return
	}

	if swing != 0 && swing != 1 {
		c.JSON(400, gin.H{
			"message": "invalid swing",
		})
		return
	}

	eco, err := strconv.Atoi(r.Eco)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid eco",
		})
		return
	}

	if eco != 0 && eco != 1 {
		c.JSON(400, gin.H{
			"message": "invalid eco",
		})
		return
	}

	onTimer, err := strconv.Atoi(r.OnTimer)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid on timer",
		})
		return
	}
	if onTimer < 0 {
		c.JSON(400, gin.H{
			"message": "invalid on timer",
		})

		return
	}

	offTimer, err := strconv.Atoi(r.OffTimer)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid off timer",
		})
		return
	}

	if offTimer < 0 {
		c.JSON(400, gin.H{
			"message": "invalid off timer",
		})

		return
	}

	json := fmt.Sprintf(`{
		"power": %d,
		"mode": %d,
		"temp": %d,
		"fan": %d,
		"swing": %d,
		"eco": %d,
		"ontimer": %d,
		"offtimer": %d
	}`, power, mode, temp, fan, swing, eco, onTimer, offTimer)

	req, err := http.NewRequest(
		"POST",
		"http://192.168.0.109/v1/post",
		bytes.NewBuffer([]byte(json)),
	)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to create request",
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to send request",
		})
		slog.Error(fmt.Sprintf("Failed to send request: %v", err))
		return
	}

	defer resp.Body.Close()

	h.airconStatus.PowerOn = (power == 1)
	h.airconStatus.Mode = Mode(mode)
	h.airconStatus.Temp = temp
	h.airconStatus.Fan = FanMode(fan)
	h.airconStatus.Swing = (swing == 1)
	h.airconStatus.Eco = (eco == 1)
	h.airconStatus.OnTimerHour = onTimer
	h.airconStatus.OffTimerHour = offTimer

	c.JSON(200, gin.H{
		"message": "success",
	})
}
