package entity

import "github.com/jinzhu/gorm"

type Host struct {
	ID                 string
	Addr               string
	HostName           string
	RunningTasks       string
	Weight             int
	Stop               bool
	LastUpdateTimeUnix int64
	Remark             string
}

func (h *Host) GetByAddr(addr string, db *gorm.DB) error {
	db.Where("addr = ?", addr)
	return db.Find(h).Error
}

func (h *Host) UpdateHeartbeatByAddr(db *gorm.DB) error {
	db.Where("addr = ?", h.Addr)
	return db.Updates(h).Error
}
