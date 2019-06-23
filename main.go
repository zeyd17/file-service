package main

import (
	"fmt"
	"net/http"

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
