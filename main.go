package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() failed with %s\n", err)
	}
	img := ""
	if _, err := os.Stat("Dockerfile.doke"); os.IsNotExist(err) {
		img = "laher/doke"
	} else if err != nil {
		log.Fatalf("cant read filesystem to find Dockerfile.doke", err)
	} else {
		//docker run --rm -it $(docker build -q -f Dockerfile.doke .) $@
		cmd := exec.Command("docker", "build", "-q", "-f", "Dockerfile.doke", ".")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("docker-build failed with %s\n", err)
		}
		img = strings.TrimSpace(string(out))
	}

	base := filepath.Base(pwd)
	cmd := exec.Command("docker")
	cmd.Args = append(cmd.Args, "run", "--rm", "-v", fmt.Sprintf("%s:/%s", pwd, base), img)
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("docker-run failed with %s\n", err)
	}
}
