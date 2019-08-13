package morph

import "github.com/go-pg/pg/v9"

type Stepable interface {
	Up(db *pg.DB) error
	Down(db *pg.DB) error
}
