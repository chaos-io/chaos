package internal

import (
	"database/sql"
	"time"
)

type Dandan struct {
	Id string
    // TODO add more fields

	CreatTime  time.Time
	UpdateTime time.Time
    DeleteTime sql.NullTime
}
