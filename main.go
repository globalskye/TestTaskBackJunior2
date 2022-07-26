package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const fileName string = "urls"
const str string = "Go"

func main() {

	urls, err := getAllUrls(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	var result int

	ch := make(chan struct{}, 5)
	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		ch <- struct{}{}
		go func(i int) {
			count, _ := getCountStringFromURL(urls[i], str)
			<-ch
			fmt.Printf("Count for %s : %v\n", strings.ReplaceAll(urls[i], "\r", ""), count)
			result += count

			wg.Done()

		}(i)
	}
	wg.Wait()
	fmt.Printf("Total : %v", result)

}
func getAllUrls(fileName string) (urls []string, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	urls = strings.Split(string(bytes), "\n")
	if len(urls) == 0 {
		return nil, errors.New("File dont have no one url")
	}
	return
}

func getCountStringFromURL(url, str string) (res int, err error) {
	url = strings.ReplaceAll(url, "\r", "")
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res = strings.Count(string(bytes), str)
	return
}
