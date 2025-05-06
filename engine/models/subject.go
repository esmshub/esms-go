package models

import (
	"slices"
	"sync"

	"go.uber.org/zap"
)

type Subject interface {
	Subscribe(Observer)
	Unsubscribe(Observer)
}

type SubjectImpl struct {
	mut       *sync.RWMutex
	observers []Observer
}

func (s *SubjectImpl) Subscribe(o Observer) {
	s.mut.Lock()
	defer s.mut.Unlock()

	observers := s.observers
	if observers == nil {
		observers = []Observer{}
	}

	if !slices.Contains(observers, o) {
		zap.L().Debug("subscribing observer", zap.Any("observer", o))
		s.observers = append(observers, o)
	} else {
		zap.L().DPanic("observer already subscribed", zap.Any("observer", o))
	}
}

func (s *SubjectImpl) Unsubscribe(observer Observer) {
	s.mut.Lock()
	defer s.mut.Unlock()

	currentObservers := s.observers
	if currentObservers != nil {
		curLen := len(currentObservers)

		currentObservers = slices.DeleteFunc(currentObservers, func(o Observer) bool {
			return o == observer
		})

		if curLen == len(currentObservers) {
			zap.L().DPanic("observer not found", zap.Any("observer", observer))
		}

		zap.L().Debug("unsubscribing observer", zap.Any("observer", observer))
		s.observers = currentObservers
	} else {
		zap.L().DPanic("no observers found")
	}
}

func (s *SubjectImpl) UnsubscribeAll() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.observers = []Observer{}
}

func NewSubject() *SubjectImpl {
	return &SubjectImpl{
		mut:       &sync.RWMutex{},
		observers: []Observer{},
	}
}
