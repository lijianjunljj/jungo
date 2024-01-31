package conf

import (
	"log"
	"time"
)

var (
	NodeName         = "local"
	ServerPort       = 4321
	ServerMaxConnNum = 1000
	ServerListenAddr = "0.0.0.0"
	CertFile         = ""
	KeyFile          = ""
	LogFlag          = log.LstdFlags
	LogLevel         = "debug"
	LogPath          = ""
	LenStackBuf      = 4096
	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian    bool   = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)
