package sqlite

import (
	"aviso/domain"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DB struct {
	Instance *sql.DB
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
        CREATE TABLE IF NOT EXISTS news (
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

func (db *DB) QueryAll() ([]domain.WebObj, error) {
	arr := []domain.WebObj{}
	query := `SELECT * FROM news ORDER BY time DESC;`
	rows, err := db.Instance.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&theme, &link, &site, &timeVal)
		if err != nil {
			return nil, err
		}
		arr = append(arr, domain.WebObj{theme, link, site, timeVal})
	}
	if rows.Err() != nil {
		return nil, err
	}

	return arr, nil
}

func (db *DB) GetByTheme(theme string) (*domain.WebObj, error) {
	query := `SELECT * FROM news
              WHERE theme=$1;`

	var site string
	var link string
	var t time.Time
	err := db.Instance.QueryRow(query, theme).Scan(&theme, &link, &site, &t)
	if err != nil {
		return nil, err
	}

	return &domain.WebObj{theme, link, site, t}, nil
}

func (db *DB) FindByTheme(theme string) ([]domain.WebObj, error) {
	arr := []domain.WebObj{}

	query := `SELECT * FROM news
              WHERE theme LIKE $1;`

	rows, err := db.Instance.Query(query, "%"+theme+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var theme string
		var link string
		var site string
		var timeVal time.Time

		err := rows.Scan(&theme, &link, &site, &timeVal)
		if err != nil {
			return nil, err
		}
		arr = append(arr, domain.WebObj{theme, link, site, timeVal})
	}
	if rows.Err() != nil {
		return nil, err
	}

	return arr, nil
}
