package scanner

import (
	"errors"
	"testing"
)

func TestResolver_Resolve_EmptyAddress(t *testing.T) {
	r := NewResolver()
	if got := r.Resolve(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestResolver_Resolve_Wildcard(t *testing.T) {
	r := NewResolver()
	if got := r.Resolve("0.0.0.0"); got != "0.0.0.0" {
		t.Errorf("expected 0.0.0.0, got %q", got)
	}
}

func TestResolver_Resolve_LookupSuccess(t *testing.T) {
	r := NewResolverWithFunc(func(addr string) ([]string, error) {
		return []string{"example.com."}, nil
	})
	got := r.Resolve("93.184.216.34")
	if got != "example.com" {
		t.Errorf("expected example.com, got %q", got)
	}
}

func TestResolver_Resolve_LookupError(t *testing.T) {
	r := NewResolverWithFunc(func(addr string) ([]string, error) {
		return nil, errors.New("lookup failed")
	})
	got := r.Resolve("1.2.3.4")
	if got != "1.2.3.4" {
		t.Errorf("expected original IP, got %q", got)
	}
}

func TestResolver_Resolve_EmptyNames(t *testing.T) {
	r := NewResolverWithFunc(func(addr string) ([]string, error) {
		return []string{}, nil
	})
	got := r.Resolve("10.0.0.1")
	if got != "10.0.0.1" {
		t.Errorf("expected original IP, got %q", got)
	}
}

func TestResolver_ResolveEntry(t *testing.T) {
	r := NewResolverWithFunc(func(addr string) ([]string, error) {
		return []string{"localhost."}, nil
	})
	e := Entry{LocalAddress: "127.0.0.1", LocalPort: 8080}
	got := r.ResolveEntry(e)
	if got != "localhost:8080" {
		t.Errorf("expected localhost:8080, got %q", got)
	}
}
