package main

var store = SimpleStore{
	Users: make(map[string]User),
}

type SimpleStore struct {
	Users map[string]User
}
