package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// EventType represents the type of audit event.
type EventType string

const (
	// Authentication events
	EventUserLogin    EventType = "USER_LOGIN"
	EventUserLogout   EventType = "USER_LOGOUT"
	EventUserRegister EventType = "USER_REGISTER"
	EventPasswordChange EventType = "PASSWORD_CHANGE"
	
	// Payment events
	EventPaymentCreated   EventType = "PAYMENT_CREATED"
	EventPaymentCompleted EventType = "PAYMENT_COMPLETED"
	EventPaymentFailed    EventType = "PAYMENT_FAILED"
	EventPaymentCanceled  EventType = "PAYMENT_CANCELED"
	
	// Subscription events
	EventSubscriptionCreated   EventType = "SUBSCRIPTION_CREATED"
	EventSubscriptionActivated EventType = "SUBSCRIPTION_ACTIVATED"
	EventSubscriptionCanceled  EventType = "SUBSCRIPTION_CANCELED"
	EventSubscriptionSuspended EventType = "SUBSCRIPTION_SUSPENDED"
	
	// Customer events
	EventCustomerCreated EventType = "CUSTOMER_CREATED"
	EventCustomerUpdated EventType = "CUSTOMER_UPDATED"
	EventCustomerDeleted EventType = "CUSTOMER_DELETED"
	
	// Admin/Staff events
	EventAdminLogin  EventType = "ADMIN_LOGIN"
	EventAdminAction EventType = "ADMIN_ACTION"
)

// Event represents an audit log entry.
type Event struct {
	EventType   EventType              `json:"event_type"`
	UserID      uint64                 `json:"user_id,omitempty"`
	Username    string                 `json:"username,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	ResourceID  uint64                 `json:"resource_id,omitempty"`
	ResourceType string                `json:"resource_type,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Status      string                 `json:"status"` // SUCCESS, FAILURE
	Error       string                 `json:"error,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Logger logs audit events.
type Logger struct {
	log zerolog.Logger
}

// New creates a new audit logger.
func New(log zerolog.Logger) *Logger {
	return &Logger{log: log.With().Str("component", "audit").Logger()}
}

// LogEvent logs an audit event.
func (l *Logger) LogEvent(ctx context.Context, event *Event) {
	if event == nil {
		return
	}
	
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	
	logBuilder := l.log.Info().
		Str("event_type", string(event.EventType)).
		Str("status", event.Status).
		Time("timestamp", event.Timestamp)
	
	if event.UserID > 0 {
		logBuilder = logBuilder.Uint64("user_id", event.UserID)
	}
	if event.Username != "" {
		logBuilder = logBuilder.Str("username", event.Username)
	}
	if event.IP != "" {
		logBuilder = logBuilder.Str("ip", event.IP)
	}
	if event.UserAgent != "" {
		logBuilder = logBuilder.Str("user_agent", event.UserAgent)
	}
	if event.ResourceID > 0 {
		logBuilder = logBuilder.
			Uint64("resource_id", event.ResourceID).
			Str("resource_type", event.ResourceType)
	}
	if event.Action != "" {
		logBuilder = logBuilder.Str("action", event.Action)
	}
	if event.Error != "" {
		logBuilder = logBuilder.Str("error", event.Error)
	}
	
	logBuilder.Interface("details", event.Details).Msg(fmt.Sprintf("Audit: %s", event.EventType))
}

// LogUserLogin logs a user login event.
func (l *Logger) LogUserLogin(ctx context.Context, userID uint64, username, ip string, success bool) {
	event := &Event{
		EventType: EventUserLogin,
		UserID:    userID,
		Username:  username,
		IP:        ip,
		Status:    "SUCCESS",
	}
	if !success {
		event.Status = "FAILURE"
		event.Error = "invalid credentials"
	}
	l.LogEvent(ctx, event)
}

// LogPaymentEvent logs payment-related events.
func (l *Logger) LogPaymentEvent(ctx context.Context, eventType EventType, paymentID, userID uint64, amount float64, status string, details map[string]interface{}) {
	event := &Event{
		EventType:    eventType,
		UserID:       userID,
		ResourceID:   paymentID,
		ResourceType: "PAYMENT",
		Status:       status,
		Details:      details,
	}
	if details == nil {
		event.Details = make(map[string]interface{})
	}
	event.Details["amount"] = amount
	l.LogEvent(ctx, event)
}

// LogSubscriptionEvent logs subscription-related events.
func (l *Logger) LogSubscriptionEvent(ctx context.Context, eventType EventType, subscriptionID, customerID, userID uint64, status string, details map[string]interface{}) {
	event := &Event{
		EventType:    eventType,
		UserID:       userID,
		ResourceID:   subscriptionID,
		ResourceType: "SUBSCRIPTION",
		Status:       status,
		Details:      details,
	}
	if details == nil {
		event.Details = make(map[string]interface{})
	}
	event.Details["customer_id"] = customerID
	l.LogEvent(ctx, event)
}

// LogAdminAction logs administrative actions.
func (l *Logger) LogAdminAction(ctx context.Context, adminID uint64, adminUsername, action string, resourceType string, resourceID uint64, details map[string]interface{}) {
	event := &Event{
		EventType:    EventAdminAction,
		UserID:       adminID,
		Username:     adminUsername,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Status:       "SUCCESS",
		Details:      details,
	}
	l.LogEvent(ctx, event)
}

// LogSecurityEvent logs security-related events (failed auth, suspicious activity, etc).
func (l *Logger) LogSecurityEvent(ctx context.Context, eventType EventType, userID uint64, ip, userAgent, details string) {
	event := &Event{
		EventType: eventType,
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
		Status:    "FAILURE",
		Error:     details,
		Details: map[string]interface{}{
			"severity": "HIGH",
		},
	}
	l.LogEvent(ctx, event)
}
