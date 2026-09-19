package migrations

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed *.sql
var FS embed.FS

// UpFS contains only the migrations that should be applied.
func UpFS() fs.FS {
	return upFS{FS: FS}
}

type upFS struct {
	fs.FS
}

func (f upFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(f.FS, name)
	if err != nil {
		return nil, err
	}

	upEntries := make([]fs.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			upEntries = append(upEntries, entry)
		}
	}
	return upEntries, nil
}
