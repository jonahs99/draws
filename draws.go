package draws

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// App is the function the user should implement
type App func(c Context, quit <-chan struct{})

// Serve serves the html
func Serve(app App, addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		startApp(app, con)
	})
	http.ListenAndServe(addr, nil)
}

func startApp(app App, con *websocket.Conn) {
	c := &context{command: make(chan string), event: make(chan string)}
	quit := make(chan struct{})

	con.SetCloseHandler(func(code int, text string) error {
		close(quit)
		return nil
	})

	go func() {
		for com := range c.command {
			if com == "BATCH" {
				batch := ""
				for subcom := range c.command {
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

	app(c, quit)
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

var connectinterval

window.addEventListener("load", function(evt) {
	connect()
})

function connect() {
	ws = null
	connectinterval = setInterval(reconnect, 1000);
	reconnect()
}

function reconnect() {
	ws = new WebSocket("{{.}}")
	ws.onopen = (evt) => {
		clearInterval(connectinterval)

		ws.onquit = function(evt) {
			console.log("quit")
			reconnect()
		}
		ws.onmessage = function(evt) {
			console.log("MESSAGE: ", evt.data)
			eval(evt.data)
		}
		return false
	}
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
