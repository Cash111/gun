package gun

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	engine := New()
	errChan := make(chan error)
	go func() {
		err := engine.Run(":8123")
		if err != nil {
			errChan <- err
		}
	}()
	err := <-errChan
	if err != nil {
		log.Fatal(err)
	}
}
