package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ToolCall struct {
	ID         string                 `json:"id,omitempty"`
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type ChatMessage struct {
	ID          string       `json:"id"`
	Role        string       `json:"role"`
	Content     string       `json:"content"`
	Thought     string       `json:"thought,omitempty"`
	ToolCalls   []ToolCall   `json:"tool_calls,omitempty"`
	ToolCallID  string       `json:"tool_call_id,omitempty"`
	ToolResults []ToolResult `json:"tool_results,omitempty"`
	CreatedAt   string       `json:"created_at"`
}

type ChatSession struct {
	ID              string        `json:"id"`
	UserID          uint32        `json:"user_id"`
	UserName        string        `json:"user_name"`
	Title           string        `json:"title"`
	Messages        []ChatMessage `json:"messages"`
	RequiresConfirm bool          `json:"requires_confirm"`
	AffectedFiles   []string      `json:"affected_files"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
}

type ToolResult struct {
	ToolName string
	Success  bool
	Output   string
	Error    string
}

func CreateSession(userID uint32, userName string, title string) (*ChatSession, error) {
	if title == "" {
		title = "新会话"
	}

	now := time.Now()
	session := &model.AssistantSession{
		UserID:          userID,
		UserName:        userName,
		Title:           title,
		RequiresConfirm: false,
		AffectedFiles:   "[]",
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}

	if err := mysql.CreateAssistantSession(session); err != nil {
		return nil, err
	}

	return convertToChatSession(session), nil
}

func GetSessions(userID uint32) ([]ChatSession, error) {
	sessions, err := mysql.GetAssistantSessionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var chatSessions []ChatSession
	for _, session := range sessions {
		chatSessions = append(chatSessions, *convertToChatSession(&session))
	}

	return chatSessions, nil
}

func GetSessionDetail(sessionID uint32, userID uint32) (*ChatSession, error) {
	session, err := mysql.GetAssistantSessionDetail(sessionID)
	if err != nil {
		return nil, err
	}

	if session.UserID != userID {
		return nil, fmt.Errorf("无权访问此会话")
	}

	return convertToChatSession(session), nil
}

func DeleteSession(sessionID uint32, userID uint32) error {
	session, err := mysql.GetAssistantSessionDetail(sessionID)
	if err != nil {
		return err
	}

	if session.UserID != userID {
		return fmt.Errorf("无权删除此会话")
	}

	return mysql.DeleteAssistantSession(sessionID)
}

func ChatWithAssistant(sessionID uint32, userID uint32, userName string, message string) (*ChatSession, error) {
	var session *model.AssistantSession
	var err error

	if sessionID > 0 {
		session, err = mysql.GetAssistantSessionDetail(sessionID)
		if err != nil {
			return nil, err
		}

		if session.UserID != userID {
			return nil, fmt.Errorf("无权访问此会话")
		}
	} else {
		now := time.Now()
		session = &model.AssistantSession{
			UserID:          userID,
			UserName:        userName,
			Title:           truncateMessage(message, 20),
			RequiresConfirm: false,
			AffectedFiles:   "[]",
			CreatedAt:       &now,
			UpdatedAt:       &now,
		}

		if err := mysql.CreateAssistantSession(session); err != nil {
			return nil, err
		}
	}

	now := time.Now()

	userMsg := &model.AssistantMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   message,
		CreatedAt: &now,
	}
	if err := mysql.CreateAssistantMessage(userMsg); err != nil {
		return nil, err
	}

	chatSession := convertToChatSession(session)
	chatSession.Messages = append(chatSession.Messages, ChatMessage{
		ID:        fmt.Sprintf("%d", userMsg.ID),
		Role:      "user",
		Content:   message,
		CreatedAt: now.Format(time.RFC3339),
	})

	llmClient := GetLLMClient()
	tools := GetToolRegistry().List()

	systemPrompt := `你是一位资深的全栈开发工程师，精通 Go/Gin 后端、UniApp 前端、MySQL 数据库。
你可以使用以下工具：
- search_order: 搜索订单信息
- search_user: 搜索用户信息
- search_device: 搜索设备信息
- search_room: 搜索房间/柜子信息
- api_tester: 测试API接口
- sql_query: 执行SQL查询（仅支持SELECT）

当用户询问数据相关问题时，优先使用工具查询真实数据。
回答要专业、简洁，使用markdown格式。`

	historyMessages := []ChatMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range chatSession.Messages {
		historyMessages = append(historyMessages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	historyMessages = append(historyMessages, ChatMessage{
		Role:    "user",
		Content: message,
	})

	content := ""
	thought := ""
	toolCalls := []ToolCall{}
	requiresConfirm := false
	affectedFiles := []string{}

	maxReActSteps := 5
	currentStep := 0

	for currentStep < maxReActSteps {
		result, err := llmClient.Chat(historyMessages, tools)
		if err != nil {
			return nil, err
		}

		if result.Content != "" {
			content += result.Content
		}

		if len(result.ToolCalls) > 0 {
			if thought == "" {
				thought = "我需要调用工具来完成这个任务"
			}

			toolCalls = append(toolCalls, result.ToolCalls...)

			for _, toolCall := range result.ToolCalls {
				toolResult, execErr := GetToolRegistry().Execute(toolCall.ToolName, toolCall.Parameters)
				if execErr != nil {
					content = fmt.Sprintf("执行工具 %s 时出错: %v\n\n%s", toolCall.ToolName, execErr, content)
					continue
				}

				if !toolResult.Success {
					content = fmt.Sprintf("工具 %s 执行失败: %s\n\n%s", toolCall.ToolName, toolResult.Error, content)
					continue
				}

				historyMessages = append(historyMessages, ChatMessage{
					Role:      "assistant",
					Content:   "",
					ToolCalls: []ToolCall{toolCall},
				})
				historyMessages = append(historyMessages, ChatMessage{
					Role:       "tool",
					Content:    toolResult.Output,
					ToolCallID: toolCall.ID,
				})
			}

			currentStep++
			continue
		}

		break
	}

	toolCallsJSON, _ := json.Marshal(toolCalls)
	affectedFilesJSON, _ := json.Marshal(affectedFiles)

	assistantMsg := &model.AssistantMessage{
		SessionID: session.ID,
		Role:      "assistant",
		Content:   content,
		Thought:   &thought,
		ToolCalls: string(toolCallsJSON),
		CreatedAt: &now,
	}
	if err := mysql.CreateAssistantMessage(assistantMsg); err != nil {
		return nil, err
	}

	session.RequiresConfirm = requiresConfirm
	session.AffectedFiles = string(affectedFilesJSON)

	if session.Title == "新会话" {
		session.Title = truncateMessage(message, 20)
	}

	session.UpdatedAt = &now
	if err := mysql.UpdateAssistantSession(session); err != nil {
		return nil, err
	}

	updatedSession, _ := mysql.GetAssistantSessionDetail(session.ID)
	return convertToChatSession(updatedSession), nil
}

func ConfirmAction(sessionID uint32, userID uint32, confirm bool) (*ChatSession, error) {
	session, err := mysql.GetAssistantSessionDetail(sessionID)
	if err != nil {
		return nil, err
	}

	if session.UserID != userID {
		return nil, fmt.Errorf("无权访问此会话")
	}

	now := time.Now()

	var content string
	var thought string
	if confirm {
		content = executeConfirmedAction(session)
		thought = "用户已确认执行，正在执行操作..."
	} else {
		content = "操作已取消。您可以继续与我对话，提出其他需求。"
		thought = "用户已取消操作"
	}

	confirmMsg := &model.AssistantMessage{
		SessionID: session.ID,
		Role:      "assistant",
		Content:   content,
		Thought:   &thought,
		CreatedAt: &now,
	}
	if err := mysql.CreateAssistantMessage(confirmMsg); err != nil {
		return nil, err
	}

	session.RequiresConfirm = false
	session.AffectedFiles = "[]"
	session.UpdatedAt = &now
	if err := mysql.UpdateAssistantSession(session); err != nil {
		return nil, err
	}

	updatedSession, _ := mysql.GetAssistantSessionDetail(session.ID)
	return convertToChatSession(updatedSession), nil
}

func executeConfirmedAction(session *model.AssistantSession) string {
	var affectedFiles []string
	if session.AffectedFiles != "" {
		json.Unmarshal([]byte(session.AffectedFiles), &affectedFiles)
	}

	result := "操作已执行成功！\n\n执行结果：\n"
	for _, file := range affectedFiles {
		result += fmt.Sprintf("- ✅ %s\n", file)
	}
	result += "\n如需进一步操作，请继续告诉我。"

	return result
}

func ChatStreamWithAssistant(sessionID uint32, userID uint32, userName string, message string, handler func(content string, toolCalls []ToolCall, done bool) error) error {
	var session *model.AssistantSession
	var err error

	if sessionID > 0 {
		session, err = mysql.GetAssistantSessionDetail(sessionID)
		if err != nil {
			return err
		}

		if session.UserID != userID {
			return fmt.Errorf("无权访问此会话")
		}
	} else {
		now := time.Now()
		session = &model.AssistantSession{
			UserID:          userID,
			UserName:        userName,
			Title:           truncateMessage(message, 20),
			RequiresConfirm: false,
			AffectedFiles:   "[]",
			CreatedAt:       &now,
			UpdatedAt:       &now,
		}

		if err := mysql.CreateAssistantSession(session); err != nil {
			return err
		}
	}

	now := time.Now()

	userMsg := &model.AssistantMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   message,
		CreatedAt: &now,
	}
	if err := mysql.CreateAssistantMessage(userMsg); err != nil {
		return err
	}

	if session.Title == "新会话" {
		session.Title = truncateMessage(message, 20)
	}

	chatSession := convertToChatSession(session)
	chatSession.Messages = append(chatSession.Messages, ChatMessage{
		ID:        fmt.Sprintf("%d", userMsg.ID),
		Role:      "user",
		Content:   message,
		CreatedAt: now.Format(time.RFC3339),
	})

	toolRegistry := GetToolRegistry()
	tools := toolRegistry.List()

	var llmMessages []ChatMessage
	for _, msg := range chatSession.Messages {
		llmMessages = append(llmMessages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	systemPrompt := `你是一位资深的全栈开发工程师，精通 Go/Gin 后端、UniApp 前端、MySQL 数据库。
你可以使用以下工具：
- search_order: 搜索订单信息
- search_user: 搜索用户信息
- search_device: 搜索设备信息
- search_room: 搜索房间/柜子信息
- api_tester: 测试API接口
- sql_query: 执行SQL查询（仅支持SELECT）

当用户询问数据相关问题时，优先使用工具查询真实数据。
回答要专业、简洁，使用markdown格式。`

	llmMessages = append([]ChatMessage{{Role: "system", Content: systemPrompt}}, llmMessages...)

	llmClient := GetLLMClient()

	var fullContent strings.Builder

	err = llmClient.ChatStream(llmMessages, tools, func(content string, toolCalls []ToolCall, done bool) error {
		if content != "" {
			fullContent.WriteString(content)
			if err := handler(content, nil, false); err != nil {
				return err
			}
		}

		if len(toolCalls) > 0 {
			for _, toolCall := range toolCalls {
				result, _ := toolRegistry.Execute(toolCall.ToolName, toolCall.Parameters)
				if result != nil {
					toolResultContent := fmt.Sprintf("工具 %s 执行结果:\n%s", toolCall.ToolName, result.Output)
					fullContent.WriteString(toolResultContent)
					if err := handler(toolResultContent, nil, false); err != nil {
						return err
					}
				}
			}
		}

		if done {
			return handler("", nil, true)
		}

		return nil
	})

	if err != nil {
		return err
	}

	finalContent := fullContent.String()

	toolCallsJSON, _ := json.Marshal([]ToolCall{})

	assistantMsg := &model.AssistantMessage{
		SessionID:   session.ID,
		Role:        "assistant",
		Content:     finalContent,
		Thought:     nil,
		ToolCalls:   string(toolCallsJSON),
		ToolResults: "[]",
		CreatedAt:   &now,
	}
	if err := mysql.CreateAssistantMessage(assistantMsg); err != nil {
		return err
	}

	session.RequiresConfirm = false
	session.AffectedFiles = "[]"
	session.UpdatedAt = &now
	if err := mysql.UpdateAssistantSession(session); err != nil {
		return err
	}

	return nil
}

func convertToChatSession(session *model.AssistantSession) *ChatSession {
	var messages []ChatMessage
	for _, msg := range session.Messages {
		var toolCalls []ToolCall
		if msg.ToolCalls != "" {
			json.Unmarshal([]byte(msg.ToolCalls), &toolCalls)
		}

		var toolResults []ToolResult
		if msg.ToolResults != "" {
			json.Unmarshal([]byte(msg.ToolResults), &toolResults)
		}

		var thought string
		if msg.Thought != nil {
			thought = *msg.Thought
		}

		var createdAt string
		if msg.CreatedAt != nil {
			createdAt = msg.CreatedAt.Format(time.RFC3339)
		}

		messages = append(messages, ChatMessage{
			ID:          fmt.Sprintf("%d", msg.ID),
			Role:        msg.Role,
			Content:     msg.Content,
			Thought:     thought,
			ToolCalls:   toolCalls,
			ToolResults: toolResults,
			CreatedAt:   createdAt,
		})
	}

	var affectedFiles []string
	if session.AffectedFiles != "" {
		json.Unmarshal([]byte(session.AffectedFiles), &affectedFiles)
	}

	var createdAt string
	if session.CreatedAt != nil {
		createdAt = session.CreatedAt.Format(time.RFC3339)
	}

	var updatedAt string
	if session.UpdatedAt != nil {
		updatedAt = session.UpdatedAt.Format(time.RFC3339)
	}

	return &ChatSession{
		ID:              fmt.Sprintf("%d", session.ID),
		UserID:          session.UserID,
		UserName:        session.UserName,
		Title:           session.Title,
		Messages:        messages,
		RequiresConfirm: session.RequiresConfirm,
		AffectedFiles:   affectedFiles,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func truncateMessage(message string, maxLen int) string {
	if len(message) <= maxLen {
		return message
	}
	return message[:maxLen] + "..."
}
