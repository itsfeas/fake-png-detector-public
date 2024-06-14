package ai

import (
	ort "github.com/yalue/onnxruntime_go"
)

type FakePngDetectorModel struct {
	*ort.Session[float32]
	ImageSize int16
}

const ImageSize = 224

var fakePngDetector FakePngDetectorModel

func GetFakePngDetectorModel() *FakePngDetectorModel {
	return &fakePngDetector
}

func (model *FakePngDetectorModel) InitializeSession(path string) error {
	err := ort.InitializeEnvironment()
	if err != nil {
		return err
	}

	inputShape := ort.NewShape(1, 3, ImageSize, ImageSize)
	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return err
	}
	defer inputTensor.Destroy()

	outputShape := ort.NewShape(2)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewSession[float32](path,
		[]string{"modelInput"}, []string{"modelOutput"},
		[]*ort.Tensor[float32]{inputTensor}, []*ort.Tensor[float32]{outputTensor})
	defer session.Destroy()

	fakePngDetector = FakePngDetectorModel{
		session,
		ImageSize,
	}
	return nil
}
