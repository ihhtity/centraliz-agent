package model

import "time"

type AssistantSession struct {
	ID              uint32             `json:"id" gorm:"size:20;primaryKey;autoIncrement;comment:主键"`
	UserID          uint32             `json:"user_id" gorm:"size:10;column:user_id;not null;index;comment:用户ID"`
	UserName        string             `json:"user_name" gorm:"size:64;column:user_name;not null;comment:用户名"`
	Title           string             `json:"title" gorm:"size:255;not null;comment:会话标题"`
	Context         string             `json:"context" gorm:"column:context;type:text;comment:会话上下文JSON"`
	RequiresConfirm bool               `json:"requires_confirm" gorm:"column:requires_confirm;default:false;comment:是否需要确认"`
	AffectedFiles   string             `json:"affected_files" gorm:"column:affected_files;type:text;comment:影响文件列表JSON"`
	CreatedAt       *time.Time         `json:"created_at" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt       *time.Time         `json:"updated_at" gorm:"column:updated_at;precision:3;comment:更新时间"`
	Messages        []AssistantMessage `json:"messages" gorm:"foreignKey:SessionID"`
}

func (AssistantSession) TableName() string {
	return "assistant_sessions"
}

type AssistantMessage struct {
	ID          uint32     `json:"id" gorm:"size:20;primaryKey;autoIncrement;comment:主键"`
	SessionID   uint32     `json:"sessionId" gorm:"size:20;column:session_id;not null;index;comment:会话ID"`
	Role        string     `json:"role" gorm:"size:20;not null;comment:角色 user/assistant"`
	Content     string     `json:"content" gorm:"type:text;comment:消息内容"`
	Thought     *string    `json:"thought" gorm:"type:text;comment:思考过程"`
	ToolCalls   string     `json:"toolCalls" gorm:"column:tool_calls;type:text;comment:工具调用JSON"`
	ToolResults string     `json:"toolResults" gorm:"column:tool_results;type:text;comment:工具执行结果JSON"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
}

func (AssistantMessage) TableName() string {
	return "assistant_messages"
}
