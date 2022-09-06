package main

import (
	"fmt"
	"github.com/kluctl/go-embed-python/pip"
	"github.com/kluctl/go-embed-python/python"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	rndName := fmt.Sprintf("pip-install-%d", rand.Uint32())
	ep, err := python.NewEmbeddedPython(rndName)
	if err != nil {
		panic(err)
	}
	defer ep.Cleanup()

	bootstrapPip(ep)

	targetDir := "./data/pip"

	err = pip.CreateEmbeddedPipPackages2(ep, "requirements.txt", targetDir)
	if err != nil {
		panic(err)
	}
}

func bootstrapPip(ep *python.EmbeddedPython) {
	getPip := downloadGetPip()
	defer os.Remove(getPip)

	cmd := ep.PythonCmd(getPip)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func downloadGetPip() string {
	resp, err := http.Get("https://bootstrap.pypa.io/get-pip.py")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic("failed to download get-pip.py: " + resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "get-pip.py")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		panic(err)
	}

	return tmpFile.Name()
}