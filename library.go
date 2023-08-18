package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	library, err := listFilesByType(LIBRARY_PATH, ACCEPTED_VIDEO_FORMAT)
	if err != nil {
		http.Error(w, "Failed to list library", http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV Library</title>
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
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #007bff;
				font-weight: bold;
			}
			.back-button a:hover {
				color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>BrewTV Library</h1>
			<ul>
				{{range .}}
				<li><a href="/library/play?path={{.}}">{{.}}</a></li>
				{{end}}
			</ul>
			<div class="back-button">
				<a href="/">Back</a>
			</div>
		</div>
	</body>
	</html>	
	`
	RenderTemplate(w, "library", html, library)
}

func LibraryPlayHandler(responseWriter http.ResponseWriter, request *http.Request) {
	relative_video_path := getURLParameter(responseWriter, request, VIDEO_PLAY_PARAMETER)
	if relative_video_path == "" {
		return
	}
	playVideo(responseWriter, request, relative_video_path)
}

func playVideo(responseWriter http.ResponseWriter, request *http.Request, relative_video_path string) {
	video_path := filepath.Join(LIBRARY_PATH, relative_video_path)
	if fileExists(video_path) {
		serveVideo(responseWriter, request, video_path, ACCEPTED_VIDEO_FORMAT)
		return
	}
	http.NotFound(responseWriter, request)
}

func getURLParameter(w http.ResponseWriter, r *http.Request, parameter_name string) string {
	if r.Method == http.MethodGet {
		return r.FormValue(parameter_name)
	}
	return ""
}

func listFilesByType(directoryPath string, file_format string) ([]string, error) {
	files := []string{}
	print(directoryPath)
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), fmt.Sprintf(".%s", file_format)) {
			files = append(files, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func serveVideo(w http.ResponseWriter, r *http.Request, file_path string, file_format string) {
	video_handle, err := os.Open(file_path)
	if err != nil {
		http.Error(w, "Failed to open the video file", http.StatusInternalServerError)
		return
	}
	defer video_handle.Close()

	w.Header().Set("Content-Type", fmt.Sprintf("video/%s", file_format))
	_, err = io.Copy(w, video_handle)
	if err != nil {
		http.Error(w, "Failed to stream the video content", http.StatusInternalServerError)
		return
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
