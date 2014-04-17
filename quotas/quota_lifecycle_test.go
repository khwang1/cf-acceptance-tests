package quotas

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"

//	. "github.com/cloudfoundry/cf-acceptance-tests/helpers"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	"fmt"
)

var _ = Describe("Quota Lifecycle", func() {
	var quotaName string

	BeforeEach(func() {
		quotaName = RandomName()
	})

	Context("Create and Destroy", func() {

		It("should succeed", func() {

			fmt.Printf("\n***************************************\nTest Running")
			Expect(Cf("quotas")).ShouldNot(Say(quotaName))

/*
			Expect(Cf("create-quota", quotaName)).To(Say("OK"))
			Expect(Cf("quotas")).Should(Say(quotaName))

			Expect(Cf("delete-quota", quotaName, "-f")).To(Say("OK"))
			Expect(Cf("quotas")).ShouldNot(Say(quotaName))
*/
		})

	})



	Context("Deleting", func() {

	})

})
