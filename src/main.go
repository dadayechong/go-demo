package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

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
	//增 接口
	r.POST("/user", func(c *gin.Context) {
		var user Users
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//
		db.Create(&user)

		c.JSON(http.StatusOK, gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"phone":   user.Phone,
			"address": user.Address,
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
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	//改 接口
	r.PUT("/user/:id", func(c *gin.Context) {
		var userData []Users
		id := c.Param("id")
		//通过id查询该条数据在数据库中是否存在
		db.Where("id = ?", id).Find(&userData)
		if len(userData) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "修改失败，找不到该用户"})
			return
		}
		//修改
		db.Model(&userData).Updates(Users{
			Name:    c.PostForm("name"),
			Phone:   c.PostForm("phone"),
			Address: c.PostForm("address"),
		})

		c.JSON(http.StatusOK, gin.H{
			"message": "修改成功",
			"data": gin.H{
				"data": userData,
			},
		})
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
		pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userData)
		c.JSON(http.StatusOK, gin.H{
			"message": "查询成功",
			"data": gin.H{
				"data": userData,
			},
		})
	})

	r.Run(":8080")

}
