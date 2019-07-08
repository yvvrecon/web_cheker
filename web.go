package main

import (
	"crypto/tls"
	"net/http"
	"time"
	//"fmt"
)

type WebStatus struct {
	Url    string
	Status int
	Err    string
    CertFrom time.Time
    CertTo time.Time
    CertDaysHave int
    IsSSL bool
}

func checkSite(url string) WebStatus {
	var webStatus WebStatus
	webStatus.Url = url

	/*
		// Перехватываем ошибку и возвращаем управление в программу.
		defer func() {
			if r := recover(); r != nil {
				//fmt.Printf("Recovered in %T", r)
			}
		}()
	*/

	// Настройка tls, чтобы игнорировать кривые сертификаты
	transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	// Настройка клиента
	client := &http.Client{
		Timeout:   (15 * time.Second),
		Transport: transCfg,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // отключаем автоматичекое следование по редиректу
			return http.ErrUseLastResponse
		},
	}

	// Настройка запроса
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		webStatus.Err = err.Error()
		return webStatus
	}
	request.Header.Set("User-Agent", "Mozilla 555 by YVV")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Host", "yp.ru")

	// Шлем запрос
	resp, err := client.Do(request)
	if err != nil {
		webStatus.Err = err.Error()
		return webStatus
	}
	defer resp.Body.Close()

	webStatus.Status = resp.StatusCode
    webStatus.IsSSL = false
    if resp.TLS != nil {
        webStatus.IsSSL = true
        webStatus.CertFrom , webStatus.CertTo, webStatus.CertDaysHave = GetCertInfo(resp)
    }
    
	return webStatus
}

// Получить информацию о сертификате
func GetCertInfo(resp *http.Response) (tStart, tEnd time.Time, daysRest int)  { 
    loc, _ := time.LoadLocation("UTC")
    now := time.Now().In(loc)
    diff := resp.TLS.PeerCertificates[0].NotAfter.Sub(now)
    return resp.TLS.PeerCertificates[0].NotBefore, resp.TLS.PeerCertificates[0].NotAfter, int(diff.Hours() / 24)
}
