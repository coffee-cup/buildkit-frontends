package main

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	gw "github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
)

func main() {
	if err := gw.RunFromEnvironment(appcontext.Context(), Build); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	// Get build args and options

	opts := c.BuildOpts()

	buildArgs := opts.Opts
	fmt.Printf("buildArgs: %+v\n", buildArgs)

	state := llb.Image("alpine").Run(llb.Shlex("ls -l /bin")).Root()
	def, err := state.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("def: %+v\n", def)

	result := client.NewResult()

	return result, nil
}
