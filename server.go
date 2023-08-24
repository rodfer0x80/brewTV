package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	DefaultPort = 80
	DebugPort   = 8080
)

func StartServer() {
	tcpAddr := getServerTCPAddr()
	loadConfig()

	setupServer()

	runServer(tcpAddr)
}

func getServerTCPAddr() net.TCPAddr {
	ip := net.ParseIP(GetLANIP())
	port := DefaultPort

	if DEBUG == 1 {
		ip = net.ParseIP("127.0.0.1")
		port = DebugPort
	}

	return net.TCPAddr{
		IP:   ip,
		Port: port,
	}
}

func runServer(tcpAddr net.TCPAddr) {
	serverAddr := tcpAddr.String()
	log.Printf("[RunServer]::Running server: %s\n", serverAddr)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Println("[RunServer]::Error starting server:", err)
	}
}

func setupServer() {
	CreateDirIfNotExists(LIBRARY_PATH)
	CreateDirIfNotExists(YTPL_PATH)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/library", LibraryHandler)
	http.HandleFunc("/library/play", LibraryPlayHandler)
	http.HandleFunc("/ytpl", YTPLVideoHandler)
	http.HandleFunc("/ytpl/play", YTPLPlayVideoHandler)
}

func loadConfig() {
	allowedMacAddrs, err := ReadAllowedMacAddresses(ALLOWED_MAC_ADDRESSES_PATH)
	if err != nil {
		log.Println("[LoadConfig]::Error loading allowed MAC address list:", err)
	}

	fmt.Println("List of allowed MAC addresses:")
	for _, mac := range allowedMacAddrs {
		fmt.Println(mac.String())
	}
}
