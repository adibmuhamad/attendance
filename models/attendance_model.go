package models

type AttendanceRequest struct {
	Base64Image string  `json:"image"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type AttendanceResponse struct {
	Time           string `json:"time"`
	RecognizedFace bool   `json:"recognizedFace"`
	ValidLocation  bool   `json:"validLocation"`
	Status         bool   `json:"status"`
}
