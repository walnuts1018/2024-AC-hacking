package handler

type AirconStatus struct {
	PowerOn      bool
	Mode         Mode
	Temp         int
	Fan          FanMode
	Swing        bool
	Eco          bool
	OnTimerHour  int
	OffTimerHour int
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
