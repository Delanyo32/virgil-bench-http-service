package service

import (
	"fmt"
	"log"
	"net/smtp"
	"time"
)

// NotificationService handles sending notifications to users.
type NotificationService struct {
	smtpHost string
	smtpPort int
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(smtpHost string, smtpPort int) *NotificationService {
	return &NotificationService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

// SendOrderConfirmation sends an order confirmation email.
// FLAW: sync-blocking-in-async -- blocking SMTP call in request handler path.
func (s *NotificationService) SendOrderConfirmation(userID, orderID int) error {
	to := fmt.Sprintf("user_%d@example.com", userID)
	subject := fmt.Sprintf("Order #%d Confirmed", orderID)
	body := fmt.Sprintf("Your order #%d has been confirmed. Thank you for your purchase!", orderID)

	return s.sendEmail(to, subject, body)
}

// SendShippingUpdate sends a shipping status notification.
func (s *NotificationService) SendShippingUpdate(userID, orderID int, trackingNumber string) error {
	to := fmt.Sprintf("user_%d@example.com", userID)
	subject := fmt.Sprintf("Order #%d Shipped", orderID)
	body := fmt.Sprintf("Your order #%d has been shipped. Tracking: %s", orderID, trackingNumber)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email via SMTP.
func (s *NotificationService) sendEmail(to, subject, body string) error {
	from := "noreply@ordersvc.example.com"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
		return fmt.Errorf("email send failed: %w", err)
	}

	return nil
}

// SendBulkNotification sends a notification to multiple users.
// FLAW: sync-blocking-in-async -- sends emails sequentially in a loop,
// blocking the caller for the sum of all SMTP round trips.
func (s *NotificationService) SendBulkNotification(userIDs []int, subject, body string) error {
	var failCount int
	for _, uid := range userIDs {
		to := fmt.Sprintf("user_%d@example.com", uid)
		if err := s.sendEmail(to, subject, body); err != nil {
			failCount++
			log.Printf("failed to notify user %d: %v", uid, err)
			continue
		}
	}

	if failCount > 0 {
		return fmt.Errorf("%d notifications failed out of %d", failCount, len(userIDs))
	}
	return nil
}

// ScheduleReminder schedules a delayed notification.
// Dead code -- not called from any handler or worker.
func (s *NotificationService) ScheduleReminder(userID int, message string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		to := fmt.Sprintf("user_%d@example.com", userID)
		_ = s.sendEmail(to, "Reminder", message) // FLAW: error swallowed
	}()
}
