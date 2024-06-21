package png_detector

import (
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

	session, err := ort.NewDynamicSession[float32, float32](modelPath, []string{"modelInput"}, []string{"modelOutput"})

	sessionWrapper := FakePngDetectorSession{
		session,
		ImageSize,
	}
	return &sessionWrapper, nil
}
