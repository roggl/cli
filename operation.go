package main

type operation int

const (
	get operation = iota
	create
	startTrack
	stopTrack
)
