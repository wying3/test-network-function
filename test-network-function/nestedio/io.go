package nestedio

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	expect "github.com/ryandgoulding/goexpect"
	"github.com/sirupsen/logrus"
	"io"
	"regexp"
	"time"
)

var testTimeout = 20 * time.Second

// Helper used to instantiate an OpenShift Client Session.
func getOcSession(pod, container, namespace string, stdinPipe *io.WriteCloser, stdoutPipe *io.Reader, timeout time.Duration, options ...expect.Option) *interactive.Oc {
	// Spawn an interactive OC shell using a goroutine (needed to avoid cross expect.Expecter interaction).  Extract the
	// Oc reference from the goroutine through a channel.  Performs basic sanity checking that the Oc session is set up
	// correctly.
	var containerOc *interactive.Oc
	ocChan := make(chan *interactive.Oc)
	var chOut <-chan error

	goExpectSpawner := interactive.NewGoExpectSpawner()
	var spawner interactive.Spawner = goExpectSpawner

	go func() {
		oc, outCh, err := interactive.SpawnOc(&spawner, pod, container, namespace, timeout, options...)
		gomega.Expect(outCh).ToNot(gomega.BeNil())
		gomega.Expect(err).To(gomega.BeNil())
		ocChan <- oc
	}()

	// Set up a go routine which reads from the error channel
	go func() {
		err := <-chOut
		gomega.Expect(err).To(gomega.BeNil())
	}()

	containerOc = <-ocChan

	gomega.Expect(containerOc).ToNot(gomega.BeNil())

	return containerOc
}

var _ = ginkgo.Describe("io", func() {
	ginkgo.When("nested io is requested", func() {
		ginkgo.It("should work", func() {
			oc := getOcSession("partner", "partner", "default", nil, nil, testTimeout, expect.Verbose(true))

			err := (*oc.GetExpecter()).Send("./a.sh\n")
			gomega.Expect(err).To(gomega.BeNil())
			(*oc.GetExpecter()).Send("ryan\n")
			output, _, err := (*oc.GetExpecter()).Expect(regexp.MustCompile(`(?m).+`), testTimeout)
			gomega.Expect(err).To(gomega.BeNil())
			logrus.Infof("output: %s", output)
		})
	})
})
