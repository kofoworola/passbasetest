package storage

import "time"

type Project struct {
	ID      string    `db:"project_id"`
	Name    string    `db:"name"`
	ApiKey  string    `db:"api_key"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}
