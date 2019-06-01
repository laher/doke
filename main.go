package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	dokeFileName     = "Dockerfile.doke"
	dokeFileTemplate = `FROM alpine 

RUN mkdir {{.WorkDir}}
RUN apk add --update make

WORKDIR {{.WorkDir}}

ENTRYPOINT ["make"]
`
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() failed with %s\n", err)
	}
	base := filepath.Base(pwd)
	if _, err := os.Stat(dokeFileName); os.IsNotExist(err) {
		fmt.Printf("%s not found - generating ...\n", dokeFileName)
		var (
			t *template.Template
			w io.WriteCloser
		)
		if t, err = template.New("dockerfile").Parse(dokeFileTemplate); err != nil {
			log.Fatal("parsing template: ", err)
		}
		if w, err = os.Create(dokeFileName); err != nil {
			log.Fatalf("cant write %s: %v", dokeFileName, err)
		}
		if err := t.Execute(w, struct{ WorkDir string }{WorkDir: "/" + base}); err != nil {
			log.Fatalf("cant write %s: %v", dokeFileName, err)
		}
		if err := w.Close(); err != nil {
			log.Fatalf("cant write %s: %v", dokeFileName, err)
		}
		fmt.Printf("Wrote a basic %s - please edit it to add your build dependencies\n", dokeFileName)
	} else if err != nil {
		log.Fatalf("cant read filesystem to find %s: %v", dokeFileName, err)
	}
	// we are trying to build the equivalent of ...
	// docker run --rm $(docker build -q -f Dockerfile.doke .) $@
	cmd := exec.Command("docker", "build", "-q", "-f", "Dockerfile.doke", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("docker-build failed with %s\n", err)
	}
	img := strings.TrimSpace(string(out))
	// fmt.Println(img) //print output for consistency

	cmd = exec.Command("docker")
	cmd.Args = append(cmd.Args, "run", "--rm", "-v", fmt.Sprintf("%s:/%s", pwd, base), img)
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("docker-run failed with %s\n", err)
	}
}
