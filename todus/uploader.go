package todus

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/websocket"
)

// Work in progress
const (
	host     string = "im.todus.cu:1756"
	port     int    = 1756
	atds_uri string = "wss://atds3.herokuapp.com"
)

var alphanumeric_Runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func parse_token(token_string string) (string, string, error) {
	type TodusClaims struct {
		Exp     float64 `json:"exp"`
		Phone   string  `json:"username"`
		Version string  `json:"version"`
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(token_string, &TodusClaims{}, nil)
	claims, ok := token.Claims.(*TodusClaims)
	if ok {
		return token_string, claims.Phone, nil
	} else {
		fmt.Println(err)
	}
	return token_string, claims.Phone, nil
}

func steal_token() (string, string, error) {
	origin := "http://localhost/"

	type req_frame struct {
		Action string `json:"accion"`
		Token  string `json:"ficha,omitempty"`
	}
	var req_token_frame = req_frame{Action: "solicitarFicha"}

	ws, err := websocket.Dial(atds_uri, "", origin)

	if err != nil {
		return "", "", errors.New("error opening websocket")
	}

	defer ws.Close()

	if err := websocket.JSON.Send(ws, req_token_frame); err != nil {
		return "", "", err
	}

	token := req_frame{}

	if err := websocket.JSON.Receive(ws, &token); err != nil {
		return "", "", err
	} else {
		return parse_token(token.Token)
	}
}

func generate_sid(length int) string {
	b := make([]rune, length)
	no_of_runes := len(alphanumeric_Runes)
	for i := range b {
		b[i] = alphanumeric_Runes[rand.Intn(no_of_runes)]
	}
	return string(b)
}

func negociate_start(response string, tls_conn *tls.Conn, authstr string, sid string) bool {

	if strings.HasPrefix(response, "<?xml version='1.0'?><stream:stream i='") &&
		strings.HasSuffix(response, "xmlns:stream='x1' f='im.todus.cu' xmlns='jc'>") {
		return true
	}

	if strings.Contains(response, "<stream:features><es xmlns='x2'><e>PLAIN</e><e>X-OAUTH2</e></es><register xmlns='http://jabber.org/features/iq-register'/></stream:features>") {
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

func Sign_url(file_size int) (string, string, error) {

	token, phone, err := steal_token()
	sid := generate_sid(5)

	if err != nil {
		return "", "", err
	}

	authstr := fmt.Sprintf("\x00%s\x00%s", phone, token)
	authstr_encoded := base64.StdEncoding.EncodeToString([]byte(authstr))
	fmt.Println(authstr)
	fmt.Println(authstr_encoded)

	conf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}

	conn, err := tls.Dial("tcp", host, conf)

	if err != nil {
		fmt.Println("SSL Error : " + err.Error())
		return "", "", err
	}

	_, err = io.WriteString(conn, "<stream:stream xmlns='jc' o='im.todus.cu' xmlns:stream='x1' v='1.0'>")

	if err != nil {
		fmt.Println("Error iniciating " + err.Error())
		return "", "", errors.New("error iniciating")
	}

	for {

		reply := make([]byte, 1024*1024)
		n, err := conn.Read(reply)

		if err != nil {
			fmt.Println("Error reading " + err.Error())
			return "", "", errors.New("error reading ")
		}

		var response string = string(reply[:n])
		fmt.Printf("Recived %d bytes\n %s\n", n, response)

		if negociate_start(response, conn, authstr_encoded, sid) {
			fmt.Println("negociating start...")
			continue
		}

		if strings.Contains(response, fmt.Sprintf("t='result' i='%s-1'>", sid)) {
			_, err1 := io.WriteString(conn, "<en xmlns='x7' u='true' max='300'/>")
			_, err2 := io.WriteString(conn, fmt.Sprintf("<iq i='%s-3' t='get'><query xmlns='todus:purl' type='0' persistent='false' size='%d' room=''></query></iq>", sid, file_size))

			if err1 != nil || err2 != nil {
				fmt.Println("Error reading " + err.Error())
				return "", "", errors.New("error reading ")
			}

			continue
		}

		if strings.HasPrefix(response, "<ed u='true' max=") {
			_, err3 := io.WriteString(conn, fmt.Sprintf("<p i='%s-4'></p>", sid))
			if err3 != nil {
				fmt.Println("Error reading " + err.Error())
				return "", "", errors.New("error reading ")
			}

			continue
		}

		if strings.HasPrefix(response, fmt.Sprintf("<iq o='%s@im.todus.cu", phone)) &&
			strings.Contains(response, "status='200'") {

			r := regexp.MustCompile(".*put='(.*)' get='(.*)' stat.*")
			res := r.FindAllStringSubmatch(response, -1)
			up := strings.ReplaceAll(res[0][1], "amp;", "")
			down := res[0][2]
			return up, down, nil
		}

		if strings.Contains(response, "<not-authorized/>") {
			return "", "", errors.New("not authorized")
		}

		if len(response) == 0 {
			return "", "", errors.New("eof")
		}

	}

}
