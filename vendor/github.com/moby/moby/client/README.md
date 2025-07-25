# Go client for the Docker Engine API

The `docker` command uses this package to communicate with the daemon. It can
also be used by your own Go applications to do anything the command-line
interface does; running containers, pulling or pushing images, etc.

For example, to list all containers (the equivalent of `docker ps --all`):

```go
package main

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func main() {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
}
```

[Full documentation is available on pkg.go.dev.](https://pkg.go.dev/github.com/moby/moby/client)
