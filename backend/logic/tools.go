package logic

import (
	"bytes"
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]interface{}
	Execute(params map[string]interface{}) (*ToolResult, error)
}

type ToolRegistry struct {
	tools map[string]Tool
}

var toolRegistry *ToolRegistry

func InitToolRegistry() {
	toolRegistry = &ToolRegistry{
		tools: make(map[string]Tool),
	}

	toolRegistry.Register(&SearchOrderTool{})
	toolRegistry.Register(&SearchUserTool{})
	toolRegistry.Register(&SearchDeviceTool{})
	toolRegistry.Register(&SearchRoomTool{})
	toolRegistry.Register(&APITesterTool{})
	toolRegistry.Register(&SQLQueryTool{})
}

func GetToolRegistry() *ToolRegistry {
	if toolRegistry == nil {
		InitToolRegistry()
	}
	return toolRegistry
}

func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) Tool {
	return r.tools[name]
}

func (r *ToolRegistry) List() []ToolDefinition {
	var definitions []ToolDefinition
	for _, tool := range r.tools {
		definitions = append(definitions, ToolDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  tool.Parameters(),
		})
	}
	return definitions
}

func (r *ToolRegistry) Execute(name string, params map[string]interface{}) (*ToolResult, error) {
	tool := r.tools[name]
	if tool == nil {
		return nil, fmt.Errorf("tool %s not found", name)
	}
	return tool.Execute(params)
}

type SearchOrderTool struct{}

func (t *SearchOrderTool) Name() string { return "search_order" }

func (t *SearchOrderTool) Description() string {
	return "搜索订单信息，支持按订单号、用户ID、状态等条件查询"
}

func (t *SearchOrderTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"order_no": map[string]interface{}{
				"type":        "string",
				"description": "订单号",
			},
			"user_id": map[string]interface{}{
				"type":        "integer",
				"description": "用户ID",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "订单状态",
				"enum":        []string{"pending", "paid", "completed", "refunded", "cancelled"},
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "页码",
				"default":     1,
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "每页数量",
				"default":     10,
			},
		},
		"required": []string{},
	}
}

func (t *SearchOrderTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	var orders []model.Order
	query := db.DB.Model(&model.Order{})

	if orderNo, ok := params["order_no"].(string); ok && orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	if userID, ok := params["user_id"].(float64); ok && userID > 0 {
		query = query.Where("user_id = ?", int(userID))
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	page := 1
	if p, ok := params["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	pageSize := 10
	if ps, ok := params["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("created_at DESC")

	var total int64
	db.DB.Model(&model.Order{}).Count(&total)

	if err := query.Find(&orders).Error; err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"orders":    orders,
	}, "", "  ")

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(data),
	}, nil
}

type SearchUserTool struct{}

func (t *SearchUserTool) Name() string { return "search_user" }

func (t *SearchUserTool) Description() string {
	return "搜索用户信息，支持按用户ID、手机号、昵称等条件查询"
}

func (t *SearchUserTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"user_id": map[string]interface{}{
				"type":        "integer",
				"description": "用户ID",
			},
			"phone": map[string]interface{}{
				"type":        "string",
				"description": "手机号",
			},
			"nickname": map[string]interface{}{
				"type":        "string",
				"description": "昵称",
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "页码",
				"default":     1,
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "每页数量",
				"default":     10,
			},
		},
		"required": []string{},
	}
}

func (t *SearchUserTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	var users []model.User
	query := db.DB.Model(&model.User{})

	if userID, ok := params["user_id"].(float64); ok && userID > 0 {
		query = query.Where("id = ?", int(userID))
	}

	if phone, ok := params["phone"].(string); ok && phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	if nickname, ok := params["nickname"].(string); ok && nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	page := 1
	if p, ok := params["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	pageSize := 10
	if ps, ok := params["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("created_at DESC")

	var total int64
	db.DB.Model(&model.User{}).Count(&total)

	if err := query.Find(&users).Error; err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"users":     users,
	}, "", "  ")

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(data),
	}, nil
}

type SearchDeviceTool struct{}

func (t *SearchDeviceTool) Name() string { return "search_device" }

func (t *SearchDeviceTool) Description() string {
	return "搜索设备信息，支持按设备ID、设备编号、状态等条件查询"
}

func (t *SearchDeviceTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"device_id": map[string]interface{}{
				"type":        "integer",
				"description": "设备ID",
			},
			"device_no": map[string]interface{}{
				"type":        "string",
				"description": "设备编号",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "设备状态",
				"enum":        []string{"online", "offline", "error", "maintenance"},
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "页码",
				"default":     1,
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "每页数量",
				"default":     10,
			},
		},
		"required": []string{},
	}
}

