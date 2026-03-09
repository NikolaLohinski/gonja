package utils

import "testing"

func TestEscapeUsesHTMLCompatibleQuoteEntities(t *testing.T) {
	got := Escape(`<tag "quote" 'apostrophe'>`)
	want := "&lt;tag &#34;quote&#34; &#39;apostrophe&#39;&gt;"
	if got != want {
		t.Fatalf("unexpected escaped output\nwant: %q\ngot:  %q", want, got)
	}
}
