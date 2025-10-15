package winfsinjector_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pivotal-cf/winfs-injector/winfsinjector"
	"github.com/pivotal-cf/winfs-injector/winfsinjector/fakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("application", func() {
	Describe("Run", func() {
		var (
			fakeReleaseCreator *fakes.ReleaseCreator
			fakeInjector       *fakes.Injector
			fakeZipper         *fakes.Zipper
			fakeExtractor      *fakes.Extractor

			inputTile  string
			outputTile string
			registry   string
			workingDir string

			app winfsinjector.Application

			err error
		)

		BeforeEach(func() {
			fakeReleaseCreator = new(fakes.ReleaseCreator)
			fakeInjector = new(fakes.Injector)
			fakeZipper = new(fakes.Zipper)
			fakeExtractor = new(fakes.Extractor)

			inputTile = "/path/to/input/tile"
			outputTile = "/path/to/output/tile"
			registry = "/path/to/docker/registry"

			workingDir, err = os.MkdirTemp("", "")
			Expect(err).ToNot(HaveOccurred())

			embedFilePath := fmt.Sprintf("%s/extracted-tile/embed", workingDir)
			err := os.MkdirAll(embedFilePath, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			winfsinjector.SetReadFile(func(path string) ([]byte, error) {
				switch filepath.Base(path) {
				case "VERSION":
					// reading VERSION file
					return []byte("9.3.6"), nil
				case "blobs.yml":
					// reading config/blobs.yml
					return []byte(`---
windows2019fs/windows2016fs-2019.0.43.tgz:
  size: 3333333333
  sha: abcdefg1234
`), nil
				case "final.yml":
					// reading config/final.yml
					return []byte(`name: windows2019fs`), nil
				default:
					return nil, errors.New("readFile called for unexpected input: " + path)
				}
			})

			err = os.MkdirAll(embedFilePath+"/windowsfs-release", os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			app = winfsinjector.NewApplication(fakeReleaseCreator, fakeInjector, fakeZipper, fakeExtractor)
		})

		AfterEach(func() {
			winfsinjector.ResetReadFile()
			winfsinjector.ResetRemoveAll()
		})

		It("unzips the tile", func() {
			err := app.Run(inputTile, outputTile, registry, workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeExtractor.ExtractCallCount()).To(Equal(1))

			inputTile, extractedTileDir := fakeExtractor.ExtractArgsForCall(0)
			Expect(inputTile).To(Equal(filepath.Join("/", "path", "to", "input", "tile")))
			Expect(extractedTileDir).To(Equal(fmt.Sprintf("%s%s", workingDir, filepath.Join("/", "extracted-tile"))))
		})

		It("creates the release", func() {
			err := app.Run(inputTile, outputTile, registry, workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseCreator.CreateReleaseCallCount()).To(Equal(1))

			releaseName, imageName, releaseDir, tarballPath, imageTag, registry, version := fakeReleaseCreator.CreateReleaseArgsForCall(0)
			Expect(releaseName).To(Equal("windows2019fs"))
			Expect(imageName).To(Equal("cloudfoundry/windows2016fs"))
			Expect(releaseDir).To(Equal(fmt.Sprintf("%s/extracted-tile/embed/windowsfs-release", workingDir)))
			Expect(tarballPath).To(Equal(fmt.Sprintf("%s/extracted-tile/releases/windows2019fs-9.3.6.tgz", workingDir)))
			Expect(imageTag).To(Equal("2019.0.43"))
			Expect(registry).To(Equal("/path/to/docker/registry"))
			Expect(version).To(Equal("9.3.6"))
		})

		It("injects the build windows release into the extracted tile", func() {
			err := app.Run(inputTile, outputTile, registry, workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseCreator.CreateReleaseCallCount()).To(Equal(1))
			Expect(fakeExtractor.ExtractCallCount()).To(Equal(1))

			Expect(fakeInjector.AddReleaseToMetadataCallCount()).To(Equal(1))
			releasePath, releaseName, releaseVersion, tileDir := fakeInjector.AddReleaseToMetadataArgsForCall(0)
			Expect(releasePath).To(Equal(fmt.Sprintf("%s/extracted-tile/releases/windows2019fs-9.3.6.tgz", workingDir)))
			Expect(releaseName).To(Equal("windows2019fs"))
			Expect(releaseVersion).To(Equal("9.3.6"))
			Expect(tileDir).To(Equal(filepath.Join(workingDir, "extracted-tile")))
		})

		It("removes the windows2016fs-release from the embed directory", func() {
			var (
				removeAllCallCount int
				removeAllPath      string
			)

			winfsinjector.SetRemoveAll(func(path string) error {
				removeAllCallCount++
				removeAllPath = path
				return nil
			})

			err := app.Run(inputTile, outputTile, registry, workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(removeAllCallCount).To(Equal(1))
			Expect(removeAllPath).To(Equal(fmt.Sprintf("%s%s", workingDir, filepath.Join("/", "extracted-tile", "embed", "windowsfs-release"))))
		})

		It("zips up the injected tile dir", func() {
			err := app.Run(inputTile, outputTile, registry, workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseCreator.CreateReleaseCallCount()).To(Equal(1))
			Expect(fakeExtractor.ExtractCallCount()).To(Equal(1))
			Expect(fakeInjector.AddReleaseToMetadataCallCount()).To(Equal(1))

			Expect(fakeZipper.ZipCallCount()).To(Equal(1))
			zipDir, zipFile := fakeZipper.ZipArgsForCall(0)
			Expect(zipDir).To(Equal(filepath.Join(workingDir, "extracted-tile")))
			Expect(zipFile).To(Equal("/path/to/output/tile"))
		})

		Context("when the image tag of release dir is malformed", func() {
			BeforeEach(func() {
				winfsinjector.SetReadFile(func(path string) ([]byte, error) {
					switch filepath.Base(path) {
					case "VERSION":
						// reading VERSION file
						return []byte(""), nil
					case "blobs.yml":
						// reading config/blobs.yml
						return []byte(`---
windows2019fs/windows2016fs-MISSING-IMAGE-TAG.tgz:
  size: 3333333333
  sha: abcdefg1234
`), nil
					case "final.yml":
						// reading config/final.yml
						return []byte(``), nil
					default:
						return nil, errors.New("readFile called for unexpected input: " + path)
					}
				})
			})

			It("returns the error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(MatchError(ContainSubstring("unable to parse tag from embedded rootfs:")))
			})
		})

		Context("when the extractor fails to unzip the tile", func() {
			BeforeEach(func() {
				fakeExtractor.ExtractReturns(errors.New("some-error"))
			})

			It("returns the error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when the injector fails to copy the release into the tile", func() {
			BeforeEach(func() {
				fakeInjector.AddReleaseToMetadataReturns(errors.New("some-error"))
			})

			It("returns the error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when there is multiple directories embedded in the tile", func() {
			BeforeEach(func() {
				embedFilePath := fmt.Sprintf("%s/extracted-tile/embed", workingDir)
				err = os.MkdirAll(embedFilePath+"/this-is-another-folder", os.ModePerm)
				Expect(err).ToNot(HaveOccurred())
			})

			It("creates a release with the windowsfs-release", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeReleaseCreator.CreateReleaseCallCount()).To(Equal(1))

				releaseName, imageName, releaseDir, tarballPath, imageTag, registry, version := fakeReleaseCreator.CreateReleaseArgsForCall(0)
				Expect(releaseName).To(Equal("windows2019fs"))
				Expect(imageName).To(Equal("cloudfoundry/windows2016fs"))
				Expect(releaseDir).To(Equal(fmt.Sprintf("%s/extracted-tile/embed/windowsfs-release", workingDir)))
				Expect(tarballPath).To(Equal(fmt.Sprintf("%s/extracted-tile/releases/windows2019fs-9.3.6.tgz", workingDir)))
				Expect(imageTag).To(Equal("2019.0.43"))
				Expect(registry).To(Equal("/path/to/docker/registry"))
				Expect(version).To(Equal("9.3.6"))
			})
		})

		Context("when windowsfs-release is already embedded in the tile", func() {
			BeforeEach(func() {
				releaseDir := fmt.Sprintf("%s/extracted-tile/releases", workingDir)
				os.MkdirAll(releaseDir, os.ModePerm)

				tarballPath := filepath.Join(releaseDir, "/windows2019fs-9.3.6.tgz")
				os.WriteFile(tarballPath, []byte("windows"), 0644)
			})

			It("does not return an error and exits", func() {
				var err error
				r, w, _ := os.Pipe()
				tmp := os.Stdout
				defer func() {
					os.Stdout = tmp
				}()
				os.Stdout = w

				err = app.Run(inputTile, outputTile, registry, workingDir)
				w.Close()

				Expect(err).ToNot(HaveOccurred())
				stdout, _ := io.ReadAll(r)
				Expect(string(stdout)).To(ContainSubstring("File system has already been injected in the tile; skipping injection"))
			})
		})

		Context("when windowsfs-release is not embedded in the tile from being previously hydrated", func() {
			BeforeEach(func() {
				embedFilePath := fmt.Sprintf("%s/extracted-tile/embed", workingDir)
				os.RemoveAll(embedFilePath + "/windowsfs-release")
			})

			It("does not return an error and exits", func() {
				var err error
				r, w, _ := os.Pipe()
				tmp := os.Stdout
				defer func() {
					os.Stdout = tmp
				}()
				os.Stdout = w
				err = app.Run(inputTile, outputTile, registry, workingDir)
				w.Close()

				Expect(err).ToNot(HaveOccurred())
				stdout, _ := io.ReadAll(r)
				Expect(string(stdout)).To(ContainSubstring("No file system found; skipping injection"))
			})
		})

		Context("when the release creator fails", func() {
			BeforeEach(func() {
				fakeReleaseCreator.CreateReleaseReturns(errors.New("some-error"))
			})

			It("returns the error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when removing the windows2016fs-release dir from the embed directory fails", func() {
			BeforeEach(func() {
				winfsinjector.SetRemoveAll(func(path string) error {
					return errors.New("remove all failed")
				})
			})

			It("returns an error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(MatchError("remove all failed"))
			})
		})

		Context("when zipping the injected tile dir fails", func() {
			BeforeEach(func() {
				fakeZipper.ZipReturns(errors.New("some-error"))
			})

			It("returns the error", func() {
				err := app.Run(inputTile, outputTile, registry, workingDir)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when input tile is not provided", func() {
			It("returns an error", func() {
				err := app.Run("", outputTile, registry, workingDir)
				Expect(err).To(MatchError("--input-tile is required"))
			})
		})

		Context("when output tile is not provided", func() {
			It("returns an error", func() {
				err := app.Run(inputTile, "", registry, workingDir)
				Expect(err).To(MatchError("--output-tile is required"))
			})
		})
	})
})
