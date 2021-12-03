package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/chats"
	"github.com/grafana/grafana/pkg/util"
	"github.com/grafana/grafana/pkg/web"
)

func (hs *HTTPServer) chatGetMessages(c *models.ReqContext) response.Response {
	cmd := chats.GetMessagesCmd{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	messages, err := hs.chatsService.GetMessages(c.Req.Context(), c.OrgId, c.UserId, cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "internal error", err)
	}

	result := make([]chats.MessageDto, 0, len(messages))
	for _, m := range messages {
		result = append(result, m.ToDTO())
	}

	return response.JSON(200, util.DynMap{
		"messages": result,
	})
}

func (hs *HTTPServer) chatSendMessage(c *models.ReqContext) response.Response {
	cmd := chats.SendMessageCmd{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	if c.SignedInUser.UserId == 0 && !c.SignedInUser.HasRole(models.ROLE_ADMIN) {
		return response.Error(http.StatusForbidden, "admin role required", nil)
	}
	message, err := hs.chatsService.SendMessage(c.Req.Context(), c.OrgId, c.SignedInUser.UserId, cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "internal error", err)
	}
	return response.JSON(200, util.DynMap{
		"message": message.ToDTO(),
	})
}
