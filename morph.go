package morph

import (
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
)

const notRan = 1
const ran = 2

type Morph struct {
	ID       uuid.NullUUID `sql:",pk,type:uuid default uuid_generate_v4()"`
	FileName string
	Status   int

	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time `sql:"default:now()"`
}

func (m *Morph) className() string {
	// File names expected in the format 0012_MorphIntoSomethingMagical
	log.Printf("Extracting name of struct from %s\n", m.FileName)

	re, err := regexp.Compile("^[0-9]+\\_(.*?)(\\.go)$")
	panicIf(err)

	match := re.FindStringSubmatch(m.FileName)
	st := match[1]

	log.Printf("Struct identified to be %s\n", st)

	return st
}

func (m *Morph) getObject(tp *map[string]interface{}) interface{} {
	n := m.className()

	log.Printf("Retrieving %s object\n", n)

	return (*tp)[n]
}

func (m *Morph) up(db *pg.DB, tp *map[string]interface{}) {
	mig := m.getObject(tp)
	migr := mig.(Stepable)

	log.Printf("Migrating %s\n", m.FileName)
	err := migr.Up(db)
	panicIf(err)

	m.Status = ran
	db.Model(m).
		Column("status").
		WherePK().
		Update()
}

func (m *Morph) down(db *pg.DB, tp *map[string]interface{}) {
	mig := m.getObject(tp)
	migr := mig.(Stepable)

	log.Printf("Rolling back %s\n", m.FileName)
	err := migr.Down(db)
	panicIf(err)

	m.Status = notRan
	db.Model(m).
		Column("status").
		WherePK().
		Update()
}

func updateMorph(db *pg.DB, fn []string) {
	for _, f := range fn {
		m := &Morph{FileName: f, Status: notRan}

		c, err := db.Model((*Morph)(nil)).
			Where("file_name = ?", m.FileName).
			Count()
		panicIf(err)

		if c != 0 {
			continue
		}

		log.Printf("Adding %s to migrations\n", f)

		err = db.Insert(m)
		panicIf(err)
	}
}

func loadFiles(dr string) []string {
	log.Println("Identifying files to be migrated")

	files, err := ioutil.ReadDir(dr)
	panicIf(err)

	var fn []string
	for _, file := range files {
		fn = append(fn, file.Name())
	}

	return fn
}

func prepare(db *pg.DB, dr string) {
	fn := loadFiles(dr)

	updateMorph(db, fn)
}

func migrate(db *pg.DB, tp *map[string]interface{}, sp int) {
	log.Println("Migrating ...")

	ms := &[]Morph{}

	var err error
	if sp == 0 {
		err = db.Model(ms).
			Where("status = ?", notRan).
			Order("file_name ASC").
			Select()
	} else {
		err = db.Model(ms).
			Where("status = ?", notRan).
			Order("file_name ASC").
			Limit(sp).
			Select()
	}

	panicIf(err)

	for _, m := range *ms {
		m.up(db, tp)
	}
}

func rollback(db *pg.DB, tp *map[string]interface{}, sp int) {
	log.Println("Rolling back ...")

	ms := &[]Morph{}

	var err error
	if sp == 0 {
		err = db.Model(ms).
			Where("status = ?", ran).
			Order("file_name DESC").
			Select()
	} else {
		err = db.Model(ms).
			Where("status = ?", ran).
			Order("file_name DESC").
			Limit(sp).
			Select()
	}

	panicIf(err)

	for _, m := range *ms {
		m.down(db, tp)
	}
}
