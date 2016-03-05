package reloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jmhodges/clock"
)

var clk = clock.Default()

// New loads the filename provided, and calls the callback.  It then spawns a
// goroutine to check for updates to that file, calling the callback again with
// any new contents. New returns the error value returned from the first call to
// callback, and discards subsequent return values.  If there is an error
// stat'ing the file or reading it, callback will be called with an error
// parameter.
func New(filename string, callback func([]byte, error) error) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	loop := func() {
		for {
			clk.Sleep(1 * time.Second)
			fmt.Println("tick")
			currentFileInfo, err := os.Stat(filename)
			if err != nil {
				fmt.Printf("err %s\n", err)
				callback(nil, err)
				continue
			}
			fmt.Printf("times are changin' %s %s \n", currentFileInfo.ModTime(), fileInfo.ModTime())
			if currentFileInfo.ModTime().After(fileInfo.ModTime()) {
				fmt.Printf("reloading\n")
				b, err := ioutil.ReadFile(filename)
				if err != nil {
					fmt.Printf("err %s\n", err)
					callback(nil, err)
					continue
				}
				fmt.Printf("success %s\n", string(b))
				fileInfo = currentFileInfo
				callback(b, nil)
			}
		}
	}
	err = callback(b, nil)
	go loop()
	return err
}
