package tracer

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressFile(source string, target string) error {

	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer targetFile.Close()

	gzipWriter := gzip.NewWriter(targetFile)

	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, sourceFile)

	return err
}