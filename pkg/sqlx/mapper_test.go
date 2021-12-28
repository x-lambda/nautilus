package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/stretchr/testify/assert"
)

var mysqldb *DB

func init() {
	Connect()
}

func Connect() {
	ctx := context.TODO()
	mysqldb = Get(ctx, "test")
}

type Schema struct {
	create string
	drop   string
}

func (s Schema) MySQL() (string, string, string) {
	return strings.Replace(s.create, `"`, `"`, -1), s.drop, `now()`
}

var defaultSchema = Schema{
	create: `
CREATE TABLE person (
	first_name text,
	last_name text,
	email text,
	added_at timestamp default now()
);
CREATE TABLE place (
	country text,
	city text NULL,
	telcode integer
);
CREATE TABLE capplace (
	COUNTRY text,
	CITY    text NULL,
	TELCODE integer
);
CREATE TABLE nullperson (
    first_name text NULL,
    last_name text NULL,
    email text NULL
);
CREATE TABLE employees (
	name text,
	id integer,
	boss_id integer
);

`,
	drop: `
drop table person;
drop table place;
drop table capplace;
drop table nullperson;
drop table employees;
`,
}

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
	AddedAt   time.Time `db:"added_at"`
}

type Person2 struct {
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	Email     sql.NullString
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

type PlacePtr struct {
	Country string
	City    *string
	TelCode int
}

type PersonPlace struct {
	Person
	Place
}

type PersonPlacePtr struct {
	*Person
	*Place
}

type EmbedConflict struct {
	FirstName string `db:"first_name"`
	Person
}

type SliceMember struct {
	Country   string
	City      sql.NullString
	TelCode   int
	People    []Person `db:"-"`
	Addresses []Place  `db:"-"`
}

type CPlace Place

func MultiExec(t *testing.T, e sqlx.Execer, query string) {
	stmts := strings.Split(query, ";\n")
	if len(strings.Trim(stmts[len(stmts)-1], "\n\t\r")) == 0 {
		stmts = stmts[:len(stmts)-1]
	}

	for _, s := range stmts {
		_, err := e.Exec(s)
		assert.Nil(t, err)
	}
}

func RunWithSchema(schema Schema, t *testing.T, test func(db *DB, t *testing.T, now string)) {
	runner := func(db *DB, t *testing.T, create string, drop string, now string) {
		defer func() {
			MultiExec(t, db, drop)
		}()

		MultiExec(t, db, create)
		test(db, t, now)
	}

	create, drop, now := schema.MySQL()
	runner(mysqldb, t, create, drop, now)
}

func loadDefaultFixture(db *DB, t *testing.T) {
	tx := db.MustBegin()
	tx.MustExec(tx.Rebind("INSERT INTO person(first_name, last_name, email) VALUES (?, ?, ?)"), "Jason", "Moiron", "foo@horizon.ai")
	tx.MustExec(tx.Rebind("INSERT INTO person(first_name, last_name, email) VALUES (?, ?, ?)"), "John", "Doe", "bar@horizon.ai")

	tx.MustExec(tx.Rebind("INSERT INTO place(country, city, telcode) VALUES (?, ?, ?)"), "United States", "New York", "1")
	tx.MustExec(tx.Rebind("INSERT INTO place(country, city, telcode) VALUES (?, ?, ?)"), "China", "Hong Kong", "852")
	tx.MustExec(tx.Rebind("INSERT INTO place(country, telcode) VALUES (?, ?)"), "Singapore", "65")

	tx.MustExec(tx.Rebind("INSERT INTO capplace(`COUNTRY`, `TELCODE`) VALUES (?, ?)"), "Sarf Efrica", "27")

	tx.MustExec(tx.Rebind("INSERT INTO employees(name, id) VALUES (?, ?)"), "Peter", "4444")
	tx.MustExec(tx.Rebind("INSERT INTO employees(name, id, boss_id) VALUES (?, ?, ?)"), "Joe", "1", "4444")
	tx.MustExec(tx.Rebind("INSERT INTO employees(name, id, boss_id) VALUES (?, ?, ?)"), "Martin", "2", "4444")

	tx.Commit()
}

func TestMissingNames(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		type PersonPlus struct {
			FirstName string `db:"first_name"`
			LastName  string `db:"last_name"`
			Email     string
		}

		// first_name text,
		// last_name  text,
		// email      text,
		// added_at   timestamp default now()

		pps := []PersonPlus{}
		err := db.Select(&pps, "SELECT * FROM person")
		if err == nil {
			t.Error("Expected missing name from Select to fail, but it did not.")
		}

		pp := PersonPlus{}
		err = db.Get(&pp, "SELECT * FROM person LIMIT 1")
		if err == nil {
			t.Error("Expected missing name Get to fail, but it did not.")
		}

		pps = []PersonPlus{}
		rows, err := db.Query("SELECT * FROM person LIMIT 1")
		if err != nil {
			t.Fatal(err)
		}
		rows.Next()

		err = sqlx.StructScan(rows, &pps)
		if err == nil {
			t.Error("Expected missing name in StructScan to fail, but it did not.")
		}
		rows.Close()

		// 这里的safe是针对MySQL查询出来的列在结构体中不存在的情况
		// try various things with unsafe set
		db = &DB{db.Unsafe()}
		pps = []PersonPlus{}
		err = db.Select(&pps, "SELECT * FROM person")
		assert.Nil(t, err)

		pp = PersonPlus{}
		err = db.Get(&pp, "SELECT * FROM person LIMIT 1")
		assert.Nil(t, err)

		pps = []PersonPlus{}
		rowsx, err := db.Queryx("SELECT * FROM person LIMIT 1")
		assert.Nil(t, err)

		rowsx.Next()
		err = sqlx.StructScan(rowsx, &pps)
		assert.Nil(t, err)
		rowsx.Close()

		nstmt, err := db.PrepareNamed(`SELECT * FROM person WHERE first_name != :name`)
		assert.Nil(t, err)

		pps = []PersonPlus{}
		err = nstmt.Select(&pps, map[string]interface{}{"name": "Jason"})
		assert.Nil(t, err)
		assert.Equal(t, len(pps), 1)
	})
}

