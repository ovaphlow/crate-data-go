package service

import (
	"encoding/json"
	"fmt"
	"time"

	"ovaphlow.com/crate/data/repository"
	"ovaphlow.com/crate/data/utility"
)

// ApplicationService 定义了应用服务操作的接口。
type ApplicationService interface {
	Create(st string, d map[string]interface{}) (string, error)
	Get(st string, f [][]string, l string) (map[string]interface{}, error)
	Update(st string, d map[string]interface{}, w string, deprecated bool) error
	Remove(st string, w string) error
}

// ApplicationServiceImpl 实现了 ApplicationService 接口。
type ApplicationServiceImpl struct {
	repo repository.RDBRepo
}

// NewApplicationService 创建一个新的 ApplicationServiceImpl 实例。
func NewApplicationService(repo repository.RDBRepo) *ApplicationServiceImpl {
	return &ApplicationServiceImpl{repo: repo}
}

// Create 创建一个新的应用服务记录。
//
// 参数:
//   - st: 服务类型。
//   - d: 应用服务数据。
//
// 返回值:
//   - string: 创建的记录ID。
//   - error: 如果创建失败，返回相应的错误。
func (s *ApplicationServiceImpl) Create(st string, d map[string]any) (string, error) {
	// id
	id, err := utility.GenerateKsuid()
	if err != nil {
		return "", err
	}
	d["id"] = id

	time_string := time.Now().Format("2006-01-02 15:04:05")

	// time
	d["event_time"] = time_string

	// state
	state := map[string]any{
		"created_at": time_string,
		"status":     "active",
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	d["data_state"] = string(stateJson)

	err = s.repo.Create(st, d)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetMany 获取多个应用服务记录。
//
// 参数:
//   - st: 服务类型。
//   - f: 查询过滤条件。
//   - l: 限制条件。
//
// 返回值:
//   - []map[string]interface{}: 应用服务数据列表。
//   - error: 如果获取失败，返回相应的错误。
func (s *ApplicationServiceImpl) GetMany(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	result, err := s.repo.Get(st, c, f, l)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []map[string]interface{}{}, nil
	}
	return result, nil
}

// Get 获取单个应用服务记录。
//
// 参数:
//   - st: 服务类型。
//   - f: 查询过滤条件。
//   - l: 限制条件。
//
// 返回值:
//   - map[string]interface{}: 应用服务数据。
//   - error: 如果获取失败，返回相应的错误。
func (s *ApplicationServiceImpl) Get(st string, f [][]string, l string) (map[string]any, error) {
	data, err := s.repo.Get(st, nil, f, l+" limit 1")
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return map[string]any{}, fmt.Errorf("记录不存在")
	}
	return data[0], nil
}

// Update 更新应用服务记录。
//
// 参数:
//   - st: schema and table。
//   - d: 更新的数据。
//   - w: 更新条件字符串
//   - deprecated: 是否标记数据弃用。
//
// 返回值:
//   - error: 如果更新失败，返回相应的错误。
func (s *ApplicationServiceImpl) Update(st string, d map[string]any, w string, deprecated bool) error {
	id, ok := d["id"].(string)
	if !ok {
		return fmt.Errorf("缺少ID")
	}

	existingData, err := s.repo.Get(st, []string{"data_state"}, [][]string{{"equal", "id", id}}, "")
	if err != nil {
		return err
	}
	if len(existingData) == 0 {
		return fmt.Errorf("记录不存在")
	}

	var state map[string]any
	err = json.Unmarshal([]byte(existingData[0]["data_state"].(string)), &state)
	if err != nil {
		return err
	}

	state["updated_at"] = time.Now().Format("2006-01-02 15:04:05")
	if deprecated {
		state["deprecated"] = true
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return err
	}
	d["data_state"] = string(stateJson)

	return s.repo.Update(st, d, w)
}

// Remove 移除应用服务记录。
//
// 参数:
//   - st: 服务类型。
//   - w: 移除条件。
//
// 返回值:
//   - error: 如果移除失败，返回相应的错误。
func (s *ApplicationServiceImpl) Remove(st string, w string) error {
	return s.repo.Remove(st, w)
}
