package main

import (
	"api/entities"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Student = entities.Student
type Teacher = entities.Teacher
type Remark = entities.Remark
type Observation = entities.Observation

var DB, ERR = sql.Open("sqlite3", "./database.db")

func errorCheck(w *http.ResponseWriter, err error, code int) bool {
	if err != nil {
		http.Error(*w, err.Error(), code)
		return true
	}
	return false
}

func handleStudent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Get all students
		rows, err := DB.Query("SELECT * FROM students")
		if errorCheck(&w, err, 500) {
			return
		}

		var students []Student
		for rows.Next() {
			var student Student
			rows.Scan(&student.Id, &student.Name, &student.Surname, &student.Class)
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
		if errorCheck(&w, err, 400) {
			return
		}
		result, err := DB.Exec("INSERT INTO students (name, surname) VALUES(?, ?)", r.Form.Get("name"), r.Form.Get("surname"))
		if errorCheck(&w, err, 500) {
			return
		}
		id, err := result.LastInsertId()
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleStudentById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/students/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/students/")

	switch r.Method {
	case "GET":
		//Get existent student
		var student Student
		err := DB.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.Id, &student.Name, &student.Surname, &student.Class)
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(student)
		return
	case "PATCH":
		//Update existent student
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("UPDATE students SET name = ?, surname = ? WHERE id = ?", r.Form.Get("name"), r.Form.Get("surname"), id)
		if errorCheck(&w, err, 500) {
			return
		}
		return
	case "DELETE":
		//Delete existent student
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("DELETE FROM students WHERE id = ?", r.Form.Get("id"))
		if errorCheck(&w, err, 500) {
			return
		}
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleTeacher(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Get all teachers
		rows, err := DB.Query("SELECT * FROM teachers")
		if errorCheck(&w, err, 500) {
			return
		}

		var teachers []Teacher
		for rows.Next() {
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
		if errorCheck(&w, err, 400) {
			return
		}
		result, err := DB.Exec("INSERT INTO teachers (name, surname) VALUES(?, ?)", r.Form.Get("name"), r.Form.Get("surname"))
		if errorCheck(&w, err, 500) {
			return
		}
		id, err := result.LastInsertId()
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleTeacherById(w http.ResponseWriter, r *http.Request) {
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
		if errorCheck(&w, err, 500) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(teacher)
		return
	case "PATCH":
		//Update existent teacher
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("UPDATE teachers SET name = ?, surname = ? WHERE id = ?", r.Form.Get("name"), r.Form.Get("surname"), id)
		if errorCheck(&w, err, 500) {
			return
		}
		return
	case "DELETE":
		//Delete existent teacher
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("DELETE FROM teachers WHERE id = ?", r.Form.Get("id"))
		if errorCheck(&w, err, 500) {
			return
		}
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleRemark(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Get all remarks
		rows, err := DB.Query("SELECT * FROM remarks")
		if errorCheck(&w, err, 500) {
			return
		}

		var remarks []Remark
		for rows.Next() {
			var remark Remark
			rows.Scan(&remark.Id, &remark.Skill, &remark.Level, &remark.Description)
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
		if errorCheck(&w, err, 400) {
			return
		}
		result, err := DB.Exec("INSERT INTO remarks (level, description) VALUES(?, ?)", r.Form.Get("level"), r.Form.Get("description"))
		if errorCheck(&w, err, 500) {
			return
		}
		id, err := result.LastInsertId()
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleRemarkById(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/remarks/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/remarks/")

	switch r.Method {
	case "GET":
		//Get existent remark
		var remark Remark
		err := DB.QueryRow("SELECT * FROM remarks WHERE id = ?", id).Scan(&remark.Id, &remark.Skill, &remark.Level, &remark.Description)
		if errorCheck(&w, err, 500) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(remark)
		return
	case "PATCH":
		//Update existent remark
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("UPDATE remarks SET level = ?, description = ? WHERE id = ?", r.Form.Get("level"), r.Form.Get("description"), id)
		if errorCheck(&w, err, 500) {
			return
		}
		return
	case "DELETE":
		//Delete existent remark
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("DELETE FROM remarks WHERE id = ?", r.Form.Get("id"))
		if errorCheck(&w, err, 500) {
			return
		}
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleObservation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Get all observations
		rows, err := DB.Query("SELECT * FROM observations")
		if errorCheck(&w, err, 500) {
			return
		}

		var observations []Observation
		for rows.Next() {
			var observation Observation
			rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
			err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Skill, &observation.Remark.Level, &observation.Remark.Description)
			if errorCheck(&w, err, 500) {
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
		if errorCheck(&w, err, 400) {
			return
		}
		result, err := DB.Exec("INSERT INTO observations (teacher, student, remark, achieved) VALUES(?, ?, ?, ?)", r.Form.Get("teacher"), r.Form.Get("student"), r.Form.Get("remark"), r.Form.Get("achieved"))
		if errorCheck(&w, err, 500) {
			return
		}
		id, err := result.LastInsertId()
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleObservationById(w http.ResponseWriter, r *http.Request) {
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
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Skill, &observation.Remark.Level, &observation.Remark.Description)
		if errorCheck(&w, err, 500) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(observation)
		return
	case "PATCH":
		//Update existent observation
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("UPDATE observations SET teacher = ?, student = ?, remark = ?, achieved = ? WHERE id = ?", r.Form.Get("teacher"), r.Form.Get("student"), r.Form.Get("remark"), r.Form.Get("achieved"), id)
		if errorCheck(&w, err, 500) {
			return
		}
		return
	case "DELETE":
		//Delete existent observation
		err := r.ParseForm()
		if errorCheck(&w, err, 400) {
			return
		}
		_, err = DB.Exec("DELETE FROM observations WHERE id = ?", r.Form.Get("id"))
		if errorCheck(&w, err, 500) {
			return
		}
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleObservationByStudentId(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/observations/student/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/observations/student/")

	switch r.Method {
	case "GET":
		//Get all observations made on student
		rows, err := DB.Query("SELECT * FROM observations where student = ?", id)
		if errorCheck(&w, err, 500) {
			return
		}

		var observations []Observation
		for rows.Next() {
			var observation Observation
			rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
			err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Skill, &observation.Remark.Level, &observation.Remark.Description)
			if errorCheck(&w, err, 500) {
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleObservationByTeacherId(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/observations/teacher/" {
		http.Error(w, "No id specified", 400)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/observations/teacher/")

	switch r.Method {
	case "GET":
		//Get all observations made by the teacher
		rows, err := DB.Query("SELECT * FROM observations where teacher = ?", id)
		if errorCheck(&w, err, 500) {
			return
		}

		var observations []Observation
		for rows.Next() {
			var observation Observation
			rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
			err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class)
			if errorCheck(&w, err, 500) {
				return
			}
			err = DB.QueryRow("SELECT * FROM remarks WHERE id = ?", observation.Remark.Id).Scan(&observation.Remark.Id, &observation.Remark.Skill, &observation.Remark.Level, &observation.Remark.Description)
			if errorCheck(&w, err, 500) {
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if !checkCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Svalutation\"")
			http.Error(w, "Authentication failed, you shall not pass", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func checkCredentials(user string, password string) bool {
	type Credentials struct {
		user     string
		password string
	}
	credentials := Credentials{}

	err := DB.QueryRow("SELECT * FROM credentials WHERE user = ?", user).Scan(&credentials.user, &credentials.password)
	if err != nil {
		log.Println("Couldn't retrieve credentials") // TODO: Manage this kind of error, it could mean the username the user provided is wrong
	}
	return bcrypt.CompareHashAndPassword([]byte(credentials.password), []byte(password)) == nil
}

func main() {
	if ERR != nil {
		log.Fatal(ERR)
	}
	defer DB.Close()

	//Teacher handlers
	http.HandleFunc("/api/teachers", auth(handleTeacher))
	http.HandleFunc("/api/teachers/", auth(handleTeacherById))

	//Student handlers
	http.HandleFunc("/api/students", auth(handleStudent))
	http.HandleFunc("/api/students/", auth(handleStudentById))

	//Remark handlers
	http.HandleFunc("/api/remarks", auth(handleRemark))
	http.HandleFunc("/api/remarks/", auth(handleRemarkById))

	//Observation handlers
	http.HandleFunc("/api/observations", auth(handleObservation))
	http.HandleFunc("/api/observations/", auth(handleObservationById))
	http.HandleFunc("/api/observations/student/", auth(handleObservationByStudentId))
	http.HandleFunc("/api/observations/teacher/", auth(handleObservationByTeacherId))
	http.ListenAndServe(":8080", nil)
}

//Note: request form only accepts content-type: application/x-www-form-urlencoded