func TestEmbeddedStructs(t *testing.T) {
	type Loop1 struct{ Person }
	type Loop2 struct{ Loop1 }
	type Loop3 struct{ Loop2 }

	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		// Warning: 在连接数据库的时候，因为做了时间的自动转换，需要在dsn中指定自动解析时间参数 ?parseTime=true
		per := []Person{}
		err := db.Select(&per, "SELECT * FROM person")
		assert.Nil(t, err)
		fmt.Println(per[0].AddedAt.Format(time.RFC3339))

		peopleAndPlaces := []PersonPlace{}
		err = db.Select(
			&peopleAndPlaces,
			`SELECT person.*, place.* FROM
             person natural join place`)
		if err != nil {
			t.Fatal(err)
		}

		for _, pp := range peopleAndPlaces {
			assert.NotEqual(t, len(pp.Person.FirstName), 0)
			assert.NotEqual(t, len(pp.Place.Country), 0)
		}

		rows, err := db.Queryx(`SELECT person.*, place.* FROM person natural join place`)
		assert.Nil(t, err)

		perp := PersonPlace{}
		rows.Next()
		err = rows.StructScan(&perp)
		assert.Nil(t, err)

		assert.NotEqual(t, len(perp.Person.FirstName), 0)
		assert.NotEqual(t, len(perp.Place.Country), 0)
		rows.Close()

		peopleAndPlacesPtrs := []PersonPlacePtr{}
		err = db.Select(&peopleAndPlacesPtrs, `SELECT person.*, place.* FROM person natural join place`)
		assert.Nil(t, err)

		for _, pp := range peopleAndPlacesPtrs {
			assert.NotEqual(t, len(pp.Person.FirstName), 0)
			assert.NotEqual(t, len(pp.Place.Country), 0)
		}

		l3s := []Loop3{}
		err = db.Select(&l3s, `SELECT * FROM person`)
		assert.Nil(t, err)

		for _, l3 := range l3s {
			assert.NotEqual(t, len(l3.Loop2.Loop1.Person.FirstName), 0)
		}

		ec := []EmbedConflict{}
		err = db.Select(&ec, `SELECT * FROM person`)
		assert.Nil(t, err)
	})
}

func TestJoinQuery(t *testing.T) {
	type Employee struct {
		Name   string
		ID     int64
		BossID sql.NullInt64 `db:"boss_id"`
	}
	type Boss Employee

	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		var employees []struct {
			Employee
			Boss `db:"boss"`
		}

		err := db.Select(&employees, `SELECT employees.*, boss.id "boss.id", boss.name "boss.name" 
				FROM employees JOIN employees AS boss ON employees.boss_id = boss.id`)
		assert.Nil(t, err)

		for _, em := range employees {
			assert.NotEqual(t, em.Employee.Name, 0)
			assert.Equal(t, em.Employee.BossID.Int64, em.Boss.ID)
		}
	})
}

