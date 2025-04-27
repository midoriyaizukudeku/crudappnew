package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actualvalue, expectedvalue T) {
	t.Helper()

	if actualvalue != expectedvalue {
		t.Errorf("wanted %v; got %v ", expectedvalue, actualvalue)
	}
}

func StringContains(t *testing.T, expected, actual string) {
	t.Helper()
	if !strings.Contains(expected, actual) {
		t.Errorf("expected %v ang got %v", expected, actual)

	}
}
