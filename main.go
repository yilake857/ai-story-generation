package main

import (
	"log"
	"net/http"
)

func Hello(wr http.ResponseWriter, r *http.Request) {
	_, err := wr.Write([]byte(`
_____.__          __    __                   .___                                     
_/ ____\  |  __ ___/  |__/  |_  ___________  __| _/______   ____ _____    _____   ______
\   __\|  | |  |  \   __\   __\/ __ \_  __ \/ __ |\_  __ \_/ __ \\__  \  /     \ /  ___/
 |  |  |  |_|  |  /|  |  |  | \  ___/|  | \/ /_/ | |  | \/\  ___/ / __ \|  Y Y  \\___ \ 
 |__|  |____/____/ |__|  |__|  \___  >__|  \____ | |__|    \___  >____  /__|_|  /____  >
								   \/           \/             \/     \/      \/     \/ 	
`))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", Hello)
	err := http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}
