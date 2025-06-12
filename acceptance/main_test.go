package acceptance_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("acceptance", func() {
	Describe("main", func() {
		var (
			winfsInjector string
			cmd           *exec.Cmd
			inputTile     string
			outputTile    string
		)
		BeforeEach(func() {
			var err error
			winfsInjector, err = gexec.Build("github.com/pivotal-cf/winfs-injector")
			Expect(err).ToNot(HaveOccurred())

			inputTile = "input-tile-path"
			outputTile = "output-tile-path"
		})

		AfterEach(func() {
			Expect(os.Remove(winfsInjector)).NotTo(HaveOccurred())
		})

		It("requires an input tile path", func() {
			cmd = exec.Command(winfsInjector, "-o", outputTile)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("--input-tile is required"))
		})

		It("requires an output tile path", func() {
			cmd = exec.Command(winfsInjector, "-i", inputTile)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("--output-tile is required"))
		})

		It("prints the tile extraction directory when the preserve-extracted flag is provided", func() {
			cmd = exec.Command(winfsInjector, "-p")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Out.Contents())).To(ContainSubstring("tile extraction directory"))
		})

		It("prints usage when the help flag is provided", func() {
			cmd = exec.Command(winfsInjector, "--help")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(`
  --input-tile, -i		path to input tile (example: /path/to/input.pivotal)
  --output-tile, -o		path to output tile (example: /path/to/output.pivotal)
  --preserve-extracted, -p	preserve the files created during the tile extraction process (useful for debugging)
  --registry, -r		path to docker registry (example: /path/to/registry, default: "https://registry.hub.docker.com")
  --help, -h			prints this usage information`))
		})
	})
})
