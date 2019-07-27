package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"

	"time"
)

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
	Instance *sql.DB
}

type QueryResult struct {
	Id    int       `json:"id"`
	Theme string    `json:"theme"`
	Link  string    `json:"link"`
	Site  string    `json:"site"`
	Time  time.Time `json:"time"`
}

func CreateDBConf(host string, port int, user, password, dbname string) *DB {
	return &DB{host, port, user, password, dbname, &sql.DB{}}
}

func (db *DB) Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DBname)

	dbinstance, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = dbinstance.Ping()
	if err != nil {
		return err
	}
	db.Instance = dbinstance
	return nil
}

func (db *DB) Init() error {
	_, err := db.Instance.Exec(`
        CREATE TABLE news (
    		id serial PRIMARY KEY,
    		theme TEXT,
    		link TEXT,
    	    site TEXT,
    	    time TIMESTAMP 
    	)`)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Insert(theme string, link string, site string, time time.Time) error {
	query := `INSERT INTO public.news (theme, link, site, time) VALUES ($1, $2, $3, $4);`
	_, err := db.Instance.Exec(query, theme, link, site, time)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) QueryAll() ([]QueryResult, error) {
	arr := []QueryResult{}
	query := `SELECT * FROM public.news;`
	rows, err := db.Instance.Query(query)
	if err != nil {
		return []QueryResult{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var id int
		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&id, &theme, &link, &site, &timeVal)
		if err != nil {
			return []QueryResult{}, err
		}
		arr = append(arr, QueryResult{id, theme, link, site, timeVal})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return arr, nil
}

func (db *DB) GetByTheme(theme string) (*QueryResult, error) {
	query := `SELECT * FROM public.news
              WHERE theme=$1;`

	var id int
	var site string
	var link string
	var timeval time.Time
	err := db.Instance.QueryRow(query, theme).Scan(&id, &theme, &link, &site, &timeval)
	if err != nil {
		return &QueryResult{}, err
	}

	return &QueryResult{id, theme, link, site, timeval}, nil
}

func (db *DB) FindByTheme(theme string) ([]QueryResult, error) {
	arr := []QueryResult{}

	query := `SELECT * FROM public.news
              WHERE theme LIKE $1;`

	rows, err := db.Instance.Query(query, "%"+theme+"%")
	if err != nil {
		return []QueryResult{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var id int
		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&id, &theme, &link, &site, &timeVal)
		if err != nil {
			return []QueryResult{}, err
		}
		arr = append(arr, QueryResult{id, theme, link, site, timeVal})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return arr, nil
}
