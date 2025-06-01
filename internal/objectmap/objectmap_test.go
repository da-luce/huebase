package objectmap

import (
	"testing"
)

type User struct {
	Name     string
	Email    string `mapto:"ContactEmail"`
	Age      int
	Location string `mapto:"Office"`
}

type Employee struct {
	Name         string
	ContactEmail string
	ID           string
	Office       string
}

func TestMapFields_CorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapFields(&user, &emp, nil, nil, "mapto")
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	fields := []struct {
		got, want, label string
	}{
		{emp.Name, user.Name, "Name"},
		{emp.ContactEmail, user.Email, "ContactEmail"},
		{emp.Office, user.Location, "Office"},
		{emp.ID, "E123", "ID"},
	}

	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s mismatch: got %q, want %q", f.label, f.got, f.want)
		}
	}

}

func TestMapFields_IncorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapFields(&user, &emp, nil, nil, "incorrecttag")
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	fields := []struct {
		got, want, label string
	}{
		{emp.Name, user.Name, "Name"},
		{emp.ContactEmail, emp.ContactEmail, "ContactEmail"}, // this will not map
		{emp.Office, emp.Office, "Office"},                   // this will not map
		{emp.ID, "E123", "ID"},                               // this should be the same
	}

	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s mismatch: got %q, want %q", f.label, f.got, f.want)
		}
	}
}

type Source struct {
	A string `mapto:"N.A"`
	B string `mapto:"N.B"`
	C string
	D string
}
type Nested struct {
	A string
	B string
}

type Dest struct {
	N Nested
	C string
	D string
}
