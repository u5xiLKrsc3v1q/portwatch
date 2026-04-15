package notifier

import (
	"errors"
	"testing"
)

// fakeNotifier records calls for assertion in tests.
type fakeNotifier struct {
	calls  []string
	errOnNth int // return error on the Nth call (1-based); 0 = never
	n      int
}

func (f *fakeNotifier) Send(title, message string) error {
	f.n++
	f.calls = append(f.calls, title+":"+message)
	if f.errOnNth > 0 && f.n == f.errOnNth {
		return errors.New("fake send error")
	}
	return nil
}

func TestMulti_SendAll(t *testing.T) {
	a := &fakeNotifier{}
	b := &fakeNotifier{}
	m := NewMulti(a, b)

	if err := m.Send("title", "body"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.calls) != 1 || len(b.calls) != 1 {
		t.Errorf("expected 1 call each, got a=%d b=%d", len(a.calls), len(b.calls))
	}
}

func TestMulti_ContinuesOnError(t *testing.T) {
	a := &fakeNotifier{errOnNth: 1}
	b := &fakeNotifier{}
	m := NewMulti(a, b)

	err := m.Send("title", "body")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// b must still have been called even though a failed
	if len(b.calls) != 1 {
		t.Errorf("expected b to be called once, got %d", len(b.calls))
	}
}

func TestMulti_ReturnsLastError(t *testing.T) {
	a := &fakeNotifier{errOnNth: 1}
	b := &fakeNotifier{errOnNth: 1}
	m := NewMulti(a, b)

	if err := m.Send("t", "m"); err == nil {
		t.Fatal("expected error from multi send")
	}
}

func TestMulti_Empty(t *testing.T) {
	m := NewMulti()
	if err := m.Send("t", "m"); err != nil {
		t.Fatalf("empty multi should not error: %v", err)
	}
}
