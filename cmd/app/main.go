package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ninehills/go-webapp-template/internal/app"
	"github.com/ninehills/go-webapp-template/pkg/version"
)

const defaultCfgFile = "./config/config.yml"

func main() {
	v := flag.Bool("v", false, "print version")
	h := flag.Bool("h", false, "print help")
	cfgFile := flag.String("c", defaultCfgFile, "config file")
	flag.Parse()
	if *v {
		fmt.Println(version.GetVersion().String())
		os.Exit(0)
	}
	if *h {
		fmt.Printf("Usage: %s [-v] [-h] [-c config file]\n", os.Args[0])
		os.Exit(0)
	}
	app.Run(*cfgFile)
}
