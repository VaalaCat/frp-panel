package models

import (
	"database/sql/driver"
	"encoding/json"
)

type GormArray[T any] []T

func (p GormArray[T]) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *GormArray[T]) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type JSON[T any] struct {
	Data T
}

func (j JSON[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSON[T]) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &j)
}
