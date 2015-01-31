package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/jaschaephraim/lrserver"
	"golang.org/x/exp/fsnotify"
)

const (
	serveDirectory = "_serve"
	script         = "<script src=\"http://localhost:35729/livereload.js\"></script>"
)

func main() {
	app := cli.NewApp()
	app.Name = "romulus"
	app.Usage = "A tool to manage romulus sites"
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Commands = []cli.Command{
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "login to romulus",
			Action: func(c *cli.Context) {
				println("username: ", c.Args().First(), "password:", c.Args().Get(1))
			},
		},
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "serve the current directory",
			Action: func(c *cli.Context) {

				// rebuild the specific file on every reload. Just listen in the root
				// directory and write new files to the _serve folder
				// make sure to inject JS into every file.
				// also check to make sure tmp files aren't reloaded

				_, err := os.Stat(serveDirectory)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Making directory")
					err := os.Mkdir(serveDirectory, 0777)
					if err != nil {
						log.Fatal(err)
					}
				}
				files, err := ioutil.ReadDir(".")
				if err != nil {
					err = fmt.Errorf("Unable to read directory: %s", err)
					fmt.Println(err)
				}
				for _, f := range files {
					name := f.Name()
					if f.IsDir() || name[0] == '.' {
						//sdaf
					} else {
						writeFileToServeFolder(name)
					}
				}

				// Create file watcher
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					log.Fatalln(err)
				}
				defer watcher.Close()

				// Add dir to watcher
				err = watcher.Watch(".")
				if err != nil {
					log.Fatalln(err)
				}

				// Start LiveReload server
				go lrserver.ListenAndServe()

				// Start goroutine that requests reload upon watcher event
				go func() {
					for {
						event := <-watcher.Event
						if event.Name[0] != '.' {
							writeFileToServeFolder(event.Name)
							lrserver.Reload(event.Name)
						}
					}
				}()

				defer fmt.Println("hi")

				// Start serving html
				// http.Handle("/assets/",
				// 	http.StripPrefix("/assets/",
				// 		http.FileServer(http.Dir("assets"))))
				http.ListenAndServe(":3000", http.FileServer(http.Dir("./"+serveDirectory)))
				// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// 	http.ServeFile(w, r, r.URL.Path[1:])
				// })

				// http.ListenAndServe(":3000", nil)
				// cool, now use this: https://github.com/jaschaephraim/lrserver
			},
		},
	}

	app.Run(os.Args)
}

func writeFileToServeFolder(name string) {
	fileData, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	match, err := regexp.MatchString(".*\\.html", name)
	if match {
		fileData = append(fileData, script...)
		// fmt.Println(fileData)
	}
	err = ioutil.WriteFile("_serve/"+name, fileData, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
