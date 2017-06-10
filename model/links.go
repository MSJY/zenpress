package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Links struct {
	LinkID          uint64 `gorm:"primary_key"`
	LinkURL         string
	LinkName        string
	LinkImage       string
	LinkTarget      string
	LinkDescription string
	LinkVisible     string //`gorm:"not null default 'Y' index VARCHAR(20)"`
	LinkOwner       uint64 //`gorm:"not null default 1 BIGINT(20)"`
	LinkRating      int
	LinkUpdated     time.Time
	LinkRel         string //XFN关系
	LinkNotes       string //XFN注释
	LinkRss         string
}

// NewLink 创建链接
func NewLink(link *Links) (db *gorm.DB) {
	db = Database.Create(link)
	return
}

// AddLink 新增链接
func AddLink(key, value, linkImage string) (db *gorm.DB) {
	db = Database.Create(&Links{LinkName: key, LinkURL: value, LinkImage: linkImage})
	return
}

// GetLink 获得链接
func GetLink(key string) (db *gorm.DB, Links Links) {
	db = Database.First(&Links, "link_name = ?", key)
	return
}

// UpdateLink 更新链接
func UpdateLink(key, value string) (db *gorm.DB, links Links) {
	links = Links{LinkURL: value}
	db = Database.Model(&links).Update("link_name", key)
	return
}

// DeleteLink 删除链接
func DeleteLink(key string) (db *gorm.DB) {
	db = Database.Delete(Links{}, "link_name = ?", key)
	return
}
