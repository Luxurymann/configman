package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Storage struct {
	path string
	data map[string]string
}
type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Storage) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}
func (s *Storage) Set(key, value string) error {
	s.data[key] = value
	// Используем метод Save, который мы создадим ниже
	return s.Save()
}
func (s *Storage) Delete(key string) error {
	delete(s.data, key)
	return s.Save()
}
func (s *Storage) Save() error {
	return SaveStorage(s.path, s.data)
}
func (s *Storage) Reload() error {
	loadedData, err := LoadStorage(s.path)
	if err != nil {
		return err
	}
	s.data = loadedData
	return nil
}
func (s *Storage) List() []Pair {
	var list []Pair
	for key, value := range s.data {
		p := Pair{Key: key, Value: value}
		list = append(list, p)
	}
	return list
}
func (s *Storage) Export(prefix string) {
	for key, value := range s.data {
		fmt.Printf("export %s%s=%q\n", prefix, strings.ToUpper(key), value)
	}
}
func LoadStorage(path string) (map[string]string, error) {
	storage := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла: %w", err)
	}
	err = json.Unmarshal(data, &storage)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при десериализации: %w", err)
	}
	return storage, nil
}
func SaveStorage(path string, storage map[string]string) error {
	jsonData, err := json.Marshal(storage)
	if err != nil {
		return fmt.Errorf("Ошибка при сериализации: %w", err)
	}
	err = os.WriteFile(path, []byte(jsonData), 0644)
	if err != nil {
		return fmt.Errorf("Ошибка при записи: %w", err)
	}
	return nil
}
func NewStorage(path string) (*Storage, error) {
	s := &Storage{
		path: path,
	}
	data, err := LoadStorage(path)
	if os.IsNotExist(err) {
		s.data = make(map[string]string)
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось создать хранилище: %w", err)
	}
	s.data = data
	return s, nil
}
