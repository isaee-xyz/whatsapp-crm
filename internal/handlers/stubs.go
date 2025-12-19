package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

// Stub handlers - to be implemented

// Contact handlers - delegating to implementations in contacts.go
func (a *App) ListContacts(r *fastglue.Request) error {
	return a.ListContactsImpl(r)
}

func (a *App) CreateContact(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) GetContact(r *fastglue.Request) error {
	return a.GetContactImpl(r)
}

func (a *App) UpdateContact(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) DeleteContact(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) AssignContact(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

// Message handlers - delegating to implementations in contacts.go
func (a *App) GetMessages(r *fastglue.Request) error {
	return a.GetMessagesImpl(r)
}

func (a *App) SendMessage(r *fastglue.Request) error {
	return a.SendMessageImpl(r)
}

func (a *App) SendTemplateMessage(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) SendMediaMessage(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) MarkMessageRead(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

// Flow handlers
func (a *App) ListFlows(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) CreateFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) GetFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) UpdateFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) DeleteFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) PublishFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) DeprecateFlow(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

// Campaign handlers - delegating to implementations in campaigns.go
func (a *App) ListCampaigns(r *fastglue.Request) error {
	return a.ListCampaignsImpl(r)
}

func (a *App) CreateCampaign(r *fastglue.Request) error {
	return a.CreateCampaignImpl(r)
}

func (a *App) GetCampaign(r *fastglue.Request) error {
	return a.GetCampaignImpl(r)
}

func (a *App) UpdateCampaign(r *fastglue.Request) error {
	return a.UpdateCampaignImpl(r)
}

func (a *App) DeleteCampaign(r *fastglue.Request) error {
	return a.DeleteCampaignImpl(r)
}

func (a *App) StartCampaign(r *fastglue.Request) error {
	return a.StartCampaignImpl(r)
}

func (a *App) PauseCampaign(r *fastglue.Request) error {
	return a.PauseCampaignImpl(r)
}

func (a *App) CancelCampaign(r *fastglue.Request) error {
	return a.CancelCampaignImpl(r)
}

func (a *App) GetCampaignProgress(r *fastglue.Request) error {
	return a.GetCampaignImpl(r) // Use get campaign for now
}

func (a *App) ImportRecipients(r *fastglue.Request) error {
	return a.AddRecipientsImpl(r)
}

func (a *App) GetCampaignRecipients(r *fastglue.Request) error {
	return a.GetCampaignRecipientsImpl(r)
}

// Chatbot handlers - delegating to implementations in chatbot.go
func (a *App) GetChatbotSettings(r *fastglue.Request) error {
	return a.GetChatbotSettingsImpl(r)
}

func (a *App) UpdateChatbotSettings(r *fastglue.Request) error {
	return a.UpdateChatbotSettingsImpl(r)
}

func (a *App) ListKeywordRules(r *fastglue.Request) error {
	return a.ListKeywordRulesImpl(r)
}

func (a *App) CreateKeywordRule(r *fastglue.Request) error {
	return a.CreateKeywordRuleImpl(r)
}

func (a *App) GetKeywordRule(r *fastglue.Request) error {
	return a.GetKeywordRuleImpl(r)
}

func (a *App) UpdateKeywordRule(r *fastglue.Request) error {
	return a.UpdateKeywordRuleImpl(r)
}

func (a *App) DeleteKeywordRule(r *fastglue.Request) error {
	return a.DeleteKeywordRuleImpl(r)
}

func (a *App) ListChatbotFlows(r *fastglue.Request) error {
	return a.ListChatbotFlowsImpl(r)
}

func (a *App) CreateChatbotFlow(r *fastglue.Request) error {
	return a.CreateChatbotFlowImpl(r)
}

func (a *App) GetChatbotFlow(r *fastglue.Request) error {
	return a.GetChatbotFlowImpl(r)
}

func (a *App) UpdateChatbotFlow(r *fastglue.Request) error {
	return a.UpdateChatbotFlowImpl(r)
}

func (a *App) DeleteChatbotFlow(r *fastglue.Request) error {
	return a.DeleteChatbotFlowImpl(r)
}

func (a *App) ListAIContexts(r *fastglue.Request) error {
	return a.ListAIContextsImpl(r)
}

func (a *App) CreateAIContext(r *fastglue.Request) error {
	return a.CreateAIContextImpl(r)
}

func (a *App) GetAIContext(r *fastglue.Request) error {
	return a.GetAIContextImpl(r)
}

func (a *App) UpdateAIContext(r *fastglue.Request) error {
	return a.UpdateAIContextImpl(r)
}

func (a *App) DeleteAIContext(r *fastglue.Request) error {
	return a.DeleteAIContextImpl(r)
}

func (a *App) ListAgentTransfers(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) CreateAgentTransfer(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) ResumeFromTransfer(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) ListChatbotSessions(r *fastglue.Request) error {
	return a.ListChatbotSessionsImpl(r)
}

func (a *App) GetChatbotSession(r *fastglue.Request) error {
	return a.GetChatbotSessionImpl(r)
}

// Analytics handlers
func (a *App) GetAnalyticsOverview(r *fastglue.Request) error {
	return a.GetDashboardStats(r)
}

func (a *App) GetMessageAnalytics(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}

func (a *App) GetChatbotAnalytics(r *fastglue.Request) error {
	return r.SendErrorEnvelope(fasthttp.StatusNotImplemented, "Not implemented yet", nil, "")
}
