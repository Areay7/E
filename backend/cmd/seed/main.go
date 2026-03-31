package main

import (
	"log"
	"math/rand"
	"time"

	"cross-border-admin/models"
)

func main() {
	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 生成测试数据
	platforms := []string{"shopee", "aliexpress", "tiktok"}
	statuses := []string{"pending", "processing", "pending_shipment", "shipped", "completed"}
	currencies := []string{"USD", "EUR", "SGD", "MYR"}
	buyers := []string{"John Doe", "Jane Smith", "李明", "王芳", "Tanaka", "Kumar"}

	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 100; i++ {
		platform := platforms[rand.Intn(len(platforms))]
		order := models.Order{
			OrderID:   platform + "-" + time.Now().Format("20060102") + "-" + randomString(8),
			Platform:  platform,
			Buyer:     buyers[rand.Intn(len(buyers))],
			Amount:    float64(rand.Intn(50000)+1000) / 100.0,
			Currency:  currencies[rand.Intn(len(currencies))],
			Status:    statuses[rand.Intn(len(statuses))],
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour),
		}

		if err := models.DB.Create(&order).Error; err != nil {
			log.Printf("创建订单失败: %v", err)
		}
	}

	log.Println("测试数据生成完成！共生成 100 条订单记录")
}

func randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
