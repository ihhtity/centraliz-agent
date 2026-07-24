package controller

import (
	"centraliz-backend/logic"
	"centraliz-backend/middleware"
	"centraliz-backend/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type ConfirmRequest struct {
	SessionID string `json:"session_id"`
	Confirm   bool   `json:"confirm"`
}

func getCurrentUser(c *gin.Context) (uint32, string) {
	return middleware.GetUserID(c), middleware.GetUsername(c)
}

type CreateSessionRequest struct {
	Title string `json:"title"`
}

func AdminCreateSession(c *gin.Context) {
	userID, userName := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	session, err := logic.CreateSession(userID, userName, req.Title)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"session": session,
	})
}

func AdminGetSessions(c *gin.Context) {
	userID, _ := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	sessions, err := logic.GetSessions(userID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"sessions": sessions,
	})
}

func AdminGetSessionDetail(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的会话ID")
		return
	}

	userID, _ := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	session, err := logic.GetSessionDetail(uint32(sessionID), userID)
	if err != nil {
		if err.Error() == "无权访问此会话" {
			response.Fail(c, 403, err.Error())
		} else if err.Error() == "record not found" {
			response.Fail(c, 404, "会话不存在")
		} else {
			response.Fail(c, 500, err.Error())
		}
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"session": session,
	})
}

func AdminDeleteSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的会话ID")
		return
	}

	userID, _ := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	err = logic.DeleteSession(uint32(sessionID), userID)
	if err != nil {
		if err.Error() == "无权删除此会话" {
			response.Fail(c, 403, err.Error())
		} else if err.Error() == "record not found" {
			response.Fail(c, 404, "会话不存在")
		} else {
			response.Fail(c, 500, err.Error())
		}
		return
	}

	response.Success(c, "删除成功")
}

func AdminChatWithAssistant(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	userID, userName := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	sessionID := uint32(0)
	if req.SessionID != "" {
		id, err := strconv.ParseUint(req.SessionID, 10, 32)
		if err == nil {
			sessionID = uint32(id)
		}
	}

	session, err := logic.ChatWithAssistant(sessionID, userID, userName, req.Message)
	if err != nil {
		if err.Error() == "无权访问此会话" {
			response.Fail(c, 403, err.Error())
		} else if strings.Contains(err.Error(), "LLM API error") {
			response.Fail(c, 503, "AI服务暂时繁忙，请稍后重试")
		} else {
			response.Fail(c, 500, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"session": session,
		},
	})
}

func AdminChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	if req.Message == "" {
		response.Fail(c, 400, "消息内容不能为空")
		return
	}

	userID, userName := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	sessionID := uint32(0)
	if req.SessionID != "" {
		id, err := strconv.ParseUint(req.SessionID, 10, 32)
		if err == nil {
			sessionID = uint32(id)
		}
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	handler := func(content string, toolCalls []logic.ToolCall, done bool) error {
		if content != "" {
			c.SSEvent("message", gin.H{
				"type":    "content",
				"content": content,
			})
			c.Writer.Flush()
		}

		if len(toolCalls) > 0 {
			c.SSEvent("message", gin.H{
				"type":       "tool_calls",
				"tool_calls": toolCalls,
			})
			c.Writer.Flush()
		}

		if done {
			c.SSEvent("message", gin.H{
				"type": "done",
			})
			c.Writer.Flush()
		}

		return nil
	}

	err := logic.ChatStreamWithAssistant(sessionID, userID, userName, req.Message, handler)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "LLM API error") {
			errorMsg = "AI服务暂时繁忙，请稍后重试"
		}
		c.SSEvent("error", gin.H{
			"error": errorMsg,
		})
		c.Writer.Flush()
	}
}

func AdminConfirmAction(c *gin.Context) {
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	sessionID, err := strconv.ParseUint(req.SessionID, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的会话ID")
		return
	}

	userID, _ := getCurrentUser(c)
	if userID == 0 {
		response.Fail(c, 401, "未授权")
		return
	}

	session, err := logic.ConfirmAction(uint32(sessionID), userID, req.Confirm)
	if err != nil {
		if err.Error() == "无权访问此会话" {
			response.Fail(c, 403, err.Error())
		} else if err.Error() == "record not found" {
			response.Fail(c, 404, "会话不存在")
		} else {
			response.Fail(c, 500, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"session": session,
		},
	})
}
