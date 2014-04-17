package quotas

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cf-acceptance-tests/helpers"
)

func TestQuotas(t *testing.T) {
	RegisterFailHandler(Fail)

	config := helpers.LoadConfig()

	fmt.Printf("\n----------------%#v\n-------------------\n", config)
	helpers.SetupAdminEnvironment(helpers.NewContext(config))

	rs := []Reporter{}

	if config.ArtifactsDirectory != "" {
		os.Setenv(
			"CF_TRACE",
			filepath.Join(
				config.ArtifactsDirectory,
				fmt.Sprintf("CATS-TRACE-%s-%d.txt", "Quotas", ginkgoconfig.GinkgoConfig.ParallelNode),
			),
		)

		rs = append(
			rs,
			reporters.NewJUnitReporter(
				filepath.Join(
					config.ArtifactsDirectory,
					fmt.Sprintf("junit-%s-%d.xml", "Quotas", ginkgoconfig.GinkgoConfig.ParallelNode),
				),
			),
		)
	}

	fmt.Printf("\n-----------------------------\nRUN TESTS\n--------------------------------\n")
	RunSpecsWithDefaultAndCustomReporters(t, "Quotas", rs)
}
