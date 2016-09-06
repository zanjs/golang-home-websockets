package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

var connid int
var conns *list.List

func Client(w http.ResponseWriter, r *http.Request) {
	html := HTML
	io.WriteString(w, html)
}

func ChatroomServer(ws *websocket.Conn) {
	defer ws.Close()

	connid++
	id := connid
	fmt.Printf("connection id: %d\n", id)

	item := conns.PushBack(ws)
	defer conns.Remove(item)

	name := fmt.Sprintf("用户%d", id)

	SendMessage(nil, fmt.Sprintf(" %s 加入\n", name))

	r := bufio.NewReader(ws)

	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Printf("disconnected id: %d\n", id)
			SendMessage(item, fmt.Sprintf("%s offline\n", name))
			break
		}

		fmt.Printf("%s: %s", name, data)

		SendMessage(item, fmt.Sprintf("%s\t> %s", name, data))
	}
}

func SendMessage(self *list.Element, data string) {
	// for _, item := range conns {
	for item := conns.Front(); item != nil; item = item.Next() {
		ws, ok := item.Value.(*websocket.Conn)
		if !ok {
			panic("item not *websocket.Conn")
		}

		if item == self {
			continue
		}

		io.WriteString(ws, data)
	}
}

func main() {

	fmt.Printf(STARTMSG)

	connid = 0
	conns = list.New()
	http.HandleFunc("/", Client)
	http.Handle("/chatroom", websocket.Handler(ChatroomServer))

	err := http.ListenAndServe(":6611", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

const STARTMSG = `Welcome chatroom server!
author: Julian
open: 127.0.0.1:6611`

const HTML = `<!doctype html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>golang websocket chatroom</title>
</head>
<body>
    <div id="log" style="height: 400px;overflow-y: scroll;border: 1px solid #CCC;">
    </div>
    <div>
        <form id="form">
            <input type="text" id="msg" size="60" />
        </form>
    </div>
</body>
<script>
window.onload = function () {
	    var hostname = window.location.hostname;
	    var hostport = window.location.port;
		if(hostport != ""){
			hostport = ':'+hostport;
		}
		var wsh = 'ws://'+hostname+hostport+'/chatroom';
        var ws = new WebSocket(wsh);

		var msg = document.getElementById("msg");
   		var log = document.getElementById("log");

        ws.onopen = function(e){
            console.log("onopen");
            console.dir(e);
        };
        ws.onmessage = function(e){
            console.log("onmessage");
            console.dir(e);
			var item = document.createElement("div");
			item.innerHTML = '<p>'+e.data+'<p>';
            log.appendChild(item);
            log.scrollTop = log.clientHeight;
        };
        ws.onclose = function(e){
            console.log("onclose");
            console.dir(e);
        };
        ws.onerror = function(e){
            console.log("onerror");
            console.dir(e);
        };
        
		document.getElementById("form").onsubmit = function(){
			if(!msg.value){
				return false;
			}

			ws.send(msg.value + "\n");
			var item = document.createElement("div")
			item.innerHTML = '<p style="color:red;">我:'+ msg.value +'<p>';
			log.appendChild(item);
			log.scrollTop = log.clientHeight;
			msg.value = "";
			return false;
		};
};
</script>
</html>`
