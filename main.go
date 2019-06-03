package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	dockerFileName     = "Dockerfile.doke"
	dockerFileTemplate = `FROM alpine 
RUN mkdir -p {{ .WorkDir }}
RUN apk add --update tzdata ca-certificates make

WORKDIR {{ .WorkDir }}

ENTRYPOINT ["make"]
`
	dockerComposeFileName = "docker-compose.doke.yaml"
	dockerComposeTemplate = `version: '3'
services:
  make:
    build:
      context: .
      dockerfile: Dockerfile.doke
    working_dir: {{ .WorkDir }}
    volumes:
      - "$PWD:{{ .WorkDir }}"
    entrypoint: make
`
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func genFileIfNotExist(name, tpl, base string) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		fmt.Printf("%s not found - generating ...\n", name)
		var (
			t *template.Template
			w io.WriteCloser
		)
		if t, err = template.New("template").Parse(tpl); err != nil {
			log.Fatal("parsing template: ", err)
		}
		if w, err = os.Create(name); err != nil {
			log.Fatalf("cant write %s: %v", name, err)
		}
		if err := t.Execute(w, struct{ WorkDir string }{WorkDir: "/" + base}); err != nil {
			log.Fatalf("cant write %s: %v", name, err)
		}
		if err := w.Close(); err != nil {
			log.Fatalf("cant write %s: %v", name, err)
		}
		fmt.Printf("Wrote a basic %s - please edit it to add your build dependencies\n", name)
	} else if err != nil {
		log.Fatalf("cant read filesystem to find %s: %v", name, err)
	}
}

func main() {
	var envVars arrayFlags
	flag.Var(&envVars, "e", "specify an env var")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() failed with %s\n", err)
	}
	base := filepath.Base(pwd)

	genFileIfNotExist(dockerFileName, dockerFileTemplate, base)
	var dockerCompose bool
	var cmd *exec.Cmd
	if _, err := exec.LookPath("docker-compose"); err != nil {
		log.Printf("docker-compose not found: using docker")
		dockerCompose = false
		if _, err := exec.LookPath("docker"); err != nil {
			log.Fatalf("docker not found: %s", err)
		}
	} else {
		dockerCompose = true
	}
	if dockerCompose {
		genFileIfNotExist(dockerComposeFileName, dockerComposeTemplate, base)
		cmd = exec.Command("docker-compose", "-f", dockerComposeFileName)
		for _, envVar := range envVars {
			cmd.Args = append(cmd.Args, "-e", envVar)
		}
		cmd.Args = append(cmd.Args, "run", "make")
	} else {
		// we are trying to build the equivalent of ...
		// docker run --rm $(docker build -q -f Dockerfile.doke .) $@
		cmd = exec.Command("docker", "build", "-q", "-f", dockerFileName, ".")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("docker-build failed with %s\n", err)
		}
		img := strings.TrimSpace(string(out))
		// fmt.Println(img) //print output for consistency

		cmd = exec.Command("docker")
		cmd.Args = append(cmd.Args, "run", "--rm", "-v", fmt.Sprintf("%s:/%s", pwd, base))
		for _, envVar := range envVars {
			cmd.Args = append(cmd.Args, "-e", envVar)
		}
		cmd.Args = append(cmd.Args, img)
	}
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("docker-run failed with %s\n", err)
	}
}
