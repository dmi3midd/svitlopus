package folder

import "time"

type Folder struct {
	Id        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	ParentId  string    `json:"parentId" db:"parent_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
