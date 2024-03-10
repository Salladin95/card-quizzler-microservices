package handlers

const (
	foldersDefaultLimit = 15
	modulesDefaultLimit = 25
)

var FolderKeysMap = map[string]bool{
	"id":         true,
	"title":      true,
	"modules":    true,
	"user_id":    true,
	"created_at": true,
	"updated_at": true,
}

var ModuleKeysMap = map[string]bool{
	"id":         true,
	"title":      true,
	"folders":    true,
	"user_id":    true,
	"created_at": true,
	"updated_at": true,
}
