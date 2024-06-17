package main

import (
	pngdetector "fake-png-detector.mod/internal/ai/png-detector"
	"fake-png-detector.mod/internal/env"
	"fmt"
)

func main() {
	if err := env.LoadEnvFile("./env/.env"); err != nil {
		fmt.Printf("Switching to using OS env variables. Err: %v\n", err)
	}

	if err := env.InitializeEnvMap("./env/.env"); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if err := pngdetector.InitializeSessionPool(); err != nil {
		fmt.Printf("%v", err)
		return
	}

	//http.HandleFunc("/huh", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "<h1>Hi there, I'm resume-gunk-api!</h1>")
	//})
	//
	//fmt.Println("Server is listening on port 8080...")
	//if err := http.ListenAndServe(":8080", api.Routes()); err != nil {
	//	fmt.Printf("Error starting the server: %v\n", err)
	//}
	fmt.Println("Hello, World!")
	pool := pngdetector.GetSessionPool()
	{
		pngdetector.GetSession(pool)
	}
	fmt.Println("Hello, World!2")
}
