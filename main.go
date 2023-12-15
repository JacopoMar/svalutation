package main

import (
	"api/entities"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/zitadel/oidc"
)

type Student = entities.Student
type Teacher = entities.Teacher
type Remark = entities.Remark
type Observation = entities.Observation

var DB, ERR = sql.Open("sqlite3", "./database.db")

func ErrorCheck(w *http.ResponseWriter, err error, code int) (bool){
	if err != nil {
		http.Error(*w, err.Error(), code)
		return true
	}
	return false
}

func HandleStudent(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			//Get all students
			rows, err := DB.Query("SELECT * FROM students")
			if ErrorCheck(&w, err, 500) {
				return
			}

			var students []Student
			for rows.Next(){
				var student Student
				rows.Scan(&student.Id, &student.Name, &student.Surname)
				students = append(students, student)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(students)
			return
		case "POST":
			//Create new student
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			result, err := DB.Exec("INSERT INTO students (name, surname) VALUES(?, ?)", r.Form.Get("name"), r.Form.Get("surname"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			id, err := result.LastInsertId()
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(id)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleStudentById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/students/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/students/")

	switch r.Method {
		case "GET":
			//Get existent student
			var student Student
			err := DB.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.Id, &student.Name, &student.Surname)
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(student)
			return
		case "PATCH":
			//Update existent student
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("UPDATE students SET name = ?, surname = ? WHERE id = ?", r.Form.Get("name"), r.Form.Get("surname"), id)
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		case "DELETE":
			//Delete existent student
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("DELETE FROM students WHERE id = ?", r.Form.Get("id"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleTeacher(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			//Get all teachers
			rows, err := DB.Query("SELECT * FROM teachers")
			if ErrorCheck(&w, err, 500) {
				return
			}

			var teachers []Teacher
			for rows.Next(){
				var teacher Teacher
				rows.Scan(&teacher.Id, &teacher.Name, &teacher.Surname)
				teachers = append(teachers, teacher)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(teachers)
			return
		case "POST":
			//Create new teacher
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			result, err := DB.Exec("INSERT INTO teachers (name, surname) VALUES(?, ?)", r.Form.Get("name"), r.Form.Get("surname"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			id, err := result.LastInsertId()
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(id)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleTeacherById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/teachers/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/teachers/")

	switch r.Method {
		case "GET":
			//Get existent teacher
			var teacher Teacher
			err := DB.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.Id, &teacher.Name, &teacher.Surname)
			if ErrorCheck(&w, err, 500) {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(teacher)
			return
		case "PATCH":
			//Update existent teacher
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("UPDATE teachers SET name = ?, surname = ? WHERE id = ?", r.Form.Get("name"), r.Form.Get("surname"), id)
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		case "DELETE":
			//Delete existent teacher
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("DELETE FROM teachers WHERE id = ?", r.Form.Get("id"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleRemark(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			//Get all remarks
			rows, err := DB.Query("SELECT * FROM remarks")
			if ErrorCheck(&w, err, 500) {
				return
			}

			var remarks []Remark
			for rows.Next(){
				var remark Remark
				rows.Scan(&remark.Id, &remark.Level, &remark.Description)
				remarks = append(remarks, remark)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(remarks)
			return
		case "POST":
			//Create new remark
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			result, err := DB.Exec("INSERT INTO remarks (level, description) VALUES(?, ?)", r.Form.Get("level"), r.Form.Get("description"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			id, err := result.LastInsertId()
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(id)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleRemarkById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/remarks/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/remarks/")
	
	switch r.Method {
		case "GET":
			//Get existent remark
			var remark Remark
			err := DB.QueryRow("SELECT * FROM remarks WHERE id = ?", id).Scan(&remark.Id, &remark.Level, &remark.Description)
			if ErrorCheck(&w, err, 500) {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(remark)
			return
		case "PATCH":
			//Update existent remark
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("UPDATE remarks SET level = ?, description = ? WHERE id = ?", r.Form.Get("level"), r.Form.Get("description"), id)
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		case "DELETE":
			//Delete existent remark
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("DELETE FROM remarks WHERE id = ?", r.Form.Get("id"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleObservation(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			//Get all observations
			rows, err := DB.Query("SELECT * FROM observations")
			if ErrorCheck(&w, err, 500) {
				return
			}

			var observations []Observation
			for rows.Next(){
				var observation Observation
				rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
				err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
				if ErrorCheck(&w, err, 500) {
					return
				}
				err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname)
				if ErrorCheck(&w, err, 500) {
					return
				}
				err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Level, &observation.Remark.Description)
				if ErrorCheck(&w, err, 500) {
					return
				}
				observations = append(observations, observation)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(observations)
			return
		case "POST":
			//Create new observation
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			result, err := DB.Exec("INSERT INTO observations (teacher, student, remark, achieved) VALUES(?, ?, ?, ?)", r.Form.Get("teacher"), r.Form.Get("student"), r.Form.Get("remark"), r.Form.Get("achieved"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			id, err := result.LastInsertId()
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(id)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleObservationById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/observations/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/observations/")
	
	switch r.Method {
		case "GET":
			//Get existent observation
			var observation Observation
			err := DB.QueryRow("SELECT * FROM observations WHERE id = ?", id).Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
			if ErrorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
			if ErrorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname)
			if ErrorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Level, &observation.Remark.Description)
			if ErrorCheck(&w, err, 500) {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(observation)
			return
		case "PATCH":
			//Update existent remark
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("UPDATE observations SET teacher = ?, student = ?, remark = ?, achieved = ? WHERE id = ?", r.Form.Get("teacher"), r.Form.Get("student"), r.Form.Get("remark"), r.Form.Get("achieved"), id)
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		case "DELETE":
			//Delete existent observation
			err := r.ParseForm()
			if ErrorCheck(&w, err, 400) {
				return
			}
			_, err = DB.Exec("DELETE FROM observations WHERE id = ?", r.Form.Get("id"))
			if ErrorCheck(&w, err, 500) {
				return
			}
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleObservationByStudentId(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/observations/student/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/observations/student/")
	
	switch r.Method {
		case "GET":
			//Get all observations made on student
			rows, err := DB.Query("SELECT * FROM observations where student = ?", id)
			if ErrorCheck(&w, err, 500) {
				return
			}

			var observations []Observation
			for rows.Next(){
				var observation Observation
				rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
				err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
				if ErrorCheck(&w, err, 500) {
					return
			  }
				err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname)
				if ErrorCheck(&w, err, 500) {
					return
				}
				err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Level, &observation.Remark.Description)
				if ErrorCheck(&w, err, 500) {
					return
				}
				observations = append(observations, observation)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(observations)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func HandleObservationByTeacherId(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/observations/teacher/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/observations/teacher/")
	
	switch r.Method {
		case "GET":
			//Get all observations made by the teacher
			rows, err := DB.Query("SELECT * FROM observations where teacher = ?", id)
			if ErrorCheck(&w, err, 500) {
				return
			}

			var observations []Observation
			for rows.Next(){
				var observation Observation
				rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
				err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
				if ErrorCheck(&w, err, 500) {
					return
				}
				err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname)
				if ErrorCheck(&w, err, 500) {
					return
				}
				err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Level, &observation.Remark.Description)
				if ErrorCheck(&w, err, 500) {
					return
				}
				observations = append(observations, observation)
			}
			rows.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(observations)
			return
		default:
			http.Error(w, "Bad Request", 400)
			return
	}
}

func main() {
	if ERR != nil	{
		log.Fatal(ERR)
	}
	defer DB.Close()

	//Teacher handlers
	http.HandleFunc("/api/teachers", HandleTeacher)
	http.HandleFunc("/api/teachers/", HandleTeacherById)
	
	//Student handlers
	http.HandleFunc("/api/students", HandleStudent)
	http.HandleFunc("/api/students/", HandleStudentById)
	
	//Remark handlers
	http.HandleFunc("/api/remarks", HandleRemark)
	http.HandleFunc("/api/remarks/", HandleRemarkById)
	
	//Observation handlers
	http.HandleFunc("/api/observations", HandleObservation)
	http.HandleFunc("/api/observations/", HandleObservationById)
	http.HandleFunc("/api/observations/student/", HandleObservationByStudentId)
  http.HandleFunc("/api/observations/teacher/", HandleObservationByTeacherId)
	http.ListenAndServe(":8080", nil)
}

//Note: request form only accepts content-type: application/x-www-form-urlencoded
