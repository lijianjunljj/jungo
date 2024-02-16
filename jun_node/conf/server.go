package conf

import (
	"time"
)

var (
	NodeName                = "local"
	IsCenter                = false
	CenterNodeHost          = "127.0.0.1"
	Cookie                  = "123abc456*"
	ServerPort              = 4321
	ServerMaxConnNum        = 1000
	ServerListenAddr        = "0.0.0.0"
	CertFile                = ""
	KeyFile                 = ""
	PendingWriteNum         = 2000
	MaxMsgLen        uint32 = 4096
	HTTPTimeout             = 10 * time.Second
	LenMsgLen               = 2
	LittleEndian     bool   = false
)
