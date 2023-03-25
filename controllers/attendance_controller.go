package controllers

import (
	"id/projects/attendance/helper"
	"id/projects/attendance/models"
	"id/projects/attendance/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type attendanceController struct {
	faceRecognitionService services.FaceRecognitionService
	geoLocationService services.GeoLocationService
}

func NewAttendanceController(frService services.FaceRecognitionService, glService services.GeoLocationService) *attendanceController {
	return &attendanceController{frService, glService}
}

func (h *attendanceController) AttendeeCheckIn(c *gin.Context) {
	var req models.AttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Unable to process request", http.StatusUnprocessableEntity, "FAILED", errorMessage)
		c.JSON(http.StatusOK, response)
		return
	}

	respFormatter := models.AttendanceResponse{}

	// Get the current time as clock-in time
	checkinTime := time.Now()
	respFormatter.Time = checkinTime.Format(time.RFC3339)

	// Validate if check-in time is later than 8:00 AM
	respFormatter.Status = checkinTime.Hour() < 8
	

	// Validate face and location
	recognizedFace, faceErr := h.faceRecognitionService.RecognizeFace(req.Base64Image)
	validLocation, geoErr := h.geoLocationService.ValidateGeolocation(c, req.Longitude, req.Latitude)

	respFormatter.RecognizedFace = recognizedFace
	respFormatter.ValidLocation = validLocation

	if faceErr != nil {
		response := helper.APIResponse(faceErr.Error(), http.StatusBadRequest, "FAILED", respFormatter)
		c.JSON(http.StatusOK, response)
		return
	}

	if geoErr != nil {
		response := helper.APIResponse(geoErr.Error(), http.StatusBadRequest, "FAILED", respFormatter)
		c.JSON(http.StatusOK, response)
		return
	}

	response := helper.APIResponse("Attendance checked in successfully", http.StatusOK, "SUCCESS", respFormatter)
	c.JSON(http.StatusOK, response)
}

func (h *attendanceController) AttendeeCheckOut(c *gin.Context) {
	var req models.AttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Unable to process request", http.StatusUnprocessableEntity, "FAILED", errorMessage)
		c.JSON(http.StatusOK, response)
		return
	}

	respFormatter := models.AttendanceResponse{}

	// Get the current time as clock-in time
	checkinTime := time.Now()
	respFormatter.Time = checkinTime.Format(time.RFC3339)
	
	// Validate if check-in time is less than 17:00 PM
	respFormatter.Status = checkinTime.Hour() > 17

	// Validate face and location
	recognizedFace, faceErr := h.faceRecognitionService.RecognizeFace(req.Base64Image)
	validLocation, geoErr := h.geoLocationService.ValidateGeolocation(c, req.Longitude, req.Latitude)

	respFormatter.RecognizedFace = recognizedFace
	respFormatter.ValidLocation = validLocation

	if faceErr != nil {
		response := helper.APIResponse(faceErr.Error(), http.StatusBadRequest, "FAILED", respFormatter)
		c.JSON(http.StatusOK, response)
		return
	}

	if geoErr != nil {
		response := helper.APIResponse(geoErr.Error(), http.StatusBadRequest, "FAILED", respFormatter)
		c.JSON(http.StatusOK, response)
		return
	}

	response := helper.APIResponse("Attendance checked out successfully", http.StatusOK, "SUCCESS", respFormatter)
	c.JSON(http.StatusOK, response)
}

