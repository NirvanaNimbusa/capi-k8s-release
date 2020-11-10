package binary_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("CF Backup Metadata inputs", func() {
	invalidArgMsgRegex := `invalid arguments, usage: %s or cat cf-metadata-another-system.json \| %s compare`

	When("Invalid arguments has been passed", func() {
		It("should report an error when invalid number of arguments in passed", func() {
			command := interactiveShell(pathToBinary, "compare", "hello", "world") // #nosec
			command.Env = []string{"CF_API_HOST=FAKE_API", "CF_CLIENT=FAKE_CLIENT", "CF_CLIENT_SECRET=FAKE_PASSWORD"}

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say(invalidArgMsgRegex, pathToBinary, pathToBinary))
		})
	})

	DescribeTable("When not all the environment variables are set",
		func(env []string, expectedError string) {
			command := exec.Command(pathToBinary) // #nosec
			command.Env = env

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say(expectedError))
		},
		Entry("CF_CLIENT_SECRET not set",
			[]string{"CF_API_HOST=FAKE_API", "CF_CLIENT=FAKE_CLIENT"},
			"'CF_CLIENT_SECRET' Environment Variable is not set"),
		Entry("CF_CLIENT not set",
			[]string{"CF_API_HOST=FAKE_API", "CF_CLIENT_SECRET=FAKE_PASSWORD"},
			"'CF_CLIENT' Environment Variable is not set"),
		Entry("CF_API_HOST not set",
			[]string{"CF_CLIENT_SECRET=FAKE_PASSWORD", "CF_CLIENT=FAKE_CLIENT"},
			"'CF_API_HOST' Environment Variable is not set"),
	)
})
