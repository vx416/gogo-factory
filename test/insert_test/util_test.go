package insert_test

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

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
