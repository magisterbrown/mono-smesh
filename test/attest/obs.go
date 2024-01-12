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
    conf := container.Config{Image: "2d8fb5f795b077287cdbfdf3b06483d23516d43308b0136586d98fa315f4e4da", Tty: true,  Cmd:  strslice.StrSlice{"--observation", "WWWWWWWWW/WWWWWWWWW/BWBW1BWBW/BBBBBBBBB/BBBBBBBBB W - - - 0"}}
    cont, err := Dock_cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")
    

   hijack, err := Dock_cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream:true, Stdout:true})

   Dock_cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
   buf, _, err:= hijack.Reader.ReadLine()

   Dock_cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{Force: true})
   fmt.Println(string(buf))
    _ = err
}
