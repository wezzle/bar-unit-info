package util

import "embed"

var repoFiles embed.FS

func InitFS(fs embed.FS) {
	repoFiles = fs
}

func IgnoreError[T any](key string, f func(string) (T, error)) T {
	v, _ := f(key)
	// TODO We should probably write these errors to somewhere
	return v
}
