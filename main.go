package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const moviesDir = "./movies"
const host = "LOCAL" // INET // LOCAL

func main() {
	config()

	IP := "127.0.0.1"
	port := "8080"

	if host == "LAN" {
		IP := getLANIP()
		port = "80"
		fmt.Println("Server LAN IP:", IP+":"+port)
	}
	if host == "INET" {
		IP = "0.0.0.0"
		port = "80"
		fmt.Println("Server INET:", IP+":"+port)
	}

	fmt.Println("Server LOCAL:", IP+":"+port)

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/movies/", streamMovieHandler)

	http.ListenAndServe(IP+":"+port, nil)
}

func config() {
	createDirIfNotExists(moviesDir)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := listMovies()
	if err != nil {
		http.Error(w, "Failed to list movies", http.StatusInternalServerError)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Cinema Movies</title>
	</head>
	<body>
		<h1>Cinema Movies</h1>
		<ul>
			{{range .}}
			<li><a href="/movies/{{.}}">{{.}}</a></li>
			{{end}}
		</ul>
	</body>
	</html>
	`

	t := template.Must(template.New("index").Parse(tmpl))
	err = t.Execute(w, movies)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func listMovies() ([]string, error) {
	var movies []string

	err := filepath.Walk(moviesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp4") {
			relPath, err := filepath.Rel(moviesDir, path)
			if err != nil {
				return err
			}
			movies = append(movies, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movies, nil
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

func streamMovieHandler(w http.ResponseWriter, r *http.Request) {
	moviePath := filepath.Join(moviesDir, filepath.Base(r.URL.Path))
	if isFileExists(moviePath) {
		serveVideo(w, r, moviePath)
		return
	}

	http.NotFound(w, r)
}

func serveVideo(w http.ResponseWriter, r *http.Request, filePath string) {
	movieFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open the movie file", http.StatusInternalServerError)
		return
	}
	defer movieFile.Close()

	w.Header().Set("Content-Type", "video/mp4")
	_, err = io.Copy(w, movieFile)
	if err != nil {
		http.Error(w, "Failed to stream the movie", http.StatusInternalServerError)
		return
	}
}

func isFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
