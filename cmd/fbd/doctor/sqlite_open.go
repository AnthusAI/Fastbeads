package doctor

import "github.com/steveyegge/fastbeads/internal/storage"

func sqliteConnString(path string, readOnly bool) string {
	return storage.SQLiteConnString(path, readOnly)
}
