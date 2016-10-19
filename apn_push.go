// httpserver
package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	apns "github.com/sideshow/apns2"
)

type App_conf struct {
	cert_file string
	topic     string
}

var client *apns.Client

type Response struct {
	ApnsID     string
	Reason     string
	StatusCode int
	Timestamp  apns.Time
}

var auth_token string = ""

func main() {
	client = apns.NewClient(tls.Certificate{}).Development() //Production()

	http.HandleFunc("/apn_push", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Please send a request body", 406)
			return
		}

		device_token := r.FormValue("token")
		topic := r.FormValue("topic")

		if auth_token == "" {
			auth_token, err = gen_token()
			if err != nil {
				http.Error(w, "error gen_token() : "+err.Error(), 500)
			}
		}
		//------------------
		notification := &apns.Notification{}
		notification.DeviceToken = device_token
		notification.Topic = topic
		notification.Authorization = "bearer " + auth_token

		notification.Payload = body

		res, err := client.Push(notification)

		if err == nil {
			rsp := Response{ApnsID: res.ApnsID, Reason: res.Reason, StatusCode: res.StatusCode, Timestamp: res.Timestamp}
			json.NewEncoder(w).Encode(rsp)
			return
		}

		if res.StatusCode == 403 {
			if res.Reason == "ExpiredProviderToken" {
				auth_token, err = gen_token()
				if err != nil {
					http.Error(w, "error gen_token() : "+err.Error(), 500)
				}
				//retry
				notification.Authorization = "bearer " + auth_token
				res, err := client.Push(notification)

				rsp := Response{ApnsID: res.ApnsID, Reason: res.Reason, StatusCode: res.StatusCode, Timestamp: res.Timestamp}
				json.NewEncoder(w).Encode(rsp)

				if err == nil {
					return
				}
			}
		}
		log.Println("Push Error:", err)
		//return

	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server start ok")
	select {}
}

func gen_token() (string, error) {
	key := `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49xCzCb63yoeRDYHajKuQpRT5J+lkCtEzX2Lr6xqpL+FmwT7hM8sG4CtpwRTOgCgYIKoZIzj0DAQehRANCAASvdynMrxUs6gqF/pIyFIPuDhITZ99ZM3kQ7hds/XlaNqwGWmYeWyqKkOPSsBEfMGhBWofC/KU2Ez2yGOdDVS41
-----END PRIVATE KEY-----`
	iss := "ABCD123132"
	kid := "FGT233DS90"

	claims := &jwt.StandardClaims{
		IssuedAt: time.Now().Unix(),
		Issuer:   iss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = kid
	skey, err := jwt.ParsePKCS8PrivateKeyFromPEM([]byte(key))
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(skey)
	if err != nil {
		return "", err
	}
	log.Println("gen token: ", tokenString)
	return tokenString, nil
}

// curl -i "127.0.0.1/apn_push?token=15323ce672ff91aeaaa68d44ef945840688f561e5568fb6bf2e0d0f78d937b6e&app_name=YouAPPNAME" -d '{"aps" : { "alert" : "Hello APNs" } }'
