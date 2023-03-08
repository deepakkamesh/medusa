package controller

import (
	"database/sql"
	"fmt"
	"sync"
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
	name TEXT NOT NULL,
	addr BLOB NOT NULL
  );`

type EventLog struct {
	Tmstmp time.Time
	Metric string  // Metric type. eg. Temperature.
	Value  float32 // value of metric.
	Room   string  // Room assigned to board.
	Name   string  // Board Name.
	Addr   []byte  // Addr of board.
}

type EventDB struct {
	db *sql.DB
	mu sync.Mutex
}

func NewEventDB() (*EventDB, error) {
	db, err := sql.Open("sqlite", ":memory:?cache=shared") // see https://github.com/mattn/go-sqlite3/blob/master/README.md#faq
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

func (e *EventDB) LogEvent(ev EventLog) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, err := e.db.Exec("INSERT INTO events VALUES(NULL,?,?,?,?,?,?);", ev.Tmstmp, ev.Metric, ev.Value, ev.Room, ev.Name, ev.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventDB) PurgeEvent(ev EventLog) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, err := e.db.Exec("DELETE from events")
	if err != nil {
		return err
	}
	return nil
}

// GetEvent returns any event from room, name in the last duration.
func (e *EventDB) GetEvent(metric string, room string, name string, last time.Duration) ([]EventLog, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	qry := "SELECT * FROM events WHERE metric = ? AND room = ? AND name = ? AND tmstmp > ? ORDER BY tmstmp"
	rows, err := e.db.Query(qry, metric, room, name, time.Now().Add(-last))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evLogs := []EventLog{}
	id := 0
	for rows.Next() {
		i := EventLog{}
		err = rows.Scan(&id, &i.Tmstmp, &i.Metric, &i.Value, &i.Room, &i.Name, &i.Addr)
		if err != nil {
			return nil, err
		}
		evLogs = append(evLogs, i)
	}
	return evLogs, nil
}

func PP(e EventLog) {
	fmt.Printf("%v %v %v %v %v %v\n", e.Tmstmp, e.Metric, e.Value, e.Room, e.Name, e.Addr)
}
