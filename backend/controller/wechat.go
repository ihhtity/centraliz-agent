package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/config"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"centraliz-backend/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`                          // 微信授权码
	Platform string `json:"platform" binding:"required,oneof=miniprogram mp"` // 平台类型：miniprogram小程序 mp公众号
	MerchsID *int32 `json:"merchsId"`                                         // 商家ID（可选）
}

type WechatUserInfoRequest struct {
	OpenID   string `json:"openId" binding:"required"`                        // 微信OpenID
	Platform string `json:"platform" binding:"required,oneof=miniprogram mp"` // 平台类型
}

type WechatBindUserRequest struct {
	OpenID   string `json:"openId" binding:"required"`   // 微信OpenID
	UsersID  int32  `json:"usersId" binding:"required"`  // 用户ID
	MerchsID int32  `json:"merchsId" binding:"required"` // 商家ID
}

type WechatUserResponse struct {
	ID        uint32  `json:"id"`
	OpenID    string  `json:"openId"`
	UnionID   *string `json:"unionId"`
	Nickname  *string `json:"nickname"`
	Avatar    *string `json:"avatar"`
	Gender    *int    `json:"gender"`
	Country   *string `json:"country"`
	Province  *string `json:"province"`
	City      *string `json:"city"`
	Platform  string  `json:"platform"`
	MerchsID  int32   `json:"merchsId"`
	UsersID   *int32  `json:"usersId"`
	Status    *string `json:"status"`
	CreatedAt string  `json:"createdAt"`
}

// WechatLogin 微信登录（小程序/公众号）
func WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 调用微信API获取用户信息
	// 小程序：使用 code 换取 openid 和 session_key
	// 公众号：使用 code 换取 openid 和 access_token
	wechatInfo, err := getWechatUserInfo(req.Code, req.Platform)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 直接返回获取到的微信用户信息给前端
	response.SuccessWithMsg(c, "获取微信用户信息成功", gin.H{
		"openId":      wechatInfo.OpenID,
		"unionId":     wechatInfo.UnionID,
		"platform":    req.Platform,
		"sessionKey":  wechatInfo.SessionKey,
		"accessToken": wechatInfo.AccessToken,
	})
}

// GetWechatUserInfo 获取微信用户信息
func GetWechatUserInfo(c *gin.Context) {
	var req WechatUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	var wechatUser model.WechatUser
	if err := db.DB.Where("open_id = ? AND platform = ?", req.OpenID, req.Platform).First(&wechatUser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}

	response.SuccessWithMsg(c, "获取成功", buildWechatUserResponse(wechatUser))
}

// BindWechatUser 绑定微信用户到系统用户
func BindWechatUser(c *gin.Context) {
	var req WechatBindUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 查询微信用户
	var wechatUser model.WechatUser
	if err := db.DB.Where("open_id = ?", req.OpenID).First(&wechatUser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}

	// 查询系统用户
	var user model.User
	if err := db.DB.Where("id = ? AND merchs_id = ?", req.UsersID, req.MerchsID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	// 检查是否已绑定
	if wechatUser.UsersID != nil {
		response.Fail(c, 400, "该微信账号已绑定其他用户")
		return
	}

	// 绑定用户
	wechatUser.UsersID = &req.UsersID
	wechatUser.MerchsID = req.MerchsID
	if err := db.DB.Save(&wechatUser).Error; err != nil {
		response.Error(c, "绑定失败")
		return
	}

	response.SuccessWithMsg(c, "绑定成功", buildWechatUserResponse(wechatUser))
}

// UnbindWechatUser 解绑微信用户
func UnbindWechatUser(c *gin.Context) {
	openID := c.Query("openId")
	if openID == "" {
		response.Fail(c, 400, "缺少openId参数")
		return
	}

	var wechatUser model.WechatUser
	if err := db.DB.Where("open_id = ?", openID).First(&wechatUser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}

	if wechatUser.UsersID == nil {
		response.Fail(c, 400, "该微信账号未绑定用户")
		return
	}

	wechatUser.UsersID = nil
	if err := db.DB.Save(&wechatUser).Error; err != nil {
		response.Error(c, "解绑失败")
		return
	}

	response.SuccessWithMsg(c, "解绑成功", nil)
}

// UpdateWechatUserInfo 更新微信用户信息
func UpdateWechatUserInfo(c *gin.Context) {
	type UpdateRequest struct {
		OpenID   string  `json:"openId" binding:"required"`
		Nickname *string `json:"nickname"`
		Avatar   *string `json:"avatar"`
		Gender   *int    `json:"gender"`
		Country  *string `json:"country"`
		Province *string `json:"province"`
		City     *string `json:"city"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	var wechatUser model.WechatUser
	if err := db.DB.Where("open_id = ?", req.OpenID).First(&wechatUser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}

	// 更新字段
	if req.Nickname != nil {
		wechatUser.Nickname = req.Nickname
	}
	if req.Avatar != nil {
		wechatUser.Avatar = req.Avatar
	}
	if req.Gender != nil {
		wechatUser.Gender = req.Gender
	}
	if req.Country != nil {
		wechatUser.Country = req.Country
	}
	if req.Province != nil {
		wechatUser.Province = req.Province
	}
	if req.City != nil {
		wechatUser.City = req.City
	}

	if err := db.DB.Save(&wechatUser).Error; err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", buildWechatUserResponse(wechatUser))
}

