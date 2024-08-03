package main

import (
	"testing"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type TestTable struct {
	gorm.Model

	Path HierarchyIdDb
}

func TestGorm(t *testing.T) {
	dsn := "sqlserver://sa:12345678@localhost:9930?database=gorm"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&TestTable{})
	if err != nil {
		panic("failed to migrate")
	}

	conn := db.Create(&TestTable{Path: HierarchyIdDb{Data: []int64{1, 2, 3, 4}}})
	if conn.Error != nil {
		panic("failed to create")
	}

	hid := TestTable{}
	conn = db.First(&TestTable{})
	if conn.Error != nil {
		panic("failed to query")
	}

	if hid.Path.Data[0] != 1 && hid.Path.Data[1] != 2 && hid.Path.Data[2] != 3 && hid.Path.Data[3] != 4 {
		panic("failed to read")
	}
}
