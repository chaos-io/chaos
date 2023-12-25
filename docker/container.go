package docker

import (
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/chaos-io/chaos/core"
	"github.com/chaos-io/chaos/logs"
)

const (
	OptionWorkingDir  = "workingDir"
	OptionEnv         = "env"
	OptionCpuSet      = "cpuset"
	OptionPorts       = "ports"
	OptionMemoryLimit = "memory"
	OptionAddHost     = "add-host"
	OptionAddDns      = "dns"
	OptionNetwork     = "network"
)

var (
	cli      *client.Client
	dictPath string
)

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithHost(Host), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("failed to connect docker due to %v", err)
	}
}

var Host = "http://172.17.0.1:2375"

// Run block to run the container, and waiting for stop
func Run(ctx context.Context, imageName, containerName string, cmd []string, options core.Options, bindPaths ...string) (int64, []byte, error) {
	logger := logs.With()
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(Host), client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = c.Close()
	}()

	cmdBuilder := strings.Builder{}
	cmdBuilder.WriteString("docker run -it --rm")

	var env []string
	if envOption, ok := options[OptionEnv].([]string); ok {
		env = envOption
	}
	for _, ev := range env {
		cmdBuilder.WriteString(" --env '" + ev + "'")
	}

	cpuset := ""
	if cpusetOption, ok := options[OptionCpuSet].(string); ok {
		cpuset = cpusetOption
	}

	var ports []string
	if portsOption, ok := options[OptionPorts].([]string); ok {
		ports = portsOption
	}
	for _, p := range ports {
		cmdBuilder.WriteString(" -p " + p)
	}

	cfg := &container.Config{
		Image: imageName,
		Cmd:   cmd,
		Env:   env,
		Tty:   true,
	}

	var mounts []mount.Mount
	for i := 0; i < len(bindPaths)-1; i += 2 {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: bindPaths[i],
			Target: bindPaths[i+1],
		})

		cmdBuilder.WriteString(" -v " + bindPaths[i] + ":" + bindPaths[i+1])
	}

	hostConfig := &container.HostConfig{
		Mounts: mounts,
		Resources: container.Resources{
			CpusetCpus: cpuset,
		},
	}

	pset, pbindings, _ := nat.ParsePortSpecs(ports)
	if len(pbindings) > 0 {
		cfg.ExposedPorts = pset
		hostConfig.PortBindings = pbindings
	}

	resp, err := c.ContainerCreate(ctx, cfg, hostConfig, nil, nil, containerName)
	if err != nil {
		return 0, nil, err
	}

	cmdBuilder.WriteString(" " + imageName)
	for _, command := range cmd {
		cmdBuilder.WriteString(" " + command)
	}

	start := time.Now()
	logger.Infow("run docker container", "containerId", resp.ID, "cmd", cmdBuilder.String())

	if err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return 0, nil, err
	}

	statusCh, errCh := c.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	logger.Debugw("container wait", "containerId", resp.ID, "duration", time.Since(start).String())

	var exitCode int64 = -1
	timeout := 0
	if timeoutOption, ok := options["timeout"].(int); ok {
		timeout = timeoutOption
	}
	if timeout > 0 {
		timeoutTimer := time.NewTimer(time.Duration(timeout) * time.Second)
		select {
		case err = <-errCh:
			timeoutTimer.Stop()
			if err != nil {
				return 0, nil, err
			}
		case status := <-statusCh:
			exitCode = status.StatusCode
			timeoutTimer.Stop()
		case <-timeoutTimer.C:
		}
	} else {
		select {
		case err = <-errCh:
			logger.Debugw("select error", "containerId", resp.ID, "error", err, "duration", time.Since(start).String())
			if err != nil {
				return 0, nil, err
			}
		case status := <-statusCh:
			logger.Debugw("select status", "containerId", resp.ID, "status", status, "duration", time.Since(start).String())
			exitCode = status.StatusCode
		}
	}

	out, err := containerLogs(c, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		logger.Warnw("failed to get the logs form container", "containerID", resp.ID, "error", err)
		return 0, nil, err
	}
	logger.Debugw("container logs", "containerId", resp.ID, "duration", time.Since(start).String())

	to := 2
	if err = c.ContainerStop(ctx, resp.ID, container.StopOptions{Timeout: &to}); err != nil {
		logger.Warn("failed to stop the container", "containerID", resp.ID, "error", err)
		return 0, out, nil
	}
	logger.Debugw("container stop", "containerId", resp.ID, "duration", time.Since(start).String())

	if err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		logger.Warn("failed to remove the container", "containerID", resp.ID, "error", err)
	}

	logger.Infow("run docker container successfully", "imageName", imageName, "containerId", resp.ID, "containerName", containerName, "exitCode", exitCode, "duration", time.Since(start).String())

	return exitCode, out, nil
}

func containerLogs(c *client.Client, id string, options types.ContainerLogsOptions) ([]byte, error) {
	out, err := c.ContainerLogs(context.Background(), id, options)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = out.Close()
	}()

	byt, err := io.ReadAll(out)
	if err != nil {
		return nil, err
	}
	return []byte(stripansi.Strip(string(byt))), nil
}
