package util

import "embed"

var repoFiles embed.FS

func InitFS(fs embed.FS) {
	repoFiles = fs
}
