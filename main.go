package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zeyd17/file-microservice/handler/api"
	"github.com/zeyd17/file-microservice/repository"
)

func main() {

	db, err := gorm.Open("sqlite3", "file.db")
	if err != nil {
		fmt.Println(err)
		panic("Db Err ")

	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "handler/swagger")
	fileServer(r, "/swagger", http.Dir(filesDir))

	fileApi := api.NewFileApi(repository.NewFileRepo(db))

	r.Route("/", func(rt chi.Router) {
		rt.Mount("/file", fileRouter(fileApi))
	})

	fmt.Println("Server listen at :8080")
	err = http.ListenAndServe(":8080", r)
	fmt.Println(err)
}

func fileRouter(fileApi *api.FileApi) http.Handler {
	r := chi.NewRouter()

	r.Get("/{id:[0-9-a-f-]+}", fileApi.Get)
	r.Post("/", fileApi.Post)
	r.Delete("/{id:[0-9-a-f-]+}", fileApi.Delete)
	r.Get("/download/{id:[0-9-a-f-]+}", fileApi.Download)

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
