package entity

import "mime/multipart"

type BindFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
