package tile_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tile Suite")
}
