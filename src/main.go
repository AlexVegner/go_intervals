package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Interval struct {
	UserId  string `json:"-"`
	Id      string `json:"id"`
	Title   string `json:"title"`
	Details string `json:"details"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}

var intervals = []Interval{
	{
		UserId:  "1",
		Id:      "1",
		Title:   "title",
		Details: "details",
		Start:   123123123,
		End:     123123123,
	},
	{
		UserId:  "2",
		Id:      "1",
		Title:   "title",
		Details: "details",
		Start:   123123124,
		End:     123123125,
	},
}

type authHeader struct {
	Authorization int    `header:"Authorization"`
	Domain        string `header:"Domain"`
}

func Filter(vs []Interval, f func(Interval) bool) []Interval {
	vsf := make([]Interval, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Find(vs []Interval, f func(Interval) bool) (Interval, error) {
	var result Interval
	for _, v := range vs {
		if f(v) {
			result = v
			return result, nil
		}
	}
	return result, errors.New("Element not found")
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello api")
	})

	v1 := router.Group("/api/v1")

	fmt.Println(intervals)

	v1.GET("/intervals", func(c *gin.Context) {
		fmt.Println(intervals)
		user := c.GetHeader("Authorization")
		if len(user) > 0 {
			filtered := Filter(intervals, func(interval Interval) bool { return interval.UserId == user })
			c.JSON(http.StatusOK, filtered)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		}
	})

	v1.POST("/intervals", func(c *gin.Context) {
		user := c.GetHeader("Authorization")
		uuid, _ := uuid.NewUUID()

		id := uuid.String()
		fmt.Println(id)
		if len(user) > 0 {
			var item Interval
			err := c.BindJSON(&item)
			if err == nil && len(item.Title) > 0 && len(item.Details) > 0 && item.Start > 0 && item.End > 0 {
				item.Id = id
				item.UserId = user
				intervals = append(intervals, item)
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		}
	})

	v1.GET("/intervals/:id", func(c *gin.Context) {
		user := c.GetHeader("Authorization")
		if len(user) > 0 {
			id := c.Param("id")
			item, err := Find(intervals, func(interval Interval) bool { return interval.UserId == user && interval.Id == id })
			if err == nil {
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		}
	})

	v1.PUT("/intervals/:id", func(c *gin.Context) {
		user := c.GetHeader("Authorization")
		if len(user) > 0 {
			id := c.Param("id")
			item, err := Find(intervals, func(interval Interval) bool { return interval.UserId == user && interval.Id == id })
			if err == nil {
				var putItem Interval
				err := c.BindJSON(&item)
				if err == nil {
					if len(putItem.Title) > 0 {
						item.Title = putItem.Title
					}
					if len(putItem.Details) > 0 {
						item.Details = putItem.Details
					}
					filtered := Filter(intervals, func(interval Interval) bool {
						return interval.UserId != user || (interval.UserId == user && interval.Id != id)
					})
					intervals = append(filtered, item)
					c.JSON(http.StatusOK, item)
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		}
	})

	v1.DELETE("/intervals/:id", func(c *gin.Context) {
		user := c.GetHeader("Authorization")
		if len(user) > 0 {
			id := c.Param("id")
			_, err := Find(intervals, func(interval Interval) bool { return interval.UserId == user && interval.Id == id })
			if err == nil {
				if err == nil {
					intervals = Filter(intervals, func(interval Interval) bool {
						return interval.UserId != user || (interval.UserId == user && interval.Id != id)
					})
					c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		}
	})

	engine := router.Run(":8080")
	if engine == nil {
		// handle your error here
	}
}
