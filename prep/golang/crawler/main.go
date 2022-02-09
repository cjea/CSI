package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"csi_prep_golang/pkg/model"
	"csi_prep_golang/pkg/pub"
	"csi_prep_golang/pkg/sub"
)

/*
                                             |               |
                            ||_|_|_|_|_|_|_|_|               |__|__|__|__|__|__|__|__|__|   <------------\
                          /                 \                                         \    URL QUEUE      \
             ____________/_                  \____________                             \_____________      \
             | subscriber |                  | publisher |                              |           |      |
                                                                                        |           |      |
WritePage :: (Page, Directory) -> ()  |     FetchPage :: URL -> Page
DequeuePage :: Queue -> Page          |     ParseHTML :: Page -> HTMLNode
NewQueue :: () -> Queue               |     GetAnchorTags :: HTMLNode -> []AnchorTag
GetDir :: Page -> Directory           |     GetURLs :: []AnchorTag -> []URL
                                      |     EnqueuePage :: (Page, Queue) -> ()
                                      |     Dequeue :: Queue -> URL
                                      |     Enqueue :: (URL, Queue) -> ()
                      Page :: { URL, Body }

Indexer flow                                     Fetcher flow
============                                     ============
NewQueue |> Dequeue |> GetDir |> WritePage       NewQueue |> Dequeue |> FetchPage |> ParseHTML |> GetAnchorTags |> GetURLs |> Enqueue(URLQueue)
                                                                                  |
                                                                                  |> TranslateURLs |> Enqueue

*/

func main() {
	fmt.Println("Running")
	pubDone := make(chan int)
	subDone := make(chan int)
	host, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Errorf("first arg needs to be a valid URL: %w", err))
		os.Exit(1)
	}
	pageQueue := model.NewPageQueue(100)
	urlQueue := model.NewURLQueue(100)
	urlQueue.Enqueue(model.URL{Value: host.String()})

	/*
		Once a second, launch a publisher and subscriber. When 2 pubs and 2 subs
		time out, then exit. Successful pub & sub restart counters.
	*/
	var pubTimeouts int32 = 0
	var subTimeouts int32 = 0
	for pubTimeouts < 2 && subTimeouts < 2 {
		go func() {
			errCh := make(chan error)
			go sub.Subscribe(pageQueue, errCh)
			err := <-errCh
			if err != nil {
				if errors.Is(err, model.ErrTimeout) {
					atomic.AddInt32(&pubTimeouts, 1)
				}
			}
		}()

		go func() {
			errCh := make(chan error)
			go pub.Publish(pageQueue, urlQueue, errCh)
			err := <-errCh
			if err != nil {
				if errors.Is(err, model.ErrTimeout) {
					atomic.AddInt32(&pubTimeouts, 1)
				}
			}
		}()
		select {
		case <-time.After(time.Second):
			pubTimeouts++
			subTimeouts++
		}
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			pub.Publish(pageQueue, urlQueue, pubDone)
			sub.Subscribe(pageQueue, ".", subDone)
			wg.Done()
		}()
	}
	wg.Wait()

	for {
		select {
		case <-pubDone:
			<-subDone
			fmt.Printf("done\n")
		case <-time.After(10 * time.Second):
			go pub.Publish(pageQueue, urlQueue, pubDone)
			go sub.Subscribe(pageQueue, ".", subDone)
		}
	}

}
