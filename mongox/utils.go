package mongox

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
)

func ObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}

// DebugCommandMonitor debug监听器
func DebugCommandMonitor() *event.CommandMonitor {
	monitor := event.CommandMonitor{}
	monitor.Started = func(ctx context.Context, event *event.CommandStartedEvent) {
		if name := event.CommandName; name != "ping" {
			log.WithField("command_name", event.CommandName).
				WithField("command", event.Command.String()).
				Info("monitor mongo command start")
		}
	}
	monitor.Succeeded = func(ctx context.Context, event *event.CommandSucceededEvent) {
		if name := event.CommandName; name != "ping" {
			log.WithField("command_name", event.CommandName).
				WithField("duration", event.Duration.String()).
				Info("monitor mongo command succeeded")
		}
	}
	monitor.Failed = func(ctx context.Context, event *event.CommandFailedEvent) {
		if name := event.CommandName; name != "ping" {
			log.WithField("command_name", event.CommandName).
				WithField("duration", event.Duration.String()).
				Info("monitor mongo command failed")
		}
	}

	return &monitor
}
