package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// type h3FileHandler struct {
// 	handler http.Handler
// }

func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./har-files/youtube")))

	mux.HandleFunc("/demo/tile", func(w http.ResponseWriter, r *http.Request) {
		log.Println("else /demo/tile 分支")
		// Small 40x40 png
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})
	mux.HandleFunc("/demo/tiles", func(w http.ResponseWriter, r *http.Request) {
		log.Println("else /demo/tiles 分支")
		io.WriteString(w, "<html><head><style>img{width:40px;height:40px;}</style></head><body>")
		for i := 0; i < 200; i++ {
			fmt.Fprintf(w, `<img src="/demo/tile?cachebust=%d">`, i)
		}
		io.WriteString(w, "</body></html>")
	})

	return mux
}

// func (h *h3FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	h.handler.ServeHTTP(w, r)
// }

func main() {
	verbose := flag.Bool("v", false, "verbose")
	www := flag.String("www", "./har-files/youtube", "www data")
	// tcp := flag.Bool("tcp", false, "also listen on TCP")
	flag.Parse()

	logger := utils.DefaultLogger

	if *verbose {
		logger.SetLogLevel(utils.LogLevelDebug)
	} else {
		logger.SetLogLevel(utils.LogLevelInfo)
	}
	logger.SetLogTimeFormat("")

	handler := setupHandler(*www)
	// quicConf := &quic.Config{}
	// certFile, keyFile := testdata.GetCertificatePaths()

	// log.Fatal(http3.ListenAndServe("localhost:6121", certFile, keyFile, handler))
	log.Fatal(http3.ListenAndServe("localhost:443", "cert.pem", "key.pem", handler))

}
