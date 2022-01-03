package greetings

import (
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name,
// checking for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("Gladys")

	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")

	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}

// Same stuff, but for Hellos.
func TestHellosName(t *testing.T) {
	msgs, err := Hellos([]string{"Thancred", "Uberdanger", "Venat"})
	length := len(msgs)

	if length != 3 || err != nil {
		t.Fatalf(`Hellos([]string{"Thancred", "Uberdanger", "Venat"}) does not return a map, %v`, err)
	}
}

func TestHellosEmpty(t *testing.T) {
	msgs, err := Hellos([]string{"", "Uberdanger", "Venat"})
	length := len(msgs)

	if length != 0 {
		t.Fatalf(`Expected empty map from Hellos([]string{"", "Uberdanger", "Venat"}), received %d instead`, length)
	}

	if err == nil {
		t.Fatalf(`Expected error from Hellos([]string{"", "Uberdanger", "Venat"}), got nil instead`)
	}
}
