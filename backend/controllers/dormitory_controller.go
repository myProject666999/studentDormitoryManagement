package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDormitories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	building := c.Query("building")
	roomNumber := c.Query("room_number")
	gender := c.Query("gender")
	status := c.Query("status")

	var dormitories []models.Dormitory
	var total int64

	query := database.DB.Model(&models.Dormitory{})

	if building != "" {
		query = query.Where("building LIKE ?", "%"+building+"%")
	}
	if roomNumber != "" {
		query = query.Where("room_number LIKE ?", "%"+roomNumber+"%")
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&dormitories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取寝室列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      dormitories,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetAvailableDormitories(c *gin.Context) {
	gender := c.Query("gender")

	var dormitories []models.Dormitory
	query := database.DB.Where("status = ?", 1)
	
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	
	if err := query.Find(&dormitories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取可用寝室列表失败",
		})
		return
	}

	availableDorms := make([]models.Dormitory, 0)
	for _, dorm := range dormitories {
		if dorm.UsedBeds < dorm.TotalBeds {
			availableDorms = append(availableDorms, dorm)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    availableDorms,
	})
}

func GetDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var dormitory models.Dormitory
	if err := database.DB.First(&dormitory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    dormitory,
	})
}

func CreateDormitory(c *gin.Context) {
	var dormitory models.Dormitory
	if err := c.ShouldBindJSON(&dormitory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if dormitory.TotalBeds == 0 {
		dormitory.TotalBeds = 4
	}

	if err := database.DB.Create(&dormitory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建寝室失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    dormitory,
	})
}

func UpdateDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var dormitory models.Dormitory
	if err := database.DB.First(&dormitory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "寝室不存在",
		})
		return
	}

	var input struct {
		Building    string `json:"building"`
		Floor       *int   `json:"floor"`
		RoomNumber  string `json:"room_number"`
		RoomType    string `json:"room_type"`
		TotalBeds   *int   `json:"total_beds"`
		Gender      string `json:"gender"`
		Status      *int   `json:"status"`
		Image       string `json:"image"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.Building != "" {
		dormitory.Building = input.Building
	}
	if input.Floor != nil {
		dormitory.Floor = *input.Floor
	}
	if input.RoomNumber != "" {
		dormitory.RoomNumber = input.RoomNumber
	}
	if input.RoomType != "" {
		dormitory.RoomType = input.RoomType
	}
	if input.TotalBeds != nil {
		dormitory.TotalBeds = *input.TotalBeds
	}
	if input.Gender != "" {
		dormitory.Gender = input.Gender
	}
	if input.Status != nil {
		dormitory.Status = *input.Status
	}
	if input.Image != "" {
		dormitory.Image = input.Image
	}
	if input.Description != "" {
		dormitory.Description = input.Description
	}

	if err := database.DB.Save(&dormitory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新寝室失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    dormitory,
	})
}

func DeleteDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var count int64
	database.DB.Model(&models.DormitoryAssignment{}).Where("dormitory_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该寝室已有学生入住，无法删除",
		})
		return
	}

	if err := database.DB.Delete(&models.Dormitory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除寝室失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