// WechatInfo 微信API返回的用户信息
type WechatInfo struct {
	OpenID       string  `json:"openid"`
	UnionID      *string `json:"unionid"`
	SessionKey   *string `json:"session_key"`
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
	ExpiresIn    *int    `json:"expires_in"`
}

// MiniprogramSessionResponse 小程序登录响应
type MiniprogramSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// MPOAuthResponse 公众号OAuth响应
type MPOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid,omitempty"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// getWechatUserInfo 调用微信API获取用户信息
func getWechatUserInfo(code, platform string) (*WechatInfo, error) {
	if config.AppConfig.Wechat == nil {
		return nil, errors.New("微信配置未设置")
	}

	switch platform {
	case "miniprogram":
		return getMiniprogramUserInfo(code)
	case "mp":
		return getMPUserInfo(code)
	default:
		return nil, errors.New("不支持的平台类型")
	}
}

// getMiniprogramUserInfo 获取小程序用户信息
// 参考文档：https://developers.weixin.qq.com/miniprogram/dev/server/API/user-login/api_code2session.html
func getMiniprogramUserInfo(code string) (*WechatInfo, error) {
	wechatConfig := config.AppConfig.Wechat.Miniprogram
	if wechatConfig.AppID == "" || wechatConfig.AppSecret == "" {
		return nil, errors.New("小程序配置未设置（AppID或AppSecret为空）")
	}

	if code == "" {
		return nil, errors.New("code不能为空")
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wechatConfig.AppID, wechatConfig.AppSecret, code)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result MiniprogramSessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// 返回原始响应体帮助调试
		return nil, fmt.Errorf("解析响应失败，原始响应: %s", string(body))
	}

	if result.ErrCode != 0 {
		// 根据微信官方错误码提供更明确的错误提示
		errorMsg := fmt.Sprintf("微信API错误: %s (code: %d)", result.ErrMsg, result.ErrCode)
		switch result.ErrCode {
		case 40029:
			errorMsg = "微信授权失败：code无效或已过期，请重新获取code后重试"
		case 40013:
			errorMsg = "微信授权失败：AppID无效，请检查配置"
		case 40001:
			errorMsg = "微信授权失败：AppSecret无效，请检查配置"
		case 40163:
			errorMsg = "微信授权失败：code已被使用，请重新获取code"
		}
		return nil, errors.New(errorMsg)
	}

	wechatInfo := &WechatInfo{
		OpenID:     result.OpenID,
		SessionKey: utils.StringPtr(result.SessionKey),
	}

	if result.UnionID != "" {
		wechatInfo.UnionID = utils.StringPtr(result.UnionID)
	}

	return wechatInfo, nil
}

