package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(10);not null""`
}

func main() {
	db:= initDb()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		if len(telephone) != 11{
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"msg": "手机号必须为11位",
			})
			return
		}

		if len(password) <6 {
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"msg":"密码不能少于6位",
			})
			return
		}

		if len(name) == 0 {
			name = randomString(10)
		}
		log.Print("注册的name:",name)

		if telephoneIsExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"msg":"用户已经存在",
			})
			return
		}

		newUser :=User{
			Name: name,
			Telephone: telephone,
			Password: string(password),
		}

		db.Create(&newUser)

		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func randomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyz")

	result := make([]byte,n)
	rand.Seed(time.Now().UnixNano())
	for i:= range result{
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func telephoneIsExist(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0{
		return true
	}
	return  false
}

func initDb() *gorm.DB {
	driverName := "mysql"
	host:="10.23.171.250"
	port:="13306"
	database:="gintestproject"
	username:="root"
	password:="my_root_password"
	charset:="utf8"
	args:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName,args)
	if err !=nil{
		panic("fail to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}


