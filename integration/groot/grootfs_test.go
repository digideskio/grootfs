package groot_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("GrootFS global flags", func() {
	Describe("logging", func() {
		Context("when setting the --log-level", func() {
			It("forwards non-human logs to stderr", func() {
				cmd := exec.Command(GrootFSBin, "--log-level", "error", "--store", StorePath, "create", "my-image")
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess).Should(gexec.Exit(1))

				Expect(err).NotTo(HaveOccurred())
				Eventually(sess.Err).Should(gbytes.Say(`"error":"invalid arguments"`))
			})
		})

		Context("when providing a --log-file", func() {
			var (
				logFile *os.File
			)

			BeforeEach(func() {
				logPath, err := ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				logFile, err = ioutil.TempFile(logPath, "mylog")
				Expect(err).NotTo(HaveOccurred())
			})

			It("forwards logs to the given file", func() {
				cmd := exec.Command(GrootFSBin, "--log-level", "debug", "--log-file", logFile.Name(), "--store", StorePath, "create", "my-image")
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess).Should(gexec.Exit(1))

				log, err := ioutil.ReadFile(logFile.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(log).To(ContainSubstring("invalid arguments"))
			})
		})
	})
})
