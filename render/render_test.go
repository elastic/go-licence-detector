package render

import "testing"

func TestCanonical(t *testing.T) {
	cases := map[string]string{
		"v1":                        "v1.0.0",
		"v1.1":                      "v1.1.0",
		"v3.2.1":                    "v3.2.1",
		"v3.2.1-meta":               "v3.2.1",
		"v1-meta":                   "v1.0.0",
		"v1.2-meta":                 "v1.2.0",
		"v1.2.3-meta":               "v1.2.3",
		"v1+meta":                   "v1.0.0",
		"v1.2+meta":                 "v1.2.0",
		"v1.2.3+meta":               "v1.2.3",
		"v1.2.3+incompatible":       "v1.2.3",
		"v1.2.3-20200707-123456abc": "v1.2.3",
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got := CanonicalVersion(in)
			if got != want {
				t.Errorf("Conversion missmatch for %q. Want: %q, Got: %q", in, want, got)
			}
		})
	}
}

func TestRevision(t *testing.T) {
	cases := map[string]string{
		"v1":                                     "",
		"v3.2.1-meta":                            "",
		"v1-meta":                                "",
		"v1.2-meta":                              "",
		"v1.2.3-meta":                            "",
		"v1+meta":                                "",
		"v1.2+meta":                              "",
		"v1.2.3+meta":                            "",
		"v1.2.3+incompatible":                    "",
		"v1.2.3-20200707-123456abc":              "123456abc",
		"v1.2.3-20200707-123456abc-meta":         "123456abc",
		"v1.2.3-20200707-123456abc+meta":         "123456abc",
		"v1.2.3-20200707-123456abc+incompatible": "123456abc",
		"v1.2.3-20200707-oops123456abc":          "",
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got := Revision(in)
			if got != want {
				t.Errorf("Conversion missmatch for %q. Want: %q, Got: %q", in, want, got)
			}
		})
	}
}
