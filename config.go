package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"time"
)

type SConfig struct {
	Sites          []string
	AuthEmailLogin string
	AuthEmailPaswd string
	Receivers      []string
	TimeoutRepeat  time.Duration
	TimeoutWait    time.Duration
	SizeChan       int
}

func (s *SConfig) Read() {
	if _, err := toml.DecodeFile(fileConfig, s); err != nil {
		panic(err)
	}
	if len(s.Sites) == 0 {
		panic(errors.New("Sites list is empty"))
	}
	if s.AuthEmailLogin == "" {
		panic(errors.New("Login for email accout is empty"))
	}
	if s.AuthEmailPaswd == "" {
		panic(errors.New("Passwd for email accout is empty"))
	}
	if len(s.Receivers) == 0 {
		panic(errors.New("Receivers list for email alarming is empty"))
	}
}
