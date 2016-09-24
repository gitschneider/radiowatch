package radiowatch

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/Sirupsen/logrus"
	"time"
	"fmt"
)

type  MysqlWriter struct {
	username string
	password string
	address  string
	database string
	port     string
	created  map[string]bool
	firstRun bool
}

func NewMysqlWriter(username, password, address, port, database string) MysqlWriter {
	return MysqlWriter{
		username : username,
		password: password,
		address: address,
		database : database,
		port: port,
		created: make(map[string]bool, 1),
		firstRun: true,
	}
}

func (m MysqlWriter) Write(ti TrackInfo) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", m.username, m.password, m.address, m.port, m.database))
	if err != nil {
		log.WithField(
		    "message",
		    err.Error(),
		).Error("Error when writing to database")
		return
	}
	defer db.Close()

	if m.firstRun {
		sql := `CREATE TABLE IF NOT EXISTS records  (
   					id  int(10) unsigned NOT NULL AUTO_INCREMENT,
				   	artist  varchar(255) COLLATE utf8_bin NOT NULL,
				   	title  varchar(255) COLLATE utf8_bin NOT NULL,
				   	station  varchar(255) COLLATE utf8_bin NOT NULL,
				   	time  int(10) unsigned NOT NULL,
			   PRIMARY KEY ( id ),
			   KEY  time  ( time )
			 ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`

		_, err := db.Exec(sql)
		if err != nil {
			log.WithField(
				"message",
				err.Error(),
			).Error("Error when writing to database")
			return
		}

		m.firstRun = false
	}

	insertQuery := `
		INSERT INTO records (artist, title, station, time)
  			SELECT ?, ?, ?, ?
    		from dual
    		WHERE NOT exists(SELECT * from records
            	WHERE id = (select max(id) from records where station = ?)
				AND artist = ?
                and title = ?);
	`
	_, err = db.Exec(
		insertQuery,
		ti.Artist,
		ti.Title,
		ti.NormalizedStationName(),
		time.Now().Unix(),
		ti.NormalizedStationName(),
		ti.Artist,
		ti.Title, )
	if err != nil {
		log.WithField(
			"message",
			err.Error(),
		).Error("Error when writing to database")
		return
	}
}
