package main

import (
	"flag"
	"os"

	"github.com/lalabuy948/linkopus/backend/cmd"
	"github.com/lalabuy948/linkopus/backend/config"
	"github.com/lalabuy948/linkopus/backend/di"
	"github.com/lalabuy948/linkopus/backend/server"
)

var (
	compress = flag.Bool("c", true, "compress server handler")

	wipe = flag.Bool("wipe", false, "")
	date = flag.String("date", "", "")
	link = flag.String("link", "", "")

	help = flag.Bool("help", false, "")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	cfg := config.ParseEnv()
	container := di.Compile(cfg)

	if *wipe {
		wipeCmd := cmd.NewWipeCommand(container)
		wipeCmd.DeleteLinkMapAndView(*date, *link)
		os.Exit(0)
	}

	for i := 0; i < 1; i++ {
		go container.WorkerService.Consume()
	}

	s := server.NewServer(container)
	s.Start(&cfg.ServerPort, compress)
}
