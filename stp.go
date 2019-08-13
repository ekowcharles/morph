package morph

import "github.com/go-pg/pg/v9"

type Stepable interface {
	Up(db *pg.DB, st []interface{}) error
	Down(db *pg.DB, st []interface{}) error
}
