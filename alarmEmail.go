package main

import "fmt"

func sendEmailGmail(config SConfig, url string) {
	sender := NewSender(config.AuthEmailLogin, config.AuthEmailPaswd)

	Subject := fmt.Sprintf("Alarm on %s", url)
	message := fmt.Sprintf(`
    <!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
    <html>
    <head>
    <meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
    </head>
    <body>Проблема с доступностью узла %s <br>
    <div><i><br>
    <br>
    Regards<br> from Eliseev Vlaimir
    <i></div>
    </body>
    </html>
    `, url)
	bodyMessage := sender.WriteHTMLEmail(config.Receivers, Subject, message)
	sender.SendMail(config.Receivers, Subject, bodyMessage)
}
