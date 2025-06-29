package config

import "sync/atomic"

// TODO: I did not like the Struct Method thing the lesson did
// the only item in the config stuct is a count variable that is synchronized, no problem with it being global
// the method being a handler didn't feel right; like it was mixing behavior that should be seperated
// unless a need arises for encapsulation, I'm doing this

// I would like a way to protect the variable, singleton, but for internal code, do I care?
var FileServerHits = atomic.Int32{}
