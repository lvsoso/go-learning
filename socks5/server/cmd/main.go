package main

import (
	"log"
	"server/socks5"
)

func main() {

	users := map[string]string{
		"admin":  "123456",
		"lvsoso": "123456",
	}

	server := socks5.SOCKS5Server{
		IP:   "localhost",
		Port: 1080,
		Config: &socks5.Config{
			AuthMethod: socks5.MethodPassword,
			PasswordChecker: func(username, password string) bool {
				wantPassword, ok := users[username]
				if !ok {
					return false
				}
				return wantPassword == password
			},
		},
	}

	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
