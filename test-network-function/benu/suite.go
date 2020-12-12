package benu

import (
	"fmt"
	expect "github.com/google/goexpect"
	"github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/testcases"
	"github.com/redhat-nfvpe/test-network-function/test-network-function/benu/configuration"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// The default test timeout.
	defaultTimeoutSeconds = 10
	// The default test timeout.
	defaultTrexServerTimeoutSeconds = 300
	testSpecName                    = "benu"
)

var (
	defaultTimeout           = time.Duration(defaultTimeoutSeconds) * time.Second
	defaultTrexServerTimeout = time.Duration(defaultTrexServerTimeoutSeconds) * time.Second

	trexServerContext *interactive.Context
	subscriberContext *interactive.Context
	dataPlaneContext  *interactive.Context

	trexServerCmd     = "oc exec -it %s --  bash -c 'cd v2.75 ;./t-rex-64 -i'"
	trexSubscriberCmd = "oc exec -it %s -c trex -- python /root/v2.75/BenuQinQ.py"
	benuDataCmd       = "oc exec -it %s -- /opt/benu/bin/cliexec -e 'show bef counter subscribers'"
	benuRegEx         = `^[\\d]+\\s+[\\d.]+\\s+(\\d+)\\s+(\\d+)`
)

var _ = ginkgo.Describe(testSpecName, func() {
	if testcases.IsInFocus(ginkgoconfig.GinkgoConfig.FocusString, testSpecName) {
		//	TimeTicker := time.NewTicker(defaultTimeout) // data is sent every 2 secsmwe are checkign eveyr 4 secs till it reaches 120 secs
		//	tickerChannel := make(chan bool)

		config, _ := configuration.GetConfig()

		ginkgo.When("a trex server is running and traffic is generated", func() {
			config := config
			//loop for 10 times to see if the packets received and transmitted are same
			//for i := 0; i <= 10; i++ {
			ginkgo.It(fmt.Sprintf("Check counter %s ", config.BNGUserPlaneContainer), func() {
				config := config
				var upstream int = 0
				var downstream int = 0
				//i := i
				defer ginkgo.GinkgoRecover()
				gomega.Expect(config).ToNot(gomega.BeNil())
				oc := getOcSession(config.BNGUserPlanePod, config.BNGUserPlaneContainer, config.Namespace, defaultTimeout, expect.Verbose(true))
				gomega.Expect(oc).ToNot(gomega.BeNil())
				gomega.Expect(oc.GetExpecter()).ToNot(gomega.BeNil())
				var dataPlane *BenuBNG
				dataPlane = NewBenuBNG(defaultTimeout, config.BNGUserPlanePod, config.Namespace, BenuBNGShowCounterCmd)
				for i := 0; i <= 10; i++ {
					time.Sleep(2 * time.Second)
					test, err := tnf.NewTest(oc.GetExpecter(), dataPlane, []reel.Handler{dataPlane}, oc.GetErrorChannel())
					gomega.Expect(err).To(gomega.BeNil())
					gomega.Expect(test).ToNot(gomega.BeNil())
					testResult, err := test.Run()
					gomega.Expect(err).To(gomega.BeNil())
					gomega.Expect(testResult).To(gomega.Equal(tnf.SUCCESS))
					gomega.Expect(dataPlane.GetResultOut()).ShouldNot(gomega.BeEmpty())
					u, d := validate(dataPlane.GetResultOut())
					gomega.Expect(u).To(gomega.Equal(d))
					gomega.Expect(u).ShouldNot(gomega.BeEquivalentTo(upstream))
					gomega.Expect(d).ShouldNot(gomega.BeEquivalentTo(downstream))
				}
			})

			//}

			//time.Sleep(120 * time.Second)
			//tickerChannel <- true

		})

	}
})

func validate(pkt string) (int, int) {
	var re = regexp.MustCompile(benuRegEx)
	/* `IP Address of this Session is set Unknown
	  Subs     Subscriber            UpStream   DownStream UpStream  DownStream AUTH
	   Id      Address               Packet     Packet     Drop Pkt   Drop Pkt  STATE
	--------- --------------------- ---------- ---------- ---------- ---------- ------
	3          10.0.0.2              213        213        0          0           AUTH`
	*/

	for _, match := range re.FindAllString(pkt, -1) {
		s := strings.Fields(match)
		upStreamCount, _ := strconv.Atoi(s[2])
		downStreamCount, _ := strconv.Atoi(s[3])
		return upStreamCount, downStreamCount

	}
	return 0, 1

}

// Helper used to instantiate an OpenShift Client Session.
func getOcSession(pod, container, namespace string, timeout time.Duration, options ...expect.Option) *interactive.Oc {
	// Spawn an interactive OC shell using a goroutine (needed to avoid cross expect.Expecter interaction).  Extract the
	// Oc reference from the goroutine through a channel.  Performs basic sanity checking that the Oc session is set up
	// correctly.
	var containerOc *interactive.Oc
	ocChan := make(chan *interactive.Oc)
	var chOut <-chan error

	goExpectSpawner := interactive.NewGoExpectSpawner()
	var spawner interactive.Spawner = goExpectSpawner

	go func(chOut <-chan error) {
		oc, chOut, err := interactive.SpawnOc(&spawner, pod, container, namespace, timeout, options...)
		gomega.Expect(chOut).ToNot(gomega.BeNil())
		gomega.Expect(err).To(gomega.BeNil())
		ocChan <- oc
	}(chOut)

	// Set up a go routine which reads from the error channel
	go func() {
		err := <-chOut
		gomega.Expect(err).To(gomega.BeNil())
	}()

	containerOc = <-ocChan

	gomega.Expect(containerOc).ToNot(gomega.BeNil())

	return containerOc
}
