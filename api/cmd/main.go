package main

import (
	pngdetector "fake-png-detector.mod/internal/ai/png-detector"
	"fake-png-detector.mod/internal/env"
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
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

	if err := pngdetector.InitializeSessionPool(); err != nil {
		fmt.Printf("%v", err)
		return
	}

	pool := pngdetector.GetSessionPool()
	{
		for i := range 10 * pool.MaxSessions {
			fmt.Printf("Get session %d!\n", i)
			pngdetector.GetSession(pool)
		}
	}
}
