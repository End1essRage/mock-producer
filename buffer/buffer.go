package buffer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
)

type Buffer struct {
	Templates map[string]api.Pattern
}

func New() *Buffer {
	return &Buffer{Templates: make(map[string]api.Pattern)}
}

func (b *Buffer) FillFromFiles(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ошибка чтения директории: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".json" {
			continue
		}

		// Формируем ключ: имя файла без .json
		key := strings.TrimSuffix(filename, filepath.Ext(filename))

		// Читаем содержимое файла
		filePath := filepath.Join(path, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
		}

		// Парсим JSON
		var pattern api.Pattern
		if err := json.Unmarshal(data, &pattern); err != nil {
			return fmt.Errorf("ошибка парсинга JSON в файле %s: %w", filePath, err)
		}

		// Сохраняем в буфер
		b.Templates[key] = pattern
	}

	return nil
}

func (b *Buffer) GetAll() []string {
	keys := make([]string, 0)
	for k, _ := range b.Templates {
		keys = append(keys, k)
	}

	return keys
}

func (b *Buffer) Get(name string) (api.Pattern, error) {
	v, ok := b.Templates[name]
	if !ok {
		return nil, fmt.Errorf("нет паттерна по ключу %s", name)
	}

	return v, nil
}
func (b *Buffer) Set(name string, pattern api.Pattern) error {
	b.Templates[name] = pattern

	return nil
}
