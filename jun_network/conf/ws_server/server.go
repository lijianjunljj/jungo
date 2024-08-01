package ws_server

import "time"

type Config struct {
	LogLevel        string
	LogPath         string
	WSAddr          string
	WSPort          string
	WSCltAddtr      string
	WSSSL           bool
	WSSSL_URL       string
	CertFile        string
	KeyFile         string
	TCPAddr         string
	HTTPTimeout     time.Duration
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32 //消息最大长度
	LenMsgLen       int
	LittleEndian    bool //字节序
}

var ServerConf Config

func init() {
	ServerConf.PendingWriteNum = 2000
	ServerConf.MaxMsgLen = 4096
	ServerConf.HTTPTimeout = 10 * time.Second
	ServerConf.LenMsgLen = 2
	ServerConf.LittleEndian = false
	ServerConf.MaxConnNum = 100
}
