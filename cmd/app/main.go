package main

import (
	"flag"
	"log"
	"os"

	"github.com/ninehills/go-webapp-template/internal/app"
	"github.com/ninehills/go-webapp-template/pkg/version"
)

const defaultCfgFile = "./config/config.yml"

func main() {
	v := flag.Bool("v", false, "print version")
	if *v {
		log.Println(version.GetVersion().String())
		os.Exit(0)
	}

	h := flag.Bool("h", false, "print help")
	if *h {
		log.Printf("Usage: %s [-v] [-h] [-c config file]\n", os.Args[0])
		os.Exit(0)
	}

	cfgFile := flag.String("c", defaultCfgFile, "config file")
	flag.Parse()

	app.Run(*cfgFile)
}
