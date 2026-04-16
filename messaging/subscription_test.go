package messaging

import (
	"errors"
	"testing"
)

func TestSubscriptionValidate(t *testing.T) {
	var nilSubscription *Subscription
	if err := nilSubscription.Validate(); !errors.Is(err, ErrNilSubscription) {
		t.Fatalf("expect ErrNilSubscription, got %v", err)
	}

	if err := (&Subscription{Topic: " \t"}).Validate(); !errors.Is(err, ErrEmptyTopic) {
		t.Fatalf("expect ErrEmptyTopic, got %v", err)
	}

	if err := (&Subscription{Topic: "topic"}).Validate(); err != nil {
		t.Fatalf("expect nil error, got %v", err)
	}
}
