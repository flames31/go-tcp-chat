package main

import "log"

func errExit(err error) {
	log.Fatalf("ERR : %v", err)
}
