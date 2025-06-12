package tile

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhoonb/archivex"
)

type Zipper struct{}

func NewZipper() Zipper {
	return Zipper{}
}

func (z Zipper) Zip(zipDir, outputFile string) error {
	zipFile := fmt.Sprintf("%s.zip", outputFile)
	zf := archivex.ZipFile{}

	err := zf.Create(zipFile)
	if err != nil {
		return err
	}

	err = zf.AddAll(zipDir, false)
	if err != nil {
		return err
	}

	err = zf.Close()
	if err != nil {
		return err
	}

	err = os.Rename(zipFile, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func (z Zipper) Unzip(zipFile, outputDir string) error {
	zipFileHandle, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("opening source file: %v", err)
	}
	defer zipFileHandle.Close()

	for _, f := range zipFileHandle.File {
		err := z.extractFile(f, outputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z Zipper) extractFile(f *zip.File, outputDir string) error {
	filePath := filepath.Join(outputDir, f.Name)

	if !strings.HasPrefix(filePath, filepath.Clean(outputDir)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("creating directory: %v", err)
		}

		return nil
	}

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating parent directory: %v", err)
	}

	dstFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating destination file: %v", err)
	}
	defer dstFile.Close()

	err = dstFile.Chmod(f.Mode().Perm())
	if err != nil {
		return fmt.Errorf("setting destination file mode: %v", err)
	}

	srcFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening source file: %v", err)
	}
	defer srcFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copying file: %v", err)
	}

	return nil
}
