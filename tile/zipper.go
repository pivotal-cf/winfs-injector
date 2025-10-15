package tile

import (
	"fmt"
	"os"

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
