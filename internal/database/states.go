package database

import (
	"errors"
	"log"
	"sync"
	"time"
)

// NOTE: I should 100% be doing redis for this, but its more fun this way
// (also more likely to be DDOSed but we dont talk about that)

type State struct {
	TTL         int64
	StateSecret string
}

type StateStore struct {
	mu     sync.Mutex
	States []State
}

func checkExpiredStates(store *StateStore) {
	store.mu.Lock()

	var validStates []State

	for _, state := range store.States {
		if time.Now().Unix() < state.TTL {
			validStates = append(validStates, state)
		}
	}

	store.States = validStates

	store.mu.Unlock()
}

func recursivelyCheckExpiredStates(store *StateStore) {
	log.Println("Checking For Expired States")
	checkExpiredStates(store)
	time.Sleep(time.Minute)

	go recursivelyCheckExpiredStates(store)
}

func NewStateStore() *StateStore {
	newStore := StateStore{
		States: []State{},
		mu:     sync.Mutex{},
	}

	go recursivelyCheckExpiredStates(&newStore)

	return &newStore
}

const GithubExpiryTime = 600 // 10 minutes

func (store *StateStore) AddNewState(secretValue string) {
	store.mu.Lock()
	store.States = append(store.States, State{
		TTL:         time.Now().Unix() + GithubExpiryTime,
		StateSecret: secretValue,
	})
	store.mu.Unlock()
}

func (store *StateStore) GetStates() []State {
	store.mu.Lock()

	states := store.States

	store.mu.Unlock()

	return states
}

var ErrStateNotValid = errors.New("that state is either not in my store or has expired")

func (store *StateStore) GetState(secretVal string) (*State, error) {
	var state *State

	store.mu.Lock()

	states := store.States

	store.mu.Unlock()

	for _, currState := range states {
		isStateAndNotExpired := currState.StateSecret == secretVal && time.Now().Unix() < currState.TTL
		if isStateAndNotExpired {
			copiedState := currState
			state = &copiedState
		}
	}

	if state == nil {
		return nil, ErrStateNotValid
	}

	return state, nil
}
