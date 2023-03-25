package services

import (
	"encoding/base64"
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"strings"

	"github.com/Kagami/go-face"
)

const (
	modelsDir = "dataset"
)

var rec *face.Recognizer

type frService struct {
}

type FaceRecognitionService interface {
	RecognizeFace(base64Image string) (bool, error)
}

func NewFaceRecognitionService() *frService {
	return &frService{}
}

func init() {
	var err error
	rec, err = face.NewRecognizer(modelsDir)
	if err != nil {
		panic(err)
	}
}

func (s *frService) RecognizeFace(base64Image string) (bool, error) {
	img, err := decodeBase64Image(base64Image)
	if err != nil {
		return false, err
	}

	faces, err := rec.Recognize(img)
	if err != nil {
		return false, err
	}
	if len(faces) != 1 {
		return false, errors.New("No face or multiple faces detected")
	}

	//TODO: Implement verification face by comparing with reference image

	return true, nil
}

func decodeBase64Image(encoded string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	return ioutil.ReadAll(reader)
}
