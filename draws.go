package draws

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const pongWait = 1 * time.Second
const pingPeriod = (pongWait * 9) / 10
const writeWait = 1 * time.Second

var upgrader = websocket.Upgrader{} // use default options

var homeTemplate = template.Must(template.New("home").Parse(client))

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
	fmt.Println("Starting app.")

	c := &context{command: make(chan string), event: make(chan string)}
	quit := make(chan struct{})

	con.SetCloseHandler(func(code int, text string) error {
		fmt.Println("Closing quit channel.")
		close(quit)
		return nil
	})

	con.SetReadDeadline(time.Now().Add(pongWait))
	con.SetPongHandler(func(string) error {
		con.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		con.Close()
	}()

	go func() {
		for {
			select {
			case com := <-c.command:
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
		}
	}()

	app(c, quit)
}
