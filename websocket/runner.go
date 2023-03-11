package websocket

import (
	"html/template"
	"log"
	"net/http"
)

// WebSocket 处理器
func echo(w http.ResponseWriter, r *http.Request) {
	// 协议升级
	c, err := Upgrade(w, r)
	if err != nil {
		log.Print("Upgrade error:", err)
		return
	}
	defer c.Close()
	// 循环处理数据，接收数据，然后返回
	for {
        message, err := c.ReadData()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		c.SendData(message)
	}
}

// index 页面处理器
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func Runner() {
	log.SetFlags(0)
	// 注册 handler
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
    log.Println("Server is running on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

// index 页面内容
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
点击 "Open" 开始一个新的WebSocket链接,
“Send" 将内容发送到服务器，
"Close" 将关闭链接。
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><textarea id="input" type="text" value="hello world!" rows="8"></textarea>
<p><button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
