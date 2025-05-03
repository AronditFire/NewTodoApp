package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AronditFire/todo-app/entity"
	"gorm.io/gorm"
)

const fileTable = "files"

type ParseRepo struct {
	db *gorm.DB
}

func NewParseRepo(db *gorm.DB) *ParseRepo {
	return &ParseRepo{db: db}
}

func (r *ParseRepo) ParseJSON(bindfile entity.BindFile) error {
	file, err := bindfile.File.Open() // open it
	if err != nil {
		return errors.New("Could not open file")
	}
	defer file.Close() // close at the end

	var fileData map[string]any
	if err := json.NewDecoder(file).Decode(&fileData); err != nil {
		return errors.New("Could not decode json")
	}

	tx := r.db.Begin() // tx launch
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	for key, value := range fileData {
		// Проверить, есть ли столбец
		if !tx.Migrator().HasColumn(fileTable, key) {
			// Выбрать SQL-тип
			sqlType := inferSQLType(value)

			stmt := fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN "%s" %s;`, fileTable, key, sqlType)
			if err := tx.Exec(stmt).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("cannot add column %q: %w", key, err)
			}
		}
	}

	if err := tx.Table(fileTable).Create(fileData).Error; err != nil {
		tx.Rollback()
		return errors.New("Could not create rows with data")
	}

	return tx.Commit().Error
}

func (r *ParseRepo) GetJsonTable() ([]map[string]any, error) {
	return nil, nil
}

func inferSQLType(value any) string {
	switch value.(type) {
	case bool:
		return "BOOLEAN"
	case float64:
		// JSON-числа в Go по умолчанию float64
		return "DOUBLE PRECISION"
	case string:
		return "TEXT"
	default:
		// на всякий случай, если придёт объект/массив/null
		return "JSONB"
	}
}
