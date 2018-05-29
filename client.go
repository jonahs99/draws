package draws

/*
	This is the client code.
	We setup a canvas and websocket connection,
	then when we receive messages we wrap them
	as function calls.
*/

const client = `
<html>
<body>

<script>

window.addEventListener("load", function(evt) {
	let canvas = document.createElement("canvas")
	let context = canvas.getContext("2d")
	document.body.appendChild(canvas)

	ws = new WebSocket("{{.}}")
	ws.addEventListener('open', (evt) => {
		console.log("OPEN")
	})
	ws.addEventListener('message', (evt) => {
		console.log("MESSAGE: ", evt.data)
		execute(canvas, context, evt.data)
	})
})

function execute(canvas, context, code) {
	f = Function('canvas,context', code)
	f(canvas, context)
}

</script>

<body>
</html>
`
