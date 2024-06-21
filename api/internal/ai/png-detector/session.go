package png_detector

import (
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
	"os"
)

const ImageSize = 224

type FakePngDetectorSession struct {
	*ort.DynamicSession[float32, float32]
	ImageSize int16
}

func InitializeSession(modelPath string) (*FakePngDetectorSession, error) {
	_, err := os.Stat(modelPath)
	if err != nil {
		return nil, err
	}

	err = ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("could not initialize ORT env\n%v\n", err)
	}

	session, err := ort.NewDynamicSession[float32, float32](modelPath, []string{"modelInput"}, []string{"modelOutput"})

	sessionWrapper := FakePngDetectorSession{
		session,
		ImageSize,
	}
	return &sessionWrapper, nil
}
