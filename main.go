package main

import (
	"log"
	"time"
)

const (
	fileConfig string = "./config/config.toml"
)

var timeRepeat *time.Timer
var config SConfig

func main() {
	log.Println("### Check site by interval ###")

	config.Read()

	// Каналы, занимающиеся постановкой и обработкой задач.
	chUrls := tasksGen(config.Sites)
	chRez := workerGen(chUrls)

	// Канал - счетчик.
	deltaWait := config.TimeoutWait * time.Second
	tickerWait := time.NewTicker(deltaWait)

	// Интервал между проверками
	timeRepeat = time.NewTimer(config.TimeoutRepeat * time.Second)
	timeRepeat.Stop()

	log.Println("==============================")

	// основной бесконечный цикл программы
	for {
		select {
		case webStatus, ok := <-chRez:
			if ok {
                    printStatusHttp(webStatus)
				if webStatus.Err != "" {
					log.Println("Send email alarm")
					sendEmailGmail(config, webStatus.Url)
				}
			}
		case <-tickerWait.C:
			log.Printf("wait %d seconds...", config.TimeoutWait)
		case <-timeRepeat.C:
			log.Printf("Restart checking after %d seconds...", config.TimeoutRepeat)
			taskAdd(config.Sites, chUrls)
		}
	}
}

// Создает поток заданий
func tasksGen(urls []string) chan string {
	out := make(chan string, config.SizeChan)
	go func() {
		for _, url := range urls {
			out <- url
		}
		timeRepeat.Reset(config.TimeoutRepeat * time.Second)
		//close(out)
	}()
	return out
}

// Добавляем задания в поток заданий
func taskAdd(urls []string, chUrls chan string) {
	go func() {
		for _, url := range urls {
			chUrls <- url
		}
		timeRepeat.Reset(config.TimeoutRepeat * time.Second)
	}()
}

// Воркер, который берет задания из потока задания и обрабатывает их.
// Результаты возвращаются в поток выполненных заданий.
func workerGen(in <-chan string) chan WebStatus {
	out := make(chan WebStatus, config.SizeChan)
	go func() {
		for url := range in {
			webStatus := checkSite(url)
			out <- webStatus
		}
		//close(out)
	}()
	return out
}


func printStatusHttp(ws WebStatus) {
    if ws.IsSSL {
        log.Printf("processed: %s, status:%d, %d days, %s", ws.Url, 
        ws.Status, ws.CertDaysHave, ws.Err)
    } else {
        log.Printf("processed: %s, status:%d, %s", ws.Url, 
        ws.Status, ws.Err)
    }
}
