package gamesocket

import (
	"testing"
)

func setupTest(t *testing.T) (func(t *testing.T), string) {
	t.Log("Setting up testing")

	// gamePool := NewGamePool()
	// lav client. Der findes mock metode til request. Kan vi finder en til writer?

	return func(t *testing.T) {
		t.Log("Tearing down testing")
	}, "HEJ"
}

func TestRegisterCharacter(t *testing.T) {
	teardown, test := setupTest(t)
	defer teardown(t)

	if test != "HEJ" {
		t.Errorf("expected test to be LOL, but got %s", test)
	}

}
