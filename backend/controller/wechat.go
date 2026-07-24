package controller

import (
	"bytes"
	"centraliz-backend/model"
	"centraliz-backend/pkg/config"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/response"
	"centraliz-backend/pkg/utils"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 获取微信小程序配置
func getMiniProgramConfig() config.WechatMiniProgramConfig {
	if config.AppConfig.Wechat == nil {
		return config.WechatMiniProgramConfig{}
	}
	return config.AppConfig.Wechat.Miniprogram
}

// 获取微信公众号配置
func getMPConfig() config.WechatMPConfig {
	if config.AppConfig.Wechat == nil {
		return config.WechatMPConfig{}
	}
	return config.AppConfig.Wechat.MP
}

// WXJSConfigResponse JS-SDK配置响应
type WXJSConfigResponse struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

// WXAccessTokenResponse 微信access_token响应
type WXAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

var tokenMutex sync.Mutex

// CachedAccessToken 缓存的access_token结构
type CachedAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	CreatedAt   int64  `json:"created_at"`
}

// getAccessTokenFilePath 获取access_token缓存文件路径
func getAccessTokenFilePath() string {
	return filepath.Join("config", "access_token.json")
}

// getAccessTokenFromFile 从文件读取access_token并检查有效性
func getAccessTokenFromFile() (string, error) {
	filePath := getAccessTokenFilePath()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", fmt.Errorf("文件内容为空")
	}

	var cached CachedAccessToken
	if err := json.Unmarshal(data, &cached); err != nil {
		return "", err
	}

	if cached.AccessToken == "" {
		return "", fmt.Errorf("access_token为空")
	}

	now := time.Now().Unix()
	expireTime := cached.CreatedAt + int64(cached.ExpiresIn) - 300

	if now >= expireTime {
		return "", fmt.Errorf("access_token已过期")
	}

	return cached.AccessToken, nil
}

// saveAccessTokenToFile 将access_token保存到文件
func saveAccessTokenToFile(token string, expiresIn int) error {
	filePath := getAccessTokenFilePath()
	dir := filepath.Dir(filePath)

	fmt.Printf("保存access_token到文件: %s\n", filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return err
	}

	cached := CachedAccessToken{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		CreatedAt:   time.Now().Unix(),
	}

	data, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return err
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return err
	}

	fmt.Printf("access_token保存成功\n")
	return nil
}

// WXJSSDKConfig 获取微信JS-SDK配置
func WXJSSDKConfig(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		response.Fail(c, 400, "URL参数不能为空")
		return
	}

	mpConfig := getMPConfig()

	// 获取access_token
	accessToken, err := getWXMPGlobalAccessToken()
	if err != nil {
		response.Fail(c, 500, "获取微信access_token失败: "+err.Error())
		return
	}

	// 获取jsapi_ticket
	ticket, err := getWXJSAPITicket(accessToken)
	if err != nil {
		response.Fail(c, 500, "获取JSAPI ticket失败: "+err.Error())
		return
	}

	// 生成签名
	timestamp := time.Now().Unix()
	nonceStr := generateNonceStr()

	// 签名格式: jsapi_ticket=sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvk&noncestr=Wm3WZYTPz0wzccnW&timestamp=1414587457&url=http://mp.weixin.qq.com?params=value
	signatureStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonceStr, timestamp, url)
	signature := sha1Hash(signatureStr)

	response.SuccessWithMsg(c, "success", WXJSConfigResponse{
		AppID:     mpConfig.AppID,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: signature,
	})
}

// getWXMPGlobalAccessToken 获取公众号全局access_token
func getWXMPGlobalAccessToken() (string, error) {
	if token, err := getAccessTokenFromFile(); err == nil {
		fmt.Printf("从缓存文件获取access_token成功\n")
		return token, nil
	}

	fmt.Printf("缓存文件不存在或已过期，准备获取新的access_token\n")

	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if token, err := getAccessTokenFromFile(); err == nil {
		fmt.Printf("加锁后从缓存文件获取access_token成功\n")
		return token, nil
	}

	mpConfig := getMPConfig()
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		mpConfig.AppID, mpConfig.AppSecret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp WXAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("%s", tokenResp.ErrMsg)
	}

	if err := saveAccessTokenToFile(tokenResp.AccessToken, tokenResp.ExpiresIn); err != nil {
		fmt.Printf("保存access_token到文件失败: %v\n", err)
	}

	return tokenResp.AccessToken, nil
}

// getWXJSAPITicket 获取JSAPI ticket
func getWXJSAPITicket(accessToken string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return "", fmt.Errorf("%s", result["errmsg"])
	}

	ticket, ok := result["ticket"].(string)
	if !ok {
		return "", fmt.Errorf("获取ticket失败")
	}

	return ticket, nil
}

// generateNonceStr 生成随机字符串
func generateNonceStr() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)
	for i := range result {
		result[i] = chars[randomInt(0, len(chars)-1)]
	}
	return string(result)
}

// randomInt 生成随机整数
func randomInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}

