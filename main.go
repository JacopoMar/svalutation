package main

import (
	"api/entities"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Student = entities.Student
type Teacher = entities.Teacher
type Remark = entities.Remark
type Observation = entities.Observation
type Class = entities.Class

var DB, DB_ERR = sql.Open("sqlite3", "./database.db")
var IGNORE sql.Null[any]

func errorCheck(w *http.ResponseWriter, err error, code int) bool {
	if err != nil {
		http.Error(*w, err.Error(), code)
		return true
	}
	return false
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM students")
	if errorCheck(&w, err, 500) {
		return
	}

	var students []Student
	for rows.Next() {
		var student Student
		rows.Scan(&student.Id, &student.Name, &student.Surname, &student.Class.Id)
		err := DB.QueryRow("SELECT * FROM classes WHERE id = ?", student.Class.Id).Scan(&student.Class.Id, &student.Class.Name)
		if errorCheck(&w, err, 500) {
			return
		}

		students = append(students, student)
	}
	rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
	return
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	result, err := DB.Exec("INSERT INTO students (name, surname, class) VALUES(?, ?, ?)", r.Form.Get("name"), r.Form.Get("surname"), r.Form.Get("class"))
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
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var student Student
	err := DB.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.Id, &student.Name, &student.Surname, &student.Class.Id)
	if errorCheck(&w, err, 500) {
		return
	}
	err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", student.Class.Id).Scan(&student.Class.Id, &student.Class.Name)
	if errorCheck(&w, err, 500) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
	return
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}

	var form_name string = r.Form.Get("name")
	if form_name != "" {
		_, err = DB.Exec("UPDATE students SET name = ? WHERE id = ?", form_name, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_surname string = r.Form.Get("surname")
	if form_surname != "" {
		_, err = DB.Exec("UPDATE students SET surname = ? WHERE id = ?", form_surname, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_class string = r.Form.Get("class")
	if form_class != "" {
		_, err = DB.Exec("UPDATE students SET class = ? WHERE id = ?", form_class, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}
	return
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	_, err = DB.Exec("DELETE FROM students WHERE id = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	return
}

func getStudentsByClass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rows, err := DB.Query("SELECT * FROM students WHERE class = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}

	var students []Student
	for rows.Next() {
		var student Student
		rows.Scan(&student.Id, &student.Name, &student.Surname, &student.Class.Id)
		err := DB.QueryRow("SELECT * FROM classes where id = ?", student.Class.Id).Scan(&student.Class.Id, &student.Class.Name)
		if errorCheck(&w, err, 500) {
			return
		}

		students = append(students, student)
	}
	rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
	return
}

func getAllTeachers(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM teachers")
	if errorCheck(&w, err, 500) {
		return
	}

	var teachers []Teacher
	for rows.Next() {
		var teacher Teacher
		rows.Scan(&teacher.Id, &teacher.Name, &teacher.Surname)

		rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", teacher.Id)
		if errorCheck(&w, err, 500) {
			return
		}

		var classes []Class
		for rows.Next() {
			var class Class
			rows.Scan(&IGNORE, &IGNORE, &class.Id)

			err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
			if errorCheck(&w, err, 500) {
				return
			}

			classes = append(classes, class)
		}
		rows.Close()
		teacher.Classes = classes

		teachers = append(teachers, teacher)
	}
	rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teachers)
	return
}

func createTeacher(w http.ResponseWriter, r *http.Request) {
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

	var classIds []int64
	err = json.Unmarshal([]byte(r.Form.Get("classes")), &classIds)
	if errorCheck(&w, err, 500) {
		return
	}

	for _, element := range classIds {
		_, err := DB.Exec("INSERT INTO classes_teachers (teacher_id, class_id) VALUES(?, ?)", id, element)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
	return
}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var teacher Teacher
	err := DB.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.Id, &teacher.Name, &teacher.Surname)
	if errorCheck(&w, err, 500) {
		return
	}

	rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", teacher.Id)
	if errorCheck(&w, err, 500) {
		return
	}
	var classes []Class
	for rows.Next() {
		var class Class
		rows.Scan(&IGNORE, &IGNORE, &class.Id)
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
		if errorCheck(&w, err, 500) {
			return
		}
		classes = append(classes, class)
	}
	rows.Close()
	teacher.Classes = classes

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher)
	return
}

func updateTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}

	var form_name string = r.Form.Get("name")
	if form_name != "" {
		_, err = DB.Exec("UPDATE teachers SET name = ? WHERE id = ?", form_name, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_surname string = r.Form.Get("surname")
	if form_surname != "" {
		_, err = DB.Exec("UPDATE teachers SET surname = ? WHERE id = ?", form_surname, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_classes string = r.Form.Get("classes")
	if form_classes != "" && form_classes != "[]" {
		_, err := DB.Exec("DELETE FROM classes_teachers WHERE teacher_id = ?", id)
		if errorCheck(&w, err, 500) {
			return
		}

		var classIds []int64
		err = json.Unmarshal([]byte(form_classes), &classIds)
		if errorCheck(&w, err, 500) {
			return
		}

		for _, element := range classIds {
			_, err := DB.Exec("INSERT INTO classes_teachers (teacher_id, class_id) VALUES(?, ?)", id, element)
			if errorCheck(&w, err, 500) {
				return
			}
		}
	}
	return
}

func deleteTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	_, err = DB.Exec("DELETE FROM teachers WHERE id = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	_, err = DB.Exec("DELETE FROM classes_teachers WHERE teacher_id = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}

	return
}

func getAllRemarks(w http.ResponseWriter, r *http.Request) {
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
}

func createRemark(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	result, err := DB.Exec("INSERT INTO remarks (skill, level, description) VALUES(?, ?, ?)", r.Form.Get("skill"), r.Form.Get("level"), r.Form.Get("description"))
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
}

func getRemark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var remark Remark
	err := DB.QueryRow("SELECT * FROM remarks WHERE id = ?", id).Scan(&remark.Id, &remark.Skill, &remark.Level, &remark.Description)
	if errorCheck(&w, err, 500) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(remark)
	return
}

func updateRemark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}

	var form_skill string = r.Form.Get("skill")
	if form_skill != "" {
		_, err = DB.Exec("UPDATE remarks SET skill = ? WHERE id = ?", form_skill, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_level string = r.Form.Get("level")
	if form_level != "" {
		_, err = DB.Exec("UPDATE remarks SET level = ? WHERE id = ?", form_level, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_description string = r.Form.Get("description")
	if form_description != "" {
		_, err = DB.Exec("UPDATE remarks SET description = ? WHERE id = ?", form_description, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}
	return
}

func deleteRemark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	_, err = DB.Exec("DELETE FROM remarks WHERE id = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	return
}

func getAllObservations(w http.ResponseWriter, r *http.Request) {
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
		rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", observation.Teacher.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		var classes []Class
		for rows.Next() {
			var class Class
			rows.Scan(&IGNORE, &IGNORE, &class.Id)

			err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
			if errorCheck(&w, err, 500) {
				return
			}
			classes = append(classes, class)
		}
		rows.Close()
		observation.Teacher.Classes = classes

		err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", observation.Student.Class.Id).Scan(&observation.Student.Class.Id, &observation.Student.Class.Name)
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
}

func createObservation(w http.ResponseWriter, r *http.Request) {
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
}

func getObservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var observation Observation
	err := DB.QueryRow("SELECT * FROM observations WHERE id = ?", id).Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
	if errorCheck(&w, err, 500) {
		return
	}
	rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", observation.Teacher.Id)
	if errorCheck(&w, err, 500) {
		return
	}
	var classes []Class
	for rows.Next() {
		var class Class
		rows.Scan(&IGNORE, &IGNORE, &class.Id)
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
		if errorCheck(&w, err, 500) {
			return
		}
		classes = append(classes, class)
	}
	rows.Close()
	observation.Teacher.Classes = classes

	err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
	if errorCheck(&w, err, 500) {
		return
	}
	err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class.Id)
	if errorCheck(&w, err, 500) {
		return
	}
	err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", observation.Student.Class.Id).Scan(&observation.Student.Class.Id, &observation.Student.Class.Name)
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
}

func updateObservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}

	var form_teacher string = r.Form.Get("teacher")
	if form_teacher != "" {
		_, err = DB.Exec("UPDATE observations SET teacher = ? WHERE id = ?", form_teacher, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_student string = r.Form.Get("student")
	if form_student != "" {
		_, err = DB.Exec("UPDATE observations SET student = ? WHERE id = ?", form_student, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_remark string = r.Form.Get("remark")
	if form_remark != "" {
		_, err = DB.Exec("UPDATE observations SET remark = ? WHERE id = ?", form_remark, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}

	var form_achieved string = r.Form.Get("achieved")
	if form_achieved != "" {
		_, err = DB.Exec("UPDATE observations SET achieved = ? WHERE id = ?", form_achieved, id)
		if errorCheck(&w, err, 500) {
			return
		}
	}
	return
}

func deleteObservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if errorCheck(&w, err, 400) {
		return
	}
	_, err = DB.Exec("DELETE FROM observations WHERE id = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	return
}

func getObservationsOnStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rows, err := DB.Query("SELECT * FROM observations where student = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	var observations []Observation
	for rows.Next() {
		var observation Observation
		rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
		rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", observation.Teacher.Id)
		if errorCheck(&w, err, 500) {
			return
		}

		var classes []Class
		for rows.Next() {
			var class Class
			rows.Scan(&IGNORE, &IGNORE, &class.Id)
			err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
			if errorCheck(&w, err, 500) {
				return
			}
			classes = append(classes, class)
		}
		rows.Close()
		observation.Teacher.Classes = classes

		err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", observation.Student.Class.Id).Scan(&observation.Student.Class.Id, &observation.Student.Class.Name)
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
}

func getObservationsByTeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rows, err := DB.Query("SELECT * FROM observations where teacher = ?", id)
	if errorCheck(&w, err, 500) {
		return
	}
	var observations []Observation
	for rows.Next() {
		var observation Observation
		rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
		rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", observation.Teacher.Id)
		if errorCheck(&w, err, 500) {
			return
		}

		var classes []Class
		for rows.Next() {
			var class Class
			rows.Scan(&IGNORE, &IGNORE, &class.Id)
			err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
			if errorCheck(&w, err, 500) {
				return
			}
			classes = append(classes, class)
		}
		rows.Close()
		observation.Teacher.Classes = classes

		err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", observation.Student.Class.Id).Scan(&observation.Student.Class.Id, &observation.Student.Class.Name)
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
}

func getObservationsByTeacherOnStudent(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("teacherId")
	studentId := r.PathValue("studentId")

	rows, err := DB.Query("SELECT * FROM observations where teacher = ? and student = ?", teacherId, studentId)
	if errorCheck(&w, err, 500) {
		return
	}

	var observations []Observation
	for rows.Next() {
		var observation Observation
		rows.Scan(&observation.Id, &observation.Teacher.Id, &observation.Student.Id, &observation.Remark.Id, &observation.Achieved, &observation.Date)
		rows, err := DB.Query("SELECT * FROM classes_teachers WHERE teacher_id = ?", observation.Teacher.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		var classes []Class
		for rows.Next() {
			var class Class
			rows.Scan(&IGNORE, &IGNORE, &class.Id)
			err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", class.Id).Scan(&class.Id, &class.Name)
			if errorCheck(&w, err, 500) {
				return
			}
			classes = append(classes, class)
		}
		rows.Close()
		observation.Teacher.Classes = classes

		err = DB.QueryRow("SELECT * FROM teachers WHERE id = ?", observation.Teacher.Id).Scan(&observation.Teacher.Id, &observation.Teacher.Name, &observation.Teacher.Surname)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM students WHERE id = ?", observation.Student.Id).Scan(&observation.Student.Id, &observation.Student.Name, &observation.Student.Surname, &observation.Student.Class.Id)
		if errorCheck(&w, err, 500) {
			return
		}
		err = DB.QueryRow("SELECT * FROM classes WHERE id = ?", observation.Student.Class.Id).Scan(&observation.Student.Class.Id, &observation.Student.Class.Name)
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
}

func enableCors(w *http.ResponseWriter) {
<<<<<<< HEAD
	(*w).Header().Set("Access-Control-Allow-Origin", "http://85.235.150.118:3000")
=======
	(*w).Header().Set("Access-Control-Allow-Origin", "85.235.150.118:3000")
>>>>>>> 43e1b5de54f0aebd2391a1d19c4ad62d3add2374
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		enableCors(&w)

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
		slog.Error("Couldn't retrieve credentials") // TODO: Manage this kind of error, it could mean the username the user provided is wrong
	}
	return bcrypt.CompareHashAndPassword([]byte(credentials.password), []byte(password)) == nil
}

func statusCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	if DB_ERR != nil {
		log.Fatal(DB_ERR)
	}
	defer DB.Close()

	slog.Info("Loaded database")

	mux := http.NewServeMux()

	// CORS OPTIONS handler
	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) { enableCors(&w); w.WriteHeader(http.StatusOK) })

	// Status handler
	mux.HandleFunc("GET /status", statusCheck)

	// Student handlers
	mux.HandleFunc("GET /api/students", auth(getAllStudents))
	mux.HandleFunc("POST /api/students", auth(createStudent))

	mux.HandleFunc("GET /api/students/{id}", auth(getStudent))
	mux.HandleFunc("PATCH /api/students/{id}", auth(updateStudent))
	mux.HandleFunc("DELETE /api/students/{id}", auth(deleteStudent))

	mux.HandleFunc("GET /api/students/class/{id}", auth(getStudentsByClass))

	// Teacher handlers
	mux.HandleFunc("GET /api/teachers", auth(getAllTeachers))
	mux.HandleFunc("POST /api/teachers", auth(createTeacher))

	mux.HandleFunc("GET /api/teachers/{id}", auth(getTeacher))
	mux.HandleFunc("PATCH /api/teachers/{id}", auth(updateTeacher))
	mux.HandleFunc("DELETE /api/teachers/{id}", auth(deleteTeacher))

	// Remark handlers
	mux.HandleFunc("GET /api/remarks", auth(getAllRemarks))
	mux.HandleFunc("POST /api/remarks", auth(createRemark))

	mux.HandleFunc("GET /api/remarks/{id}", auth(getRemark))
	mux.HandleFunc("PATCH /api/remarks/{id}", auth(updateRemark))
	mux.HandleFunc("DELETE /api/remarks/{id}", auth(deleteRemark))

	// Observation handlers
	mux.HandleFunc("GET /api/observations", auth(getAllObservations))
	mux.HandleFunc("POST /api/observations", auth(createObservation))

	mux.HandleFunc("GET /api/observations/{id}", auth(getObservation))
	mux.HandleFunc("PATCH /api/observations/{id}", auth(updateObservation))
	mux.HandleFunc("DELETE /api/observations/{id}", auth(deleteObservation))

	mux.HandleFunc("GET /api/observations/student/{id}", auth(getObservationsOnStudent))
	mux.HandleFunc("GET /api/observations/teacher/{id}", auth(getObservationsByTeacher))
	mux.HandleFunc("GET /api/observations/teacher/{teacherId}/student/{studentId}", auth(getObservationsByTeacherOnStudent))

	slog.Info("Starting server")
	http.ListenAndServe(":8080", mux)
}

// Note: request form only accepts content-type: application/x-www-form-urlencoded
