package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go-students-api/internal/student"
)

type Handler struct {
	Service *student.Service
}

func NewHandler(s *student.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// CreateStudent
func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student student.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	student.CreatedBy = userID
	student.UpdatedBy = userID

	log.Printf("Creating student: %+v", student)

	if err := h.Service.CreateStudent(&student); err != nil {
		log.Printf("Error creating student: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

// GetStudent By ID
func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/students/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

// UpdateStudent
func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var student student.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	student.UpdatedBy = userID

	if err := h.Service.UpdateStudent(&student); err != nil {
		log.Printf("Error updating student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

// DeleteStudent By ID
func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/students/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteStudent(id); err != nil {
		log.Printf("Error deleting student: %v", err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllStudents
func (h *Handler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.Service.GetAllStudents()
	if err != nil {
		log.Printf("Error retrieving students: %v", err)
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
}
