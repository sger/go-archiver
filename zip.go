package archiver

import "sync"

// Zip struct
type Zip struct{}

var instance *Zip
var once sync.Once

// GetInstance provides global access to Zip struct
// and ensures that our type only gets
// initialized exactly once.
func GetInstance() *Zip {
	once.Do(func() {
		instance = &Zip{}
	})
	return instance
}
