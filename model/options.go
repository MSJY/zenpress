package model

import (
	"github.com/jinzhu/gorm"
)

type Options struct {
	OptionID    uint64 `gorm:"primary_key"`
	OptionName  string
	OptionValue string //`gorm:"not null LONGTEXT"`
	Autoload    string
}

/*
这是一种向选项数据库表中添加有名称的选项/值对的安全方法。
如果所需选项已存在，add_option()不添加内容。
选项被保存后，可通过get_option()来访问选项，
通过update_option()来修改选项，还可以通过delete_option()删除该选项。
*/

// AddOption 新增选项
func AddOption(key, value string, autoload ...string) (db *gorm.DB) {
	var iAutoload = "yes"
	if len(autoload) > 0 {
		if !((autoload[0] != "yes") && (autoload[0] != "no")) {
			iAutoload = autoload[0]
		}
	}
	db = Database.Create(&Options{OptionName: key, OptionValue: value, Autoload: iAutoload})
	return
}

// GetOption 获得选项
func GetOption(key string) (db *gorm.DB, options Options) {
	db = Database.First(&options, "option_name = ?", key)
	return
}

// UpdateOption 更新选项
func UpdateOption(key, value string) (db *gorm.DB, options Options) {
	options = Options{OptionValue: value}
	db = Database.Model(&options).Update("option_name", key)
	return
}

// DeleteOption 删除选项
func DeleteOption(key string) (db *gorm.DB) {
	/*
		var opt Options
		db, opt = GetOption(key)
		db = Database.Delete(opt)
	*/
	db = Database.Delete(Options{}, "option_name = ?", key)
	return
}
