package example

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartSimpleServe(wg *sync.WaitGroup) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)

	mux := http.NewServeMux()

	uploadDir := "./upload"

	mux.Handle("/", &HandlerRoot{})
	mux.Handle("/upload", &HandlerUpload{uploadDir})
	mux.Handle("/upload-list", &HandlerUploadList{uploadDir})

	srv := http.Server{
		Addr:         "localhost:9000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	go func() {
		log.Println("Starting server ...")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	close(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutdown", err)
	}

	log.Println("Server shutdown")
	wg.Done()
}

func isDirExist(path string) bool {
	dir, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !dir.IsDir() {
		return false
	}

	return true
}

func CreateDir(path string) error {
	return os.Mkdir(path, 0777)
}
