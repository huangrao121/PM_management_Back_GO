package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func SavetoLocalWithNewName(c *gin.Context, jsonName string, folderName string) (string, error) {
	file, header, err := c.Request.FormFile(jsonName)
	if err != nil {
		log.Error("Happened error getting image from request. Error", err)
		return "", err
	}
	defer file.Close()

	//Create UUID string
	uid := uuid.New().String()

	//Get original file name
	newFilename := fmt.Sprintf("%s_%s", uid, header.Filename)
	dst := fmt.Sprintf("public/%s/%s", folderName, newFilename)
	if save_err := c.SaveUploadedFile(header, dst); save_err != nil {
		log.Error("Happened error saving image to local disk. Error", err)
		return "", save_err
	}
	return newFilename, nil
}
