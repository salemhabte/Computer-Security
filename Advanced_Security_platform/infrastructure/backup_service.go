package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	domain "security/domain"
	"security/config"
)

type BackupService struct {
	lastRun time.Time
}

func NewBackupService() domain.IBackupService {
	return &BackupService{}
}

func (b *BackupService) RunBackup() error {
	now := time.Now()
	filename := fmt.Sprintf("backup-%s.txt", now.Format("20060102-150405"))
	path := filepath.Join(config.BACKUP_PATH, filename)
	if err := os.MkdirAll(config.BACKUP_PATH, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte("backup placeholder at "+now.String()), 0600); err != nil {
		return err
	}
	b.lastRun = now
	return nil
}

func (b *BackupService) LastRun() time.Time {
	return b.lastRun
}

