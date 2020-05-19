package docker

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/moby/term"
	"go.uber.org/zap"
)

type Service struct {
	log *zap.SugaredLogger
	cli *client.Client
}

func New(log *zap.SugaredLogger, host string) *Service {
	svs := &Service{log: log}
	if host == "tcp://127.0.0.1:2376" && runtime.GOOS == "darwin" {
		fmt.Println("⚠️ ️当前使用的是本地docker环境")
		envClient, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err.Error())
		}
		svs.cli = envClient
		return svs
	}

	tlsCert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		panic(err.Error())
	}

	tlsConfig := tlsconfig.ClientDefault()
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.Certificates = []tls.Certificate{tlsCert}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	newClient, err := client.NewClientWithOpts(client.WithHost(host), client.WithHTTPClient(httpClient))
	if err != nil {
		panic(err.Error())
	}
	svs.cli = newClient
	return svs
}

func (s *Service) Pull(image string) error {
	s.log.Debugf("docker pull image: %s", image)
	body, err := s.cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer body.Close()
	inFd, is := term.GetFdInfo(body)

	var buf = &bytes.Buffer{}
	if err := jsonmessage.DisplayJSONMessagesStream(body, buf, inFd, is, nil); err != nil {
		s.log.Warnf("docker pull image display result: %s", err.Error())
	}
	buf.Reset()
	return nil
}

func (s *Service) Run(image, name string, cmd, ports []string) (id string, err error) {
	s.log.Debugf("docker run image: %s", image)

	if runtime.GOOS != "darwin" {
		if err = s.Pull(image); err != nil {
			s.log.Errorf("docker pull image %s error: %s", image, err.Error())
		}
	} else {
		fmt.Println("⚠️ 当前跳过 docker pull️ image:", image)
	}

	config := &container.Config{Image: image, Cmd: cmd}
	hostConfig := &container.HostConfig{NetworkMode: container.NetworkMode("host"), RestartPolicy: container.RestartPolicy{Name: "no"}}
	networkConnect := &network.NetworkingConfig{}

	if len(ports) > 0 {
		hostConfig.NetworkMode = container.NetworkMode("default")
		config.ExposedPorts, hostConfig.PortBindings, err = nat.ParsePortSpecs(ports)
		if err != nil {
			return
		}
	}

	// 创建容器
	containerCreateBody, err := s.cli.ContainerCreate(context.Background(), config, hostConfig, networkConnect, name)
	if err != nil {
		return
	}

	for _, warning := range containerCreateBody.Warnings {
		s.log.Warn("docker create container: ", warning)
	}
	s.log.Debugf("docker create container success, Id: %s", containerCreateBody.ID)

	// 启动容器
	err = s.cli.ContainerStart(context.Background(), containerCreateBody.ID, types.ContainerStartOptions{})
	if err != nil {
		return
	}

	time.Sleep(2 * time.Second)
	inspect, err := s.cli.ContainerInspect(context.Background(), containerCreateBody.ID)
	if err != nil {
		return
	}

	// 检查容器启动是否成功
	if !inspect.State.Running || inspect.State.ExitCode != 0 {
		var responseBody io.ReadCloser
		responseBody, err = s.cli.ContainerLogs(context.Background(), containerCreateBody.ID,
			types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
		if err != nil {
			return
		}

		// 打印容器启动失败的日志信息，并返回错误
		defer responseBody.Close()
		_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, responseBody)
		err = fmt.Errorf("failed to start %s, please check config", name)
	}
	return containerCreateBody.ID, err
}
