package portmanager_test

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/chaos-io/chaos/test/portmanager"
)

func Test_UDPTCP(t *testing.T) {
	lsn, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lsn.Close()

	port := lsn.Addr().(*net.TCPAddr).Port

	udp, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	require.NoError(t, err)
	defer udp.Close()
}

func TestPortManager_SinglePort(t *testing.T) {
	pm := portmanager.New(t)

	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", pm.GetPort()))
	require.NoError(t, err)
	require.NoError(t, lsn.Close())
}

func TestPortManager_PortExhaustion(t *testing.T) {
	pm := portmanager.New(t)

	for i := 0; i < 1000000; i++ {
		port, err := pm.TryGetPort()
		if err != nil {
			return
		}

		t.Logf("allocated port %d", port)
	}

	t.Fatalf("TryGetPort should fail")
}

func TestPortManager_Concurrent(t *testing.T) {
	var takenPorts sync.Map

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pm := portmanager.New(t)
			for i := 0; i < 100; i++ {
				port := pm.GetPort()

				_, loaded := takenPorts.LoadOrStore(port, struct{}{})
				assert.Falsef(t, loaded, "port %d is taken", port)
			}
		}()
	}

	wg.Wait()
}
func ExamplePortManager() {
	var t *testing.T

	pm := portmanager.New(t)

	port := pm.GetPort()
	uiPort := pm.GetPort(8080)

	_ = port
	go log.Fatalf("failed to start UI: %v", http.ListenAndServe(fmt.Sprintf(":%d", uiPort), nil))
}