// sha1Hash SHA1哈希
func sha1Hash(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// WXScanImage 图片扫码识别
func WXScanImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Fail(c, 400, "请选择图片")
		return
	}

	// 读取图片文件
	f, err := file.Open()
	if err != nil {
		response.Fail(c, 500, "读取图片失败")
		return
	}
	defer f.Close()

	// 获取access_token
	accessToken, err := getWXMPGlobalAccessToken()
	if err != nil {
		response.Fail(c, 500, "获取微信access_token失败")
		return
	}

	// 调用微信图片识别API
	result, err := scanQRCodeFromImage(accessToken, f)
	if err != nil {
		response.Fail(c, 500, "识别失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "识别成功", gin.H{
		"result": result,
	})
}

// scanQRCodeFromImage 调用微信二维码识别API
func scanQRCodeFromImage(accessToken string, imageData io.Reader) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cv/img/qrcode?access_token=%s", accessToken)

	// 创建multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("img", "qrcode.jpg")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, imageData)
	if err != nil {
		return "", err
	}

	writer.Close()

	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return "", fmt.Errorf("%s", result["errmsg"])
	}

	// 解析识别结果
	if codeResults, ok := result["code_results"].([]interface{}); ok && len(codeResults) > 0 {
		if firstResult, ok := codeResults[0].(map[string]interface{}); ok {
			if data, ok := firstResult["data"].(string); ok {
				return data, nil
			}
		}
	}

	return "", fmt.Errorf("未识别到二维码")
}

// MiniProgramUserInfo 微信小程序用户信息
type MiniProgramUserInfo struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
}

