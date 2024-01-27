package conf

var NodeConf struct {
	NodePort       string
	NodeListenAddr string
}

func init() {
	NodeConf.NodePort = "4321"
	NodeConf.NodeListenAddr = "0.0.0.0"
}
