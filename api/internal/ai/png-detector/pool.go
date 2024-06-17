package png_detector

import "sync"

type SessionPool struct {
	sessions    []*FakePngDetectorSession
	queue       []int8
	modelPath   string
	mut         *sync.Mutex
	cond        *sync.Cond
	maxSessions int8
}

func InitializePool(maxSessions int8, modelPath string) error {
	newPool := SessionPool{
		sessions:    make([]*FakePngDetectorSession, maxSessions),
		queue:       make([]int8, maxSessions),
		modelPath:   modelPath,
		mut:         new(sync.Mutex),
		cond:        new(sync.Cond),
		maxSessions: maxSessions,
	}
	newPool.mut.Lock()
	for i := range maxSessions {
		newSession, err := InitializeSession(modelPath)
		if err != nil {
			return err
		}
		newPool.sessions[i] = newSession
		newPool.queue = append(newPool.queue, i)
	}
	newPool.mut.Unlock()
	return nil
}

func GetSession(pool *SessionPool) *FakePngDetectorSession {
	pool.mut.Lock()
	for len(pool.queue) == 0 {
		pool.cond.Wait()
	}
	sessionId := pool.queue[0]
	session := pool.sessions[]
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