func TestJoinQueryNamedPointerStructs(t *testing.T) {
	type Employee struct {
		Name   string
		ID     int64
		BossID sql.NullInt64 `db:"boss_id"`
	}
	type Boss Employee

	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		var employees []struct {
			Emp1  *Employee `db:"emp1"`
			Emp2  *Employee `db:"emp2"`
			*Boss `db:"boss"`
		}

		err := db.Select(&employees, `SELECT emp.name "emp1.name", emp.id "emp1.id", emp.boss_id "emp1.boss_id",
				emp.name "emp2.name", emp.id "emp2.id", emp.boss_id "emp2.boss_id", 
       			boss.id "boss.id", boss.name "boss.name" 
				FROM employees AS emp JOIN employees AS boss ON emp.boss_id = boss.id`)
		assert.Nil(t, err)

		for _, em := range employees {
			assert.NotEqual(t, em.Emp1.Name, 0)
			assert.NotEqual(t, em.Emp2.Name, 0)

			assert.Equal(t, em.Emp1.BossID.Int64, em.Boss.ID)
			assert.Equal(t, em.Emp2.BossID.Int64, em.Boss.ID)
		}
	})
}

func TestSelectSliceMapTime(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		rows, err := db.Queryx("SELECT * FROM person")
		assert.Nil(t, err)

		for rows.Next() {
			_, err = rows.SliceScan()
			assert.Nil(t, err)
		}

		rows, err = db.Queryx("SELECT * FROM person")
		assert.Nil(t, err)

		for rows.Next() {
			m := map[string]interface{}{}
			err = rows.MapScan(m)
			assert.Nil(t, err)
		}
	})
}

func TestNilReceiver(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		var p *Person
		err := db.Get(p, "SELECT * FROM person LIMIT 1")
		assert.NotNil(t, err)

		var pp *[]Person
		err = db.Select(pp, "SELECT * FROM person")
		assert.NotNil(t, err)
	})
}

