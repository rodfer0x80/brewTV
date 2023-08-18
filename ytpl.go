package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
)

func extractVideoID(url string) string {
	match := regexp.MustCompile(`\?v=([^&]+)`).FindStringSubmatch(url)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

func getUUIDVideoPath(directoryPath string) string {
	return filepath.Join(directoryPath, uuid.New().String()+".mp4")
}

func ytdlVideoExec(videoPath string, videoID string) {
	exec.Command("yt-dlp", "--no-playlist", "--format", "mp4", "--output", videoPath, videoID).Run()
}

func downloadVideo(conWrite http.ResponseWriter, request *http.Request, url string) string {
	videoPath := getUUIDVideoPath(YTPL_PATH)
	ytdlVideoExec(videoPath, extractVideoID(url))
	return videoPath
}

func getVideoURL(conWrite http.ResponseWriter, request *http.Request) string {
	if request.Method == http.MethodPost {
		return request.FormValue("url")
	}
	return ""
}

func cleanupDownload(videoPath string) {
	os.Remove(videoPath)
	// if err := os.Remove(videoPath); err != nil {
	// fmt.Printf("Error deleting video file: %v\n", err)
	// }
}

func streamVideo(conWrite http.ResponseWriter, videoPath string) {
	videoFile, err := os.Open(videoPath)
	if err != nil {
		fmt.Printf("Error opening video file: %v\n", err)
	}
	defer videoFile.Close()
	conWrite.Header().Set("Content-Type", "video/mp4")
	io.Copy(conWrite, videoFile)
}

func YTPLVideoHandler(w http.ResponseWriter, request *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV YTPL</title>
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
			form {
				margin-top: 20px;
			}
			label {
				display: block;
				margin-bottom: 10px;
				font-weight: bold;
				color: #fff;
			}
			.center-input {
				display: flex;
				justify-content: center;
				align-items: center;
			}
			input[type="text"] {
				width: 100%;
				padding: 10px;
				border: none;
				border-radius: 4px;
				background-color: #333;
				color: #fff;
			}
			button[type="submit"] {
				background-color: #007bff;
				color: #fff;
				border: none;
				border-radius: 4px;
				padding: 12px 20px;
				cursor: pointer;
				margin-top: 10px;
			}
			button[type="submit"]:hover {
				background-color: #0056b3;
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
			<h1>BrewTV YTPL</h1>
			<form action="/ytpl/play" method="post">
				<label for="url">URL:</label>
				<div class="center-input">
					<input type="text" id="url" name="url" placeholder="https://www.youtube.com/watch?v=0123456789" required>
				</div>
				<button type="submit">Play</button>
			</form>
			<div class="back-button">
				<a href="/">Back</a>
			</div>
		</div>
	</body>
	</html>	
	`
	RenderTemplate(w, "ytpl", html, make([]string, 0))
}

func YTPLPlayVideoHandler(w http.ResponseWriter, r *http.Request) {
	url := getVideoURL(w, r)
	if url == "" {
		return
	}
	videoPath := downloadVideo(w, r, url)
	defer cleanupDownload(videoPath)
	streamVideo(w, videoPath)
}
