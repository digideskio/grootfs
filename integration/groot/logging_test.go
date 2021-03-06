package groot_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/grootfs/integration/runner"
	"code.cloudfoundry.org/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Logging", func() {
	It("forwards human logs to stdout", func() {
		cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", "my-image")
		sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess).Should(gexec.Exit(1))

		Expect(err).NotTo(HaveOccurred())
		Eventually(sess.Out).Should(gbytes.Say("invalid arguments"))
	})

	It("re-logs the nested unpack commands logs", func() {
		imgPath, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(filepath.Join(imgPath, "unreadable-file"), []byte("foo bar"), 0644)).To(Succeed())

		logBuffer := gbytes.NewBuffer()
		cmd := runner.CreateCmd{
			GrootFSBin: GrootFSBin,
			StorePath:  StorePath,
			Spec: groot.CreateSpec{
				ID:    "random-id",
				Image: imgPath,
			},
			LogLevel: lager.DEBUG,
			LogFile:  logBuffer,
		}
		_, err = cmd.Run()
		Expect(err).NotTo(HaveOccurred())

		Expect(logBuffer).To(gbytes.Say("namespaced-unpacking.unpack"))
	})
})
