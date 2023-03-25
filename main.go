package main

import (
	"id/projects/attendance/controllers"
	"id/projects/attendance/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	faceRecognitionService := services.NewFaceRecognitionService()
	geoLocationService := services.NewGeoLocationService()

	attendanceController := controllers.NewAttendanceController(faceRecognitionService, geoLocationService)

	attendance := r.Group("/api/v1/attendance")
	{
		attendance.POST("/check-in", attendanceController.AttendanceCheckIn)
		attendance.POST("/check-out", attendanceController.AttendanceCheckOut)
	}

	r.Run(":8080")
}