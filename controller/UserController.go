package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/git_test_project/common"
	"github.com/git_test_project/model"
	"github.com/git_test_project/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	/*name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")*/
	var requestUser model.User
	ctx.ShouldBind(&requestUser)
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "手机号必须为11位",
		})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "密码不能少于6位",
		})
		return
	}

	if len(name) == 0 {
		name = utils.RandomString(10)
	}
	log.Print("注册的name:", name)

	if telephoneIsExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "用户已经存在",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 500,
			"msg":  "加密失败",
		})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}

	DB.Create(&newUser)

	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	/*var userLogin model.User
	if err:= ctx.ShouldBind(&userLogin); err!=nil{
		log.Panicf("login param failed, err: %v",err)
		return
	}

	//获取参数
	telephone := userLogin.Telephone
	password := userLogin.Password*/
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	log.Printf("telephone:%s,password:%s", telephone, password)
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "手机号必须为11位",
		})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "密码不能少于6位",
		})
		return
	}

	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "用户不存在",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户或者密码错误",
		})
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 500,
			"msg":  "系统异常",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": token,
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": user},
	})
}

func telephoneIsExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
