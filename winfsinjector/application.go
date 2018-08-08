package winfsinjector

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	readFile  = ioutil.ReadFile
	removeAll = os.RemoveAll
)

type Application struct {
	injector       injector
	releaseCreator releaseCreator
	zipper         zipper
}

//go:generate counterfeiter -o ./fakes/injector.go --fake-name Injector . injector

type injector interface {
	AddReleaseToMetadata(releasePath, releaseName, releaseVersion, extractedTileDir string) error
}

//go:generate counterfeiter -o ./fakes/zipper.go --fake-name Zipper . zipper

type zipper interface {
	Zip(dir, zipFile string) error
	Unzip(zipFile, dest string) error
}

//go:generate counterfeiter -o ./fakes/release_creator.go --fake-name ReleaseCreator . releaseCreator

type releaseCreator interface {
	CreateRelease(imageName, releaseDir, tarballPath, imageTagPath, versionDataPath string) error
}

func NewApplication(releaseCreator releaseCreator, injector injector, zipper zipper) Application {
	return Application{
		injector:       injector,
		releaseCreator: releaseCreator,
		zipper:         zipper,
	}
}

func (a Application) Run(inputTile, outputTile, workingDir string) error {
	if inputTile == "" {
		return errors.New("--input-tile is required")
	}

	if outputTile == "" {
		return errors.New("--output-tile is required")
	}

	extractedTileDir := filepath.Join(workingDir, "extracted-tile")
	err := a.zipper.Unzip(inputTile, extractedTileDir)
	if err != nil {
		return err
	}

	releaseDir := filepath.Join(extractedTileDir, "embed", "windowsfs-release")
	releaseVersion, err := a.extractReleaseVersion(releaseDir)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("git", "config", "core.filemode", "false")
		cmd.Dir = releaseDir
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unable to fix file permissions for windows: %s, %s", stdoutStderr, err)
		}

		cmd = exec.Command("git", "submodule", "foreach", "git", "config", "core.filemode", "false")
		cmd.Dir = releaseDir
		stdoutStderr, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unable to fix file permissions for windows: %s, %s", stdoutStderr, err)
		}
	}

	// Dependent on what the tile metadata expects, p-windows-runtime-2016/jobs/windows1803fs.yml
	releaseName := "windows1803fs"
	imageName := "cloudfoundry/windows2016fs"
	imageTagPath := filepath.Join(releaseDir, "src", "code.cloudfoundry.org", "windows2016fs", "1803", "IMAGE_TAG")
	tarballPath := filepath.Join(extractedTileDir, "releases", fmt.Sprintf("%s-%s.tgz", releaseName, releaseVersion))

	err = a.releaseCreator.CreateRelease(imageName, releaseDir, tarballPath, imageTagPath, filepath.Join(releaseDir, "VERSION"))
	if err != nil {
		return err
	}

	err = a.injector.AddReleaseToMetadata(tarballPath, releaseName, releaseVersion, extractedTileDir)
	if err != nil {
		return err
	}

	err = removeAll(filepath.Join(extractedTileDir, "embed", "windowsfs-release"))
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
