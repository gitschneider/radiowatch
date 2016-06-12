package radiowatch

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"time"
)

type  MysqlWriter struct {
	username string
	password string
	address  string
	database string
	port     string
	created  map[string]bool
}

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "Error when writing to database: %v\n", err.Error())
}

func NewMysqlWriter(username, password, address, port, database string) MysqlWriter {
	return MysqlWriter{
		username : username,
		password: password,
		address: address,
		database : database,
		port: port,
		created: make(map[string]bool, 1),
	}
}

func (m MysqlWriter) Write(ti TrackInfo) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", m.username, m.password, m.address, m.port, m.database))
	if err != nil {
		handleErr(err)
		return
	}
	defer db.Close()

	if !m.created[ti.NormalizedStationName()] {
		_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `radiowatch`.`%v` ( `id` INT UNSIGNED NOT NULL AUTO_INCREMENT , `artist` VARCHAR(255) NOT NULL , `title` VARCHAR(255) NOT NULL , `station` VARCHAR(255) NOT NULL , `time` INT UNSIGNED NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB; ", ti.NormalizedStationName()))
		if err != nil {
			handleErr(err)
			return
		}
		m.created[ti.NormalizedStationName()] = true
	}

	insertQuery := fmt.Sprintf(`
INSERT INTO %[1]v.%[2]v (artist, title, station, time)
  SELECT ?, ?, ?, ?
    from dual
    WHERE NOT exists(SELECT * from %[1]v.%[2]v
                        WHERE id = (select max(id) from %[1]v.%[2]v)
                          AND artist = ?
                          and title = ?);
	`, m.database, ti.NormalizedStationName())
	_, err = db.Exec(
		insertQuery,
		ti.Artist,
		ti.Title,
		ti.Station,
		time.Now().Unix(),
		ti.Artist,
		ti.Title,)
	if err != nil {
		handleErr(err)
		return
	}
}
