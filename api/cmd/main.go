package main

import (
	pngdetector "fake-png-detector.mod/internal/ai/png-detector"
	"fake-png-detector.mod/internal/env"
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
	"sync"
	"time"
)

func main() {
	//if err := env.LoadEnvFile("./env/.env"); err != nil {
	//	fmt.Printf("Switching to using OS env variables. Err: %v\n", err)
	//}

	if err := env.InitializeEnvMap("./env/.env"); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	envMap := *env.GetEnvMap()
	ort.SetSharedLibraryPath(envMap["ORT_LIB_PATH"])
	err := ort.InitializeEnvironment()
	if err != nil {
		fmt.Printf("could not initialize ORT env\n%v\n", err)
		return
	}

	if err := pngdetector.InitializeSessionPool(); err != nil {
		fmt.Printf("%v", err)
		return
	}

	pool := pngdetector.GetSessionPool()
	var wait sync.WaitGroup
	for i := range 10 * pool.MaxSessions {
		go func() {
			wait.Add(1)
			defer wait.Done()
			_, clean := pngdetector.GetSession(pool)
			defer clean()
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Got session %d!\n", i)
		}()
	}
	wait.Wait()
}
