// backend/main.go
package main

import (
    "fmt"
    "log"
    "mime"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "github.com/gorilla/handlers"
)


func main() {
    mux := http.NewServeMux()

    // Serve video and subtitles from API
    mux.HandleFunc("/video", videoHandler)
    mux.HandleFunc("/sub", subtitleHandler)
    mux.HandleFunc("/favicon.ico", faviconHandler)

    // Get the current working directory
    cwd, err1 := os.Getwd()
    if err1 != nil {
        log.Fatal(err1)
    }

    log.Println("Current working directory:", cwd)
    // Build the full path to the "build" directory
    buildPath := filepath.Join(cwd, "../frontend/build")

    // Serve static files from the build folder
    fs := http.FileServer(http.Dir(buildPath))
    mux.Handle("/", fs)

    // Start the server
    log.Println("Serving files from", buildPath)

    // CORS options
    headers := handlers.AllowedHeaders([]string{"Content-Type"})
    methods := handlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})
    origins := handlers.AllowedOrigins([]string{"*"}) // Allow all origins for testing

    fmt.Println("Server started at :80")
    log.Fatal(http.ListenAndServe(":80", handlers.CORS(origins, headers, methods)(mux)))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./favicon.ico") // Adjust path as necessary
}

func subtitleHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request for subtitles from: %s", r.RemoteAddr)

    subtitlePath := "video/The.Greatest.Hits.2024.vtt" // Change this to your subtitle file path
    w.Header().Set("Content-Type", "text/vtt")
    http.ServeFile(w, r, subtitlePath)

    log.Println("Completed serving subtitles")
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
    // Path to the video file
    videoPath := "video/The.Greatest.Hits.2024.mkv"

    log.Printf("Attempting to stream video: %s", videoPath)

    // Open the file
    file, err := os.Open(videoPath)
    if err != nil {
        http.Error(w, "Video not found", http.StatusNotFound)
        log.Printf("Error opening file: %v", err)
        return
    }
    defer file.Close()

    // Get file info
    fi, err := file.Stat()
    if err != nil {
        http.Error(w, "Cannot obtain file info", http.StatusInternalServerError)
        log.Printf("Error getting file info: %v", err)
        return
    }

    log.Printf("Streaming video: %s", videoPath)

    // Set headers
    size := fi.Size()
    rangeHeader := r.Header.Get("Range")

    log.Printf("Range header: %s", rangeHeader)

    // Set Content-Type
    ext := filepath.Ext(videoPath)
    mimeType := mime.TypeByExtension(ext)
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }
    w.Header().Set("Content-Type", mimeType)
    w.Header().Set("Accept-Ranges", "bytes")

    log.Printf("Content-Type: %s", mimeType)

    if rangeHeader == "" {
        w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
        w.WriteHeader(http.StatusOK)
        // Serve the entire file
        http.ServeContent(w, r, videoPath, fi.ModTime(), file)
        return
    }

    // Parse Range header
    var start, end int64
    fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)

    if end == 0 || end >= size {
        end = size - 1
    }

    // Set headers for partial content
    w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
    w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
    w.WriteHeader(http.StatusPartialContent)

    // Seek to the start position
    file.Seek(start, 0)

    // Serve the requested range
    http.ServeContent(w, r, videoPath, fi.ModTime(), file)
}
