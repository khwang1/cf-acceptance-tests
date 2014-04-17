package helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"

	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"fmt"
)

var AdminUserContext cf.UserContext
var RegularUserContext cf.UserContext

type SuiteContext interface {
	Setup()
	Teardown()

	AdminUserContext() cf.UserContext
	RegularUserContext() cf.UserContext
}

func SetupEnvironment(context SuiteContext) {
	var originalCfHomeDir, currentCfHomeDir string

	BeforeEach(func() {
		fmt.Printf("\n---------------\nSetupEnvironment:BeforeEach...\n---------------------------\n")
		AdminUserContext = context.AdminUserContext()
		RegularUserContext = context.RegularUserContext()

		context.Setup()

		fmt.Printf("\n---------------\nSetupEnvironment:BeforeEach:set up space as regular user...\n---------------------------\n")
		cf.AsUser(AdminUserContext, func() {
			setUpSpaceWithUserAccess(RegularUserContext)
		})

		fmt.Printf("\n---------------\nSetupEnvironment:BeforeEach:target space as regular user...\n---------------------------\n")
		originalCfHomeDir, currentCfHomeDir = cf.InitiateUserContext(RegularUserContext)
		cf.TargetSpace(RegularUserContext)
	})

	AfterEach(func() {
		fmt.Printf("\n---------------\nSetupEnvironment:AfterEach...\n---------------------------\n")
		cf.RestoreUserContext(RegularUserContext, originalCfHomeDir, currentCfHomeDir)

		context.Teardown()
	})
}

func SetupAdminEnvironment(context SuiteContext) {
	var originalCfHomeDir, currentCfHomeDir string

	BeforeEach(func() {
		fmt.Printf("\n---------------\nSetupAdminEnvironment:BeforeEach...\n---------------------------\n")
		AdminUserContext = context.AdminUserContext()

		context.Setup()
	})

	AfterEach(func() {
			fmt.Printf("\n---------------\nSetupAdminEnvironment:AfterEach...\n---------------------------\n")
			cf.RestoreUserContext(RegularUserContext, originalCfHomeDir, currentCfHomeDir)

			context.Teardown()
	})

}

func setUpSpaceWithUserAccess(uc cf.UserContext) {
	Expect(cf.Cf("create-space", "-o", uc.Org, uc.Space)).To(ExitWith(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceManager")).To(ExitWith(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceDeveloper")).To(ExitWith(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceAuditor")).To(ExitWith(0))
}
