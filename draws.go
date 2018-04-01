package draws

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Canvas is the interface type
type Canvas struct {
	App     func(Context)
	open    chan bool
	close   chan bool
	started bool
}

// New returns a new canvas
func New(app func(Context)) Canvas {
	return Canvas{
		App:   app,
		open:  make(chan bool),
		close: make(chan bool),
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
	if c.started {
		return
	}

	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	c.startApp(con)
}

func (c *Canvas) startApp(con *websocket.Conn) {
	ctx := &context{command: make(chan string), event: make(chan string)}

	go func() {
		for com := range ctx.command {
			if com == "BATCH" {
				batch := ""
				for subcom := range ctx.command {
					if subcom == "DRAW" {
						break
					}
					batch += subcom + "\n"
				}
				com = batch[:len(batch)-1]
			}
			err := con.WriteMessage(websocket.TextMessage, []byte(com))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}()

	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}

			evt := event{}
			err = json.Unmarshal(message, &evt)
			if err != nil {
				log.Println("unmarshal:", err)
			}
			ctx.mousex, ctx.mousey = evt.X, evt.Y
		}
	}()

	con.SetCloseHandler(func(code int, text string) error {
		fmt.Println("CLOSE")
		return nil
	})

	c.started = true
	c.App(ctx)
}

type event struct {
	X, Y float64
}

func forward(con *websocket.Conn, ctx *context) {
	go func() {
		for com := range ctx.command {
			fmt.Println("COMMAND: ", com)
			err := con.WriteMessage(websocket.TextMessage, []byte(com))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}()
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}

			evt := event{}
			err = json.Unmarshal(message, &evt)
			if err != nil {
				log.Println("unmarshal:", err)
			}
			ctx.mousex, ctx.mousey = evt.X, evt.Y
		}
	}()
}

var upgrader = websocket.Upgrader{} // use default options

var homeTemplate = template.Must(template.New("home").Parse(
	`<html>
<body>
<script>
var ws

var canvas = document.createElement("canvas")
var context = canvas.getContext("2d")
document.body.appendChild(canvas)

window.addEventListener("load", function(evt) {
	reconnect()
})

function reconnect() {
	ws = null
	let connectinterval = setInterval(() => {
		ws = new WebSocket("{{.}}")
		ws.onopen = (evt) => {
			clearInterval(connectinterval)

			ws.onclose = function(evt) {
				console.log("CLOSE")
				reconnect()
			}
			ws.onmessage = function(evt) {
				console.log("MESSAGE: ", evt.data)
				eval(evt.data)
			}
			return false
		}
	}, 1000)
}

function sendMouseEvents() {
	let rect = canvas.getBoundingClientRect()
	canvas.addEventListener("mousemove", (evt) => {
		msg = {x: evt.clientX - rect.left, y: evt.clientY - rect.top}
		ws.send(JSON.stringify(msg))
	})
}
</script>
<body>
</html>`))
