package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GeoLocation struct {
	Latitude  float64
	Longitude float64
}

type IPAPIResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Replace these values with the latitude and longitude of your allowed area's center
const (
	centerLatitude  = 40.7128
	centerLongitude = -74.0060
)

// Replace this value with the maximum distance in meters that is considered valid
const maxAllowedDistance = 100.0

type glService struct {
}

type GeoLocationService interface {
	ValidateGeolocation(c *gin.Context, longitude, latitude float64) (bool, error)
}

func NewGeoLocationService() *glService {
	return &glService{}
}

func (s *glService) ValidateGeolocation(c *gin.Context, longitude, latitude float64) (bool, error) {
	userLocation := GeoLocation{Latitude: latitude, Longitude: longitude}
	centerLocation := GeoLocation{Latitude: centerLatitude, Longitude: centerLongitude}

	distance := calculateDistance(userLocation, centerLocation)
	if distance > maxAllowedDistance {
		formattedDistance := formatDistanceWithThousandsSeparator(distance)
		message := fmt.Sprint("Location is more than ", formattedDistance, " meters from Office")
		return false, errors.New(message)
	}

	// Verify the user's IP-based geolocation
	clientIP := getClientIP(c)
	ipLocation, err := getIPLocation(clientIP)

	if err != nil {
		return false, err
	}

	ipDistance := calculateDistance(userLocation, ipLocation)
	if ipDistance > maxAllowedDistance {
		formattedDistance := formatDistanceWithThousandsSeparator(distance)
		message := fmt.Sprint("IP Location is more than ", formattedDistance, " meters from Office")
		return false, errors.New(message)
	}

	return true, nil
}

func formatDistanceWithThousandsSeparator(distance float64) string {
	intPart := int(math.Floor(distance))
	intStr := strconv.Itoa(intPart)
	buf := &bytes.Buffer{}

	for i, v := range intStr {
		if i > 0 && (len(intStr)-i)%3 == 0 {
			buf.WriteString(",")
		}
		buf.WriteRune(v)
	}

	return buf.String()
}

func getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Forwarded-For")
	}
	return clientIP
}

// Implement the getIPLocation function
func getIPLocation(ip string) (GeoLocation, error) {
	// Only on development,  is the IPv6 address for localhost, which indicates that you are running the application on your local machine
	// replace with "" or change to hardcoded public IP address 8:8:8:8
	if ip == "::1" {
		ip = ""
	}
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return GeoLocation{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GeoLocation{}, err
	}

	var ipAPIResponse IPAPIResponse
	err = json.Unmarshal(body, &ipAPIResponse)
	if err != nil {
		return GeoLocation{}, err
	}

	return GeoLocation{Latitude: ipAPIResponse.Lat, Longitude: ipAPIResponse.Lon}, nil
}

func calculateDistance(loc1, loc2 GeoLocation) float64 {
	const earthRadius = 6371e3 // Earth's radius in meters

	lat1 := loc1.Latitude * (math.Pi / 180)
	lat2 := loc2.Latitude * (math.Pi / 180)
	deltaLat := (loc2.Latitude - loc1.Latitude) * (math.Pi / 180)
	deltaLng := (loc2.Longitude - loc1.Longitude) * (math.Pi / 180)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
