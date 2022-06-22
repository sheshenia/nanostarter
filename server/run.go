package server

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", ":8085", "http service address")
	useEmbed = flag.Bool("embed", true, "serve static client with embed Go")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
)

//go:embed client/assets
var assets embed.FS

//go:embed client/index.html
var mainPage []byte

func Run(ctx context.Context) error {
	flag.Parse()

	//output address in blue color
	colored := fmt.Sprintf("\x1b[%dm%s%s\x1b[0m", 34, "http://localhost", *addr)
	log.Println("Serving Nanostarter:", colored)
	/*if _, err := exec.Command("bash", "-c", `sed -i -e 's/__PORT__\s\?=\s\?":[0-9]\{4\}"/__PORT__ = "`+*addr+`"/' ./client/dist/index.html`).Output(); err != nil {
		log.Println(err)
	}*/

	mux := http.NewServeMux()

	// -embed=true
	// by default in production, we use embed mode.
	// You need nothing to change just use make command:
	// make my
	// it will build static Vuejs to client/dist
	// and copy client/dist to server/client
	// embed server/client contents to var assets
	// and index.html to var mainPage, using mainPage as template
	//
	// -embed=false
	// use in dev mode, be sure to run app in default addr :8085
	// in non linux system change in ./client/index.html :
	// window.__GOOS__ = "linux";
	// window.__GOARCH__ = "amd64";
	// to:
	// window.__GOOS__ = "darwin";
	// window.__GOARCH__ = "arm64"; //if not amd64
	//
	// serving static content in ./client with:
	// npm run dev
	// in this case usually open client on http://localhost:3000
	if *useEmbed {
		// replacing default port in index.html to *addr
		templatePort := regexp.MustCompile(`(?i)window\.__PORT__\s?=\s?(?:"|'):\d{4,}(?:"|')`)
		mainPage = templatePort.ReplaceAll(mainPage, []byte(fmt.Sprintf(`window.__PORT__ = "%s"`, *addr)))

		//runtime.GOOS, runtime.GOARCH for default commands on client
		bytes.Replace(mainPage, []byte(`"linux";`), []byte(fmt.Sprintf(`"%s";`, runtime.GOOS)), 1)
		bytes.Replace(mainPage, []byte(`"amd64";`), []byte(fmt.Sprintf(`"%s";`, runtime.GOARCH)), 1)

		//serving main page
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write(mainPage); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		contentStatic, err := fs.Sub(assets, "client")
		if err != nil {
			log.Println(err)
		}
		//serving assets: *.js, *.css and other
		mux.Handle("/assets/", http.FileServer(http.FS(contentStatic))) //http.FileServer(http.Dir("./client/dist")))
	} else {
		// in -embed=false mode serve all static content from ./client/dist
		// if for some reason you used npm run build, and don't want to serve static with npm run dev
		mux.Handle("/", http.FileServer(http.Dir("./client/dist"))) //http.FileServer(http.Dir("./client/dist")))
	}

	//each command has its own websocket handler
	for _, cmnd := range logCommands {
		log.Printf("websocket  endpoint: /%s\n", cmnd.alias)
		mux.Handle(cmnd.pattern(), cmnd)
	}

	//serving common Log simple commands terminal
	log.Println("http  endpoint     : /command")
	mux.HandleFunc("/command", simpleCommandHandler)

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	//graceful shutdown, when signal from context
	go func() {
		<-ctx.Done()
		log.Println("Shutdown Nanostarter!")
		if err := server.Shutdown(ctx); err != nil {
			return
		}
	}()

	return server.ListenAndServe()
}

// simpleCommandHandler handles simple linux commands that don't require logging
func simpleCommandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	cmdText := r.FormValue("cmd")
	if cmdText == "" {
		http.Error(w, "no command to exec", http.StatusBadRequest)
		return
	}
	log.Println("try to parse, cmdText:", cmdText)
	cmd := NewCommandFromString(cmdText)

	// give 3 seconds to finish the command
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	c := exec.CommandContext(ctx, "bash", "-c", cmd.String())
	c.Stdout = w
	c.Stderr = w
	if err := c.Run(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)
