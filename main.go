package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Demo struct {
	Id    int64 `gorm:"primaryKey,autoIncrement"`
	Key   string
	Ctime int64
	Utime int64
}

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		//Addr: "localhost:31379",
		Addr: "k8s-demo-redis-service:6379",
	})
}

func InitDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:30002)/k8s_demo"))
	db, err := gorm.Open(mysql.Open("root:root@tcp(k8s-demo-mysql-service:3308)/k8s_demo"))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Demo{})
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	redisClient := InitRedis()
	mysqlClient := InitDB()
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	server.GET("/set_redis", func(ctx *gin.Context) {
		err := redisClient.Set(ctx, "test_redis_key", "test_redis_value", -1).Err()
		if err != nil {
			ctx.String(http.StatusOK, "set redis error")
			return
		}
		ctx.String(http.StatusOK, "set redis success, key: test_redis_key, val: test_redis_value")
	})
	server.GET("/get_redis", func(ctx *gin.Context) {
		s, err := redisClient.Get(ctx, "test_redis_key").Result()
		if err != nil {
			ctx.String(http.StatusOK, "get redis error")
			return
		}
		ctx.String(http.StatusOK, fmt.Sprintf("get redis success, key: %s, val: %s", "test_redis_key", s))
	})

	server.GET("/set_mysql", func(ctx *gin.Context) {
		demo := Demo{
			Key:   "test_mysql_key",
			Ctime: time.Now().UnixMilli(),
			Utime: time.Now().UnixMilli(),
		}
		err := mysqlClient.WithContext(ctx).Create(&demo).Error
		if err != nil {
			ctx.String(http.StatusOK, "set mysql error")
			return
		}
		ctx.String(http.StatusOK, "set mysql success")
	})

	server.GET("/get_mysql", func(ctx *gin.Context) {
		var d Demo
		err := mysqlClient.WithContext(ctx).First(&d).Error
		if err != nil {
			fmt.Println(err)
			ctx.String(http.StatusOK, "get mysql error")
			return
		}
		ctx.String(http.StatusOK, "get mysql success, id: "+strconv.FormatInt(d.Id, 10))
	})

	if err := server.Run(":8082"); err != nil {
		panic(err)
	}

	fmt.Printf("web server start....")
}