func TestNamedQuery(t *testing.T) {
	var schema = Schema{
		create: `
				CREATE TABLE place(id integer PRIMARY KEY, name text NULL);
				CREATE TABLE person(first_name text NULL, last_name text NULL, email text NULL);
				CREATE TABLE placeperson(first_name text NULL, last_name text NULL, email text NULL, place_id integer NULL);
				CREATE TABLE jsperson(FIRST text NULL, last_name text NULL, EMAIL text NULL);
				`,
		drop: `
				drop table person;
				drop table jsperson;
				drop table place;
				drop table placeperson;
			`,
	}

	RunWithSchema(schema, t, func(db *DB, t *testing.T, now string) {
		type Person struct {
			FirstName sql.NullString `db:"first_name"`
			LastName  sql.NullString `db:"last_name"`
			Email     sql.NullString
		}

		p := Person{
			FirstName: sql.NullString{String: "ben", Valid: true},
			LastName:  sql.NullString{String: "doe", Valid: true},
			Email:     sql.NullString{String: "ben@doe.com", Valid: true},
		}

		q1 := `INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)`
		_, err := db.NamedExec(q1, p)
		assert.Nil(t, err)

		p2 := &Person{}
		rows, err := db.NamedQuery("SELECT * FROM person WHERE first_name=:first_name", p)
		assert.Nil(t, err)
		for rows.Next() {
			err = rows.StructScan(p2)
			assert.Nil(t, err)

			assert.Equal(t, p2.FirstName.String, "ben")
			assert.Equal(t, p2.LastName.String, "doe")
		}

		old := *db.Mapper

		type JSONPerson struct {
			FirstName sql.NullString `json:"FIRST"`
			LastName  sql.NullString `json:"last_name"`
			Email     sql.NullString
		}

		jp := JSONPerson{
			FirstName: sql.NullString{String: "ben", Valid: true},
			LastName:  sql.NullString{String: "smith", Valid: true},
			Email:     sql.NullString{String: "ben@smith.com", Valid: true},
		}

		db.Mapper = reflectx.NewMapperFunc("json", strings.ToUpper)

		pdb := func(s string, db *DB) string {
			if db.DriverName() == "mysql" {
				return strings.Replace(s, `"`, "`", -1)
			}
			return s
		}

		q1 = `INSERT INTO jsperson(FIRST, last_name, EMAIL) VALUES (:FIRST, :last_name, :EMAIL)`
		_, err = db.NamedExec(pdb(q1, db), jp)
		assert.Nil(t, err)

		check := func(t *testing.T, rows *sqlx.Rows) {
			jp = JSONPerson{}

			for rows.Next() {
				err = rows.StructScan(&jp)
				assert.Nil(t, err)
				assert.Equal(t, jp.FirstName.String, "ben")
				assert.Equal(t, jp.LastName.String, "smith")
				assert.Equal(t, jp.Email.String, "ben@smith.com")
			}
		}

		ns, err := db.PrepareNamed(pdb(`SELECT * FROM jsperson WHERE FIRST=:FIRST AND last_name=:last_name AND EMAIL=:EMAIL`, db))
		assert.Nil(t, err)
		rows, err = ns.Queryx(jp)
		assert.Nil(t, err)
		check(t, rows)

		rows, err = db.NamedQuery(pdb(`SELECT * FROM jsperson WHERE FIRST=:FIRST AND last_name=:last_name AND EMAIL=:EMAIL`, db), jp)
		assert.Nil(t, err)
		check(t, rows)

		db.Mapper = &old

		type Place struct {
			ID   int            `db:"id"`
			Name sql.NullString `db:"name"`
		}

		type PlacePerson struct {
			FirstName sql.NullString `db:"first_name"`
			LastName  sql.NullString `db:"last_name"`
			Email     sql.NullString
			Place     Place `db:"place"`
		}

		pl := Place{Name: sql.NullString{String: "myplace", Valid: true}}
		pp := PlacePerson{
			FirstName: sql.NullString{String: "ben", Valid: true},
			LastName:  sql.NullString{String: "doe", Valid: true},
			Email:     sql.NullString{String: "ben@doe.com", Valid: true},
		}

		q2 := `INSERT INTO place(id, name) VALUES (1, :name)`
		_, err = db.NamedQuery(q2, pl)
		assert.Nil(t, err)

		id := 1
		pp.Place.ID = id
		q3 := `INSERT INTO placeperson (first_name, last_name, email, place_id) VALUES (:first_name, :last_name, :email, :place.id)`
		_, err = db.NamedExec(q3, pp)
		assert.Nil(t, err)

		pp2 := &PlacePerson{}
		rows, err = db.NamedQuery(`SELECT first_name, last_name, email, place.id AS "place.id", place.name AS "place.name"
					FROM placeperson INNER JOIN place ON place.id = placeperson.place_id WHERE place.id=:place.id`, pp)
		assert.Nil(t, err)

		for rows.Next() {
			err = rows.StructScan(pp2)
			assert.Nil(t, err)

			assert.Equal(t, pp2.FirstName.String, "ben")
			assert.Equal(t, pp2.LastName.String, "doe")
			assert.Equal(t, pp2.Place.Name.String, "myplace")
			assert.Equal(t, pp2.Place.ID, pp.Place.ID)
		}
	})
}

func TestNilInserts(t *testing.T) {
	var schema = Schema{
		create: `CREATE TABLE tt (id integer, value text NULL DEFAULT NULL);`,
		drop:   "drop table tt;",
	}

	RunWithSchema(schema, t, func(db *DB, t *testing.T, now string) {
		type TT struct {
			ID    int
			Value *string
		}

		var v, v2 TT
		r := db.Rebind

		db.MustExec(r(`INSERT INTO tt(id) VALUES (1)`))
		db.Get(&v, r(`SELECT * FROM tt`))

		assert.Equal(t, v.ID, 1)
		assert.Nil(t, v.Value)

		v.ID = 2
		db.NamedExec(`INSERT INTO tt(id, value ) VALUES (:id, :value)`, v)
		db.Get(&v2, r(`SELECT * FROM tt WHERE id=2`))

		assert.Equal(t, v.ID, v2.ID)
		assert.Nil(t, v2.Value)
	})
}

func TestScanError(t *testing.T) {
	var schema = Schema{
		create: `CREATE TABLE kv (k text, v integer);`,
		drop:   `drop table kv;`,
	}

	RunWithSchema(schema, t, func(db *DB, t *testing.T, now string) {
		type WrongTypes struct {
			K int
			V string
		}

		_, err := db.Exec(db.Rebind("INSERT INTO kv (k, v) VALUES (?, ?)"), "hi", 1)
		assert.Nil(t, err)

		rows, err := db.Queryx("SELECT * FROM kv")
		assert.Nil(t, err)

		for rows.Next() {
			var wt WrongTypes
			err := rows.StructScan(&wt)
			assert.NotNil(t, err)
		}
	})
}

