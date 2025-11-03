// Package handlers provides handlers for calendar service
package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/venexene/calendar/internal"
)

// TestServerHandle handles test request
func TestServerHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Server definetly works",
	})
}

// AddHandle handles add requests
func AddHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Date   string `form:"date" json:"date" binding:"required"`
		Event  string `form:"event" json:"event" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date is required",
		})
		return
	}

	if request.Event == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event is required",
		})
		return
	}

	err := calendarDB.Add(request.UserID, request.Date, request.Event)
	if err != nil {
		if err.Error() == "Invalid date" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "New event created successfully",
	})
}

// UpdateHandle handles update requests
func UpdateHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID   string `form:"user_id" json:"user_id" binding:"required"`
		Date     string `form:"date" json:"date" binding:"required"`
		Event    string `form:"event" json:"event" binding:"required"`
		NewEvent string `form:"new_event" json:"new_event" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date is required",
		})
		return
	}

	if request.Event == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event is required",
		})
		return
	}

	if request.NewEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "New event is required",
		})
		return
	}

	err := calendarDB.Update(request.UserID, request.Date, request.Event, request.NewEvent)
	if err != nil {
		if err.Error() == "Invalid date" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Event updated successfully",
	})
}

// DeleteHandle handles delete requests
func DeleteHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Date   string `form:"date" json:"date" binding:"required"`
		Event  string `form:"event" json:"event" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date is required",
		})
		return
	}

	if request.Event == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event is required",
		})
		return
	}

	err := calendarDB.Delete(request.UserID, request.Date, request.Event)
	if err != nil {
		if err.Error() == "Invalid date" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Event deleted successfully",
	})
}

// DayEventsHandle handles requests to get events by date
func DayEventsHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Day    string `form:"day" json:"day" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Day == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Day is required",
		})
		return
	}

	events, err := calendarDB.GetEventsByDay(request.UserID, request.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": request.UserID,
		"day":     request.Day,
		"events":  events,
		"count":   len(events),
	})
}

// WeekEventsHandle handles requests to get events by week
func WeekEventsHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Week   string `form:"week" json:"week" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Week == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Week is required",
		})
		return
	}

	events, err := calendarDB.GetEventsByWeek(request.UserID, request.Week)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": request.UserID,
		"week":    request.Week,
		"events":  events,
		"count":   len(events),
	})
}

// MonthEventsHandle handles requests to get events by month
func MonthEventsHandle(c *gin.Context) {
	db, exists := c.Get("calendar")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Calendar not available",
		})
		return
	}
	calendarDB := db.(*calendar.Calendar)

	var request struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Month  string `form:"month" json:"month" binding:"required"`
	}

	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data: " + err.Error(),
			})
			return
		}
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID is required",
		})
		return
	}

	if request.Month == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Month is required",
		})
		return
	}

	events, err := calendarDB.GetEventsByMonth(request.UserID, request.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": request.UserID,
		"month":   request.Month,
		"events":  events,
		"count":   len(events),
	})
}

// LoggingMiddleware provides middleware logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		log.Printf("[%s] %s %s %s %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration,
		)
	}
}

// CalendarMiddleware adds calendar to context
func CalendarMiddleware(calendarDB *calendar.Calendar) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("calendar", calendarDB)
		c.Next()
	}
}
