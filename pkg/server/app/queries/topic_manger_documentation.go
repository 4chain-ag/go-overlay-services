package queries

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

type TopicManagerDocumentationProvider interface {
	GetTopicManagerDocumentation(ctx context.Context) error
}

type TopicManagerDocumentationHandler struct {
	provider TopicManagerDocumentationProvider
}

func (t *TopicManagerDocumentationHandler) Handle(c *fiber.Ctx) error {
	// TODO: Add custom validation logic.
	err := t.provider.GetTopicManagerDocumentation(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseNonOK)
	}
	return c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK)
}

func NewTopicManagerDocumentationHandler(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationHandler {
	if provider == nil {
		panic("topic manager documentation provider is nil")
	}
	return &TopicManagerDocumentationHandler{
		provider: provider,
	}
}
