package files

import (
	"bufio"
	"github.com/samber/do"
	"os"
	"path/filepath"

	"github.com/cewuandy/wonderland/internal/domain"
)

type fileRepo struct {
	path string
}

func (f *fileRepo) ListFilesContent() (map[string][]*string, error) {
	var (
		dirEntry []os.DirEntry
		file     *os.File
		result   = map[string][]*string{}
		err      error
	)

	dirEntry, err = os.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dirEntry); i++ {
		toolName := dirEntry[i].Name()
		file, err = os.Open(filepath.Join(f.path, dirEntry[i].Name()))
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			result[toolName] = append(result[toolName], &line)
		}
		_ = file.Close()
	}

	return result, nil
}

func NewFileRepo(i *do.Injector) (domain.FileRepo, error) {
	return &fileRepo{do.MustInvokeNamed[string](i, "material_path")}, nil
}
