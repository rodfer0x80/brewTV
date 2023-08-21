package main

import "flag"

const DEBUG = 0

func main() {
	startFlag := flag.Bool("start", false, "Run the server")
	setupFlag := flag.Bool("setup", false, "Configure allowed MAC addresses")

	flag.Parse()

	if *startFlag {
		StartServer()
	}

	if *setupFlag {
		ConfigureAllowedMacAddresses()
	}
}
