package png_detector

import (
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
	"os"
)

const ImageSize = 224

type FakePngDetectorSession struct {
	*ort.Session[float32]
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

	inputShape := ort.NewShape(1, 3, ImageSize, ImageSize)
	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return nil, err
	}
	defer inputTensor.Destroy()

	outputShape := ort.NewShape(2)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewSession[float32](modelPath,
		[]string{"modelInput"}, []string{"modelOutput"},
		[]*ort.Tensor[float32]{inputTensor}, []*ort.Tensor[float32]{outputTensor})
	defer session.Destroy()

	sessionWrapper := FakePngDetectorSession{
		session,
		ImageSize,
	}
	return &sessionWrapper, nil
}
