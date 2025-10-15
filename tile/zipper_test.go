package tile_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/pivotal-cf/winfs-injector/tile"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Zipper", func() {
	Describe("Zip", func() {
		var (
			zipper  tile.Zipper
			srcDir  string
			zipFile *os.File
		)

		BeforeEach(func() {
			zipper = tile.NewZipper()

			var err error
			srcDir, err = os.MkdirTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			zipFile, err = os.CreateTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(srcDir, "top-level-file"), []byte("foo"), os.FileMode(0644))
			Expect(err).NotTo(HaveOccurred())

			err = os.Mkdir(filepath.Join(srcDir, "second-level-dir"), os.FileMode(0755))
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(srcDir, "second-level-dir", "second-level-file"), []byte("bar"), os.FileMode(0644))
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(srcDir)
			Expect(err).NotTo(HaveOccurred())

			err = zipFile.Close()
			Expect(err).NotTo(HaveOccurred())

			err = os.RemoveAll(zipFile.Name())
			Expect(err).NotTo(HaveOccurred())
		})

		It("zips the specified directory and creates a zip at the specified path", func() {
			err := zipper.Zip(srcDir, zipFile.Name())
			Expect(err).NotTo(HaveOccurred())

			actualZip, err := zip.OpenReader(zipFile.Name())
			Expect(err).NotTo(HaveOccurred())

			Expect(actualZip.File).To(HaveLen(3))

			fileAssertions := map[string]string{
				"top-level-file":   "foo",
				"second-level-dir": "",
				filepath.Join("second-level-dir", "second-level-file"): "bar",
			}

			for _, f := range actualZip.File {
				openedFile, err := f.Open()
				Expect(err).NotTo(HaveOccurred())

				defer openedFile.Close()

				fileContents, err := io.ReadAll(openedFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(fileContents)).To(Equal(fileAssertions[f.Name]))
			}
		})

		Context("failure cases", func() {
			Context("when an intermediate dir in the destination path does not exist", func() {
				It("returns an error", func() {
					err := zipper.Zip(srcDir, "/path/to/non-existing/dir")
					Expect(err).To(MatchError(ContainSubstring("/path/to/non-existing/dir")))
				})
			})

			Context("when the source dir does not exist", func() {
				It("returns an error", func() {
					err := zipper.Zip("/path/to/non-existing/dir", zipFile.Name())
					Expect(err).To(MatchError(ContainSubstring("/path/to/non-existing/dir")))
				})
			})
		})
	})
})
