package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"text/template"
)

const DEBUG = 0

const LIBRARY_PATH = "/opt/brewTV/library"
const YTPL_PATH = "/opt/brewTV/ytpl"
const ACCEPTED_VIDEO_FORMAT = "mp4"
const VIDEO_PLAY_PARAMETER = "path"
const YT_URL_TEMPLATE_URL = "https://www.youtube.com/watch?v="

func main() {
	// Setup server interface and port
	tcpAddr := &net.TCPAddr{
		IP:   net.ParseIP(getLANIP()),
		Port: 80,
	}
	if DEBUG == 1 {
		tcpAddr = &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 8080,
		}
	}
	// Setup working volumes
	createDirIfNotExists(LIBRARY_PATH)
	createDirIfNotExists(YTPL_PATH)
	// Assign paths to handler functions
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/library", LibraryHandler)
	http.HandleFunc("/library/play", LibraryPlayHandler)
	http.HandleFunc("/ytpl", YTPLVideoHandler)
	http.HandleFunc("/ytpl/play", YTPLPlayVideoHandler)
	// Start server
	fmt.Printf("Running server: %s\n", tcpAddr.String())
	err := http.ListenAndServe(tcpAddr.String(), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getLANIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Failed to get local IP:", err)
		os.Exit(1)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	fmt.Println("Unable to determine LAN IP. Exiting...")
	os.Exit(1)
	return ""
}

func createDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func RenderTemplate(w http.ResponseWriter, name string, html string, data []string) {
	if template.Must(template.New(name).Parse(html)).Execute(w, data) != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #000;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background-color: #1a1a1a;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0px 2px 10px rgba(255, 255, 255, 0.1);
				width: 400px;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #fff;
			}
			ul {
				list-style: none;
				padding: 0;
				margin-top: 20px;
			}
			li {
				margin-bottom: 10px;
			}
			a {
				text-decoration: none;
				color: #007bff;
				font-weight: bold;
				font-size: 18px; /* Adjust the font size */
				margin-left: 10px; /* Add margin for spacing */
			}
			a:hover {
				color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>BrewTV</h1>
			<ul>
				<li><a href="/library">Library</a></li>
				<li><a href="/ytpl">YouTube</a></li>
			</ul>
		</div>
	</body>
	</html>	
	`
	RenderTemplate(w, "index", html, make([]string, 0))
}
