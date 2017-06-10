package model

type RolePermissions struct {
	ID           uint64 `gorm:"primary_key"`
	RoleID       uint64
	PermissionID uint64
}
