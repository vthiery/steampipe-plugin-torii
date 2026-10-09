package torii

import (
	"encoding/json"
	"testing"
)

// Torii sends contract owners as a numeric user id on some tenants and as a
// name on others. Both must decode, and both must survive a round trip into
// the string column.
func TestContractOwnerAcceptsStringAndNumber(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{"numeric id", `{"id":1,"owner":4815162342}`, "4815162342"},
		{"string name", `{"id":1,"owner":"ada@example.com"}`, "ada@example.com"},
		{"null", `{"id":1,"owner":null}`, ""},
		{"absent", `{"id":1}`, ""},
		{"numeric string", `{"id":1,"owner":"42"}`, "42"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var c Contract
			if err := json.Unmarshal([]byte(tc.body), &c); err != nil {
				t.Fatalf("unmarshal %s: %v", tc.body, err)
			}
			if string(c.Owner) != tc.want {
				t.Errorf("owner = %q, want %q", c.Owner, tc.want)
			}
		})
	}
}

func TestContractOwnerRejectsNonScalar(t *testing.T) {
	var c Contract
	if err := json.Unmarshal([]byte(`{"owner":{"email":"ada@example.com"}}`), &c); err == nil {
		t.Fatal("expected an error for an object owner, got nil")
	}
}
