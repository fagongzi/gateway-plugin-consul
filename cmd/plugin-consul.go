package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"runtime"

	"github.com/fagongzi/gateway-plugin-consul/pkg/conf"
	"github.com/fagongzi/gateway-plugin-consul/pkg/server"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/fagongzi/gateway/pkg/util"
)

var (
	cpus       = flag.Int("cpus", 1, "use cpu nums")
	logFile    = flag.String("log-file", "", "which file to record log, if not set stdout to use.")
	logLevel   = flag.String("log-level", "info", "log level.")
	configFile = flag.String("config", "", "config file")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*cpus)

	util.InitLog(*logFile)
	util.SetLogLevel(*logLevel)

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.PanicErrorf(err, "read config file <%s> failure.", *configFile)
	}

	cnf := &conf.Conf{}
	err = json.Unmarshal(data, cnf)
	if err != nil {
		log.PanicErrorf(err, "parse config file <%s> failure.", *configFile)
	}

	server := server.NewServer(cnf)
	server.Start()
}
