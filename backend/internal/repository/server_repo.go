package repository

import (
	"errors"
	"fmt"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// ServerRepository 服务器数据仓库
type ServerRepository struct {
	store *JSONStore
}

// ServerData 服务器数据结构（用于JSON存储）
type ServerData struct {
	Servers []model.Server `json:"servers"`
}

// NewServerRepository 创建服务器仓库
func NewServerRepository(store *JSONStore) *ServerRepository {
	return &ServerRepository{
		store: store,
	}
}

// ListServers 获取服务器列表
func (r *ServerRepository) ListServers() ([]model.Server, error) {
	var data ServerData
	if err := r.store.Load("servers.json", &data); err != nil {
		return nil, err
	}
	return data.Servers, nil
}

// GetServer 根据ID获取服务器
func (r *ServerRepository) GetServer(id string) (*model.Server, error) {
	servers, err := r.ListServers()
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.ID == id {
			return &server, nil
		}
	}

	return nil, errors.New("服务器不存在")
}

// CreateServer 创建服务器
func (r *ServerRepository) CreateServer(server *model.Server) error {
	var data ServerData
	r.store.Load("servers.json", &data)

	data.Servers = append(data.Servers, *server)

	return r.store.Save("servers.json", &data)
}

// UpdateServer 更新服务器
func (r *ServerRepository) UpdateServer(server *model.Server) error {
	var data ServerData
	if err := r.store.Load("servers.json", &data); err != nil {
		return err
	}

	found := false
	for i, s := range data.Servers {
		if s.ID == server.ID {
			data.Servers[i] = *server
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("服务器不存在: %s", server.ID)
	}

	return r.store.Save("servers.json", &data)
}

// DeleteServer 删除服务器
func (r *ServerRepository) DeleteServer(id string) error {
	var data ServerData
	if err := r.store.Load("servers.json", &data); err != nil {
		return err
	}

	found := false
	newServers := make([]model.Server, 0, len(data.Servers)-1)
	for _, s := range data.Servers {
		if s.ID != id {
			newServers = append(newServers, s)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("服务器不存在: %s", id)
	}

	data.Servers = newServers
	return r.store.Save("servers.json", &data)
}

// UpdateServerOrder 更新服务器排序
func (r *ServerRepository) UpdateServerOrder(order []string) error {
	var data ServerData
	if err := r.store.Load("servers.json", &data); err != nil {
		return err
	}

	// 创建ID到服务器的映射
	serverMap := make(map[string]model.Server)
	for _, server := range data.Servers {
		serverMap[server.ID] = server
	}

	// 按照新顺序重新组织服务器列表
	orderedServers := make([]model.Server, 0, len(order))
	for _, id := range order {
		if server, exists := serverMap[id]; exists {
			orderedServers = append(orderedServers, server)
			delete(serverMap, id)
		}
	}

	// 将剩余未在order中的服务器追加到末尾
	for _, server := range serverMap {
		orderedServers = append(orderedServers, server)
	}

	data.Servers = orderedServers
	return r.store.Save("servers.json", &data)
}
