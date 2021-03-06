package integration

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("proxy", func() {
	var proxyURL, apiRegexp string

	BeforeEach(func() {
		proxyURL = "127.0.0.1:9999"
		apiRegexp = "(?:https:\\/\\/)*" + apiURL
	})

	Context("V2 Legacy", func() {
		It("handles a proxy", func() {
			session := helpers.CFWithEnv(map[string]string{"http_proxy": proxyURL}, "api", apiURL)
			Eventually(session).Should(Say("Error performing request: Get %s\\/v2\\/info: http: error connecting to proxy http:\\/\\/%s", apiRegexp, proxyURL))
			Eventually(session).Should(Exit(1))
		})
	})

	Context("V3", func() {
		It("handles a proxy", func() {
			session := helpers.CFWithEnv(map[string]string{"http_proxy": proxyURL}, "run-task", "app", "echo")
			Eventually(session.Err).Should(Say("Get %s: http: error connecting to proxy http:\\/\\/%s", apiRegexp, proxyURL))
			Eventually(session).Should(Exit(1))
		})
	})
})
