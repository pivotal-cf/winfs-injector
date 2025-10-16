package winfsinjector

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

var (
	readFile  = os.ReadFile
	removeAll = os.RemoveAll
)

type Application struct {
	injector       injector
	releaseCreator releaseCreator
	zipper         zipper
	extractor      extractor
}

//go:generate counterfeiter -generate

//counterfeiter:generate -o ./fakes/injector.go --fake-name Injector . injector

type injector interface {
	AddReleaseToMetadata(releasePath, releaseName, releaseVersion, extractedTileDir string) error
}

//counterfeiter:generate -o ./fakes/zipper.go --fake-name Zipper . zipper

type zipper interface {
	Zip(dir, zipFile string) error
}

//counterfeiter:generate -o ./fakes/extractor.go --fake-name Extractor . extractor

type extractor interface {
	Extract(src, dest string) error
}

//counterfeiter:generate -o ./fakes/release_creator.go --fake-name ReleaseCreator . releaseCreator

type releaseCreator interface {
	CreateRelease(releaseName, imageName, releaseDir, tarballPath, imageTag, registry, version string) error
}

func NewApplication(releaseCreator releaseCreator, injector injector, zipper zipper, extractor extractor) Application {
	return Application{
		injector:       injector,
		releaseCreator: releaseCreator,
		zipper:         zipper,
		extractor:      extractor,
	}
}

func (a Application) Run(inputTile, outputTile, registry, workingDir string) error {
	if inputTile == "" {
		return errors.New("--input-tile is required")
	}

	if outputTile == "" {
		return errors.New("--output-tile is required")
	}

	extractedTileDir := filepath.Join(workingDir, "extracted-tile")
	err := a.extractor.Extract(inputTile, extractedTileDir)
	if err != nil {
		return err
	}

	injectedReleaseTarball := filepath.Join(extractedTileDir, "releases/windows*fs*")
	matches, _ := filepath.Glob(injectedReleaseTarball)
	if len(matches) > 0 {
		fmt.Println("File system has already been injected in the tile; skipping injection")
		return nil
	}

	embeddedReleaseDir := filepath.Join(extractedTileDir, "embed/windowsfs-release")
	if _, err := os.Stat(embeddedReleaseDir); os.IsNotExist(err) {
		fmt.Println("No file system found; skipping injection")
		return nil
	}

	releaseVersion, err := a.extractReleaseVersion(embeddedReleaseDir)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("git", "config", "core.filemode", "false")
		cmd.Dir = embeddedReleaseDir
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unable to fix file permissions for windows: %s, %s", stdoutStderr, err)
		}

		cmd = exec.Command("git", "submodule", "foreach", "git", "config", "core.filemode", "false")
		cmd.Dir = embeddedReleaseDir
		stdoutStderr, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unable to fix file permissions for windows: %s, %s", stdoutStderr, err)
		}
	}

	releaseName, err := a.extractReleaseName(embeddedReleaseDir)
	if err != nil {
		return err
	}

	imageName := "cloudfoundry/windows2016fs"
	imageTag, err := a.determineImageTag(embeddedReleaseDir)
	if err != nil {
		return err
	}

	tarballPath := filepath.Join(extractedTileDir, "releases", fmt.Sprintf("%s-%s.tgz", releaseName, releaseVersion))

	err = a.releaseCreator.CreateRelease(releaseName, imageName, embeddedReleaseDir, tarballPath, imageTag, registry, releaseVersion)
	if err != nil {
		return err
	}

	err = a.injector.AddReleaseToMetadata(tarballPath, releaseName, releaseVersion, extractedTileDir)
	if err != nil {
		return err
	}

	err = removeAll(embeddedReleaseDir)
	if err != nil {
		return err
	}

	return a.zipper.Zip(extractedTileDir, outputTile)
}

func (a Application) extractReleaseVersion(releaseDir string) (string, error) {
	rawReleaseVersion, err := readFile(filepath.Join(releaseDir, "VERSION"))
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(rawReleaseVersion), "\n"), nil
}

func (a Application) extractReleaseName(releaseDir string) (string, error) {
	contents, err := readFile(filepath.Join(releaseDir, "config", "final.yml"))
	if err != nil {
		return "", err
	}

	type NameFile struct {
		Name string `yaml:"name"`
	}

	var f NameFile
	err = yaml.Unmarshal(contents, &f)
	if err != nil {
		return "", err
	}

	return f.Name, nil
}

func (a Application) determineImageTag(releaseDir string) (string, error) {
	var (
		blobs       = map[string]interface{}{}
		blobPath    = filepath.Join(releaseDir, "config", "blobs.yml")
		blobPattern = regexp.MustCompile(`windows.*fs\/windows.*fs-(\d+\.\d+\.\d+)\.tgz`)
	)

	data, err := readFile(blobPath)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(data, &blobs)
	if err != nil {
		return "", err
	}

	for key, _ := range blobs {
		matches := blobPattern.FindStringSubmatch(key)

		if len(matches) == 2 {
			return matches[1], nil
		}
	}

	return "", errors.New("unable to parse tag from embedded rootfs: Please confirm that you are using the appropriate winfs-injector version for this tile")
}
