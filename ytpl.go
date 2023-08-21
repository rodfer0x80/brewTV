package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
)

const YTPL_PATH = "/opt/brewTV/ytpl"
const YT_URL_TEMPLATE_URL = "https://www.youtube.com/watch?v="

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
}

func streamVideo(conWrite http.ResponseWriter, videoPath string) {
	videoFile, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Error opening video file: %v\n", err)
	}
	defer videoFile.Close()
	conWrite.Header().Set("Content-Type", "video/mp4")
	io.Copy(conWrite, videoFile)
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
