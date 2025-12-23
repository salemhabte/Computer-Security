package controller

import (
	"net/http"

	domain "security/domain"

	"github.com/gin-gonic/gin"
)

type BackupController struct {
	svc domain.IBackupService
}

func NewBackupController(svc domain.IBackupService) *BackupController {
	return &BackupController{svc: svc}
}

func (bc *BackupController) RunBackup(c *gin.Context) {
	if err := bc.svc.RunBackup(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "backup triggered",
		"lastRun": bc.svc.LastRun(),
	})
}

