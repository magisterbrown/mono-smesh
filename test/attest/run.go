package main
import (
    "context"
    "fmt"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/strslice"
    "github.com/docker/docker/client"
)

func main() {
    Dock_cli, _ := client.NewClientWithOpts(client.FromEnv)
    conf := container.Config{Image: agent.Image, Tty: true,  Cmd:  strslice.StrSlice{command}}
    cont, err := Dock_cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")
    conf := container.Config{Image: "39c14dd445e9", Tty: true,  Cmd:  strslice.StrSlice{"sadsad dsds dds"}}
    cont, _ := Dock_cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")

    hijack, _ := Dock_cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream:true, Stdout:true})

    Dock_cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
    buf, _, _:= hijack.Reader.ReadLine()
    fmt.Println(string(buf))

    Dock_cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{Force: true})
    fmt.Println(cont.ID)

}
