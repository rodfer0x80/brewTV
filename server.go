package main

import (
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
	CreateDirIfNotExists(MUSIC_PATH)
	CreateDirIfNotExists(TV_PATH)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/tv", TVHandler)
	http.HandleFunc("/tv/play", TVPlayHandler)
	http.HandleFunc("/music", MusicHandler)
	http.HandleFunc("/music/play", MusicPlayHandler)
	http.HandleFunc("/ytpl", YTPLVideoHandler)
	http.HandleFunc("/ytpl/play", YTPLPlayVideoHandler)
}

func loadConfig() {
	allowedMacAddrs, err := ReadAllowedMacAddresses(ALLOWED_MAC_ADDRESSES_PATH)
	if err != nil {
		log.Println("[LoadConfig]::Error loading allowed MAC address list:", err)
	}

	log.Println("List of allowed MAC addresses:")
	for _, mac := range allowedMacAddrs {
		log.Println(mac.String())
	}
}
