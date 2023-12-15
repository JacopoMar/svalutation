package entities

import (
	"time"
)

type Observation struct {
	Id       int64
	Teacher  Teacher
	Student  Student
	Remark   Remark
	Achieved bool
	Date     time.Time
}