func TestMultiInsert(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		q := db.Rebind(`INSERT INTO employees (name, id) VALUES (?, ?), (?, ?)`)
		db.MustExec(q, "name1", 400, "name2", 500)
	})
}

func TestUsage(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)

		slicemembers := []SliceMember{}
		err := db.Select(&slicemembers, "SELECT * FROM place ORDER BY telcode ASC")
		assert.Nil(t, err)

		people := []Person{}
		err = db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
		assert.Nil(t, err)

		jason, john := people[0], people[1]
		assert.Equal(t, jason.FirstName, "Jason")
		assert.Equal(t, jason.LastName, "Moiron")
		assert.Equal(t, jason.Email, "foo@horizon.ai")

		assert.Equal(t, john.FirstName, "John")
		assert.Equal(t, john.LastName, "Doe")
		assert.Equal(t, john.Email, "bar@horizon.ai")

		jason = Person{}
		err = db.Get(&jason, db.Rebind("SELECT * FROM person WHERE first_name=?"), "Jason")
		assert.Nil(t, err)
		assert.Equal(t, jason.FirstName, "Jason")

		err = db.Get(&jason, db.Rebind("SELECT * FROM person WHERE first_name=?"), "Foobar")
		assert.Equal(t, err, sql.ErrNoRows)

		stmt1, err := db.Preparex(db.Rebind("SELECT * FROM person WHERE first_name=?"))
		assert.Nil(t, err)

		jason = Person{}
		row := stmt1.QueryRowx("DoesNotExist")
		row.Scan(&jason)
		row = stmt1.QueryRowx("DoesNotExist")
		row.Scan(&jason)

		err = stmt1.Get(&jason, "DoesNotExist User")
		assert.NotNil(t, err)

		err = stmt1.Get(&jason, "DoesNotExist User 2")
		assert.NotNil(t, err)

		stmt2, err := db.Preparex(db.Rebind("SELECT * FROM person WHERE first_name=?"))
		assert.Nil(t, err)

		jason = Person{}
		tx, err := db.Beginx()
		assert.Nil(t, err)
		tstmt2 := tx.Stmtx(stmt2)
		row2 := tstmt2.QueryRowx("Jason")
		err = row2.StructScan(&jason)
		assert.Nil(t, err)
		tx.Commit()

		places := []*Place{}
		err = db.Select(&places, "SELECT telcode FROM place ORDER BY telcode ASC")
		assert.Nil(t, err)

		usa, singsing, honkers := places[0], places[1], places[2]
		assert.Equal(t, usa.TelCode, 1)
		assert.Equal(t, honkers.TelCode, 852)
		assert.Equal(t, singsing.TelCode, 65)

		placesptr := []PlacePtr{}
		err = db.Select(&placesptr, "SELECT * FROM place ORDER BY telcode ASC")
		assert.Nil(t, err)

		places2 := []Place{}
		err = db.Select(&places2, "SELECT * FROM place ORDER BY telcode ASC")
		assert.Nil(t, err)

		usa, singsing, honkers = &places2[0], &places2[1], &places2[2]
		p := Place{}
		err = db.Select(&p, "SELECT * FROM place ORDER BY telcode ASC")
		assert.NotNil(t, err)

		pl := []Place{}
		err = db.Select(pl, "SELECT * FROM place ORDER BY telcode ASC")
		assert.NotNil(t, err)

		assert.Equal(t, usa.TelCode, 1)
		assert.Equal(t, honkers.TelCode, 852)
		assert.Equal(t, singsing.TelCode, 65)

		stmt, err := db.Preparex(db.Rebind("SELECT country, telcode FROM place WHERE telcode > ? ORDER BY telcode ASC"))
		assert.Nil(t, err)

		places = []*Place{}
		err = stmt.Select(&places, 10)
		assert.Nil(t, err)
		assert.Equal(t, len(places), 2)

		singsing, honkers = places[0], places[1]
		assert.Equal(t, singsing.TelCode, 65)
		assert.Equal(t, honkers.TelCode, 852)

		rows, err := db.Queryx("SELECT * FROM place")
		assert.Nil(t, err)
		place := Place{}
		for rows.Next() {
			err = rows.StructScan(&place)
			assert.Nil(t, err)
		}

		rows, err = db.Queryx("SELECT * FROM place")
		assert.Nil(t, err)
		m := map[string]interface{}{}
		for rows.Next() {
			err = rows.MapScan(m)
			assert.Nil(t, err)

			_, ok := m["country"]
			assert.True(t, ok)
		}

		rows, err = db.Queryx("SELECT * FROM place")
		assert.Nil(t, err)
		for rows.Next() {
			s, err := rows.SliceScan()
			assert.Nil(t, err)
			assert.Equal(t, len(s), 3)
		}

		_, err = db.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first, :last, :email)", map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "bensmith@allblacks.nz",
		})
		assert.Nil(t, err)

		rows, err = db.NamedQuery("SELECT * FROM person WHERE first_name=:first", map[string]interface{}{"first": "Bin"})
		assert.Nil(t, err)
		ben := &Person{}
		for rows.Next() {
			err = rows.StructScan(ben)
			assert.Nil(t, err)

			assert.Equal(t, ben.FirstName, "Bin")
			assert.Equal(t, ben.LastName, "Smuth")
		}

		ben.FirstName = "Ben"
		ben.LastName = "Smith"
		ben.Email = "binsmuth@allblacks.nz"
		_, err = db.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", ben)
		assert.Nil(t, err)

		rows, err = db.NamedQuery("SELECT * FROM person WHERE first_name=:first_name", ben)
		assert.Nil(t, err)
		for rows.Next() {
			err = rows.StructScan(ben)
			assert.Nil(t, err)

			assert.Equal(t, ben.FirstName, "Ben")
			assert.Equal(t, ben.LastName, "Smith")
		}

		person := &Person{}
		err = db.Get(person, "SELECT * FROM person WHERE first_name=$1", "does-not-exist")
		assert.NotNil(t, err)

		stmt, err = db.Preparex(db.Rebind("SELECT * FROM person WHERE first_name=?"))
		assert.Nil(t, err)
		rows, err = stmt.Queryx("Ben")
		assert.Nil(t, err)
		for rows.Next() {
			err = rows.StructScan(ben)
			assert.Nil(t, err)

			assert.Equal(t, ben.FirstName, "Ben")
			assert.Equal(t, ben.LastName, "Smith")
		}

		john = Person{}
		stmt, err = db.Preparex(db.Rebind("SELECT * FROM person WHERE first_name=?"))
		assert.Nil(t, err)
		err = stmt.Get(&john, "John")
		assert.Nil(t, err)

		db.MapperFunc(strings.ToUpper)
		rsa := CPlace{}
		err = db.Get(&rsa, "SELECT * FROM capplace;")
		assert.Nil(t, err)
		db.MapperFunc(strings.ToLower)

		dbCopy := sqlx.NewDb(db.DB.DB, db.DriverName())
		dbCopy.MapperFunc(strings.ToUpper)
		err = dbCopy.Get(&rsa, "SELECT * FROM capplace;")
		assert.Nil(t, err)

		err = db.Get(&rsa, "SELECT * FROM cappplace;")
		assert.NotNil(t, err)

		rows, err = db.Queryx("SELECT email FROM person ORDER BY email ASC;")
		assert.Nil(t, err)
		// ignore scan

		var count int
		err = db.Get(&count, "SELECT count(*) FROM person;")
		assert.Nil(t, err)

		var addedAt time.Time
		err = db.Get(&addedAt, "SELECT added_at FROM person LIMIT 1;")
		assert.Nil(t, err)

		var addedAts []time.Time
		err = db.Select(&addedAts, "SELECT added_at FROM person;")
		assert.Nil(t, err)

		var pcount *int
		err = db.Get(&pcount, "SELECT count(*) FROM person;")
		assert.Nil(t, err)
		assert.Equal(t, *pcount, count)

		sdest := []string{}
		err = db.Select(&sdest, "SELECT first_name FROM person ORDER BY first_name ASC;")
		assert.Nil(t, err)

		expected := []string{"Ben", "Bin", "Jason", "John"}
		for i, got := range sdest {
			assert.Equal(t, got, expected[i])
		}

		var nsdest []sql.NullString
		err = db.Select(&nsdest, "SELECT city FROM place ORDER BY city ASC")
		assert.Nil(t, err)
		assert.Equal(t, nsdest[0].String, "")
		assert.Equal(t, nsdest[1].String, "Hong Kong")
		assert.Equal(t, nsdest[2].String, "New York")
	})
}
