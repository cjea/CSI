package sub

import (
	"csi_prep_golang/pkg/model"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func Subscribe(pageQueue model.PageQueue, dir string, done chan error) {
	go WritePage(pageQueue.Dequeue(), dir, done)
}

func DequeuePage(queue model.PageQueue) model.Page {
	page := model.Page{}
	select {
	case page = <-queue.Q:
		return page
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return page
	}
}

func WritePage(page model.Page, dir string, done chan error) {
	// parsed := page.ParseHTML()
	// if parsed.Node == nil {
	// 	fmt.Printf("couldn't parse page with URL %s\n", page.URL.Value)
	// 	done <- 1
	// }
	u, err := url.Parse(page.URL.Value)
	if err != nil {
		done <- fmt.Errorf("couldn't parse URL '%s'\n", page.URL.Value)
		return
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		done <- fmt.Errorf("Error creating '%s': %w\n", dir, err)
		return
	}
	filename := filepath.Join(dir, u.Path)
	err = os.WriteFile(filename, page.Data, 0666)
	if err != nil {
		done <- fmt.Errorf("Error writing to path '%s': %w\n", filename, err)
		return
	}
	done <- nil
}
