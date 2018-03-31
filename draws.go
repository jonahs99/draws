package draws

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Canvas is the interface type
type Canvas struct {
	App func(Context)

	register chan *websocket.Conn
}

// New returns a new canvas
func New(app func(Context)) Canvas {
	return Canvas{
		App: app,
	}
}

// Serve serves the html
func (c *Canvas) Serve(addr string) {
	http.HandleFunc("/", c.home)
	http.HandleFunc("/ws", c.ws)
	http.ListenAndServe(addr, nil)
}

func (c *Canvas) home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func (c *Canvas) ws(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.startApp(con)
}

func (c *Canvas) startApp(con *websocket.Conn) {
	ctx := &context{make(chan string)}
	go forwardCommands(con, ctx.command)
	c.App(ctx)
}

func forwardCommands(con *websocket.Conn, command chan string) {
	for com := range command {
		fmt.Println("COMMAND: ", com)
		err := con.WriteMessage(websocket.TextMessage, []byte(com))
		if err != nil {
			log.Println("write:", err)
		}
	}
}

var upgrader = websocket.Upgrader{} // use default options

var homeTemplate = template.Must(template.New("home").Parse(
	`<html>
<body>
<script>
var c = document.createElement("canvas")
c.id = "canvas"
c.width = 600
c.height = 600
c.style.background = "#333"
document.body.appendChild(c)

ctx = c.getContext("2d")

window.addEventListener("load", function(evt) {
	ws = new WebSocket("{{.}}")
	ws.onopen = function(evt) {
		console.log("OPEN")
		ws.send("open")
	}
	ws.onmessage = function(evt) {
		console.log("MESSAGE: ", evt.data)

		canvas = c
		context = ctx

		eval(evt.data)
	}
	return false
})

</script>
<body>
</html>`))
