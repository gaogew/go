package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var palette = []color.Color{color.RGBA{R: 2, G: 4, B: 1}, color.RGBA{R: 45, B: 12}, color.RGBA{A: 19, B: 12}}

const (
	whiteIndex = 0
	blackIndex = 1
	otherIndex = 2
)

func main() {
	httpServer()
	go syncHTTP()
	fmt.Printf("sleep now %s\n", time.Now())
	time.Sleep(time.Second * 3)
	fmt.Printf("run now %s\n", time.Now())

	for i := 0; i < 10; i++ {
		go fetchGo("http://localhost:8000")

	}
}

/**
 * animation
 */
func animFunc() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w)
		}
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 5
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
			img.SetColorIndex(size+int(x*size+0.6), size+int(y*size+0.6), otherIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	_ = gif.EncodeAll(out, &anim)
}

// fetch url
func fetchURLFunc() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

//并发fetch
func fectchAll() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // launch a goroutine
	}

	for range os.Args[0:] {
		fmt.Println(<-ch) //从通道ch接收
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // 发送到通道ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

func httpServer() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Printf("result: %s, %s\n", r.Method, r.RemoteAddr)
}

func fetchGo(url string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error %q", err)
	}
}

var mu sync.Mutex
var count int

func syncHTTP() {
	http.HandleFunc("/", handler3)
}

func handler3(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	_, _ = fmt.Printf("URL.Path = %q\n", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	_, _ = fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

func handler2(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		_, _ = fmt.Fprintf(w, "Header[%s] = %q\n", k, v)
	}
	_, _ = fmt.Fprintf(w, "Host = %q\n", r.Host)
	_, _ = fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		_, _ = fmt.Fprintf(w, "Form [%q] = %q\n", k, v)
	}

}
