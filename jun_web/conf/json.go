package conf

var Server struct {
	LogLevel         string
	LogPath          string
	HttpPort         int
	TokenAdminExp    string
	TokenAdminSecret string
	ConsolePort      int
	ProfilePath      string
	JwtSecret        string
	JwtValidTime     string
	EnableJwt        bool
}

func init() {
	Server.HttpPort = 8080
	Server.JwtSecret = "SHGDJKSADFSJKKSFKGFHiD"
	Server.JwtValidTime = "24h"
}
