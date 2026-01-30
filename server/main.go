package main

import (
	"net/http"
	"timetrace/database"
	"timetrace/models"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.Connect()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api/v1")
	{
		// Health Check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Categories
		api.GET("/categories", getCategories)
		api.POST("/categories", createCategory)
		api.DELETE("/categories/:id", deleteCategory)

		// Activities
		api.GET("/activities", getActivities)
		api.POST("/activities", createActivity)
		api.POST("/activities/:id/pin", pinActivity)
		api.DELETE("/activities/:id", deleteActivity)

		// Timer
		api.GET("/timer/current", getCurrentTimer)
		api.POST("/timer/start", startTimer)
		api.POST("/timer/stop", stopTimer)

		// Time Entries
		api.GET("/time-entries", getTimeEntries)
		api.POST("/time-entries", createTimeEntry)
		api.PATCH("/time-entries/:id", updateTimeEntry)
		api.DELETE("/time-entries/:id", deleteTimeEntry)
		
		// Reports
		api.GET("/reports/daily", getDailyReport)
	}

	r.Run(":8080")
}

// --- Handlers (To be refactored into separate files) ---

func getCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func createCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Defaults
	if category.UserID == 0 {
		category.UserID = 1
	}
	
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func deleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func getActivities(c *gin.Context) {
	var activities []models.Activity
	database.DB.Preload("Category").Find(&activities)
	c.JSON(http.StatusOK, activities)
}

func createActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if activity.UserID == 0 {
		activity.UserID = 1
	}
	if err := database.DB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func deleteActivity(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Activity{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func pinActivity(c *gin.Context) {
	id := c.Param("id")
	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}
	
	// Toggle pin or set to 1? PRD says POST /pin, usually implies pinning. 
	// Let's assume toggle or set true. For now set to 1 (pinned).
	// But let's check body if we want to support unpin.
	// For MVP simplicity: Toggle
	newPin := 1
	if activity.Pinned == 1 {
		newPin = 0
	}
	
	database.DB.Model(&activity).Update("pinned", newPin)
	c.JSON(http.StatusOK, gin.H{"pinned": newPin})
}

func getCurrentTimer(c *gin.Context) {
	var entry models.TimeEntry
	// Find entry with end_time IS NULL
	result := database.DB.Preload("Activity").Preload("Category").Where("user_id = ? AND end_time IS NULL", 1).First(&entry)
	if result.Error != nil {
		c.JSON(http.StatusOK, nil) // No running timer
		return
	}
	c.JSON(http.StatusOK, entry)
}

func startTimer(c *gin.Context) {
	var req struct {
		ActivityID uint   `json:"activity_id"`
		Note       string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	// 1. Stop current running timer if exists
	var current models.TimeEntry
	if err := tx.Where("user_id = ? AND end_time IS NULL", 1).First(&current).Error; err == nil {
		now := time.Now()
		if err := tx.Model(&current).Update("end_time", now).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop current timer"})
			return
		}
	}

	// 2. Get Activity to fill CategoryID
	var activity models.Activity
	if err := tx.First(&activity, req.ActivityID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activity not found"})
		return
	}

	// 3. Create new entry
	newEntry := models.TimeEntry{
		UserID:     1,
		CategoryID: activity.CategoryID,
		ActivityID: req.ActivityID,
		StartTime:  time.Now(),
		Note:       req.Note,
	}

	if err := tx.Create(&newEntry).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start timer"})
		return
	}

	tx.Commit()
	// Reload to get preloads if needed, or just return
	c.JSON(http.StatusOK, newEntry)
}

func stopTimer(c *gin.Context) {
	var req struct {
		EndTime *time.Time `json:"end_time"`
	}
	c.ShouldBindJSON(&req) // Optional

	var current models.TimeEntry
	if err := database.DB.Where("user_id = ? AND end_time IS NULL", 1).First(&current).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "No running timer"}) // Or 204
		return
	}

	endTime := time.Now()
	if req.EndTime != nil {
		endTime = *req.EndTime
	}

	if err := database.DB.Model(&current).Update("end_time", endTime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop timer"})
		return
	}
	
	// Return the full object
	current.EndTime = &endTime
	c.JSON(http.StatusOK, current)
}

