package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"student-dormitory-management/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWorkers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	workerNumber := c.Query("worker_number")
	name := c.Query("name")
	specialty := c.Query("specialty")

	var workers []models.MaintenanceWorker
	var total int64

	query := database.DB.Model(&models.MaintenanceWorker{})

	if workerNumber != "" {
		query = query.Where("worker_number LIKE ?", "%"+workerNumber+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if specialty != "" {
		query = query.Where("specialty LIKE ?", "%"+specialty+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取维修工列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      workers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetAllWorkers(c *gin.Context) {
	var workers []models.MaintenanceWorker
	if err := database.DB.Where("status = ?", 1).Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取维修工列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    workers,
	})
}

func CreateWorker(c *gin.Context) {
	var worker models.MaintenanceWorker
	if err := c.ShouldBindJSON(&worker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if worker.Password == "" {
		worker.Password = "123456"
	}

	if err := database.DB.Create(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建维修工失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    worker,
	})
}

func UpdateWorker(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var worker models.MaintenanceWorker
	if err := database.DB.First(&worker, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "维修工不存在",
		})
		return
	}

	var input struct {
		WorkerNumber string `json:"worker_number"`
		Name         string `json:"name"`
		Gender       string `json:"gender"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
		Specialty    string `json:"specialty"`
		Password     string `json:"password"`
		Status       *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.WorkerNumber != "" {
		worker.WorkerNumber = input.WorkerNumber
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
	if input.Status != nil {
		worker.Status = *input.Status
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
		worker.Password = hashedPassword
	}

	if err := database.DB.Save(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新维修工失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    worker,
	})
}

func DeleteWorker(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	if err := database.DB.Delete(&models.MaintenanceWorker{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除维修工失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
