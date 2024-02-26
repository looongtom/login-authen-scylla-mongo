package models

import (
	"github.com/gocql/gocql"
	"time"
)

type HistoryConnection struct {
	Id       gocql.UUID `db:"id"`
	Username string     `db:"username"`
	LoginAt  time.Time  `db:"login_at"`
	LogoutAt time.Time  `db:"logout_at"`
}
