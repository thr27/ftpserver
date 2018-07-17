// ftpserver allows to create your own FTP(S) server
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	//"github.com/fclairamb/ftpserver/sample"
	"github.com/thr27/ftpserver/sample"
	"github.com/thr27/ftpserver/server"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	ftpServer *server.FtpServer
)

func confFileContent() []byte {
	str := `# ftpserver configuration file
#
# These are all the config parameters with their default values. If not present,
# Max number of control connections to accept
# max_connections = 0
max_connections = 10
[server]
# Address to listen on
# listen_host = "0.0.0.0"
# Port to listen on
listen_port = 2121
# Public host to expose in the passive connection
# public_host = ""
# Idle timeout time
# idle_timeout = 900
# Data port range from 10000 to 15000
# [dataPortRange]
# start = 2122
# end = 2200
[server.dataPortRange]
start = 2122
end = 2200
[[users]]
user="admin"
pass="123456"
dir="shared"
[[users]]
user="test"
pass="test"
dir="shared"
`
	return []byte(str)
}
func main() {
	var err error
	flag.Parse()
	confFile := "settings.toml"
	if err = ioutil.WriteFile(confFile, confFileContent(), 0644); err != nil {
		log15.Error("Couldn't create config file", "action", "conf_file.could_not_create", "confFile", confFile)
	}
	//drv, err := sample.NewSampleDriver("", confFile)
	drv := sample.NewSampleDriver()
	//if err != nil {
	//	log15.Error("Problem creating driver", "err", err)
	//}
	ftpServer = server.NewFtpServer(drv)

	go signalHandler()

	err = ftpServer.ListenAndServe()
	if err != nil {
		log15.Error("Problem listening", "err", err)
	}
}

func signalHandler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	for {
		switch <-ch {
		case syscall.SIGTERM:
			ftpServer.Stop()
			break
		}
	}
}
