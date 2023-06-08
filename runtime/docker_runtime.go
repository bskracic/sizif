package runtime

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerRuntime struct {
	cli   *docker.Client
	specs Specs
}

func NewDockerRuntime() *DockerRuntime {
	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &DockerRuntime{
		cli: cli,
	}
}

func (dr *DockerRuntime) Prepare(s Specs) string {
	return dr.createContainer(fmt.Sprintf("sizif-%s-%v", s.Lang, time.Now().Unix()), s.Image)
}

func (dr *DockerRuntime) createContainer(name string, image string) string {
	ctx := context.Background()
	resp, err := dr.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		Cmd:   []string{"sleep", "infinity"},
	}, nil, nil, nil, name)
	if err != nil {
		panic(err)
	}
	if err := dr.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func (dr *DockerRuntime) Exec(id string, cmd []string, ch chan *ExecResult) error {

	config := types.ExecConfig{AttachStdout: true, AttachStderr: true, Cmd: cmd}
	ctx := context.Background()

	respCreate, err := dr.cli.ContainerExecCreate(ctx, id, config)
	if err != nil {
		log.Fatalln("cannot create exec: ", err)
		return err
	}

	respExec, err := dr.cli.ContainerExecAttach(ctx, respCreate.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatalln("cannot exec attach: ", err)
	}
	defer respExec.Close()

	// Read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, respExec.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			log.Fatalln(err)
			return err
		}
		break

	case <-ctx.Done():
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		log.Fatalln(err)
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := dr.cli.ContainerExecInspect(ctx, respCreate.ID)
	if err != nil {
		log.Fatalln("cannot exec inspect ubicu se majkemi", err)
	}

	ch <- &ExecResult{Stdout: string(stdout), Stderr: string(stderr), ExitCode: res.ExitCode}
	return err
}

func (dr *DockerRuntime) CopyFile(id string, content string, filename string, dst string) error {
	f, err := generateFileContent(filename, content)
	if err != nil {
		return err
	}
	ctx := context.Background()

	return dr.cli.CopyToContainer(ctx, id, dst, f, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
}

func (dr *DockerRuntime) CreateDir(id string, dirpath string) error {

	cmd := fmt.Sprintf("mkdir %s", dirpath)
	config := types.ExecConfig{Cmd: []string{"bash", "-c", cmd}}
	ctx := context.Background()
	respCreate, err := dr.cli.ContainerExecCreate(ctx, id, config)
	if err != nil {
		return err
	}
	// Listen for an event and return only after exec finishes
	msgs, errs := dr.cli.Events(ctx, types.EventsOptions{})
	err = dr.cli.ContainerExecStart(ctx, respCreate.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-errs:
			return err
		case msg := <-msgs:
			if msg.Action == "exec_die" && msg.Actor.ID == id {
				return nil
			}
		}
	}
}

func (dr *DockerRuntime) Kill(id string) {
	ctx := context.Background()
	if err := dr.cli.ContainerKill(ctx, id, ""); err != nil {
		log.Default().Printf("error killing %s: %v\n", id, err)
	}

	err := dr.cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		log.Default().Printf("error removing %s: %v", id, err)
	}
}
