package models

import "time"

type File struct {
	Id        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Mime      string    `json:"mime" db:"mime"`
	Size      int       `json:"size" db:"size"`
	ParentId  string    `json:"parentId" db:"parent_id"`
	FileId    string    `json:"fileId" db:"file_id"`
	MessageId int       `json:"MessageId" db:"message_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
