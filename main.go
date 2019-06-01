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
		img = "laher/doke-base"
	} else if err != nil {
		log.Fatalf("cant read filesystem to find Dockerfile.doke", err)
	} else {
		// we are trying to build the equivalent of ...
		// docker run --rm $(docker build -q -f Dockerfile.doke .) $@
		cmd := exec.Command("docker", "build", "-q", "-f", "Dockerfile.doke", ".")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("docker-build failed with %s\n", err)
		}
		img = strings.TrimSpace(string(out))
		// fmt.Println(img) //print output for consistency
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
