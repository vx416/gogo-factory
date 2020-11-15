package test

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

func gormSqlite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
}

func AllEmployees(db *sql.DB, driver string) ([]*Employee, map[int64]*Employee, error) {
	data := make([]*Employee, 0, 1)
	xDB := sqlx.NewDb(db, driver)
	err := xDB.Select(&data, "select * from employees order by id")
	if err != nil {
		return nil, nil, err
	}
	dataMap := make(map[int64]*Employee)
	for i := range data {
		dataMap[data[i].ID] = data[i]
	}

	return data, dataMap, nil
}

func AllProjects(db *sql.DB, driver string) ([]*Project, error) {
	xDB := sqlx.NewDb(db, driver)
	data := make([]*Project, 0, 1)
	err := xDB.Select(&data, "select * from projects order by id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AllEmployeesProjects(db *sql.DB, driver string) ([]*EmployeesProjects, error) {
	xDB := sqlx.NewDb(db, driver)
	data := make([]*EmployeesProjects, 0, 1)
	err := xDB.Select(&data, "select * from employees_projects order by project_id, employee_id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AllTasks(db *sql.DB, driver string) ([]*Task, error) {
	xDB := sqlx.NewDb(db, driver)
	data := make([]*Task, 0, 1)
	err := xDB.Select(&data, "select * from tasks order by id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AllDomains(db *sql.DB, driver string) ([]*Domain, error) {
	xDB := sqlx.NewDb(db, driver)
	data := make([]*Domain, 0, 1)
	err := xDB.Select(&data, "select * from domains order by id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AllSpecialties(db *sql.DB, driver string) ([]*Specialty, error) {
	xDB := sqlx.NewDb(db, driver)
	data := make([]*Specialty, 0, 1)
	err := xDB.Select(&data, "select * from specialties order by id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Clear(db *sql.DB) error {
	var err error
	tables := []string{"employees", "projects", "tasks", "domains", "specialties", "employees_projects"}
	for _, table := range tables {
		_, err = db.Exec("DELETE FROM " + table)
	}
	return err
}

func readSchema(fileName string) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	bytes, err := ioutil.ReadFile(filepath.Join(dir, "./schema/"+fileName))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func initSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	schema, err := readSchema("sqlite.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// func initPGDB() *sql.DB {
// 	return
// }

// func initMySQLDB() *sql.DB {
// 	return
// }
