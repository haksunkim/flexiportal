package entity

import "time"

type Post struct {
	Id        uint
	Guid      string
	Title     string
	Content   string
	CreatedBy    uint
	CreatedAt *time.Time
	StrCreatedAt string
	ModifiedBy uint
	ModifiedAt *time.Time
	StrModifiedBy string
	Session Session
}
type Posts []Post
