package main

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

var sysCtx = &types.SystemContext{
	RegistriesDirPath:        "",
	SystemRegistriesConfPath: "",
	BigFilesTemporaryDir:     "/var/tmp",
	OSChoice:                 "linux",
	ArchitectureChoice:       "arm64",
}

func getSource(imgref string) (types.ImageSource, error) {
	ref, err := alltransports.ParseImageName(imgref)
	if err != nil {
		return nil, err
	}

	src, err := ref.NewImageSource(context.Background(), sysCtx)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func run(imgref string) error {
	src, err := getSource(imgref)
	if err != nil {
		return err
	}

	data, _, err := src.GetManifest(context.Background(), nil)
	if err != nil {
		return err
	}

	index, err := manifest.OCI1IndexFromManifest(data)
	if err != nil {
		return err
	}

	for _, m := range index.Manifests {
		data, _, err = src.GetManifest(context.Background(), &m.Digest)
		if err != nil {
			return err
		}

		result, err := manifest.OCI1FromManifest(data)
		if err != nil {
			return err
		}

		fmt.Println("Result: ", result.ConfigInfo().Digest)
	}

	return nil
}

func main() {
	imgrefs := map[string]string{
		"directory":     fmt.Sprintf("dir:%s", "containers"),
		"remote":        fmt.Sprintf("docker://%s", "quay.io/fedora/fedora-minimal"),
		"local storage": fmt.Sprintf("containers-storage:%s", "localhost/multi"),
	}

	for key, ref := range imgrefs {
		fmt.Println("Resolving", key, "manifest list")
		err := run(ref)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----------")
		fmt.Println()
	}

}
