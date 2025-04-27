package code_runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
)

type ResultData struct {
	CallbackToken string `json:"callback_token"`
	Status        string `json:"Status"`
}

func runCodeInsideContainer(run Run) string {
	timeLimit := fmt.Sprintf("%.3f", float64(run.TimeLimitMs+5000)/float64(1000))

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return "Compilation failed"
	}
	defer cli.Close()

	problemDirOnHost := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER_SRC"), fmt.Sprintf("%d", run.PproblemID))
	problemDirOnContainer := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", run.PproblemID))

	config := &container.Config{
		Image: "go-code-runner",
		Cmd:   []string{""},
		Env: []string{
			"TIME_LIMIT=" + timeLimit,
		},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(problemDirOnHost, fmt.Sprintf("%d.go", run.SubmissionID)),
				Target:   "/mnt/problem/code.go",
				ReadOnly: true,
			},
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(problemDirOnHost, "input.txt"),
				Target:   "/mnt/problem/input.txt",
				ReadOnly: true,
			},
		},
		Resources: container.Resources{
			CPUCount: 1,
			Memory:   int64(512) * 1024 * 1024, // for build code
			NanoCPUs: 1_000_000_000,
			Ulimits:  []*units.Ulimit{{Name: "nofile", Soft: 1024, Hard: 1024}},
		},

		NetworkMode: "none", // No network
		CapDrop:     []string{"ALL"},
		SecurityOpt: []string{"no-new-privileges"},
		AutoRemove:  false,
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		return "Compilation failed"
	}
	defer cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return "Compilation failed"
	}

	attachResp, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		fmt.Printf("Error attaching to container: %v\n", err)
		return "Compilation failed"
	}
	defer attachResp.Close()

	var compileOutput string
	dScanner := bufio.NewScanner(attachResp.Reader)
	for dScanner.Scan() {
		line := dScanner.Text()
		if len(line) > 8 {
			line = line[8:]
		}
		compileOutput += line
	}
	if err := dScanner.Err(); err != nil {
		fmt.Printf("Error reading compile container output: %v\n", err)
		return "Compilation failed"
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		fmt.Printf("Error waiting for container: %v\n", err)
		return "Compilation failed"
	case status := <-statusCh:
		if status.StatusCode != 0 {
			fmt.Println("Compilation failed", status.Error)
			return "Compilation failed"
		}
	}

	if strings.Contains(compileOutput, "failed") {
		fmt.Println("Compilation failed")
		return "Compilation failed"
	}

	_, err = cli.ContainerUpdate(ctx, resp.ID, container.UpdateConfig{
		Resources: container.Resources{
			Memory: int64(run.MemoryLimitMb+6) * 1024 * 1024, // 6 mb for container
		},
	})
	if err != nil {
		fmt.Printf("error updating container memory: %v", err)
		return "Compilation failed"
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return "Compilation failed"
	}

	statusCh, errCh = cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		fmt.Printf("Error waiting for container: %v\n", err)
		return "Compilation failed"
	case status := <-statusCh:
		if status.StatusCode != 0 {
			fmt.Println("Compilation failed", status.Error)
			return "Compilation failed"
		}
	}
	logs, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		fmt.Println("error reading logs for container", err)
		return "Compilation failed"
	}
	defer logs.Close()
	isFirstLine := true
	var output string
	dScanner = bufio.NewScanner(logs)
	for dScanner.Scan() {
		line := dScanner.Text()
		if isFirstLine {
			isFirstLine = false
			continue // skip first line
		}
		if len(line) > 8 {
			line = line[8:]
		}
		output += line
	}

	if err := dScanner.Err(); err != nil {
		fmt.Println("error reading exec output", err)
		return "Compilation failed"
	}

	charLimit := 200
	if len(output) < charLimit {
		charLimit = len(output)
	}

	if strings.Contains(output[:charLimit], "Wrong answer") {
		return "Wrong answer"
	} else if strings.Contains(output[:charLimit], "Time limit exceeded") {
		return "Time limit exceeded"
	} else if strings.Contains(output[:charLimit], "Memory limit exceeded") {
		return "Memory limit exceeded"
	} else if strings.Contains(output[:charLimit], "Runtime error") {
		return "Runtime error"
	}

	file, err := os.Open(filepath.Join(problemDirOnContainer, "output.txt"))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "Compilation failed"
	}
	defer file.Close()

	fScanner := bufio.NewScanner(file)
	var main_output string
	for fScanner.Scan() {
		line := fScanner.Text()
		main_output += line
	}

	if output == main_output {
		return "Accepted"
	}

	return "Wrong answer"

}

func sendRunCallBack(result string, run Run) {
	resultData := ResultData{
		CallbackToken: run.CallbackToken,
		Status:        result,
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post("http://judge:80/code/callback", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}
