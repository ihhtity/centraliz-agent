package main

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/config"
	"centraliz-backend/pkg/db"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库连接
	db.InitDB()

	// 生成20条测试数据
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 20; i++ {
		order := generateRandomOrder()
		if err := db.DB.Create(&order).Error; err != nil {
			fmt.Printf("创建订单 %d 失败: %v\n", i+1, err)
		} else {
			fmt.Printf("创建订单 %d 成功: ID=%d\n", i+1, order.ID)
		}
	}

	fmt.Println("数据生成完成!")
}

// generateRandomOrder 生成随机订单数据
func generateRandomOrder() model.Order {
	now := time.Now()
	
	// 随机生成今天之前的开始时间（最近30天内）
	daysAgo := rand.Intn(30)
	startTime := now.Add(-time.Duration(daysAgo)*24*time.Hour - time.Duration(rand.Intn(24))*time.Hour - time.Duration(rand.Intn(60))*time.Minute)
	
	// 随机生成持续时间（1小时到24小时）
	duration := int64(rand.Intn(24) + 1)
	endTime := startTime.Add(time.Duration(duration) * time.Hour)

	// 随机金额（10-1000元）
	price := float64(rand.Intn(9900)+100) / 100
	
	// 随机状态
	statuses := []string{"0", "1", "3", "4", "5"}
	status := statuses[rand.Intn(len(statuses))]
	
	// 随机数量
	amount := float64(rand.Intn(10) + 1)

	// 生成订单编号
	orderCode := fmt.Sprintf("ORD%s%04d", now.Format("20060102150405"), rand.Intn(10000))
	
	// 商品名称列表
	productNames := []string{"商品A", "商品B", "商品C", "商品D", "商品E", "服务套餐", "会员充值", "活动门票"}
	name := productNames[rand.Intn(len(productNames))]

	return model.Order{
		MerchsID:   14,
		UsersID:    int32(rand.Intn(100) + 1),
		DevicesID:  int32(rand.Intn(50) + 1),
		RoomsID:    int32(rand.Intn(30) + 1),
		GroupsID:   int32(rand.Intn(20) + 1),
		Name:       &name,
		Code:       &orderCode,
		Status:     &status,
		Amount:     &amount,
		Duration:   &duration,
		Price:      &price,
		StartTime:  &startTime,
		EndTime:    &endTime,
		UserPhone:  strPtr(fmt.Sprintf("1%d", rand.Intn(9000000000)+1000000000)),
		MerchPhone: strPtr("13800138000"),
		CreatedAt:  &startTime,
	}
}

// strPtr 将字符串转换为指针
func strPtr(s string) *string {
	return &s
}
