package controllers

import (
	"net/http"
	"student-dormitory-management/database"
	"student-dormitory-management/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNotices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	title := c.Query("title")

	var notices []models.Notice
	var total int64

	query := database.DB.Model(&models.Notice{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取公告列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      notices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "公告不存在",
		})
		return
	}

	notice.ViewCount++
	database.DB.Save(&notice)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    notice,
	})
}

func CreateNotice(c *gin.Context) {
	userID := c.GetUint("user_id")
	username := c.GetString("username")

	var notice models.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	notice.AuthorID = userID
	notice.AuthorName = username
	notice.ViewCount = 0

	if err := database.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建公告失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    notice,
	})
}

func UpdateNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "公告不存在",
		})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Image   string `json:"image"`
		Status  *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if input.Title != "" {
		notice.Title = input.Title
	}
	if input.Content != "" {
		notice.Content = input.Content
	}
	if input.Image != "" {
		notice.Image = input.Image
	}
	if input.Status != nil {
		notice.Status = *input.Status
	}

	if err := database.DB.Save(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新公告失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    notice,
	})
}

func DeleteNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID格式错误",
		})
		return
	}

	if err := database.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除公告失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
