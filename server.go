package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func StartServer() {
	tcp_addr := &net.TCPAddr{
		IP:   net.ParseIP(GetLANIP()),
		Port: 80,
	}
	if DEBUG == 1 {
		tcp_addr = &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 8080,
		}
	}
	loadConfig()
	setupServer()
	runServer(*tcp_addr)
}

func runServer(tcp_addr net.TCPAddr) {
	log.Printf("[RunServer]::Running server: %s\n", tcp_addr.String())
	err := http.ListenAndServe(tcp_addr.String(), nil)
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
	allowed_mac_addrs, err := ReadAllowedMacAddresses(ALLOWED_MAC_ADDRESSES_PATH)
	if err != nil {
		log.Println("[RunServer]::Error loading allowed MAC address list")
	}
	fmt.Println("List of allowed MAC addresses:")
	for _, mac := range allowed_mac_addrs {
		fmt.Println(mac.String())
	}
}
