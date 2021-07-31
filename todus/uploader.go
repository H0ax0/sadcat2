package todus

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Work in progress
const (
	host        string = "im.todus.cu"
	port        int    = 1756
	auth_stream string = `<stream:features><es xmlns='x2'><e>PLAIN</e><e>X-OAUTH2</e></es>
						<register xmlns='http://jabber.org/features/iq-register'/></stream:features>`
)

func steal_token() (string, string, error) {
	return "", "", nil
}

func sign_url(file_size int) (string, error) {

	phone, token, err := steal_token()

	if err != nil {
		fmt.Print("blear")
	}

	authstr := fmt.Sprintf("\x00%s\x00%s", phone, token)

	conf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}

	conn, err := tls.Dial("tcp", "", conf)

	if err != nil {
		fmt.Println("SSL Error : " + err.Error())
		return "", errors.New("SSL Error")
	}

	_, err = io.WriteString(conn, "<stream:stream xmlns='jc' o='im.todus.cu' xmlns:stream='x1' v='1.0'>")

	if err != nil {
		fmt.Println("Error iniciating " + err.Error())
		return "", errors.New("error iniciating")
	}

	reply := make([]byte, 1024*1024)
	n, err := conn.Read(reply)
	if err != nil {
		fmt.Println("Error reading " + err.Error())
		return "", errors.New("error reading ")
	}
	fmt.Printf("Recived %d bytes", n)

	var response string = string(reply[:n])

	for {
		if negociate_start(response, conn, authstr, phone) {
			continue
		}

		if strings.Contains(response, fmt.Sprintf("t='result' i='%s-1'>", phone)) {
			_, err1 := io.WriteString(conn, "<en xmlns='x7' u='true' max='300'/>")
			_, err2 := io.WriteString(conn, fmt.Sprintf("<iq i='%s-3' t='get'><query xmlns='todus:purl' type='0' persistent='false' size='%d' room=''></query></iq>", phone, file_size))

			if err1 != nil || err2 != nil {
				fmt.Println("Error reading " + err.Error())
				return "", errors.New("error reading ")
			}
		}
		if strings.HasPrefix(response, "<ed u='true' max=") {
			_, err3 := io.WriteString(conn, fmt.Sprintf("<p i='%s-4'></p>", phone))
			if err3 != nil {
				fmt.Println("Error reading " + err.Error())
				return "", errors.New("error reading ")
			}
		}

		if strings.Contains(response, fmt.Sprintf("t='result' i='%s-2'>", phone)) &&
			strings.Contains(response, "status='200'") {
			r, _ := regexp.Compile(".*du='(.*)' stat.*")
			res := r.FindString(response)
			return strings.ReplaceAll(res, "amp;", ""), nil
		}

	}

}

func negociate_start(response string, tls_conn *tls.Conn,
	authstr string, sid string) bool {

	if strings.HasPrefix(response, "<?xml version='1.0'?><stream:stream i='") &&
		strings.HasSuffix(response, "xmlns:stream='x1' f='im.todus.cu' xmlns='jc'>") {
		return true
	}

	if response == auth_stream {
		_, err := io.WriteString(tls_conn, fmt.Sprintf("<ah xmlns='ah:ns' e='PLAIN'>%s</ah>", authstr))
		if err != nil {
			fmt.Println("error")
		}
		return true
	}

	if response == "<ok xmlns='x2'/>" {
		_, err := io.WriteString(tls_conn, "<stream:stream xmlns='jc' o='im.todus.cu' xmlns:stream='x1' v='1.0'>")
		if err != nil {
			fmt.Println("error")
		}
		return true
	}

	if strings.Contains(response, "<stream:features><b1 xmlns='x4'/>") {
		_, err := io.WriteString(tls_conn, fmt.Sprintf("<iq i='%s-1' t='set'><b1 xmlns='x4'></b1></iq>", sid))
		if err != nil {
			fmt.Println("error")
		}
		return true
	}

	return false
}