func getTimeEntries(c *gin.Context) {
	// Query params from/to
	fromStr := c.Query("from")
	toStr := c.Query("to")
	
	db := database.DB.Preload("Activity").Preload("Category").Where("user_id = ?", 1)
	
	if fromStr != "" {
		// Parse
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			db = db.Where("start_time >= ?", t)
		}
	}
	if toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			// Add 1 day to include the end date fully if it's just a date
			t = t.Add(24 * time.Hour)
			db = db.Where("start_time < ?", t)
		}
	}
	
	var entries []models.TimeEntry
	db.Order("start_time desc").Find(&entries)
	c.JSON(http.StatusOK, entries)
}

func createTimeEntry(c *gin.Context) {
	var entry models.TimeEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validation: Start < End
	if entry.EndTime != nil && !entry.EndTime.After(entry.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}
	
	if entry.UserID == 0 {
		entry.UserID = 1
	}

	// Fetch category if not provided but activity is? 
	// Ideally frontend sends both, or we look it up.
	if entry.CategoryID == 0 && entry.ActivityID != 0 {
		var act models.Activity
		if err := database.DB.First(&act, entry.ActivityID).Error; err == nil {
			entry.CategoryID = act.CategoryID
		}
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func updateTimeEntry(c *gin.Context) {
	id := c.Param("id")
	var entry models.TimeEntry
	if err := database.DB.First(&entry, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	var req models.TimeEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update allowed fields
	updates := make(map[string]interface{})
	if req.ActivityID != 0 {
		updates["activity_id"] = req.ActivityID
		// Should update category too? Yes if changed.
		var act models.Activity
		if err := database.DB.First(&act, req.ActivityID).Error; err == nil {
			updates["category_id"] = act.CategoryID
		}
	}
	if !req.StartTime.IsZero() {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = req.EndTime
	}
	if req.Note != "" {
		updates["note"] = req.Note
	}
	
	if err := database.DB.Model(&entry).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func deleteTimeEntry(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.TimeEntry{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func getDailyReport(c *gin.Context) {
	dateStr := c.Query("date") // YYYY-MM-DD
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	
	start, _ := time.Parse("2006-01-02", dateStr)
	end := start.Add(24 * time.Hour)
	
	var entries []models.TimeEntry
	database.DB.Preload("Category").Preload("Activity").
		Where("user_id = ? AND start_time >= ? AND start_time < ?", 1, start, end).
		Find(&entries)
		
	// Aggregate
	type Stat struct {
		Name     string `json:"name"`
		Duration float64 `json:"duration"` // Seconds
		Color    string `json:"color"`
	}
	
	catStats := make(map[string]*Stat)
	actStats := make(map[string]*Stat)
	
	totalDuration := 0.0
	
	for _, e := range entries {
		var duration float64
		if e.EndTime != nil {
			duration = e.EndTime.Sub(e.StartTime).Seconds()
		} else {
			// Running
			duration = time.Since(e.StartTime).Seconds()
		}
		
		totalDuration += duration
		
		// Category
		cName := "Unknown"
		cColor := ""
		if e.Category.Name != "" {
			cName = e.Category.Name
			cColor = e.Category.Color
		}
		if _, ok := catStats[cName]; !ok {
			catStats[cName] = &Stat{Name: cName, Duration: 0, Color: cColor}
		}
		catStats[cName].Duration += duration
		
		// Activity
		aName := "Unknown"
		if e.Activity.Name != "" {
			aName = e.Activity.Name
		}
		if _, ok := actStats[aName]; !ok {
			actStats[aName] = &Stat{Name: aName, Duration: 0}
		}
		actStats[aName].Duration += duration
	}
	
	cList := make([]Stat, 0, len(catStats))
	for _, v := range catStats {
		cList = append(cList, *v)
	}
	aList := make([]Stat, 0, len(actStats))
	for _, v := range actStats {
		aList = append(aList, *v)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"date": dateStr,
		"total_duration": totalDuration,
		"by_category": cList,
		"by_activity": aList,
	})
}
