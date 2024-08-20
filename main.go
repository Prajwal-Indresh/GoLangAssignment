package main

import (
	"log"
	"net/http"

	"go-students-api/internal/config"
	"go-students-api/internal/database"
	"go-students-api/internal/student"
	httptransport "go-students-api/internal/transport/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	studService := student.NewService(database.GetDB())
	handler := httptransport.NewHandler(studService)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		httptransport.Login(w, r, []byte(cfg.JWTSecret))
	})
	http.HandleFunc("/students", httptransport.JWTAuthMiddleware(handler.GetAllStudents, []byte(cfg.JWTSecret)))
	http.HandleFunc("/students/", httptransport.JWTAuthMiddleware(handler.GetStudentByID, []byte(cfg.JWTSecret)))
	http.HandleFunc("/students/update", httptransport.JWTAuthMiddleware(handler.UpdateStudent, []byte(cfg.JWTSecret)))
	http.HandleFunc("/students/delete/", httptransport.JWTAuthMiddleware(handler.DeleteStudent, []byte(cfg.JWTSecret)))
	http.HandleFunc("/students/create", httptransport.JWTAuthMiddleware(handler.CreateStudent, []byte(cfg.JWTSecret)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
