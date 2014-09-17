package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
	"github.com/smtc/justcms/database"
)

var (
	// reserve table name
	NotAllowTables = []string{"account", "role", "post"}
)

type Table struct {
	Id        int64
	Name      string    `sql:"size:45;not null;unique"`
	Alias     string    `sql:"size:45"`
	Des       string    `Sql:"size:512"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
	Columns   []Column
}

func getTableDB() *gorm.DB {
	return database.GetDB("table")
}

func (t *Table) Exist() bool {
	for _, nat := range NotAllowTables {
		if t.Name == nat {
			return true
		}
	}
	db := getTableDB()
	var count int
	db.Model(Table{}).Where("id!=? and name=?", t.Id, t.Name).Count(&count)
	return count > 0
}

func (t *Table) Get(id int64) error {
	db := getTableDB()
	return db.First(t, id).Error
}

func (t *Table) GetColumns() error {
	db := getTableDB()
	t.Columns = nil
	return db.Model(t).Related(&t.Columns).Order("name").Error
}

func (t *Table) Refresh() error {
	db := getTableDB()
	err := db.First(t, t.Id).Error
	if err != nil {
		return err
	}
	return t.GetColumns()
}

func (t *Table) Save() error {
	var (
		db     = getTableDB()
		isNew  = false
		column Column
	)
	if t.Exist() {
		return fmt.Errorf("table '%v' is existed", t.Name)
	}
	if t.Id == 0 {
		t.CreatedAt = time.Now()
		isNew = true
	}
	t.EditAt = time.Now()
	if err := db.Save(t).Error; err != nil {
		return err
	}

	if isNew {
		column.Name = "id"
		column.Alias = "id"
		column.Type = database.BIGINT
		column.PrimaryKey = true
		column.TableId = t.Id
		column.NotNull = true
		column.Save()
	}
	return nil
}

func (t *Table) Delete() error {
	db := getTableDB()
	if err := db.Where("table_id = ?", t.Id).Delete(Column{}).Error; err != nil {
		return err
	}
	return db.Delete(t).Error
}

func TableDelete(where string) {
	db := getTableDB()
	db.Where(where).Delete(&Table{})
}

func TableList(tbls *[]Table) error {
	db := getTableDB()
	err := db.Find(tbls).Error
	return err
}

func (t *Table) Field(name string) Column {
	var column Column
	for _, c := range t.Columns {
		if c.Name == name {
			column = c
		}
	}
	return column
}

func (t Table) MarshalJSON() ([]byte, error) {
	j, _ := goutils.ToJsonOnly(t)
	return []byte(j), nil
}

func (t *Table) CreateTable() error {
	var sql []string
	sql = append(sql, fmt.Sprintf("CREATE table `%v` (", t.Name))
	return nil
}

func (t *Table) MigrateTable() error {
	return nil
}
