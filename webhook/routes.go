package webhook

import (
	"github.com/gin-gonic/gin"
)

func WebhookRouter(webhookRepo Repository, router *gin.Engine) {

	webhookService := NewService(webhookRepo)
	webhookController := NewWebhookController(webhookService)

	webhookIncomingRouter := router.Group("/webhook/incoming")
	{
		webhookIncomingRouter.POST("", webhookController.Create)
		webhookIncomingRouter.GET("/:id", webhookController.Find)
		webhookIncomingRouter.POST("/:token", webhookController.Run)
		webhookIncomingRouter.PATCH("/:id", webhookController.Update)
		webhookIncomingRouter.DELETE("/:id", webhookController.Delete)
	}

	webhooksIncomingRouter := router.Group("/webhooks/incoming")
	{
		webhooksIncomingRouter.GET("", webhookController.List)
	}

	webhookOutgoingRouter := router.Group("/webhook/outgoing")
	{
		webhookOutgoingRouter.POST("", webhookController.Create)
		webhookOutgoingRouter.GET("/:id", webhookController.Find)
		webhookOutgoingRouter.PATCH("/:id", webhookController.Update)
		webhookOutgoingRouter.DELETE("/:id", webhookController.Delete)
	}

	webhooksOutgoingRouter := router.Group("/webhooks/uutgoing")
	{
		webhooksOutgoingRouter.GET("", webhookController.List)
	}
}
