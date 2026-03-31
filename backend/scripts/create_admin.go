package main

import (
	"log"
	"time"

	"cross-border-admin/config"
	"cross-border-admin/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("../config/config.yaml"); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 检查是否已存在admin用户
	var count int64
	models.DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("admin用户已存在，跳过创建")
		return
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建admin用户
	now := time.Now()
	admin := models.User{
		Username:       "admin",
		Password:       string(hashedPassword),
		Nickname:       "系统管理员",
		Email:          "admin@example.com",
		Status:         1,
		Role:           "admin",
		LoginFailCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := models.DB.Create(&admin).Error; err != nil {
		log.Fatalf("创建admin用户失败: %v", err)
	}

	log.Println("admin用户创建成功！")
	log.Println("用户名: admin")
	log.Println("密码: admin123")
}
