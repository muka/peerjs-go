package emitter

import (
	"github.com/chuckpreslar/emission"
)

//EventHandler wrap an event callback
type EventHandler func(interface{})

// NewEmitter initializes an Emitter
func NewEmitter() Emitter {
	return Emitter{
		emitter: emission.NewEmitter(),
	}
}

// Emitter exposes an EventEmitter-like interface
type Emitter struct {
	emitter *emission.Emitter
}

//Emit emits an event with contextual data
func (p *Emitter) Emit(event string, data interface{}) {
	// log.Printf("EMIT %s %++v", event, data)
	p.emitter.Emit(event, data)
}

//On register a function. Note that the pointer to the function need to be
//the same to be removed with Off
func (p *Emitter) On(event string, handler EventHandler) {
	p.emitter.On(event, handler)
}

//Off remove a listener function, pointer of the function passed must match with the one
func (p *Emitter) Off(event string, handler EventHandler) {
	p.emitter.Off(event, handler)
}
