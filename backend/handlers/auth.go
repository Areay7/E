package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"cross-border-admin/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	captchaStore = make(map[string]*CaptchaData)
	captchaMutex sync.RWMutex
	jwtSecret    = []byte("your-secret-key-change-in-production")
)

type CaptchaData struct {
	Code      string
	ExpireAt  time.Time
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha"`
	CaptchaID string `json:"captchaId"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	User     *models.User `json:"user"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	// 启动定时清理过期验证码
	go cleanExpiredCaptcha()
	return &AuthHandler{}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	// 检查是否需要验证码（登录失败3次以上）
	if user.LoginFailCount >= 3 {
		if req.Captcha == "" || req.CaptchaID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入验证码", "needCaptcha": true})
			return
		}

		// 验证验证码
		if !verifyCaptcha(req.CaptchaID, req.Captcha) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误", "needCaptcha": true})
			return
		}
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 密码错误，增加失败次数
		user.LoginFailCount++
		user.LoginFailTime = &time.Time{}
		*user.LoginFailTime = time.Now()
		models.DB.Save(&user)

		needCaptcha := user.LoginFailCount >= 3
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":       "用户名或密码错误",
			"needCaptcha": needCaptcha,
		})
		return
	}

	// 登录成功，重置失败次数
	now := time.Now()
	user.LoginFailCount = 0
	user.LastLoginAt = &now
	user.LastLoginIP = c.ClientIP()
	models.DB.Save(&user)

	// 生成JWT token
	token, err := generateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  &user,
	})
}

// GetCaptcha 获取验证码
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	// 生成验证码ID和代码
	captchaID := generateID()
	code := generateCaptchaCode()

	// 存储验证码
	captchaMutex.Lock()
	captchaStore[captchaID] = &CaptchaData{
		Code:     code,
		ExpireAt: time.Now().Add(5 * time.Minute),
	}
	captchaMutex.Unlock()

	// 生成验证码图片
	img := generateCaptchaImage(code)

	// 设置响应头
	c.Header("Content-Type", "image/png")
	c.Header("X-Captcha-ID", captchaID)

	// 输出图片
	png.Encode(c.Writer, img)
}

// GetCurrentUser 获取当前登录用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

// generateToken 生成JWT token
func generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// generateCaptchaCode 生成验证码
func generateCaptchaCode() string {
	const chars = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	code := make([]byte, 4)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

// generateCaptchaImage 生成验证码图片
func generateCaptchaImage(code string) image.Image {
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景色
	bgColor := color.RGBA{240, 240, 240, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 绘制干扰线
	for i := 0; i < 5; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(width)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(height)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(width)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(height)))
		drawLine(img, int(x1.Int64()), int(y1.Int64()), int(x2.Int64()), int(y2.Int64()), color.RGBA{200, 200, 200, 255})
	}

	// 绘制文字（简化版）
	colors := []color.RGBA{
		{50, 50, 200, 255},
		{200, 50, 50, 255},
		{50, 200, 50, 255},
		{200, 100, 50, 255},
	}

	for i, ch := range code {
		x := 20 + i*25
		y := 25
		drawChar(img, x, y, string(ch), colors[i%len(colors)])
	}

	return img
}

// drawLine 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// drawChar 绘制字符（简化版，使用像素点）
func drawChar(img *image.RGBA, x, y int, ch string, col color.Color) {
	// 简化的字符绘制，使用5x7像素矩阵
	patterns := map[string][][]int{
		"0": {{1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}},
		"1": {{0, 1, 0}, {1, 1, 0}, {0, 1, 0}, {0, 1, 0}, {0, 1, 0}, {0, 1, 0}, {1, 1, 1}},
		"2": {{1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 1, 1}},
		"3": {{1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {1, 1, 1}},
		"4": {{1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}},
		"5": {{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {1, 1, 1}},
		"6": {{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}},
		"7": {{1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {0, 1, 0}, {0, 1, 0}, {1, 0, 0}, {1, 0, 0}},
		"8": {{1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}},
		"9": {{1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {1, 1, 1}},
		"A": {{0, 1, 0}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}},
		"B": {{1, 1, 0}, {1, 0, 1}, {1, 0, 1}, {1, 1, 0}, {1, 0, 1}, {1, 0, 1}, {1, 1, 0}},
		"C": {{0, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}, {0, 1, 1}},
		"D": {{1, 1, 0}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 0}},
		"E": {{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 1, 1}},
		"F": {{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}},
	}

	pattern, ok := patterns[ch]
	if !ok {
		pattern = patterns["0"]
	}

	scale := 2
	for row, line := range pattern {
		for colIdx, pixel := range line {
			if pixel == 1 {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.Set(x+colIdx*scale+dx, y+row*scale+dy, col)
					}
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// verifyCaptcha 验证验证码
func verifyCaptcha(captchaID, code string) bool {
	captchaMutex.RLock()
	data, exists := captchaStore[captchaID]
	captchaMutex.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(data.ExpireAt) {
		captchaMutex.Lock()
		delete(captchaStore, captchaID)
		captchaMutex.Unlock()
		return false
	}

	// 删除已使用的验证码
	captchaMutex.Lock()
	delete(captchaStore, captchaID)
	captchaMutex.Unlock()

	return strings.EqualFold(data.Code, code)
}

// generateID 生成随机ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// cleanExpiredCaptcha 定时清理过期验证码
func cleanExpiredCaptcha() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		captchaMutex.Lock()
		now := time.Now()
		for id, data := range captchaStore {
			if now.After(data.ExpireAt) {
				delete(captchaStore, id)
			}
		}
		captchaMutex.Unlock()
	}
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 解析token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证信息"})
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", uint(claims["user_id"].(float64)))
			c.Set("username", claims["username"].(string))
			c.Set("role", claims["role"].(string))
		}

		c.Next()
	}
}
