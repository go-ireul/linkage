package main

import "ireul.com/orm"

// DB wrapped orm.DB
type DB struct {
	*orm.DB
}

// Record record
type Record struct {
	orm.Model
	Name string `json:"name" orm:"unique_index"`
	URL  string `json:"url"`
}

// NewDB create a new DB
func NewDB(c string, debug bool) (db *DB, err error) {
	var d *orm.DB
	if d, err = orm.Open("mysql", c); err != nil {
		return
	}
	db = &DB{DB: d}
	db.LogMode(debug)
	err = db.AutoMigrate()
	return
}

// AutoMigrate migrate all models
func (d *DB) AutoMigrate() error {
	return d.DB.AutoMigrate(Record{}).Error
}
