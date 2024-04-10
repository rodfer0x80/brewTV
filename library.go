package main

import (
	"fmt"
	"io"
  "log"
  "compress/gzip"
  "net/http"
	"os"
  "time"
  "strconv"
	"path/filepath"
	"strings"
)

const LIBRARY_PATH = "/opt/brewTV/library"
const TV_PATH = "/opt/brewTV/library/tv"
const MUSIC_PATH = "/opt/brewTV/library/music"
const ACCEPTED_TV_FORMAT = "mp4"
const ACCEPTED_MUSIC_FORMAT = "mp3"
const VIDEO_PLAY_PARAMETER = "path"
const CACHE_MAX_AGE = 360000 // 100h

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
    videoHandle, err := os.Open(file_path)
    if err != nil {
        http.Error(w, "Failed to open the video file", http.StatusInternalServerError)
        log.Fatal("Failed to open video file")
        return
    }
    defer videoHandle.Close()

    fileStat, err := videoHandle.Stat()
    if err != nil {
        http.Error(w, "Failed to get video file information", http.StatusInternalServerError)
        log.Fatal("Failed to get video file information")
        return
    }

    w.Header().Set("Content-Type", "video/"+file_format)

    w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(CACHE_MAX_AGE))
    w.Header().Set("Expires", time.Now().Add(time.Duration(CACHE_MAX_AGE)*time.Second).Format(http.TimeFormat))

    w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

    _, err = io.Copy(w, videoHandle)
    if err != nil {
        http.Error(w, "Failed to stream the video content", http.StatusInternalServerError)
        log.Fatal("Failed to stream video content")
        return
    }
}


func serveCompressedContent(w http.ResponseWriter, r *http.Request, content []byte, modTime time.Time, content_type string) {
    w.Header().Set("Content-Type", "video/"+content_type)
    w.Header().Set("Cache-Control", "max-age=3600") // Cache for 1 hour

    // Check if client supports gzip
    if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()
        gz.Write(content)
    } else {
        // If client doesn't support gzip, serve uncompressed content
        w.Write(content)
    }

    w.WriteHeader(http.StatusOK)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
