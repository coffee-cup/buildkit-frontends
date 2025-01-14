package main

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/buildkit/client/llb"
)

func writellb() {
	dt, err := createLLBState().Marshal(context.TODO(), llb.LinuxAmd64)
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)
}

func createLLBState() llb.State {

	fmt.Fprintf(os.Stderr, "\n\nBuild function called\n\n")

	// Get build args
	// Base image
	base := llb.Image("ubuntu:noble")

	// Install curl
	state := base.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y curl")).
		Run(llb.Shlex("rm -rf /var/lib/apt/lists/*")).
		Root()

	// Set environment variables
	state = state.AddEnv("GIT_SSL_CAINFO", "/etc/ssl/certs/ca-certificates.crt").
		AddEnv("MISE_DATA_DIR", "/mise").
		AddEnv("MISE_CONFIG_DIR", "/mise").
		AddEnv("MISE_INSTALL_PATH", "/usr/local/bin/mise").
		AddEnv("PATH", "/mise/shims:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	// Install mise
	// state = state.Run(llb.Shlex("curl https://mise.run | sh")).Root()

	state = state.Run(llb.Shlex("sh -c 'curl -fsSL https://mise.run | sh'")).Root()

	// Set working directory
	state = state.Dir("/app")

	// Create mise config
	miseConfig := `[tools.node]
version = "23.5.0"
[tools.npm]
version = "11.0.0"
`

	state = state.File(llb.Mkdir("/etc/mise", 0755, llb.WithParents(true)))
	state = state.File(llb.Mkfile("/etc/mise/config.toml", 0644, []byte(miseConfig)))

	// Trust and install mise tools
	state = state.Run(llb.Shlex("mise trust -a")).
		Run(llb.Shlex("mise install")).
		Root()

	// Copy the application files
	src := llb.Local("context")
	state = state.File(llb.Copy(src, ".", ".", &llb.CopyInfo{
		CopyDirContentsOnly: true,
	}))

	// Run npm ci
	state = state.Run(llb.Shlex("npm ci")).Root()

	return state
}
