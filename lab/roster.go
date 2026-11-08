package lab

import (
	"fmt"
	"net/mail"
	"sort"
	"strings"
)

type Recipient struct {
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Topics []string `json:"topics"`
}

type Input struct {
	Topic    string      `json:"topic"`
	Previous []Recipient `json:"previous"`
	Current  []Recipient `json:"current"`
	Required []string    `json:"required"`
}

type Plan struct {
	Topic      string   `json:"topic"`
	Added      []string `json:"added"`
	Removed    []string `json:"removed"`
	Retained   []string `json:"retained"`
	Deliveries []string `json:"deliveries"`
}

// BuildPlan compares subscribed email endpoints. It performs no delivery.
func BuildPlan(input Input) (Plan, error) {
	if !validToken(input.Topic) {
		return Plan{}, fmt.Errorf("topic must contain 1..64 lowercase letters, digits or hyphens")
	}
	previous, _, err := selected(input.Previous, input.Topic)
	if err != nil {
		return Plan{}, fmt.Errorf("previous roster: %w", err)
	}
	current, ids, err := selected(input.Current, input.Topic)
	if err != nil {
		return Plan{}, fmt.Errorf("current roster: %w", err)
	}
	if len(input.Required) > 200 {
		return Plan{}, fmt.Errorf("at most 200 required teams")
	}
	for _, id := range input.Required {
		if !validToken(id) || !ids[id] {
			return Plan{}, fmt.Errorf("required team %q is missing or not subscribed", id)
		}
	}
	plan := Plan{Topic: input.Topic, Added: []string{}, Removed: []string{}, Retained: []string{}, Deliveries: []string{}}
	for endpoint := range current {
		plan.Deliveries = append(plan.Deliveries, endpoint)
		if previous[endpoint] {
			plan.Retained = append(plan.Retained, endpoint)
		} else {
			plan.Added = append(plan.Added, endpoint)
		}
	}
	for endpoint := range previous {
		if !current[endpoint] {
			plan.Removed = append(plan.Removed, endpoint)
		}
	}
	sort.Strings(plan.Added)
	sort.Strings(plan.Removed)
	sort.Strings(plan.Retained)
	sort.Strings(plan.Deliveries)
	return plan, nil
}

func selected(roster []Recipient, topic string) (map[string]bool, map[string]bool, error) {
	if len(roster) > 200 {
		return nil, nil, fmt.Errorf("at most 200 recipients")
	}
	endpoints := make(map[string]bool)
	ids := make(map[string]bool)
	subscribed := make(map[string]bool)
	for _, recipient := range roster {
		if !validToken(recipient.ID) || ids[recipient.ID] {
			return nil, nil, fmt.Errorf("invalid or duplicate team ID %q", recipient.ID)
		}
		ids[recipient.ID] = true
		endpoint, err := canonicalEmail(recipient.Email)
		if err != nil {
			return nil, nil, fmt.Errorf("team %q: %w", recipient.ID, err)
		}
		if len(recipient.Topics) > 32 {
			return nil, nil, fmt.Errorf("team %q has more than 32 topics", recipient.ID)
		}
		seenTopics := make(map[string]bool)
		for _, entry := range recipient.Topics {
			if !validToken(entry) || seenTopics[entry] {
				return nil, nil, fmt.Errorf("team %q has invalid or duplicate topics", recipient.ID)
			}
			seenTopics[entry] = true
		}
		if seenTopics[topic] {
			endpoints[endpoint] = true
			subscribed[recipient.ID] = true
		}
	}
	return endpoints, subscribed, nil
}

func canonicalEmail(value string) (string, error) {
	if len(value) > 254 || strings.TrimSpace(value) != value {
		return "", fmt.Errorf("email must be a bare address up to 254 bytes")
	}
	for _, char := range value {
		if char < 33 || char > 126 {
			return "", fmt.Errorf("email must use printable ASCII without spaces")
		}
	}
	address, err := mail.ParseAddress(value)
	if err != nil || address.Name != "" || address.Address != value {
		return "", fmt.Errorf("invalid bare email address")
	}
	index := strings.LastIndexByte(value, '@')
	if index < 1 || index == len(value)-1 {
		return "", fmt.Errorf("email needs a local part and domain")
	}
	return value[:index+1] + strings.ToLower(value[index+1:]), nil
}

func validToken(value string) bool {
	if len(value) == 0 || len(value) > 64 {
		return false
	}
	for _, char := range value {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}
