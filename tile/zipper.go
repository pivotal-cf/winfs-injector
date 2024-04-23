package tile

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
)

type Zipper struct{}

func NewZipper() Zipper {
	return Zipper{}
}

func (z Zipper) Zip(zipDir, outputFile string) error {
	zipFile := fmt.Sprintf("%s.zip", outputFile)

	err := z.CreateZip(zipDir, zipFile)

	if err != nil {
		return err
	}

	err = os.Rename(zipFile, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func (a Zipper) CreateZip(zipDir string, zipFile string) error {

	_, err := os.Stat(zipDir)

	if err != nil {
		return err
	}

	destinationZip, err := os.Create(zipFile)

	if err != nil {
		return err
	}

	defer destinationZip.Close()

	zipWriter := zip.NewWriter(destinationZip)

	defer zipWriter.Close()

	return filepath.Walk(zipDir, func(filePath string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if zipDir == filePath {
			return nil
		}

		relPath := strings.TrimPrefix(filePath, zipDir+string(os.PathSeparator))

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)

		if err != nil {
			return err
		}

		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		zippedFile, err := zipWriter.CreateHeader(header)

		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)

			if err != nil {
				return err
			}

			defer file.Close()

			_, err = io.Copy(zippedFile, file)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (z Zipper) Unzip(zipFile, outputDir string) error {
	return archiver.DefaultZip.Unarchive(zipFile, outputDir)
}
