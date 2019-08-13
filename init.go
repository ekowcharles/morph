package morph

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

var typeRegistry *map[string]interface{}

type Morpher struct {
	typeRegistry *map[string]interface{}

	db     *pg.DB
	Folder string
}

func Init(db *pg.DB, tp []interface{}) *Morpher {
	m := &Morpher{}
	m.db = db
	m.Folder = getEnv("MIGRATION_FOLDER", "./morph")

	m.typeRegistry = createTypeRegistry(tp)
	m.createSchema()

	return m
}

func (m *Morpher) Migrate() {
	m.MigrateWithStep(0)
}

func (m *Morpher) MigrateWithStep(sp int) {
	prepare(m.db, m.Folder)

	migrate(m.db, m.typeRegistry, sp)
}

func (m *Morpher) Rollback() {
	m.RollbackWithStep(0)
}

func (m *Morpher) RollbackWithStep(sp int) {
	prepare(m.db, m.Folder)

	rollback(m.db, m.typeRegistry, sp)
}

func (m *Morpher) createSchema() {
	_, err := m.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	panicIf(err)

	for _, model := range []interface{}{
		(*Morph)(nil),
	} {
		m.db.CreateTable((model), &orm.CreateTableOptions{
			IfNotExists: true,
		})
	}

	qs := []string{
		"DROP INDEX IF EXISTS morph_status_idx",
		"CREATE INDEX morph_status_idx ON morphs (status)",
		"DROP INDEX IF EXISTS morph_status_filename_idx",
		"CREATE UNIQUE INDEX morph_status_filename_idx ON morphs (status, file_name)",
		"DROP INDEX IF EXISTS morph_filename_idx",
		"CREATE UNIQUE INDEX morph_filename_idx ON morphs (file_name)",
	}

	for _, q := range qs {
		_, err := m.db.Exec(q)

		panicIf(err)
	}
}
