package model

type TermRelationships struct {
	ObjectID       uint64 `gorm:"primary_key"` //对应文章ID/链接ID
	TermTaxonomyID uint64 //对应分类方法ID
	TermOrder      int    //排序
}
