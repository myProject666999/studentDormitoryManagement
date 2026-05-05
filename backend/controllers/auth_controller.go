package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"student-dormitory-management/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var admin models.Admin
	var student models.Student
	var worker models.MaintenanceWorker
	var foundUser *models.User
	var userInfo interface{}

	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err == nil {
		foundUser = &admin.User
		userInfo = admin
	} else if err := database.DB.Where("username = ?", req.Username).First(&student).Error; err == nil {
		foundUser = &student.User
		userInfo = student
	} else if err := database.DB.Where("username = ?", req.Username).First(&worker).Error; err == nil {
		foundUser = &worker.User
		userInfo = worker
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	if foundUser.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "账号已被禁用",
		})
		return
	}

	if !utils.CheckPassword(req.Password, foundUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	token, err := utils.GenerateToken(foundUser.ID, foundUser.Username, foundUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user":  userInfo,
			"role":  foundUser.Role,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var userInfo interface{}
	switch role {
	case "admin":
		var admin models.Admin
		database.DB.Where("id = ?", userID).First(&admin)
		userInfo = admin
	case "student":
		var student models.Student
		database.DB.Where("id = ?", userID).First(&student)
		userInfo = student
	case "worker":
		var worker models.MaintenanceWorker
		database.DB.Where("id = ?", userID).First(&worker)
		userInfo = worker
	default:
		var user models.User
		database.DB.Where("id = ?", userID).First(&user)
		userInfo = user
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    userInfo,
	})
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	switch role {
	case "admin":
		var admin models.Admin
		if err := database.DB.First(&admin, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		if !utils.CheckPassword(req.OldPassword, admin.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "原密码错误",
			})
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败",
			})
			return
		}
		admin.Password = hashedPassword
		database.DB.Save(&admin)

	case "student":
		var student models.Student
		if err := database.DB.First(&student, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		if !utils.CheckPassword(req.OldPassword, student.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "原密码错误",
			})
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败",
			})
			return
		}
		student.Password = hashedPassword
		database.DB.Save(&student)

	case "worker":
		var worker models.MaintenanceWorker
		if err := database.DB.First(&worker, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		if !utils.CheckPassword(req.OldPassword, worker.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "原密码错误",
			})
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败",
			})
			return
		}
		worker.Password = hashedPassword
		database.DB.Save(&worker)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户角色",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
	})
}
