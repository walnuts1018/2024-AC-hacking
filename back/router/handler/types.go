package handler

type AirconStatus struct {
	PowerOn      bool    `json:"power"`
	Mode         Mode    `json:"mode"`
	Temp         int     `json:"temp"`
	Fan          FanMode `json:"fan"`
	Swing        bool    `json:"swing"`
	Eco          bool    `json:"eco"`
	OnTimerHour  int     `json:"ontimer"`
	OffTimerHour int     `json:"offtimer"`
}

type Mode int

const (
	Heat Mode = 0
	Dry  Mode = 1
	Cool Mode = 2
	Fan  Mode = 3
)

type FanMode int

const (
	Auto   FanMode = 0
	Low    FanMode = 1
	Medium FanMode = 2
	High   FanMode = 3
)
