package png_detector

import "sync"

type SessionPool struct {
	sessions    []*FakePngDetectorSession
	queue       []int8
	modelPath   string
	mut         *sync.Mutex
	maxSessions int8
}
