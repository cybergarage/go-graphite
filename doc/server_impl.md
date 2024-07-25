# How to implement a Graphite-compatible server?

This section describes how to implement a Graphite-compatible server using the go-graphite.

## Creating your Graphite-com-compatible server

### STEP1: Inheritancing the base server

At first, inherit the base server of the go-graphite to implement your Redis-compatible server as the following:

```
import (
	"github.com/cybergarage/go-graphite/net/graphite"
)

type Server struct {
	*graphite.Server
}

func NewServer() *Server {
    server := &Server{
        Server: graphite.NewServer(),
    }
    return server
}
```
### STEP2: Implementing your Carbon command handler

Next, implement your Carbon command handler according to the [CarbonListener (PlainTextRequestListener)](../net/graphite/carbon_listener.go) interface of the go-graphite as the following:

```
func (server *Server) InsertMetricsRequestReceived([]*Metrics, error) {
    // Implement your command handler here
    ....
}
```

### STEP3: Implementing your Render command handler

Next, implement your Render command handler according to the [RenderListener (RenderRequestListener)](../net/graphite/render_listener.go) interface of the go-graphite as the following:

```
func (server *Server) FindMetricsRequestReceived(*Query, error) ([]*Metrics, error) {
    // Implement your command handler here
    ....
}

func (server *Server) QueryMetricsRequestReceived(*Query, error) ([]*Metrics, error) {
    // Implement your command handler here
    ....
}
```

### STEP4: Setting your user command handlers

Next, set your user command handler to your server using `Server::SetCarbonListener()` and `Server::SetRenderListener()` as the following:

```
func NewServer() *Server {
	server: = &Server{
		Server:    redis.NewServer(),
	}
    server.SetCarbonListener(server)
    server.SetRenderListener(server)
    return server
}
```

### STEP5: Starting server

Finally, start your compatible server using `Server::Start()` as the following:

```
server := NewServer()
....
server::Start()
```
