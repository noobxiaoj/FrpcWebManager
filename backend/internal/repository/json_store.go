package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// JSONStore JSON 文件存储
type JSONStore struct {
	dataDir string
	mu      sync.RWMutex
}

// NewJSONStore 创建 JSON 存储
func NewJSONStore(dataDir string) (*JSONStore, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &JSONStore{
		dataDir: dataDir,
	}, nil
}

// Load 加载 JSON 文件
func (s *JSONStore) Load(filename string, v interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.dataDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在,返回默认值
			return nil
		}
		return err
	}

	return json.Unmarshal(data, v)
}

// Save 保存 JSON 文件
func (s *JSONStore) Save(filename string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.dataDir, filename)

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Delete 删除文件
func (s *JSONStore) Delete(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.dataDir, filename)
	return os.Remove(filePath)
}

// Exists 检查文件是否存在
func (s *JSONStore) Exists(filename string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.dataDir, filename)
	_, err := os.Stat(filePath)
	return err == nil
}

// GetDataDir 获取数据目录
func (s *JSONStore) GetDataDir() string {
	return s.dataDir
}
