package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRepairRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	title := c.Query("title")
	status := c.Query("status")
	studentName := c.Query("student_name")

	var requests []models.RepairRequest
	var total int64

	query := database.DB.Model(&models.RepairRequest{}).Preload("Student").Preload("Dormitory").Preload("Worker")

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if studentName != "" {
		query = query.Joins("JOIN students ON students.id = repair_requests.student_id").
			Where("students.name LIKE ?", "%"+studentName+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报修列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      requests,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetMyRepairRequests(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var requests []models.RepairRequest
	var total int64

	query := database.DB.Model(&models.RepairRequest{}).Preload("Student").Preload("Dormitory").Preload("Worker")

	if role == "student" {
		query = query.Where("student_id = ?", userID)
	} else if role == "worker" {
		query = query.Where("worker_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报修列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      requests,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetRepairRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var request models.RepairRequest
	if err := database.DB.Preload("Student").Preload("Dormitory").Preload("Worker").First(&request, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "报修记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    request,
	})
}

func CreateRepairRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "只有学生可以创建报修",
		})
		return
	}

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Attachment  string `json:"attachment"`
		WorkerID    *uint  `json:"worker_id"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var student models.Student
	if err := database.DB.First(&student, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生信息不存在",
		})
		return
	}

	if student.DormitoryID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "您还没有分配寝室，无法报修",
		})
		return
	}

	request := models.RepairRequest{
		Title:       input.Title,
		Description: input.Description,
		Attachment:  input.Attachment,
		StudentID:   userID,
		DormitoryID: *student.DormitoryID,
		WorkerID:    input.WorkerID,
		Status:      "pending",
		Priority:    input.Priority,
	}

	if request.Priority == "" {
		request.Priority = "normal"
	}

	if err := database.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建报修失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    request,
	})
}

func UpdateRepairRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var request models.RepairRequest
	if err := database.DB.First(&request, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "报修记录不存在",
		})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Attachment  string `json:"attachment"`
		WorkerID    *uint  `json:"worker_id"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.Title != "" {
		request.Title = input.Title
	}
	if input.Description != "" {
		request.Description = input.Description
	}
	if input.Attachment != "" {
		request.Attachment = input.Attachment
	}
	if input.WorkerID != nil {
		request.WorkerID = input.WorkerID
	}
	if input.Status != "" {
		request.Status = input.Status
	}
	if input.Priority != "" {
		request.Priority = input.Priority
	}

	if err := database.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新报修失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    request,
	})
}

func DeleteRepairRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	if err := database.DB.Delete(&models.RepairRequest{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除报修失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
