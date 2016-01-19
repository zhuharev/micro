package gorp

import (
	"database/sql"
	"github.com/lann/squirrel"
	"github.com/zhuharev/micro/db"
	"reflect"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type GorpDb struct {
	db *gorp.DbMap
}

func New(driver, config string) (*GorpDb, error) {
	if conn, e := sql.Open(driver, config); e != nil {
		return nil, e
	} else {
		e = conn.Ping()
		if e != nil {
			return nil, e
		}
		dbmap := &gorp.DbMap{Db: conn}
		switch driver {
		case "sqlite3":
			dbmap.Dialect = gorp.SqliteDialect{}
		case "mysql":
			dbmap.Dialect = gorp.MySQLDialect{}

		//todo check syntax
		case "postgres":
			dbmap.Dialect = gorp.PostgresDialect{}
		default:
			panic("can't find dialect")
		}
		db := &GorpDb{}
		db.SetDb(dbmap)
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
		return db, nil
	}
}

func (db *GorpDb) SetDb(dbmap *gorp.DbMap) {
	db.db = dbmap
}

func (db *GorpDb) Registry(i interface{}, args ...interface{}) error {
	db.db.AddTable(i).SetKeys(true, "Id")
	for i := range args {
		db.db.AddTable(args[i]).SetKeys(true, "Id")
	}
	return db.db.CreateTablesIfNotExists()
	//return nil
}

func (db *GorpDb) AddTable(i interface{}) error {
	//panic if Id field not exist
	db.db.AddTable(i).SetKeys(true, "Id")
	return db.db.CreateTablesIfNotExists()
}

func (db *GorpDb) Save(id int64, obj interface{}) error {
	if id == 0 {
		return db.db.Insert(obj)
	} else {
		_, e := db.db.Update(obj)
		return e
	}
	return nil
}

func (db *GorpDb) Read(id int64, notPtrInterface interface{}) (interface{}, error) {
	return db.db.Get(notPtrInterface, id)
}

func (gorpdb *GorpDb) ReadMany(in interface{}, opts *db.Filter) ([]interface{}, error) {
	t := reflect.TypeOf(in)
	//sliceType := reflect.SliceOf(t)
	//sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	tmap, e := gorpdb.db.TableFor(t, false)
	if e != nil {
		return nil, e
	}
	tname := tmap.TableName
	query := squirrel.Select("*").From(tname)

	var (
		limit  = opts.Limit
		offset = opts.Offset
		page   = opts.Page
	)

	if page != 0 && limit != 0 {
		offset = (page - 1) * limit
	}

	if limit != 0 {
		query = query.Limit(limit)
	}
	if offset != 0 {
		query = query.Offset(offset)
	}
	if opts.Order != "" {
		query = query.OrderBy(opts.Order)
	}
	if opts.WhereEq != nil {
		for _, v := range opts.WhereEq {
			query = query.Where(v)
		}
	}
	qstring, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	//typeSliceInterface := reflect.ValueOf([]interface{}{}).Type()
	//slIface := sliceValue.Convert(typeSliceInterface).Interface().([]interface{})
	//var re []interface{}

	res, err := gorpdb.db.Select(in, qstring, args...)
	return res, err
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
