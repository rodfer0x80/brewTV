package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const LIBRARY_PATH = "/opt/brewTV/library"
const TV_PATH = "/opt/brewTV/library/tv"
const MUSIC_PATH = "/opt/brewTV/library/music"
const ACCEPTED_TV_FORMAT = "mp4"
const ACCEPTED_MUSIC_FORMAT = "mp3"
const VIDEO_PLAY_PARAMETER = "path"

func TVPlayHandler(responseWriter http.ResponseWriter, request *http.Request) {
	relative_video_path := getURLParameter(responseWriter, request, VIDEO_PLAY_PARAMETER)
	if relative_video_path == "" {
		return
	}
	playVideo(responseWriter, request, TV_PATH, relative_video_path)
}

func MusicPlayHandler(responseWriter http.ResponseWriter, request *http.Request) {
	relative_video_path := getURLParameter(responseWriter, request, VIDEO_PLAY_PARAMETER)
	if relative_video_path == "" {
		return
	}
	playVideo(responseWriter, request, MUSIC_PATH, relative_video_path)
}

func playVideo(responseWriter http.ResponseWriter, request *http.Request, dir_path string, relative_video_path string) {
	video_path := filepath.Join(dir_path, relative_video_path)
	if fileExists(video_path) {
		serveVideo(responseWriter, request, video_path, ACCEPTED_TV_FORMAT)
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
