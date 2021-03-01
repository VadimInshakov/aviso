package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"

	"time"
)

type DB struct {
	Instance *sql.DB
}

type QueryResult struct {
	Theme string    `json:"theme"`
	Link  string    `json:"link"`
	Site  string    `json:"site"`
	Time  time.Time `json:"time"`
}

func New(dbname string) (*DB, error) {
	dbinstance, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	err = dbinstance.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{Instance: dbinstance}, err
}

func (db *DB) Init() error {
	_, err := db.Instance.Exec(`
        CREATE TABLE news (
    		theme TEXT PRIMARY KEY,
    		link TEXT,
    	    site TEXT,
    	    time TIMESTAMP 
    	)`)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Insert(theme string, link string, site string, t time.Time) error {
	query := `INSERT INTO news (theme, link, site, time) VALUES ($1, $2, $3, $4);`
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	_, err = db.Instance.Exec(query, theme, link, site, t.In(loc))
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) QueryAll() ([]QueryResult, error) {
	arr := []QueryResult{}
	query := `SELECT * FROM news;`
	rows, err := db.Instance.Query(query)
	if err != nil {
		return []QueryResult{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&theme, &link, &site, &timeVal)
		if err != nil {
			return []QueryResult{}, err
		}
		arr = append(arr, QueryResult{theme, link, site, timeVal})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return arr, nil
}

func (db *DB) GetByTheme(theme string) (*QueryResult, error) {
	query := `SELECT * FROM news
              WHERE theme=$1;`

	var site string
	var link string
	var timeval time.Time
	err := db.Instance.QueryRow(query, theme).Scan(&theme, &link, &site, &timeval)
	if err != nil {
		return &QueryResult{}, err
	}

	return &QueryResult{theme, link, site, timeval}, nil
}

func (db *DB) FindByTheme(theme string) ([]QueryResult, error) {
	arr := []QueryResult{}

	query := `SELECT * FROM news
              WHERE theme LIKE $1;`

	rows, err := db.Instance.Query(query, "%"+theme+"%")
	if err != nil {
		return []QueryResult{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&theme, &link, &site, &timeVal)
		if err != nil {
			return []QueryResult{}, err
		}
		arr = append(arr, QueryResult{theme, link, site, timeVal})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return arr, nil
}
