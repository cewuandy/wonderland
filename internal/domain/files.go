package domain

type FileRepo interface {
	ListFilesContent() (map[string][]*string, error)
}
