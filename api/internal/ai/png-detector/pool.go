package png_detector

import (
	"fake-png-detector.mod/internal/env"
	"fmt"
	"strconv"
	"sync"
)

type SessionPool struct {
	sessions    []*FakePngDetectorSession
	queue       []int8
	modelPath   string
	mut         *sync.Mutex
	cond        *sync.Cond
	MaxSessions int8
}

var sessionPool SessionPool

func GetSessionPool() *SessionPool {
	return &sessionPool
}

func InitializeSessionPool() error {
	envMap := *env.GetEnvMap()
	maxSessions, err := strconv.Atoi(envMap["MAX_PNG_DETECTOR_SESSIONS"])
	modelPath := envMap["FAKE_PNG_DETECTOR_MODEL"]

	if err != nil {
		return fmt.Errorf("error while parsing max ort sessions: %v\n", err)
	}
	newPool, err := NewPool(int8(maxSessions), modelPath)
	if err != nil {
		return fmt.Errorf("err during pool instanstiation: %v\n", err)
	}
	sessionPool = *newPool
	return nil
}

func NewPool(maxSessions int8, modelPath string) (*SessionPool, error) {
	newPool := SessionPool{
		sessions:    make([]*FakePngDetectorSession, maxSessions),
		queue:       make([]int8, maxSessions),
		modelPath:   modelPath,
		mut:         new(sync.Mutex),
		cond:        new(sync.Cond),
		MaxSessions: maxSessions,
	}
	newPool.mut.Lock()
	for i := range maxSessions {
		newSession, err := InitializeSession(modelPath)
		if err != nil {
			return nil, err
		}
		newPool.sessions[i] = newSession
		newPool.queue = append(newPool.queue, i)
	}
	newPool.mut.Unlock()
	return &newPool, nil
}

func GetSession(pool *SessionPool) *FakePngDetectorSession {
	pool.mut.Lock()
	for len(pool.queue) == 0 {
		pool.cond.Wait()
	}
	sessionId := pool.queue[0]
	session := pool.sessions[0]
	pool.queue = pool.queue[1:]
	pool.mut.Unlock()
	defer FinishSession(pool, sessionId)
	return session
}

func FinishSession(pool *SessionPool, idx int8) {
	pool.mut.Lock()
	pool.queue = append(pool.queue, idx)
	pool.cond.Broadcast()
	pool.mut.Unlock()
}
