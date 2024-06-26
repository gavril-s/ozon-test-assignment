package comment

import (
	graphModel "ozon-test-assignment/graph/model"
	"sync"
)

type SubscriptionManager struct {
	mu          sync.Mutex
	subscribers map[int]map[int]chan *graphModel.Comment
	lastId      int
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[int]map[int]chan *graphModel.Comment),
		lastId:      0,
	}
}

func (sm *SubscriptionManager) AddSubscriber(postID int) (int, chan *graphModel.Comment) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.subscribers[postID] == nil {
		sm.subscribers[postID] = make(map[int]chan *graphModel.Comment)
	}

	sm.lastId++
	id := sm.lastId
	ch := make(chan *graphModel.Comment, 1)
	sm.subscribers[postID][id] = ch
	return id, ch
}

func (sm *SubscriptionManager) RemoveSubscriber(postID int, id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.subscribers[postID] != nil {
		close(sm.subscribers[postID][id])
		delete(sm.subscribers[postID], id)
		if len(sm.subscribers[postID]) == 0 {
			delete(sm.subscribers, postID)
		}
	}
}

func (sm *SubscriptionManager) BroadcastComment(postID int, comment *graphModel.Comment) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if subs, ok := sm.subscribers[postID]; ok {
		for _, ch := range subs {
			ch <- comment
		}
	}
}
