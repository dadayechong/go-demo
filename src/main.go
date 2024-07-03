package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	//导入mysql驱动
	"gorm.io/driver/mysql"
)

func main() {

	//定义数据库连接地址
	dsn := "root:12345678@tcp(127.0.0.1:3306)/CRUD?charset=utf8mb4&parseTime=True&loc=Local"

	//连接本地数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败~")
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	type Users struct {
		gorm.Model
		Name    string `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
		Phone   string `gorm:"type:varchar(255);not null" json:"phone" binding:"required"`
		Address string `gorm:"type:varchar(255);not null" json:"address" binding:"required"`
	}
	// 迁移模式
	db.AutoMigrate(&Users{})

	//创建一个gin应用
	r := gin.Default()

	// CORS 中间件 解决跨域问题
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	})

	//增 接口
	r.POST("/user", func(c *gin.Context) {
		var userData Users
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//
		db.Create(&userData)

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建成功",
			"data": userData,
		})
	})

	//删 接口
	r.DELETE("/user/:id", func(c *gin.Context) {
		var userData []Users
		id := c.Param("id")
		if err := c.ShouldBindUri(&id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//通过id查询该条数据在数据库中是否存在
		db.Where("id = ?", id).Find(&userData)
		if len(userData) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "删除失败，找不到该用户"})
			return
		}
		db.Where("id = ?", id).Delete(&userData)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除成功",
		})
	})

	//改 接口
	r.PUT("/user/:id", func(c *gin.Context) {
		var userData Users
		id := c.Param("id")
		// 使用ShouldBindJSON来解析JSON请求体到结构体中
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 更新数据库
		if err := db.Model(&Users{}).Where("id = ?", id).Updates(userData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "找不到用户"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": userData})

	})
	//根据id、name、phone查询用户
	r.GET("/user", func(c *gin.Context) {
		var userData []Users
		if c.Query("id") != "" {
			db.Where("id = ?", c.Query("id")).Find(&userData)
		} else if c.Query("name") != "" {
			db.Where("name = ?", c.Query("name")).Find(&userData)
		} else if c.Query("phone") != "" {
			db.Where("phone = ?", c.Query("phone")).Find(&userData)
		} else {
			db.Find(&userData)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": userData,
		})
	})

	//列表分页查询接口
	r.GET("/userlist", func(c *gin.Context) {
		var userData []Users
		//获取url的分页参数
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userData)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": userData,
		})
	})

	//获取全部分页数据
	r.GET("/userlistall", func(c *gin.Context) {
		var userData []Users
		db.Find(&userData)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": len(userData),
		})
	})

	r.Run(":8080")

}