func (t *SearchDeviceTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	var devices []model.Device
	query := db.DB.Model(&model.Device{})

	if deviceID, ok := params["device_id"].(float64); ok && deviceID > 0 {
		query = query.Where("id = ?", int(deviceID))
	}

	if deviceNo, ok := params["device_no"].(string); ok && deviceNo != "" {
		query = query.Where("device_no LIKE ?", "%"+deviceNo+"%")
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	page := 1
	if p, ok := params["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	pageSize := 10
	if ps, ok := params["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("created_at DESC")

	var total int64
	db.DB.Model(&model.Device{}).Count(&total)

	if err := query.Find(&devices).Error; err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"devices":   devices,
	}, "", "  ")

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(data),
	}, nil
}

type SearchRoomTool struct{}

func (t *SearchRoomTool) Name() string { return "search_room" }

func (t *SearchRoomTool) Description() string {
	return "搜索房间/柜子信息，支持按房间ID、房间号、状态等条件查询"
}

func (t *SearchRoomTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"room_id": map[string]interface{}{
				"type":        "integer",
				"description": "房间ID",
			},
			"room_no": map[string]interface{}{
				"type":        "string",
				"description": "房间号",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "房间状态",
				"enum":        []string{"empty", "occupied", "maintenance"},
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "页码",
				"default":     1,
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "每页数量",
				"default":     10,
			},
		},
		"required": []string{},
	}
}

func (t *SearchRoomTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	var rooms []model.Room
	query := db.DB.Model(&model.Room{})

	if roomID, ok := params["room_id"].(float64); ok && roomID > 0 {
		query = query.Where("id = ?", int(roomID))
	}

	if roomNo, ok := params["room_no"].(string); ok && roomNo != "" {
		query = query.Where("room_no LIKE ?", "%"+roomNo+"%")
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	page := 1
	if p, ok := params["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	pageSize := 10
	if ps, ok := params["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("created_at DESC")

	var total int64
	db.DB.Model(&model.Room{}).Count(&total)

	if err := query.Find(&rooms).Error; err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"rooms":     rooms,
	}, "", "  ")

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(data),
	}, nil
}

type APITesterTool struct{}

func (t *APITesterTool) Name() string { return "api_tester" }

func (t *APITesterTool) Description() string {
	return "测试API接口，发送HTTP请求并返回响应结果"
}

func (t *APITesterTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"method": map[string]interface{}{
				"type":        "string",
				"description": "HTTP方法",
				"enum":        []string{"GET", "POST", "PUT", "DELETE"},
				"default":     "GET",
			},
			"url": map[string]interface{}{
				"type":        "string",
				"description": "API URL",
			},
			"headers": map[string]interface{}{
				"type":        "object",
				"description": "请求头",
			},
			"body": map[string]interface{}{
				"type":        "object",
				"description": "请求体",
			},
		},
		"required": []string{"url"},
	}
}

func (t *APITesterTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	method := "GET"
	if m, ok := params["method"].(string); ok && m != "" {
		method = strings.ToUpper(m)
	}

	url, ok := params["url"].(string)
	if !ok || url == "" {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    "url is required",
		}, nil
	}

	var body []byte
	if b, ok := params["body"]; ok && b != nil {
		var err error
		body, err = json.Marshal(b)
		if err != nil {
			return &ToolResult{
				ToolName: t.Name(),
				Success:  false,
				Error:    fmt.Sprintf("marshal body error: %v", err),
			}, nil
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    fmt.Sprintf("create request error: %v", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/json")

	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    fmt.Sprintf("http request error: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output: fmt.Sprintf("Status: %s\nHeaders: %v\nBody: %s",
			resp.Status,
			resp.Header,
			string(respBody)),
	}, nil
}

type SQLQueryTool struct{}

func (t *SQLQueryTool) Name() string { return "sql_query" }

func (t *SQLQueryTool) Description() string {
	return "执行SQL查询，返回查询结果（仅支持SELECT查询）"
}

func (t *SQLQueryTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"sql": map[string]interface{}{
				"type":        "string",
				"description": "SQL查询语句（仅支持SELECT）",
			},
		},
		"required": []string{"sql"},
	}
}

func (t *SQLQueryTool) Execute(params map[string]interface{}) (*ToolResult, error) {
	sql, ok := params["sql"].(string)
	if !ok || sql == "" {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    "sql is required",
		}, nil
	}

	sql = strings.TrimSpace(strings.ToUpper(sql))
	if !strings.HasPrefix(sql, "SELECT") {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    "仅支持SELECT查询",
		}, nil
	}

	rows, err := db.DB.Raw(sql).Rows()
	if err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return &ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Error:    err.Error(),
		}, nil
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"columns": columns,
		"count":   len(results),
		"results": results,
	}, "", "  ")

	return &ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(data),
	}, nil
}
