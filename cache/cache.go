package cache

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type CacheItem struct {
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	cacheFile = "cache/data.json"
	ttl       = time.Minute
	mutex     sync.Mutex
)

// GetData проверяет кэш и возвращает данные либо из файла, либо генерирует новые
func GetData(generate func() any) (any, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// Проверяем, существует ли файл и не истёк ли TTL
	if fi, err := os.Stat(cacheFile); err == nil && time.Since(fi.ModTime()) < ttl {
		content, err := os.ReadFile(cacheFile)
		if err == nil {
			var item CacheItem
			if json.Unmarshal(content, &item) == nil {
				return item.Data, nil
			}
		}
	}

	// Генерируем новые данные
	data := generate()
	item := CacheItem{
		Data:      data,
		CreatedAt: time.Now(),
	}

	// Сохраняем в файл
	content, _ := json.MarshalIndent(item, "", "  ")
	os.WriteFile(cacheFile, content, 0644)

	return data, nil
}
