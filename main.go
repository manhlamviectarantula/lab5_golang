package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define the Student struct
type Student struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Major string `json:"major"`
}

var DB *gorm.DB

func initDB() {
	dsn := "root:my-secret-pw@tcp(127.0.0.1:3308)/student_management?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&Student{})
}

func main() {
	initDB()

	router := gin.Default()

	router.GET("/get-students", func(c *gin.Context) {
		var students []Student
		DB.Find(&students)
		c.JSON(http.StatusOK, students)
	})

	router.GET("/get-student-detail/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}

		var student Student
		if err := DB.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusOK, student)
	})

	router.POST("/add-student", func(c *gin.Context) {
		var newStudent Student
		if err := c.ShouldBindJSON(&newStudent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		DB.Create(&newStudent)
		c.JSON(http.StatusCreated, newStudent)
	})

	router.PUT("/update-student/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}

		var student Student
		if err := DB.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}

		if err := c.ShouldBindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		DB.Save(&student)
		c.JSON(http.StatusOK, student)
	})

	router.DELETE("/delete-student/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}

		var student Student
		if err := DB.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}

		DB.Delete(&student)
		c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
	})

	router.Run(":8080") // Start the server on port 8080
}
