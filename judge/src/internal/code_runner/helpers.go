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

func runCodeInsideContainer(timeLimitMs, problemID, submissionID int) string {
	timeLimit := fmt.Sprintf("%.2f", float64(timeLimitMs)/float64(1000))

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return "Compilation failed"
	}
	defer cli.Close()

	problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problemID))

	config := &container.Config{
		Image: "go-code-runner",
		Cmd:   []string{timeLimit},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(problemDir, fmt.Sprintf("%d.go", submissionID)),
				Target:   "/mnt/problem/code.go",
				ReadOnly: true,
			},
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(problemDir, "input.txt"),
				Target:   "/mnt/problem/input.txt",
				ReadOnly: true,
			},
		},
		Resources: container.Resources{
			CPUCount: 1,
			Memory:   256 * 1024 * 1024,
			NanoCPUs: 1_000_000_000,
			Ulimits:  []*units.Ulimit{{Name: "nofile", Soft: 1024, Hard: 1024}},
		},

		NetworkMode: "none", // No network
		CapDrop:     []string{"ALL"},
		SecurityOpt: []string{"no-new-privileges"},
		AutoRemove:  true, // Equivalent to --rm
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		return "Compilation failed"
	}

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

	isFirstLine := true
	var output string
	dScanner := bufio.NewScanner(attachResp.Reader)
	for dScanner.Scan() {
		line := dScanner.Text()
		if isFirstLine {
			isFirstLine = false
			if len(line) > 8 {
				line = line[8:]
			}
		}
		output += line
	}
	if err := dScanner.Err(); err != nil {
		fmt.Printf("Error reading container output: %v\n", err)
		return "Compilation failed"
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		fmt.Printf("Error waiting for container: %v\n", err)
		return "Compilation failed"
	case status := <-statusCh:
		if status.StatusCode != 0 {
			fmt.Printf("Compilation failed")
			return "Compilation failed"
		}
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
	} else if strings.Contains(output[:charLimit], "Compilation failed") {
		return "Compilation failed"
	} else if strings.Contains(output[:charLimit], "Runtime error") {
		return "Runtime error"
	}

	file, err := os.Open(filepath.Join(problemDir, "output.txt"))
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

	resp, err := http.Post("http://localhost:80/code/callback", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}
