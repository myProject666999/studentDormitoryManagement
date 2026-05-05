package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"student-dormitory-management/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	studentNumber := c.Query("student_number")
	name := c.Query("name")
	class := c.Query("class")

	var students []models.Student
	var total int64

	query := database.DB.Model(&models.Student{})

	if studentNumber != "" {
		query = query.Where("student_number LIKE ?", "%"+studentNumber+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if class != "" {
		query = query.Where("class LIKE ?", "%"+class+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Dormitory").Offset(offset).Limit(pageSize).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      students,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if student.Password == "" {
		student.Password = "123456"
	}

	if err := database.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建学生失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    student,
	})
}

func UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生不存在",
		})
		return
	}

	var input struct {
		StudentNumber string `json:"student_number"`
		Name          string `json:"name"`
		Gender        string `json:"gender"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		Class         string `json:"class"`
		Major         string `json:"major"`
		Password      string `json:"password"`
		Status        *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.StudentNumber != "" {
		student.StudentNumber = input.StudentNumber
	}
	if input.Name != "" {
		student.Name = input.Name
	}
	if input.Gender != "" {
		student.Gender = input.Gender
	}
	if input.Phone != "" {
		student.Phone = input.Phone
	}
	if input.Email != "" {
		student.Email = input.Email
	}
	if input.Class != "" {
		student.Class = input.Class
	}
	if input.Major != "" {
		student.Major = input.Major
	}
	if input.Status != nil {
		student.Status = *input.Status
	}

	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败",
			})
			return
		}
		student.Password = hashedPassword
	}

	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新学生失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    student,
	})
}

func DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	if err := database.DB.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除学生失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var updatedUser interface{}

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

		var input struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		if input.Name != "" {
			admin.Name = input.Name
		}
		if input.Phone != "" {
			admin.Phone = input.Phone
		}
		if input.Email != "" {
			admin.Email = input.Email
		}

		database.DB.Save(&admin)
		updatedUser = admin

	case "student":
		var student models.Student
		if err := database.DB.First(&student, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		var input struct {
			Name   string `json:"name"`
			Gender string `json:"gender"`
			Phone  string `json:"phone"`
			Email  string `json:"email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		if input.Name != "" {
			student.Name = input.Name
		}
		if input.Gender != "" {
			student.Gender = input.Gender
		}
		if input.Phone != "" {
			student.Phone = input.Phone
		}
		if input.Email != "" {
			student.Email = input.Email
		}

		database.DB.Save(&student)
		updatedUser = student

	case "worker":
		var worker models.MaintenanceWorker
		if err := database.DB.First(&worker, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
			})
			return
		}

		var input struct {
			Name      string `json:"name"`
			Gender    string `json:"gender"`
			Phone     string `json:"phone"`
			Email     string `json:"email"`
			Specialty string `json:"specialty"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		if input.Name != "" {
			worker.Name = input.Name
		}
		if input.Gender != "" {
			worker.Gender = input.Gender
		}
		if input.Phone != "" {
			worker.Phone = input.Phone
		}
		if input.Email != "" {
			worker.Email = input.Email
		}
		if input.Specialty != "" {
			worker.Specialty = input.Specialty
		}

		database.DB.Save(&worker)
		updatedUser = worker
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    updatedUser,
	})
}
