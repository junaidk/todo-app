package main

import (
	"fmt"
	"github.com/swaggo/http-swagger"
	"io"
	"log"
	"net/http"
	"os"
	"todo/api"
	"todo/config"
	_ "todo/docs"
)

func init() {

	config.DataStorePath = "/tmp/test.json"
	config.CacheDirPath = "/tmp/Cache"
	config.LogFilePath = "/tmp/logFile.log"

	createDirectory(config.CacheDirPath)
}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	api.Init()

	openLogFile(config.LogFilePath)
	routes := api.NewRouter()

	routes.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":8080", logRequest(routes)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		log.Printf("%s | %s | %d | %s \n", r.RemoteAddr, r.Method, sw.status, r.URL)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		multi := io.MultiWriter(lf, os.Stdout)
		log.SetOutput(multi)
	}
}

// https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	return n, err
}

func createDirectory(dirName string) bool {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}

	if src.Mode().IsRegular() {
		fmt.Println(dirName, "already exist as a file!")
		return false
	}

	return false
}
