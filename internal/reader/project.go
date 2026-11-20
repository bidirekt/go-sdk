package reader

import (
	"os"
	"path/filepath"
)

type Project struct {
	Dir string
}

func DirToProject(dir string) (Project, bool) {
	for candidate := dir; ; candidate = filepath.Dir(candidate) {
		if info, err := os.Stat(filepath.Join(candidate, "go.mod")); err == nil && !info.IsDir() {
			return Project{Dir: dir}, true
		}
		if filepath.Dir(candidate) == candidate {
			return Project{}, false
		}
	}
}
