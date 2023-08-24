package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
)

const (
	YTPL_PATH           = "/opt/brewTV/ytpl"
	YT_URL_TEMPLATE_URL = "https://www.youtube.com/watch?v="
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

func ytdlVideoExec(videoPath, videoID string) error {
	cmd := exec.Command("yt-dlp", "--no-playlist", "--format", "mp4", "--output", videoPath, videoID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing yt-dlp: %v\n%s", err, output)
	}
	return nil
}

func downloadVideo(url string) (string, error) {
	videoPath := getUUIDVideoPath(YTPL_PATH)
	videoID := extractVideoID(url)
	if videoID == "" {
		return "", fmt.Errorf("invalid video URL")
	}
	if err := ytdlVideoExec(videoPath, videoID); err != nil {
		return "", err
	}
	return videoPath, nil
}

func getVideoURL(request *http.Request) string {
	if request.Method == http.MethodPost {
		return request.FormValue("url")
	}
	return ""
}

func cleanupDownload(videoPath string) {
	if err := os.Remove(videoPath); err != nil {
		log.Printf("Error cleaning up video file: %v\n", err)
	}
}

func streamVideo(responseWriter http.ResponseWriter, videoPath string) {
	videoFile, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Error opening video file: %v\n", err)
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer videoFile.Close()

	responseWriter.Header().Set("Content-Type", "video/mp4")
	_, err = io.Copy(responseWriter, videoFile)
	if err != nil {
		log.Printf("Error streaming video: %v\n", err)
	}
}

func YTPLPlayVideoHandler(responseWriter http.ResponseWriter, request *http.Request) {
	url := getVideoURL(request)
	if url == "" {
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	videoPath, err := downloadVideo(url)
	if err != nil {
		log.Printf("Error downloading video: %v\n", err)
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cleanupDownload(videoPath)

	streamVideo(responseWriter, videoPath)
}
