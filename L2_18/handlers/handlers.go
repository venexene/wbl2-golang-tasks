package handlers 

import (
	"net/http"
    "time"
    "fmt"
    "strings"

	"github.com/gin-gonic/gin"

    "github.com/venexene/calendar/internal"
)

func TestServerHandle(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "Server definetly works",
    })
}


func AddHandler(c *gin.Context) {
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
			"error": "user_id is required",
		})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date is required",
		})
		return
	}

	if request.Event == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "event is required",
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
		"result": "Event created successfully",
	})
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %s %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration,
		)
	}
}

func CalendarMiddleware(calendarDB *calendar.Calendar) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("calendar", calendarDB)
		c.Next()
	}
}
