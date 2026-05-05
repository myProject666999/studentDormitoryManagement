package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRepairRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	workerName := c.Query("worker_name")

	var records []models.RepairRecord
	var total int64

	query := database.DB.Model(&models.RepairRecord{}).Preload("RepairRequest").Preload("Worker")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if workerName != "" {
		query = query.Joins("JOIN maintenance_workers ON maintenance_workers.id = repair_records.worker_id").
			Where("maintenance_workers.name LIKE ?", "%"+workerName+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取维修记录列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetMyRepairRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var records []models.RepairRecord
	var total int64

	query := database.DB.Model(&models.RepairRecord{}).Preload("RepairRequest").Preload("Worker")

	if role == "worker" {
		query = query.Where("worker_id = ?", userID)
	} else if role == "student" {
		query = query.Joins("JOIN repair_requests ON repair_requests.id = repair_records.repair_request_id").
			Where("repair_requests.student_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取维修记录列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetRepairRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var record models.RepairRecord
	if err := database.DB.Preload("RepairRequest").Preload("Worker").First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "维修记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    record,
	})
}

func CreateRepairRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "worker" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足",
		})
		return
	}

	var input struct {
		RepairRequestID uint    `json:"repair_request_id" binding:"required"`
		Description     string  `json:"description"`
		Cost            float64 `json:"cost"`
		Remark          string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var repairRequest models.RepairRequest
	if err := database.DB.First(&repairRequest, input.RepairRequestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "报修记录不存在",
		})
		return
	}

	workerID := userID
	if repairRequest.WorkerID != nil && role != "admin" {
		workerID = *repairRequest.WorkerID
	}

	now := time.Now()
	record := models.RepairRecord{
		RepairRequestID: input.RepairRequestID,
		WorkerID:        workerID,
		Description:     input.Description,
		Status:          "in_progress",
		StartTime:       &now,
		Cost:            input.Cost,
		Remark:          input.Remark,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建维修记录失败",
		})
		return
	}

	repairRequest.Status = "in_progress"
	database.DB.Save(&repairRequest)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    record,
	})
}

func UpdateRepairRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var record models.RepairRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "维修记录不存在",
		})
		return
	}

	var input struct {
		Description string  `json:"description"`
		Status      string  `json:"status"`
		Cost        float64 `json:"cost"`
		Remark      string  `json:"remark"`
		IsCompleted bool    `json:"is_completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.Description != "" {
		record.Description = input.Description
	}
	if input.Status != "" {
		record.Status = input.Status
	}
	if input.Cost > 0 {
		record.Cost = input.Cost
	}
	if input.Remark != "" {
		record.Remark = input.Remark
	}

	if input.IsCompleted && record.EndTime == nil {
		now := time.Now()
		record.EndTime = &now
		record.Status = "completed"

		var repairRequest models.RepairRequest
		database.DB.First(&repairRequest, record.RepairRequestID)
		repairRequest.Status = "completed"
		database.DB.Save(&repairRequest)
	}

	if err := database.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新维修记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    record,
	})
}

func DeleteRepairRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	if err := database.DB.Delete(&models.RepairRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除维修记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
