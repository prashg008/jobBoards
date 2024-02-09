// models/notification.go
package models

import (
	"fmt"

	"github.com/lib/pq"
)

// Notification represents a notification to be sent to users
type Notification struct {
	UserID     uint
	Message    string
	Channels   pq.StringArray         `gorm:"type:text[]"` // Channels to send the notification (e.g., "email", "sms")
	Additional map[string]interface{} `gorm:"type:json"`   // Additional data to be sent with the notification (e.g., "email_template", "sms_template")
}

// SendNotification sends the notification to the user through the specified channels
func (n *Notification) SendNotification() error {
	// Iterate through channels and send the notification
	for _, channel := range n.Channels {
		switch channel {
		case "email":
			err := n.sendEmail()
			if err != nil {
				return fmt.Errorf("failed to send email notification: %v", err)
			}
		case "sms":
			err := n.sendSMS()
			if err != nil {
				return fmt.Errorf("failed to send SMS notification: %v", err)
			}
		// Add more cases for other channels as needed
		default:
			return fmt.Errorf("unsupported notification channel: %s", channel)
		}
	}
	return nil
}

// sendEmail sends the notification via email
func (n *Notification) sendEmail() error {
	// Print email notification message
	fmt.Printf("Sending email notification to user %d: %s\n", n.UserID, n.Message)
	return nil
}

// sendSMS sends the notification via SMS
func (n *Notification) sendSMS() error {
	// Print SMS notification message
	fmt.Printf("Sending SMS notification to user %d: %s\n", n.UserID, n.Message)
	return nil
}
