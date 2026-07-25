package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// NotificationChannelType defines supported notification channels
type NotificationChannelType string

const (
	ChannelSlack   NotificationChannelType = "slack"
	ChannelDiscord NotificationChannelType = "discord"
	ChannelJira    NotificationChannelType = "jira"
	ChannelWebhook NotificationChannelType = "webhook"
)

// NotificationService sends notifications to external services
type NotificationService struct {
	configured map[string]bool
	mu         sync.RWMutex
	logger     *zap.Logger
}

// NewNotificationService creates a new notification service
func NewNotificationService(logger *zap.Logger) *NotificationService {
	return &NotificationService{
		configured: make(map[string]bool),
		logger:     logger,
	}
}

// Configure sets up a notification channel
func (ns *NotificationService) Configure(channelType NotificationChannelType, config map[string]string) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.configured[string(channelType)] = true
	ns.logger.Info("Notification channel configured", zap.String("type", string(channelType)))

	return nil
}

// NotifyFinding sends a finding notification
func (ns *NotificationService) NotifyFinding(finding *models.UniversalFinding) error {
	ns.mu.RLock()
	channels := make([]string, 0, len(ns.configured))
	for channel := range ns.configured {
		channels = append(channels, channel)
	}
	ns.mu.RUnlock()

	for _, channel := range channels {
		switch NotificationChannelType(channel) {
		case ChannelSlack:
			if err := ns.sendSlackNotification(finding); err != nil {
				ns.logger.Error("Failed to send Slack notification", zap.Error(err))
			}
		case ChannelDiscord:
			if err := ns.sendDiscordNotification(finding); err != nil {
				ns.logger.Error("Failed to send Discord notification", zap.Error(err))
			}
		case ChannelJira:
			if err := ns.sendJiraNotification(finding); err != nil {
				ns.logger.Error("Failed to send Jira notification", zap.Error(err))
			}
		}
	}

	return nil
}

// sendSlackNotification sends a notification to Slack
func (ns *NotificationService) sendSlackNotification(finding *models.UniversalFinding) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf("🚨 New Finding: %s", finding.Title),
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*\n%s\n*Severity:* %s", finding.Title, finding.Description, finding.Severity),
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// TODO: Send to Slack webhook
	ns.logger.Debug("Slack notification prepared", zap.String("finding", finding.ID))

	return nil
}

// sendDiscordNotification sends a notification to Discord
func (ns *NotificationService) sendDiscordNotification(finding *models.UniversalFinding) error {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title": finding.Title,
				"description": finding.Description,
				"color": getSeverityColor(finding.Severity),
				"fields": []map[string]string{
					{
						"name":  "Severity",
						"value": finding.Severity,
					},
					{
						"name":  "Category",
						"value": finding.Category,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// TODO: Send to Discord webhook
	ns.logger.Debug("Discord notification prepared", zap.String("finding", finding.ID))

	return nil
}

// sendJiraNotification sends a notification/issue to Jira
func (ns *NotificationService) sendJiraNotification(finding *models.UniversalFinding) error {
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{"key": "HLFA"},
			"issuetype": map[string]string{"name": "Bug"},
			"summary": finding.Title,
			"description": finding.Description,
			"priority": map[string]string{"name": finding.Severity},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// TODO: Send to Jira API
	ns.logger.Debug("Jira notification prepared", zap.String("finding", finding.ID))

	return nil
}

func getSeverityColor(severity string) int {
	switch severity {
	case "CRITICAL":
		return 15158332 // Red
	case "HIGH":
		return 15105570 // Orange
	case "MEDIUM":
		return 15263976 // Yellow
	default:
		return 3447003 // Blue
	}
}
