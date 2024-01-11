package jun_log

import (
	"github.com/lijianjunljj/jungo/jun_config"
)

func Init() {
	logger, err := New(jun_config.LogLevel, jun_config.LogPath, jun_config.LogFlag)
	if err != nil {
		panic(err)
	}
	Export(logger)
	//defer logger.Close()
}
