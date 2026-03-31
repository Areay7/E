package main

import (
	"log"

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

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 更新admin用户密码
	result := models.DB.Model(&models.User{}).Where("username = ?", "admin").Update("password", string(hashedPassword))
	if result.Error != nil {
		log.Fatalf("更新密码失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Println("未找到admin用户")
		return
	}

	log.Println("admin密码已更新为: admin123")
}
