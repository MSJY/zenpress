package model

type Sitemeta struct {
	MetaID    uint64 `gorm:"primary_key"`
	SiteID    uint64 `gorm:"not null default 0 index BIGINT(20)"`
	MetaKey   string `gorm:"index VARCHAR(255)"`
	MetaValue string `gorm:"LONGTEXT"`
}
