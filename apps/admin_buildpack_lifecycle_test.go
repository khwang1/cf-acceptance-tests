package apps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	catsHelpers "github.com/cloudfoundry/cf-acceptance-tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	archive_helpers "github.com/pivotal-golang/archiver/extractor/test_helper"
	. "github.com/vito/cmdtest/matchers"
)

var _ = Describe("An application using an admin buildpack", func() {
	var (
		appName       string
		BuildpackName string

		appPath string

		buildpackPath        string
		buildpackArchivePath string
	)

	matchingFilename := func(appName string) string {
		return fmt.Sprintf("simple-buildpack-please-match-%s", appName)
	}

	BeforeEach(func() {
		AsUser(catsHelpers.AdminUserContext, func() {
			BuildpackName = RandomName()
			appName = RandomName()

			tmpdir, err := ioutil.TempDir(os.TempDir(), "matching-app")
			Expect(err).ToNot(HaveOccurred())

			appPath = tmpdir

			tmpdir, err = ioutil.TempDir(os.TempDir(), "matching-buildpack")
			Expect(err).ToNot(HaveOccurred())

			buildpackPath = tmpdir
			buildpackArchivePath = path.Join(buildpackPath, "buildpack.zip")

			archive_helpers.CreateZipArchive(buildpackArchivePath, []archive_helpers.ArchiveFile{
				{
					Name: "bin/compile",
					Body: `#!/usr/bin/env bash


echo "Staging with Simple Buildpack"

sleep 2
`,
				},
				{
					Name: "bin/detect",
					Body: fmt.Sprintf(`#!/bin/bash

if [ -f "${1}/%s" ]; then
  echo Simple
else
  echo no
  exit 1
fi
`, matchingFilename(appName)),
				},
				{
					Name: "bin/release",
					Body: `#!/usr/bin/env bash

cat <<EOF
---
config_vars:
  PATH: bin:/usr/local/bin:/usr/bin:/bin
  FROM_BUILD_PACK: "yes"
default_process_types:
  web: while true; do { echo -e 'HTTP/1.1 200 OK\r\n'; echo "hi from a simple admin buildpack"; } | nc -l \$PORT; done
EOF
`,
				},
			})

			_, err = os.Create(path.Join(appPath, matchingFilename(appName)))
			Expect(err).ToNot(HaveOccurred())

			_, err = os.Create(path.Join(appPath, "some-file"))
			Expect(err).ToNot(HaveOccurred())

			createBuildpack := Cf("create-buildpack", BuildpackName, buildpackArchivePath, "0")
			Expect(createBuildpack).To(Say("Creating"))
			Expect(createBuildpack).To(Say("OK"))
			Expect(createBuildpack).To(Say("Uploading"))
			Expect(createBuildpack).To(Say("OK"))
		})
	})

	AfterEach(func() {
		AsUser(catsHelpers.AdminUserContext, func() {
			Expect(Cf("delete-buildpack", BuildpackName, "-f")).To(Say("OK"))
		})
	})

	Context("when the buildpack is detected", func() {
		It("is used for the app", func() {
			push := Cf("push", appName, "-p", appPath)
			Expect(push).To(Say("Staging with Simple Buildpack"))
			Expect(push).To(Say("App started"))
		})
	})

	Context("when the buildpack fails to detect", func() {
		BeforeEach(func() {
			err := os.Remove(path.Join(appPath, matchingFilename(appName)))
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails to stage", func() {
			Expect(Cf("push", appName, "-p", appPath)).To(Say("Staging error"))
		})
	})

	Context("when the buildpack is deleted", func() {
		BeforeEach(func() {
			AsUser(catsHelpers.AdminUserContext, func() {
				Expect(Cf("delete-buildpack", BuildpackName, "-f")).To(Say("OK"))
			})
		})

		It("fails to stage", func() {
			Expect(Cf("push", appName, "-p", appPath)).To(Say("Staging error"))
		})
	})

	Context("when the buildpack is disabled", func() {
		BeforeEach(func() {
			AsUser(catsHelpers.AdminUserContext, func() {
				var response QueryResponse

				ApiRequest("GET", "/v2/buildpacks?q=name:"+BuildpackName, &response)

				Expect(response.Resources).To(HaveLen(1))

				buildpackGuid := response.Resources[0].Metadata.Guid

				ApiRequest(
					"PUT",
					"/v2/buildpacks/"+buildpackGuid,
					nil,
					`{"enabled":false}`,
				)
			})
		})

		It("fails to stage", func() {
			Expect(Cf("push", appName, "-p", appPath)).To(Say("Staging error"))
		})
	})
})