// WXSessionResponse 微信会话信息
type WXSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WXUserInfo 微信用户信息
type WXUserInfo struct {
	ID        uint32 `json:"id"`
	OpenID    string `json:"openid"`
	GOpenID   string `json:"gopenid"`
	UnionID   string `json:"unionid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"headimgurl"`
	Platform  string `json:"platform"`
	Gender    int    `json:"sex"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
}

// WXMPAccessToken 公众号access_token响应
type WXMPAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// WXMPUserInfo 公众号用户信息
type WXMPUserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
}

// WXLogin 微信登录（小程序/H5）
func WXLogin(c *gin.Context) {
	// WXLoginRequest 微信登录请求
	type WXLoginRequest struct {
		ID       uint32 `json:"id"`                          // 用户ID
		Code     string `json:"code" binding:"required"`     // 登录凭证code
		Platform string `json:"platform" binding:"required"` // 平台类型: miniprogram/mp
	}
	var req WXLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	var unionID, openID string
	var userInfo WXUserInfo

	if req.Platform == "miniprogram" {
		// 微信小程序登录
		session, err := getWXSession(req.Code)
		if err != nil {
			response.Fail(c, 401, "获取微信会话失败: "+err.Error())
			return
		}

		unionID = session.UnionID
		openID = session.OpenID

		userInfo = WXUserInfo{
			OpenID:  openID,
			UnionID: unionID,
		}
	} else if req.Platform == "mp" {
		// 微信公众号(H5)登录
		token, err := getWXMPAccessToken(req.Code)
		if err != nil {
			response.Fail(c, 401, "获取公众号token失败: "+err.Error())
			return
		}

		unionID = token.UnionID
		openID = token.OpenID

		// 获取用户信息
		mpUserInfo, err := getWXMPUserInfo(token.AccessToken, token.OpenID)
		if err != nil {
			response.Fail(c, 401, "获取公众号用户信息失败: "+err.Error())
			return
		}

		userInfo = WXUserInfo{
			GOpenID:   mpUserInfo.OpenID,
			UnionID:   mpUserInfo.UnionID,
			Nickname:  mpUserInfo.Nickname,
			AvatarURL: mpUserInfo.HeadImgURL,
			Gender:    mpUserInfo.Sex,
			Country:   mpUserInfo.Country,
			Province:  mpUserInfo.Province,
			City:      mpUserInfo.City,
			Language:  "",
		}
	} else {
		response.Fail(c, 400, "不支持的平台类型")
		return
	}

	// 通过UnionID查找或创建用户
	userInfo.ID = req.ID
	userInfo.Platform = req.Platform
	user, err := findOrCreateUserByUnionID(userInfo)
	if err != nil {
		response.Error(c, "创建或查找用户失败: "+err.Error())
		return
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(uint32(user.ID), user.Account, "user")
	if err != nil {
		response.Error(c, "生成token失败")
		return
	}

	response.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user":  user,
	})
}

// getWXSession 获取微信小程序会话
func getWXSession(code string) (*WXSessionResponse, error) {
	wxConfig := getMiniProgramConfig()
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wxConfig.AppID, wxConfig.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var session WXSessionResponse
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}

	if session.ErrCode != 0 {
		return nil, fmt.Errorf("%s", session.ErrMsg)
	}

	return &session, nil
}

// getWXMPAccessToken 获取公众号access_token
func getWXMPAccessToken(code string) (*WXMPAccessToken, error) {
	mpConfig := getMPConfig()
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		mpConfig.AppID, mpConfig.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token WXMPAccessToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	if token.ErrCode != 0 {
		return nil, fmt.Errorf("%s", token.ErrMsg)
	}

	return &token, nil
}

// getWXMPUserInfo 获取公众号用户信息
func getWXMPUserInfo(accessToken, openID string) (*WXMPUserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo WXMPUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	if userInfo.ErrCode != 0 {
		return nil, fmt.Errorf("%s", userInfo.ErrMsg)
	}

	return &userInfo, nil
}

// findOrCreateUserByUnionID 通过UnionID查找或创建用户
func findOrCreateUserByUnionID(userInfo WXUserInfo) (*model.User, error) {
	// 先通过UnionID查找用户
	var user model.User
	if userInfo.UnionID != "" || userInfo.ID != 0 {
		if err := db.DB.Where("union_id = ? OR id = ?", userInfo.UnionID, userInfo.ID).First(&user).Error; err == nil {
			// 找到用户，更新信息
			if userInfo.Platform != "mp" {
				user.OpenID = userInfo.OpenID
			} else {
				user.Name = userInfo.Nickname
				user.AvatarURL = &userInfo.AvatarURL
				user.GOpenID = userInfo.GOpenID
			}
			user.UnionID = userInfo.UnionID
			if err := db.DB.Save(&user).Error; err != nil {
				return nil, err
			}
			return &user, nil
		}
	}

	// 没有找到，创建新用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(utils.GenerateRandomPassword(12)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	username := userInfo.Nickname
	if username == "" {
		username = "wx" + userInfo.OpenID[:8]
	}

	newUser := model.User{
		Name:      username,
		Account:   username,
		Password:  string(hashedPassword),
		AvatarURL: &userInfo.AvatarURL,
		UnionID:   userInfo.UnionID,
		GOpenID:   userInfo.GOpenID,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

// WXMPLoginCallback 微信公众号登录回调（H5端）
func WXMPLoginCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state != "wxlogin" {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	// 获取公众号access_token
	token, err := getWXMPAccessToken(code)
	if err != nil {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	// 获取用户信息
	mpUserInfo, err := getWXMPUserInfo(token.AccessToken, token.OpenID)
	if err != nil {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	userInfo := WXUserInfo{
		OpenID:    mpUserInfo.OpenID,
		UnionID:   mpUserInfo.UnionID,
		Nickname:  mpUserInfo.Nickname,
		AvatarURL: mpUserInfo.HeadImgURL,
		Gender:    mpUserInfo.Sex,
		Country:   mpUserInfo.Country,
		Province:  mpUserInfo.Province,
		City:      mpUserInfo.City,
		Language:  "",
	}

	// 通过UnionID查找或创建用户
	user, err := findOrCreateUserByUnionID(userInfo)
	if err != nil {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	// 更新微信用户信息
	err = updateWechatUser(user.ID, token.OpenID, token.UnionID, "", token.AccessToken, token.RefreshToken, "mp", userInfo)
	if err != nil {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	// 生成JWT令牌
	jwtToken, err := jwt.GenerateToken(uint32(user.ID), user.Account, "user")
	if err != nil {
		c.Redirect(http.StatusFound, "/pages/user/index/index")
		return
	}

	// 重定向到前端页面，携带token
	redirectURL := fmt.Sprintf("/pages/user/index/callback?token=%s&user_id=%d", jwtToken, user.ID)
	c.Redirect(http.StatusFound, redirectURL)
}

// updateWechatUser 更新微信用户信息
func updateWechatUser(userID uint32, openID, unionID, sessionKey, accessToken, refreshToken, platform string, userInfo WXUserInfo) error {
	var wechatUser model.WechatUser
	err := db.DB.Where("open_id = ?", openID).First(&wechatUser).Error

	if err != nil {
		// 创建新记录
		wechatUser = model.WechatUser{
			OpenID:       openID,
			UnionID:      unionID,
			SessionKey:   sessionKey,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Nickname:     userInfo.Nickname,
			Avatar:       userInfo.AvatarURL,
			Gender:       userInfo.Gender,
			Country:      userInfo.Country,
			Province:     userInfo.Province,
			City:         userInfo.City,
			Language:     userInfo.Language,
			Platform:     platform,
			UsersID:      int32(userID),
			Status:       "0",
			ExpiredAt:    time.Now().Add(time.Hour * 2),
		}
		return db.DB.Create(&wechatUser).Error
	}

	// 更新记录
	wechatUser.UnionID = unionID
	wechatUser.SessionKey = sessionKey
	wechatUser.AccessToken = accessToken
	wechatUser.RefreshToken = refreshToken
	if userInfo.Nickname != "" {
		wechatUser.Nickname = userInfo.Nickname
	}
	if userInfo.AvatarURL != "" {
		wechatUser.Avatar = userInfo.AvatarURL
	}
	if userInfo.Gender != 0 {
		wechatUser.Gender = userInfo.Gender
	}
	wechatUser.Country = userInfo.Country
	wechatUser.Province = userInfo.Province
	wechatUser.City = userInfo.City
	wechatUser.Language = userInfo.Language
	wechatUser.Platform = platform
	wechatUser.ExpiredAt = time.Now().Add(time.Hour * 2)

	return db.DB.Save(&wechatUser).Error
}
