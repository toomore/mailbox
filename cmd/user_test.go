package cmd

import "testing"

func TestFormattedEmails(t *testing.T) {
	in := []string{
		" Foo.Bar+promo@example.com,foo.bar@example.com ",
		"foo.bar+newsletter@example.com",
		"bar@example.com",
	}
	got := formattedEmails(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 unique emails, got %d (%v)", len(got), got)
	}
	if got[0] != "foobar@example.com" {
		t.Fatalf("unexpected first email: %s", got[0])
	}
	if got[1] != "bar@example.com" {
		t.Fatalf("unexpected second email: %s", got[1])
	}
}

func TestBuildUnsubscribeUpdateQueryWithoutGroup(t *testing.T) {
	query, args := buildUnsubscribeUpdateQuery([]string{"a@example.com", "b@example.com"}, "")
	expectQuery := "UPDATE user SET alive=0 WHERE email_uni IN (?,?)"
	if query != expectQuery {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(args) != 2 {
		t.Fatalf("unexpected args length: %d", len(args))
	}
}

func TestBuildUnsubscribeUpdateQueryWithGroup(t *testing.T) {
	query, args := buildUnsubscribeUpdateQuery([]string{"a@example.com"}, "weekly")
	expectQuery := "UPDATE user SET alive=0 WHERE email_uni IN (?) AND groups=?"
	if query != expectQuery {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(args) != 2 {
		t.Fatalf("unexpected args length: %d", len(args))
	}
	if args[1] != "weekly" {
		t.Fatalf("unexpected group arg: %v", args[1])
	}
}
