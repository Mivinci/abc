package store

import (
	"errors"
	"time"
)

var (
	ErrDup      = errors.New("key exists")
	ErrExpired  = errors.New("key has expired")
	ErrNotFound = errors.New("key does not exists")
)

const Forever = -1

type Store interface {
	Get(string) (interface{}, error)
	Add(string, interface{}, time.Duration) error
	Set(string, interface{}, time.Duration) error
	Refresh(string, time.Duration) error
	Remove(string) error
	GetAndRemove(string) (interface{}, error)
}

func NewMemoryStore() Store {
	return &memoryStore{c: make(map[string]*entry)}
}

type entry struct {
	dl  int64
	val interface{}
}

type EvictFunc func(key string, val interface{})

type memoryStore struct {
	c       map[string]*entry
	OnEvict EvictFunc
}

func (s *memoryStore) Get(key string) (val interface{}, err error) {
	if ent, ok := s.c[key]; ok {
		if ent.dl < 0 || ent.dl > time.Now().Unix() {
			return ent.val, nil
		}

		delete(s.c, key)
		if s.OnEvict != nil {
			s.OnEvict(key, ent.val)
		}

		return nil, ErrExpired
	}
	return nil, ErrNotFound
}

func (s *memoryStore) Add(key string, value interface{}, d time.Duration) error {
	if _, ok := s.c[key]; ok {
		return ErrDup
	}
	var dl int64
	if d < 0 {
		dl = -1
	} else {
		dl = time.Now().Add(d).Unix()
	}
	s.c[key] = &entry{dl, value}
	return nil
}

func (s *memoryStore) Set(key string, value interface{}, d time.Duration) error {
	if ent, ok := s.c[key]; ok {
		ent.val = value
		if d >= 0 {
			ent.dl = time.Now().Add(d).Unix()
		}
		return nil
	}
	return s.Add(key, value, d)
}

func (s *memoryStore) Refresh(key string, d time.Duration) error {
	if ent, ok := s.c[key]; ok {
		ent.dl = time.Now().Add(d).Unix()
		return nil
	}
	return ErrNotFound
}

func (s *memoryStore) Remove(key string) error {
	if ent, ok := s.c[key]; ok {
		delete(s.c, key)
		if s.OnEvict != nil {
			s.OnEvict(key, ent.val)
		}
		return nil
	}
	return ErrNotFound
}

func (s *memoryStore) GetAndRemove(key string) (value interface{}, err error) {
	if value, err = s.Get(key); err != nil {
		return
	}
	delete(s.c, key)
	if s.OnEvict != nil {
		s.OnEvict(key, value)
	}
	return
}
