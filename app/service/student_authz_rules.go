package service

import (
	"api-students/app/model"
	"api-students/helper"
)

// CanAccessStudent: pemilik data selalu boleh; selain itu butuh permission :any.
func CanAccessStudent(current model.AuthUser, ownerID int, perms *helper.PermissionSet, anyPermission string) bool {
	if current.UserID == ownerID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}