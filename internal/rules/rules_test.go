package rules

import (
	"testing"
)

func TestEvaluatePriorityRules_UrgentKeywords(t *testing.T) {
	rules := DefaultRules()

	tests := []struct {
		name     string
		msg      MessageMetadata
		expected int
	}{
		{
			name: "Urgent keyword in subject",
			msg: MessageMetadata{
				Sender:  "someone@example.com",
				Subject: "URGENT: Please review",
				Snippet: "This is a normal message",
			},
			expected: 1,
		},
		{
			name: "ASAP keyword in snippet",
			msg: MessageMetadata{
				Sender:  "someone@example.com",
				Subject: "Meeting request",
				Snippet: "Please respond ASAP",
			},
			expected: 1,
		},
		{
			name: "No urgent keywords",
			msg: MessageMetadata{
				Sender:  "someone@example.com",
				Subject: "Weekly update",
				Snippet: "Here is the weekly status",
			},
			expected: 0,
		},
		{
			name: "Invoice keyword (case insensitive)",
			msg: MessageMetadata{
				Sender:  "billing@company.com",
				Subject: "Your Invoice for December",
				Snippet: "Payment details enclosed",
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePriorityRules(rules, tt.msg)
			if result != tt.expected {
				t.Errorf("EvaluatePriorityRules() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestEvaluatePriorityRules_VIPSenders(t *testing.T) {
	rules := DefaultRules()
	rules.PriorityRules.VIPSenders = []string{
		"boss@company.com",
		"ceo@company.com",
	}

	tests := []struct {
		name     string
		sender   string
		expected int
	}{
		{
			name:     "VIP sender - exact match",
			sender:   "boss@company.com",
			expected: 1,
		},
		{
			name:     "VIP sender - with display name",
			sender:   "John Boss <boss@company.com>",
			expected: 1,
		},
		{
			name:     "Non-VIP sender",
			sender:   "colleague@company.com",
			expected: 0,
		},
		{
			name:     "VIP sender - case insensitive",
			sender:   "CEO@COMPANY.COM",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MessageMetadata{
				Sender:  tt.sender,
				Subject: "Regular subject",
				Snippet: "Regular message",
			}
			result := EvaluatePriorityRules(rules, msg)
			if result != tt.expected {
				t.Errorf("EvaluatePriorityRules() = %d, want %d for sender %s", result, tt.expected, tt.sender)
			}
		})
	}
}

func TestEvaluatePriorityRules_VIPDomains(t *testing.T) {
	rules := DefaultRules()
	rules.PriorityRules.VIPDomains = []string{
		"importantclient.com",
		"partner.io",
	}

	tests := []struct {
		name     string
		sender   string
		expected int
	}{
		{
			name:     "VIP domain",
			sender:   "contact@importantclient.com",
			expected: 1,
		},
		{
			name:     "VIP domain - different user",
			sender:   "support@partner.io",
			expected: 1,
		},
		{
			name:     "Non-VIP domain",
			sender:   "random@example.com",
			expected: 0,
		},
		{
			name:     "VIP domain - case insensitive",
			sender:   "sales@IMPORTANTCLIENT.COM",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MessageMetadata{
				Sender:  tt.sender,
				Subject: "Regular subject",
				Snippet: "Regular message",
			}
			result := EvaluatePriorityRules(rules, msg)
			if result != tt.expected {
				t.Errorf("EvaluatePriorityRules() = %d, want %d for sender %s", result, tt.expected, tt.sender)
			}
		})
	}
}

func TestEvaluatePriorityRules_NilRules(t *testing.T) {
	msg := MessageMetadata{
		Sender:  "urgent@example.com",
		Subject: "URGENT",
		Snippet: "ASAP",
	}

	result := EvaluatePriorityRules(nil, msg)
	if result != 0 {
		t.Errorf("EvaluatePriorityRules(nil, msg) = %d, want 0", result)
	}
}
