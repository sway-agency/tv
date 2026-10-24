package tv

import "testing"

func TestValidatePackage(t *testing.T) {
	for _, name := range []string{"com.netflix.ninja", "net.oqee.androidtv.store"} {
		if ValidatePackage(name) != nil {
			t.Fatalf("%q refused", name)
		}
	}
	for _, name := range []string{"", "netflix", "com.a;reboot", "com.a b", "$(id).x", "com..x", "-com.x"} {
		if ValidatePackage(name) == nil {
			t.Fatalf("%q accepted", name)
		}
	}
}

func TestValidateText(t *testing.T) {
	if ValidateText("Stranger Things 4") != nil {
		t.Fatal("plain text refused")
	}
	for _, text := range []string{"", "a;b", "a%sb", "a\tb", "émile", string(make([]byte, MaxTextLength+1))} {
		if ValidateText(text) == nil {
			t.Fatalf("%q accepted", text)
		}
	}
}

func TestValidateComponent(t *testing.T) {
	if validateComponent("com.netflix.ninja/.MainActivity") != nil {
		t.Fatal("valid component refused")
	}
	if validateComponent("com.x/.A;reboot") == nil {
		t.Fatal("component injection accepted")
	}
}

func TestParseKey(t *testing.T) {
	key, err := ParseKey("home")
	if err != nil || key.Code != 3 {
		t.Fatalf("home key: %+v %v", key, err)
	}
	_, err = ParseKey("KEYCODE_POWER")
	if err == nil {
		t.Fatal("raw keycode name accepted")
	}
}