// getMPUserInfo 获取公众号用户信息
// 参考文档：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
func getMPUserInfo(code string) (*WechatInfo, error) {
	wechatConfig := config.AppConfig.Wechat.MP
	if wechatConfig.AppID == "" || wechatConfig.AppSecret == "" {
		return nil, errors.New("公众号配置未设置（AppID或AppSecret为空）")
	}

	if code == "" {
		return nil, errors.New("code不能为空")
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wechatConfig.AppID, wechatConfig.AppSecret, code)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result MPOAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败，原始响应: %s", string(body))
	}

	if result.ErrCode != 0 {
		errorMsg := fmt.Sprintf("微信API错误: %s (code: %d)", result.ErrMsg, result.ErrCode)
		switch result.ErrCode {
		case 40029:
			errorMsg = "微信授权失败：code无效或已过期，请重新获取code后重试"
		case 40013:
			errorMsg = "微信授权失败：AppID无效，请检查配置"
		case 40001:
			errorMsg = "微信授权失败：AppSecret无效，请检查配置"
		case 40163:
			errorMsg = "微信授权失败：code已被使用，请重新获取code"
		}
		return nil, errors.New(errorMsg)
	}

	wechatInfo := &WechatInfo{
		OpenID:       result.OpenID,
		AccessToken:  utils.StringPtr(result.AccessToken),
		RefreshToken: utils.StringPtr(result.RefreshToken),
		ExpiresIn:    &result.ExpiresIn,
	}

	if result.UnionID != "" {
		wechatInfo.UnionID = utils.StringPtr(result.UnionID)
	}

	return wechatInfo, nil
}

// refreshAccessToken 刷新公众号访问令牌
func refreshAccessToken(refreshToken string) (*WechatInfo, error) {
	wechatConfig := config.AppConfig.Wechat.MP
	if wechatConfig.AppID == "" || wechatConfig.AppSecret == "" {
		return nil, errors.New("公众号配置未设置")
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s",
		wechatConfig.AppID, refreshToken)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result MPOAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %s (code: %d)", result.ErrMsg, result.ErrCode)
	}

	return &WechatInfo{
		OpenID:       result.OpenID,
		AccessToken:  utils.StringPtr(result.AccessToken),
		RefreshToken: utils.StringPtr(result.RefreshToken),
		ExpiresIn:    &result.ExpiresIn,
	}, nil
}

// getMPUserInfoByAccessToken 使用access_token获取公众号用户信息
func getMPUserInfoByAccessToken(accessToken, openID string) (*WechatInfo, error) {
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		url.QueryEscape(accessToken), url.QueryEscape(openID))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		OpenID     string `json:"openid"`
		Nickname   string `json:"nickname"`
		Sex        int    `json:"sex"`
		Province   string `json:"province"`
		City       string `json:"city"`
		Country    string `json:"country"`
		HeadImgURL string `json:"headimgurl"`
		UnionID    string `json:"unionid,omitempty"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %s (code: %d)", result.ErrMsg, result.ErrCode)
	}

	wechatInfo := &WechatInfo{
		OpenID:     result.OpenID,
		SessionKey: utils.StringPtr(""),
	}

	if result.UnionID != "" {
		wechatInfo.UnionID = utils.StringPtr(result.UnionID)
	}

	return wechatInfo, nil
}

// buildWechatUserResponse 构建微信用户响应
func buildWechatUserResponse(user model.WechatUser) WechatUserResponse {
	resp := WechatUserResponse{
		ID:        user.ID,
		OpenID:    user.OpenID,
		UnionID:   user.UnionID,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Country:   user.Country,
		Province:  user.Province,
		City:      user.City,
		Platform:  user.Platform,
		MerchsID:  user.MerchsID,
		UsersID:   user.UsersID,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	return resp
}


