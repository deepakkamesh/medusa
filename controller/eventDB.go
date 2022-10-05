package controller

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

const createSQL string = `
  CREATE TABLE IF NOT EXISTS events (
  id INTEGER NOT NULL PRIMARY KEY,
  tmstmp DATETIME NOT NULL,
  metric TEXT NOT NULL,
	value REAL NOT NULL,
	room TEXT NOT NULL,
	addr BLOB NOT NULL
  );`

type eventLog struct {
	tmstmp time.Time
	metric string
	value  float32
	room   string
	addr   []byte
}

type EventDB struct {
	db *sql.DB
}

func NewEventDB() (*EventDB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createSQL); err != nil {
		return nil, err
	}
	return &EventDB{
		db: db,
	}, nil
}

func (e *EventDB) LogEvent(ev eventLog) error {

	_, err := e.db.Exec("INSERT INTO events VALUES(NULL,?,?,?,?,?);", ev.tmstmp, ev.metric, ev.value, ev.room, ev.addr)
	if err != nil {
		return err
	}
	return nil
}

// returns any event from room in the last secs.
func (e *EventDB) GetEvent(metric string, room string, last time.Duration) ([]eventLog, error) {
	qry := "SELECT * FROM events WHERE metric = ? AND room = ? AND tmstmp > ? ORDER BY tmstmp"
	rows, err := e.db.Query(qry, metric, room, time.Now().Add(-last))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evLogs := []eventLog{}
	id := 0
	for rows.Next() {
		i := eventLog{}
		err = rows.Scan(&id, &i.tmstmp, &i.metric, &i.value, &i.room, &i.addr)
		if err != nil {
			return nil, err
		}
		evLogs = append(evLogs, i)
	}
	return evLogs, nil
}

func PP(e eventLog) {
	fmt.Printf("%v %v %v %v %v\n", e.tmstmp, e.metric, e.value, e.room, e.addr)
}
