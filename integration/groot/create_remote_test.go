package groot_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/grootfs/integration"
	"code.cloudfoundry.org/grootfs/store"
	"code.cloudfoundry.org/grootfs/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	specsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var _ = Describe("Create with remote images", func() {
	var imageURL string

	Context("when using the default registry", func() {
		BeforeEach(func() {
			imageURL = "docker:///cfgarden/empty:v0.1.0"
		})

		It("creates a root filesystem based on the image provided", func() {
			bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)

			Expect(path.Join(bundle.RootFSPath(), "hello")).To(BeARegularFile())
		})

		It("saves the image.json to the bundle folder", func() {
			bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)

			imageJsonPath := path.Join(bundle.Path(), "image.json")
			Expect(imageJsonPath).To(BeARegularFile())

			imageJsonReader, err := os.Open(imageJsonPath)
			Expect(err).ToNot(HaveOccurred())
			var imageJson specsv1.Image
			Expect(json.NewDecoder(imageJsonReader).Decode(&imageJson)).To(Succeed())

			Expect(imageJson.Created).To(Equal("2016-08-03T16:50:55.797615406Z"))
			Expect(imageJson.RootFS.DiffIDs).To(Equal([]string{
				"sha256:3355e23c079e9b35e4b48075147a7e7e1850b99e089af9a63eed3de235af98ca",
			}))
		})

		Describe("OCI image caching", func() {
			It("caches the image in the store", func() {
				integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)

				blobPath := path.Join(
					StorePath, "cache", "blobs",
					"sha256-6c1f4533b125f8f825188c4f4ff633a338cfce0db2813124d3d518028baf7d7a",
				)
				Expect(blobPath).To(BeARegularFile())
			})

			It("uses the cached image from the store", func() {
				_, err := integration.CreateBundleWSpec(GrootFSBin, StorePath, DraxBin, groot.CreateSpec{
					ID:    "random-id",
					Image: imageURL,
					UIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{NamespaceID: 0, HostID: os.Getuid(), Size: 1},
						groot.IDMappingSpec{NamespaceID: 1, HostID: 100000, Size: 65000},
					},
					GIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{NamespaceID: 0, HostID: os.Getgid(), Size: 1},
						groot.IDMappingSpec{NamespaceID: 1, HostID: 100000, Size: 65000},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// change the cache
				blobPath := path.Join(
					StorePath, "cache", "blobs",
					"sha256-6c1f4533b125f8f825188c4f4ff633a338cfce0db2813124d3d518028baf7d7a",
				)

				buffer := bytes.NewBuffer([]byte{})
				tarWriter := tar.NewWriter(buffer)
				Expect(tarWriter.WriteHeader(&tar.Header{
					Name: "i-hacked-your-cache",
					Mode: 0666,
					Size: int64(len([]byte("cache-hit!"))),
				})).To(Succeed())
				_, err = tarWriter.Write([]byte("cache-hit!"))
				Expect(err).NotTo(HaveOccurred())
				Expect(tarWriter.Close()).To(Succeed())

				blob, err := os.OpenFile(blobPath, os.O_WRONLY, 0666)
				Expect(err).NotTo(HaveOccurred())
				gzip := gzip.NewWriter(blob)
				blobContents, err := ioutil.ReadAll(buffer)
				Expect(err).NotTo(HaveOccurred())
				_, err = gzip.Write(blobContents)
				Expect(err).NotTo(HaveOccurred())
				Expect(gzip.Close()).To(Succeed())

				bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id-2", 0)
				Expect(path.Join(bundle.RootFSPath(), "i-hacked-your-cache")).To(BeARegularFile())
			})

			Context("when the image has opaque white outs", func() {
				BeforeEach(func() {
					imageURL = "docker:///cfgarden/opq-whiteout-busybox"
				})

				It("empties the folder contents but keeps the dir", func() {
					bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)

					whiteoutedDir := path.Join(bundle.RootFSPath(), "var")
					Expect(whiteoutedDir).To(BeADirectory())
					contents, err := ioutil.ReadDir(whiteoutedDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(contents).To(BeEmpty())
				})
			})

			Context("when image size exceeds disk quota", func() {
				BeforeEach(func() {
					imageURL = "docker:///cfgarden/empty:v0.1.1"
				})

				Context("when the image is not accounted for in the quota", func() {
					It("succeeds", func() {
						cmd := exec.Command(GrootFSBin, "--store", StorePath, "--drax-bin", DraxBin, "create", "--disk-limit-size-bytes", "10", "--exclude-image-from-quota", imageURL, "random-id")
						sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, 12*time.Second).Should(gexec.Exit(0))
					})
				})

				Context("when the image is accounted for in the quota", func() {
					It("returns an error", func() {
						cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", "--disk-limit-size-bytes", "10", imageURL, "random-id")

						sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())
						Eventually(sess, 12*time.Second).Should(gexec.Exit(1))

						Eventually(sess).Should(gbytes.Say("layers exceed disk quota"))
					})
				})
			})

			Describe("Unpacked layer caching", func() {
				It("caches the unpacked image in a subvolume with snapshots", func() {
					integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)

					layerSnapshotPath := filepath.Join(StorePath, "volumes", "sha256:3355e23c079e9b35e4b48075147a7e7e1850b99e089af9a63eed3de235af98ca")
					Expect(ioutil.WriteFile(layerSnapshotPath+"/injected-file", []byte{}, 0666)).To(Succeed())

					bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id-2", 0)
					Expect(path.Join(bundle.RootFSPath(), "hello")).To(BeARegularFile())
					Expect(path.Join(bundle.RootFSPath(), "injected-file")).To(BeARegularFile())
				})

				Describe("when unpacking the image fails", func() {
					It("deletes the layer volume cache", func() {
						// write an invalid cached layer blob
						blobPath := path.Join(
							StorePath, "cache", "blobs",
							"sha256-6c1f4533b125f8f825188c4f4ff633a338cfce0db2813124d3d518028baf7d7a",
						)
						Expect(os.MkdirAll(path.Join(StorePath, "cache", "blobs"), 0755)).To(Succeed())
						Expect(ioutil.WriteFile(blobPath, []byte("corrupted"), 0666)).To(Succeed())

						layerSnapshotPath := filepath.Join(StorePath, "volumes", "sha256:3355e23c079e9b35e4b48075147a7e7e1850b99e089af9a63eed3de235af98ca")
						cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", imageURL, "random-id-2")
						sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
						Expect(err).ToNot(HaveOccurred())
						Eventually(sess, 12*time.Second).Should(gexec.Exit(1))

						// layer should be cleaned up
						Expect(layerSnapshotPath).ToNot(BeAnExistingFile())
					})
				})
			})
		})

		Context("when the image has a version 1 manifest schema", func() {
			BeforeEach(func() {
				imageURL = "docker:///cfgarden/empty:schemaV1"
			})

			It("creates a root filesystem based on the image provided", func() {
				bundle := integration.CreateBundle(GrootFSBin, StorePath, DraxBin, imageURL, "random-id", 0)
				Expect(path.Join(bundle.RootFSPath(), "allo")).To(BeAnExistingFile())
				Expect(path.Join(bundle.RootFSPath(), "hello")).To(BeAnExistingFile())
			})
		})
	})

	Context("when a private registry is used", func() {
		var (
			fakeRegistry *testhelpers.FakeRegistry
		)

		BeforeEach(func() {
			dockerHubUrl, err := url.Parse("https://registry-1.docker.io")
			Expect(err).NotTo(HaveOccurred())
			fakeRegistry = testhelpers.NewFakeRegistry(dockerHubUrl)
			Expect(fakeRegistry.Start()).To(Succeed())

			imageURL = fmt.Sprintf("docker://%s/cfgarden/empty:v0.1.1", fakeRegistry.Addr())
		})

		AfterEach(func() {
			fakeRegistry.Stop()
		})

		It("fails to fetch the image", func() {
			_, err := integration.CreateBundleWSpec(GrootFSBin, StorePath, DraxBin, groot.CreateSpec{
				ID:    "random-id",
				Image: imageURL,
			})

			Eventually(err).Should(MatchError("This registry is insecure. To pull images from this registry, please use the --insecure-registry option."))
		})

		Context("when it's provided as a valid insecure registry", func() {
			It("should create a root filesystem based on the image provided by the private registry", func() {
				cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", "--insecure-registry", fakeRegistry.Addr(), imageURL, "random-id")
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess, 12*time.Second).Should(gexec.Exit(0))

				rootFSPath := strings.TrimSpace(string(sess.Out.Contents())) + "/rootfs"
				Expect(path.Join(rootFSPath, "hello")).To(BeARegularFile())

				Expect(fakeRegistry.RequestedBlobs()).To(HaveLen(3))
			})
		})
	})

	Context("when the image does not exist", func() {
		It("returns a useful error", func() {
			_, err := integration.CreateBundleWSpec(GrootFSBin, StorePath, DraxBin, groot.CreateSpec{
				ID:    "random-id",
				Image: "docker:///cfgaren/sorry-not-here",
			})

			Eventually(err).Should(MatchError("Image does not exist or you do not have permissions to see it."))
		})
	})

	Context("when downloading the layer fails", func() {
		var (
			fakeRegistry *testhelpers.FakeRegistry
			layerID      string
		)

		BeforeEach(func() {
			dockerHubUrl, err := url.Parse("https://registry-1.docker.io")
			Expect(err).NotTo(HaveOccurred())
			fakeRegistry = testhelpers.NewFakeRegistry(dockerHubUrl)

			layerID = "6c1f4533b125f8f825188c4f4ff633a338cfce0db2813124d3d518028baf7d7a"
			fakeRegistry.WhenGettingBlob(layerID, 1, func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte("i-am-groot"))
			})

			Expect(fakeRegistry.Start()).To(Succeed())

			imageURL = fmt.Sprintf("docker://%s/cfgarden/empty:v0.1.0", fakeRegistry.Addr())
		})

		AfterEach(func() {
			fakeRegistry.Stop()
		})

		It("does not leak corrupted state", func() {
			cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", "--insecure-registry", fakeRegistry.Addr(), imageURL, "random-id")
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess, 12*time.Second).Should(gexec.Exit(1))
			Expect(path.Join(StorePath, "cache", "blobs", fmt.Sprintf("sha256-%s", layerID))).NotTo(BeARegularFile())

			// Can still be used succesfully later
			cmd = exec.Command(GrootFSBin, "--store", StorePath, "create", "--insecure-registry", fakeRegistry.Addr(), imageURL, "random-id")
			sess, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess, 12*time.Second).Should(gexec.Exit(0))
			Expect(path.Join(StorePath, "cache", "blobs", fmt.Sprintf("sha256-%s", layerID))).To(BeARegularFile())
		})
	})

	Context("when a layer is corrupted", func() {
		var (
			fakeRegistry *testhelpers.FakeRegistry
			layerID      string
		)

		BeforeEach(func() {
			dockerHubUrl, err := url.Parse("https://registry-1.docker.io")
			Expect(err).NotTo(HaveOccurred())
			fakeRegistry = testhelpers.NewFakeRegistry(dockerHubUrl)

			tarBuffer := bytes.NewBuffer([]byte{})
			tarWriter := tar.NewWriter(tarBuffer)
			Expect(tarWriter.WriteHeader(&tar.Header{
				Name: "a_file",
				Size: 11,
			})).To(Succeed())
			_, err = tarWriter.Write([]byte("hello-world"))
			Expect(err).NotTo(HaveOccurred())
			Expect(tarWriter.Close()).To(Succeed())
			tarContents := tarBuffer.Bytes()

			layerID = "6c1f4533b125f8f825188c4f4ff633a338cfce0db2813124d3d518028baf7d7a"
			fakeRegistry.WhenGettingBlob(layerID, 2, func(rw http.ResponseWriter, req *http.Request) {
				gzipWriter := gzip.NewWriter(rw)
				_, err = gzipWriter.Write(tarContents)
				Expect(err).NotTo(HaveOccurred())
				Expect(gzipWriter.Close()).To(Succeed())
			})

			Expect(fakeRegistry.Start()).To(Succeed())

			imageURL = fmt.Sprintf("docker://%s/cfgarden/empty:v0.1.0", fakeRegistry.Addr())
		})

		AfterEach(func() {
			fakeRegistry.Stop()
		})

		It("does not leak corrupted state", func() {
			cmd := exec.Command(GrootFSBin, "--store", StorePath, "create", "--insecure-registry", fakeRegistry.Addr(), imageURL, "random-id")
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess, 12*time.Second).Should(gexec.Exit(1))

			Expect(path.Join(StorePath, "cache", "blobs", fmt.Sprintf("sha256-%s", layerID))).NotTo(BeARegularFile())
			Expect(path.Join(StorePath, store.BUNDLES_DIR_NAME, "random-id")).NotTo(BeADirectory())
		})
	})
})
