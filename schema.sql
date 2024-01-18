CREATE table IF NOT EXISTS "classes" (
  "id" INTEGER NOT NULL UNIQUE,
  "name" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE table IF NOT EXISTS "classes_teachers" (
  "id" INTEGER NOT NULL UNIQUE,
  "teacher_id" INTEGER NOT NULL,
  "class_id" INTEGER NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT),
  FOREIGN KEY("teacher_id") REFERENCES "teachers"("id"),
  FOREIGN KEY("class_id") REFERENCES "classes"("id")
);

CREATE TABLE IF NOT EXISTS "credentials" (
  "user" TEXT NOT NULL UNIQUE,
  "password" TEXT NOT NULL,
  PRIMARY KEY("user")
);

CREATE table IF NOT EXISTS "observations" (
  "id" INTEGER NOT NULL UNIQUE,
  "teacher" INTEGER NOT NULL,
  "student" INTEGER NOT NULL,
  "remark" INTEGER NOT NULL,
  "achieved" INTEGER NOT NULL,
  "date" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY("id" AUTOINCREMENT),
  FOREIGN KEY("student") REFERENCES "students",
  FOREIGN KEY("teacher") REFERENCES "teachers",
  FOREIGN KEY("remark") REFERENCES "remarks"
);

CREATE table IF NOT EXISTS "remarks" (
  "id" INTEGER NOT NULL UNIQUE,
  "skill" INTEGER NOT NULL,
  "level" INTEGER NOT NULL,
  "description" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE table IF NOT EXISTS "students" (
  "id" INTEGER NOT NULL UNIQUE,
  "name" TEXT NOT NULL,
  "surname" TEXT NOT NULL,
  "class" INTEGER NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE table IF NOT EXISTS "teachers" (
  "id" INTEGER NOT NULL UNIQUE,
  "name" TEXT NOT NULL,
  "surname" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);