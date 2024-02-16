package conf

import (
	"encoding/json"
	"flag"
	"github.com/lijianjunljj/jungo/jun_log"
	"os"
)

var configFile = flag.String("node", "", "")

var NodeConf struct {
	NodeName       string
	IsCenter       bool
	CenterNodeHost string
	Cookie         string
}

func init() {
	flag.Parse()
	jun_log.Debug("开始初始化节点配置%s", *configFile)
	data, err := os.ReadFile(*configFile)
	if err != nil {
		jun_log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &NodeConf)
	if err != nil {
		jun_log.Fatal("%v", err)
	}
	NodeName = NodeConf.NodeName
	IsCenter = NodeConf.IsCenter
	CenterNodeHost = NodeConf.CenterNodeHost
	Cookie = NodeConf.Cookie
	jun_log.Debug("初始化节点配置成功：%v", NodeConf)
}
