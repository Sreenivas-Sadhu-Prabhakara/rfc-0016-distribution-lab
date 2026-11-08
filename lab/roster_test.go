package lab

import (
	"reflect"
	"testing"
)

func member(id, email string, topics ...string) Recipient {
	return Recipient{ID: id, Email: email, Topics: topics}
}

func TestPlanAddsRemovesAndDeduplicates(t *testing.T) {
	input := Input{
		Topic:    "api",
		Previous: []Recipient{member("old", "old@example.invalid", "api"), member("shared", "ops@EXAMPLE.INVALID", "api")},
		Current:  []Recipient{member("new", "new@example.invalid", "api"), member("shared", "ops@example.invalid", "api"), member("alias", "ops@example.invalid", "api"), member("cloud", "cloud@example.invalid", "cloud")},
		Required: []string{"new", "alias"},
	}
	plan, err := BuildPlan(input)
	if err != nil {
		t.Fatal(err)
	}
	want := Plan{Topic: "api", Added: []string{"new@example.invalid"}, Removed: []string{"old@example.invalid"}, Retained: []string{"ops@example.invalid"}, Deliveries: []string{"new@example.invalid", "ops@example.invalid"}}
	if !reflect.DeepEqual(plan, want) {
		t.Fatalf("got %+v, want %+v", plan, want)
	}
}

func TestRequiredMembershipAndInvalidRosters(t *testing.T) {
	cases := []Input{
		{Topic: "api", Required: []string{"absent"}},
		{Topic: "api", Current: []Recipient{member("team", "a@example.invalid", "cloud")}, Required: []string{"team"}},
		{Topic: "api", Current: []Recipient{member("a", "a@example.invalid", "api"), member("a", "b@example.invalid", "api")}},
		{Topic: "api", Current: []Recipient{member("a", "Name <a@example.invalid>", "api")}},
		{Topic: "api", Current: []Recipient{member("a", "a@example.invalid", "api", "api")}},
		{Topic: "API"},
	}
	for _, input := range cases {
		if _, err := BuildPlan(input); err == nil {
			t.Fatalf("accepted malformed roster %+v", input)
		}
	}
}

func TestContactCorrectionAndLocalPartCase(t *testing.T) {
	plan, err := BuildPlan(Input{Topic: "api", Previous: []Recipient{member("team", "A@example.invalid", "api")}, Current: []Recipient{member("team", "a@example.invalid", "api")}})
	if err != nil || len(plan.Added) != 1 || len(plan.Removed) != 1 {
		t.Fatalf("%+v, %v", plan, err)
	}
}
