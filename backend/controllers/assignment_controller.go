package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAssignments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	studentName := c.Query("student_name")
	studentNumber := c.Query("student_number")
	building := c.Query("building")
	roomNumber := c.Query("room_number")

	var assignments []models.DormitoryAssignment
	var total int64

	query := database.DB.Model(&models.DormitoryAssignment{}).Preload("Student").Preload("Dormitory")

	if studentName != "" {
		query = query.Joins("JOIN students ON students.id = dormitory_assignments.student_id").
			Where("students.name LIKE ?", "%"+studentName+"%")
	}
	if studentNumber != "" {
		query = query.Joins("JOIN students ON students.id = dormitory_assignments.student_id").
			Where("students.student_number LIKE ?", "%"+studentNumber+"%")
	}
	if building != "" {
		query = query.Joins("JOIN dormitories ON dormitories.id = dormitory_assignments.dormitory_id").
			Where("dormitories.building LIKE ?", "%"+building+"%")
	}
	if roomNumber != "" {
		query = query.Joins("JOIN dormitories ON dormitories.id = dormitory_assignments.dormitory_id").
			Where("dormitories.room_number LIKE ?", "%"+roomNumber+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取寝室安排列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      assignments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var assignment models.DormitoryAssignment
	if err := database.DB.Preload("Student").Preload("Dormitory").First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室安排不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    assignment,
	})
}

func GetMyAssignment(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "只有学生可以查询自己的寝室安排",
		})
		return
	}

	var assignment models.DormitoryAssignment
	if err := database.DB.Where("student_id = ?", userID).Preload("Dormitory").First(&assignment).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "暂无寝室安排",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    assignment,
	})
}

func CreateAssignment(c *gin.Context) {
	var input struct {
		StudentID   uint `json:"student_id" binding:"required"`
		DormitoryID uint `json:"dormitory_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var existingAssignment models.DormitoryAssignment
	if err := database.DB.Where("student_id = ?", input.StudentID).First(&existingAssignment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该学生已有寝室安排",
		})
		return
	}

	var dormitory models.Dormitory
	if err := database.DB.First(&dormitory, input.DormitoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室不存在",
		})
		return
	}

	if dormitory.UsedBeds >= dormitory.TotalBeds {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该寝室床位已满",
		})
		return
	}

	assignment := models.DormitoryAssignment{
		StudentID:      input.StudentID,
		DormitoryID:    input.DormitoryID,
		AssignmentDate: time.Now(),
		Status:         1,
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建寝室安排失败",
		})
		return
	}

	dormitory.UsedBeds++
	database.DB.Save(&dormitory)

	var student models.Student
	database.DB.First(&student, input.StudentID)
	student.DormitoryID = &dormitory.ID
	database.DB.Save(&student)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    assignment,
	})
}

func UpdateAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var assignment models.DormitoryAssignment
	if err := database.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室安排不存在",
		})
		return
	}

	var input struct {
		DormitoryID *uint `json:"dormitory_id"`
		Status      *int  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.DormitoryID != nil && *input.DormitoryID != assignment.DormitoryID {
		var oldDormitory models.Dormitory
		database.DB.First(&oldDormitory, assignment.DormitoryID)
		oldDormitory.UsedBeds--
		database.DB.Save(&oldDormitory)

		var newDormitory models.Dormitory
		if err := database.DB.First(&newDormitory, *input.DormitoryID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "新寝室不存在",
			})
			return
		}

		if newDormitory.UsedBeds >= newDormitory.TotalBeds {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "新寝室床位已满",
			})
			return
		}

		newDormitory.UsedBeds++
		database.DB.Save(&newDormitory)

		assignment.DormitoryID = *input.DormitoryID

		var student models.Student
		database.DB.First(&student, assignment.StudentID)
		student.DormitoryID = input.DormitoryID
		database.DB.Save(&student)
	}

	if input.Status != nil {
		assignment.Status = *input.Status
	}

	if err := database.DB.Save(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新寝室安排失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    assignment,
	})
}

func DeleteAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var assignment models.DormitoryAssignment
	if err := database.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室安排不存在",
		})
		return
	}

	var dormitory models.Dormitory
	database.DB.First(&dormitory, assignment.DormitoryID)
	dormitory.UsedBeds--
	database.DB.Save(&dormitory)

	var student models.Student
	database.DB.First(&student, assignment.StudentID)
	student.DormitoryID = nil
	database.DB.Save(&student)

	if err := database.DB.Delete(&models.DormitoryAssignment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除寝室安排失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
