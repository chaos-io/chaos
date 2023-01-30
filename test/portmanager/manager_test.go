package portmanager_test

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chaos-io/chaos/test/portmanager"
	"github.com/chaos-io/chaos/test/yatest"
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

func TestPortManager_PythonCompat(t *testing.T) {
	binaryPath, err := yatest.BinaryPath("library/go/test/portmanager/burn_ports/burn_ports")
	require.NoError(t, err)

	pyProcess := exec.Command(binaryPath)
	pyProcess.Stderr = os.Stderr

	var stdout bytes.Buffer
	pyProcess.Stdout = &stdout

	require.NoError(t, pyProcess.Start())

	pm := portmanager.New(t)

	goPorts := map[int]struct{}{}
	for i := 0; i < 1000; i++ {
		goPorts[pm.GetPort()] = struct{}{}
	}

	require.NoError(t, pyProcess.Wait())

	pyPorts := map[int]struct{}{}
	for _, line := range strings.Split(stdout.String(), "\n") {
		if line == "" {
			continue
		}

		port, err := strconv.Atoi(line)
		require.NoError(t, err)

		pyPorts[port] = struct{}{}
	}

	for goPort := range goPorts {
		if _, ok := pyPorts[goPort]; ok {
			t.Errorf("found port conflict: %d", goPort)
		}
	}

	for pyPort := range pyPorts {
		if _, ok := goPorts[pyPort]; ok {
			t.Errorf("found port conflict: %d", pyPort)
		}
	}
}

func ExamplePortManager() {
	var t *testing.T

	pm := portmanager.New(t)

	port := pm.GetPort()
	uiPort := pm.GetPort(8080)

	_ = port
	go log.Fatalf("failed to start UI: %v", http.ListenAndServe(fmt.Sprintf(":%d", uiPort), nil))
}
