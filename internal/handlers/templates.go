package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/zerodha/fastglue"
)

// TemplateRequest represents the request body for creating/updating a template
type TemplateRequest struct {
	WhatsAppAccount string        `json:"whatsapp_account" validate:"required"` // WhatsApp account name
	Name            string        `json:"name" validate:"required"`
	DisplayName     string        `json:"display_name"`
	Language        string        `json:"language" validate:"required"`
	Category        string        `json:"category" validate:"required"` // MARKETING, UTILITY, AUTHENTICATION
	HeaderType      string        `json:"header_type"`                  // TEXT, IMAGE, DOCUMENT, VIDEO, NONE
	HeaderContent   string        `json:"header_content"`
	BodyContent     string        `json:"body_content" validate:"required"`
	FooterContent   string        `json:"footer_content"`
	Buttons         []interface{} `json:"buttons"`
	SampleValues    []interface{} `json:"sample_values"`
}

// TemplateResponse represents the response for a template
type TemplateResponse struct {
	ID              uuid.UUID     `json:"id"`
	WhatsAppAccount string        `json:"whatsapp_account"` // WhatsApp account name
	MetaTemplateID  string        `json:"meta_template_id"`
	Name            string        `json:"name"`
	DisplayName     string        `json:"display_name"`
	Language        string        `json:"language"`
	Category        string        `json:"category"`
	Status          string        `json:"status"`
	HeaderType      string        `json:"header_type"`
	HeaderContent   string        `json:"header_content"`
	BodyContent     string        `json:"body_content"`
	FooterContent   string        `json:"footer_content"`
	Buttons         []interface{} `json:"buttons"`
	SampleValues    []interface{} `json:"sample_values"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
}

// ListTemplates returns all templates for the organization
func (a *App) ListTemplates(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	// Optional filters
	accountName := string(r.RequestCtx.QueryArgs().Peek("account")) // Filter by account name
	status := string(r.RequestCtx.QueryArgs().Peek("status"))
	category := string(r.RequestCtx.QueryArgs().Peek("category"))

	query := a.DB.Where("organization_id = ?", orgID)

	if accountName != "" {
		query = query.Where("whats_app_account = ?", accountName)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var templates []models.Template
	if err := query.Order("created_at DESC").Find(&templates).Error; err != nil {
		a.Log.Error("Failed to list templates", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to list templates", nil, "")
	}

	response := make([]TemplateResponse, len(templates))
	for i, t := range templates {
		response[i] = templateToResponse(t)
	}

	return r.SendEnvelope(map[string]interface{}{
		"templates": response,
	})
}

// CreateTemplate creates a new message template
func (a *App) CreateTemplate(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req TemplateRequest
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Validate required fields
	if req.WhatsAppAccount == "" || req.Name == "" || req.Language == "" || req.Category == "" || req.BodyContent == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "whatsapp_account, name, language, category, and body_content are required", nil, "")
	}

	// Verify account belongs to organization
	var account models.WhatsAppAccount
	if err := a.DB.Where("name = ? AND organization_id = ?", req.WhatsAppAccount, orgID).First(&account).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "WhatsApp account not found", nil, "")
	}

	// Normalize template name (lowercase, underscores)
	templateName := normalizeTemplateName(req.Name)

	// Check if template with same name exists for this account
	var existingTemplate models.Template
	if err := a.DB.Where("organization_id = ? AND whats_app_account = ? AND name = ?", orgID, req.WhatsAppAccount, templateName).First(&existingTemplate).Error; err == nil {
		return r.SendErrorEnvelope(fasthttp.StatusConflict, "Template with this name already exists", nil, "")
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Name
	}

	template := models.Template{
		OrganizationID:  orgID,
		WhatsAppAccount: req.WhatsAppAccount,
		Name:            templateName,
		DisplayName:     displayName,
		Language:        req.Language,
		Category:        strings.ToUpper(req.Category),
		Status:          "DRAFT", // Local draft until submitted to Meta
		HeaderType:      strings.ToUpper(req.HeaderType),
		HeaderContent:   req.HeaderContent,
		BodyContent:     req.BodyContent,
		FooterContent:   req.FooterContent,
		Buttons:         convertToJSONBArray(req.Buttons),
		SampleValues:    convertToJSONBArray(req.SampleValues),
	}

	if err := a.DB.Create(&template).Error; err != nil {
		a.Log.Error("Failed to create template", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create template", nil, "")
	}

	return r.SendEnvelope(templateToResponse(template))
}

// GetTemplate returns a single template
func (a *App) GetTemplate(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	idStr, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || idStr == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Missing template ID", nil, "")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid template ID", nil, "")
	}

	var template models.Template
	if err := a.DB.Where("id = ? AND organization_id = ?", id, orgID).First(&template).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Template not found", nil, "")
	}

	return r.SendEnvelope(templateToResponse(template))
}

// UpdateTemplate updates a message template
func (a *App) UpdateTemplate(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	idStr, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || idStr == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Missing template ID", nil, "")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid template ID", nil, "")
	}

	var template models.Template
	if err := a.DB.Where("id = ? AND organization_id = ?", id, orgID).First(&template).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Template not found", nil, "")
	}

	// Cannot edit approved templates (Meta doesn't allow)
	if template.Status == "APPROVED" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Cannot edit approved templates", nil, "")
	}

	var req TemplateRequest
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Update fields
	if req.DisplayName != "" {
		template.DisplayName = req.DisplayName
	}
	if req.Language != "" {
		template.Language = req.Language
	}
	if req.Category != "" {
		template.Category = strings.ToUpper(req.Category)
	}
	if req.HeaderType != "" {
		template.HeaderType = strings.ToUpper(req.HeaderType)
	}
	template.HeaderContent = req.HeaderContent
	if req.BodyContent != "" {
		template.BodyContent = req.BodyContent
	}
	template.FooterContent = req.FooterContent
	if req.Buttons != nil {
		template.Buttons = convertToJSONBArray(req.Buttons)
	}
	if req.SampleValues != nil {
		template.SampleValues = convertToJSONBArray(req.SampleValues)
	}

	if err := a.DB.Save(&template).Error; err != nil {
		a.Log.Error("Failed to update template", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update template", nil, "")
	}

	return r.SendEnvelope(templateToResponse(template))
}

// DeleteTemplate deletes a message template
func (a *App) DeleteTemplate(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	idStr, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || idStr == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Missing template ID", nil, "")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid template ID", nil, "")
	}

	var template models.Template
	if err := a.DB.Where("id = ? AND organization_id = ?", id, orgID).First(&template).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Template not found", nil, "")
	}

	// If template exists on Meta, delete it there too
	if template.MetaTemplateID != "" {
		var account models.WhatsAppAccount
		if err := a.DB.Where("name = ? AND organization_id = ?", template.WhatsAppAccount, orgID).First(&account).Error; err == nil {
			// Delete from Meta API
			go a.deleteTemplateFromMeta(&account, template.Name)
		}
	}

	if err := a.DB.Delete(&template).Error; err != nil {
		a.Log.Error("Failed to delete template", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete template", nil, "")
	}

	return r.SendEnvelope(map[string]string{"message": "Template deleted successfully"})
}

// SubmitTemplate submits a template to Meta for approval
func (a *App) SubmitTemplate(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	idStr, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || idStr == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Missing template ID", nil, "")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid template ID", nil, "")
	}

	var template models.Template
	if err := a.DB.Where("id = ? AND organization_id = ?", id, orgID).First(&template).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Template not found", nil, "")
	}

	// Check if already submitted and not rejected
	if template.MetaTemplateID != "" && template.Status != "REJECTED" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Template already submitted to Meta", nil, "")
	}

	// Get the WhatsApp account
	var account models.WhatsAppAccount
	if err := a.DB.Where("name = ? AND organization_id = ?", template.WhatsAppAccount, orgID).First(&account).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "WhatsApp account not found", nil, "")
	}

	// For rejected templates, delete the old one first then create new
	if template.Status == "REJECTED" && template.MetaTemplateID != "" {
		a.Log.Info("Deleting rejected template before resubmission", "template", template.Name)
		a.deleteTemplateFromMeta(&account, template.Name)
		// Clear the old meta template ID
		template.MetaTemplateID = ""
	}

	// Submit template to Meta
	metaTemplateID, submitErr := a.submitTemplateToMeta(&account, &template)
	if submitErr != nil {
		a.Log.Error("Failed to submit template to Meta", "error", submitErr)
		return r.SendErrorEnvelope(fasthttp.StatusBadGateway, "Failed to submit template to Meta: "+submitErr.Error(), nil, "")
	}
	template.MetaTemplateID = metaTemplateID

	// Update template status
	template.Status = "PENDING"
	if err := a.DB.Save(&template).Error; err != nil {
		a.Log.Error("Failed to update template after submission", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Template submitted but failed to update local record", nil, "")
	}

	return r.SendEnvelope(map[string]interface{}{
		"message":          "Template submitted to Meta for approval",
		"meta_template_id": metaTemplateID,
		"status":           "PENDING",
		"template":         templateToResponse(template),
	})
}

// submitTemplateToMeta submits a template to Meta's API
func (a *App) submitTemplateToMeta(account *models.WhatsAppAccount, template *models.Template) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates",
		account.APIVersion, account.BusinessID)

	// Build components array
	components := []map[string]interface{}{}

	// Body component (required) - add first as per Meta's expected order
	body := map[string]interface{}{
		"type": "BODY",
		"text": template.BodyContent,
	}
	// Add examples if there are variables in body
	if strings.Contains(template.BodyContent, "{{") {
		bodyExamples := extractExamplesForComponent(template.SampleValues, "body")
		if len(bodyExamples) > 0 {
			body["example"] = map[string]interface{}{
				"body_text": [][]string{bodyExamples},
			}
		} else {
			// Count variables in body and return error if no samples provided
			varCount := strings.Count(template.BodyContent, "{{")
			if varCount > 0 {
				return "", fmt.Errorf("sample values are required for template variables. Found %d variable(s) in body but no sample values provided", varCount)
			}
		}
	}
	components = append(components, body)

	// Header component
	if template.HeaderType != "" && template.HeaderType != "NONE" {
		header := map[string]interface{}{
			"type":   "HEADER",
			"format": template.HeaderType,
		}
		switch template.HeaderType {
		case "TEXT":
			header["text"] = template.HeaderContent
			// Add example if there are variables
			if strings.Contains(template.HeaderContent, "{{") {
				headerExamples := extractExamplesForComponent(template.SampleValues, "header")
				if len(headerExamples) > 0 {
					header["example"] = map[string]interface{}{
						"header_text": headerExamples,
					}
				}
			}
		case "IMAGE", "VIDEO", "DOCUMENT":
			// For media headers, we need example handle (media ID)
			if template.HeaderContent != "" {
				header["example"] = map[string]interface{}{
					"header_handle": []string{template.HeaderContent},
				}
			}
		}
		components = append(components, header)
	}

	// Footer component
	if template.FooterContent != "" {
		components = append(components, map[string]interface{}{
			"type": "FOOTER",
			"text": template.FooterContent,
		})
	}

	// Buttons component
	if len(template.Buttons) > 0 {
		buttons := []map[string]interface{}{}
		for _, btn := range template.Buttons {
			if btnMap, ok := btn.(map[string]interface{}); ok {
				btnType, _ := btnMap["type"].(string)
				btnType = strings.ToUpper(btnType)
				btnText, _ := btnMap["text"].(string)

				// Skip buttons without text
				if btnText == "" {
					continue
				}

				button := map[string]interface{}{}

				switch btnType {
				case "QUICK_REPLY":
					button["type"] = "QUICK_REPLY"
					button["text"] = btnText
				case "URL":
					btnURL, _ := btnMap["url"].(string)
					if btnURL == "" {
						continue // Skip URL buttons without URL
					}
					button["type"] = "URL"
					button["text"] = btnText
					button["url"] = btnURL
					// Add example if URL has variable
					if strings.Contains(btnURL, "{{") {
						if example, ok := btnMap["example"].(string); ok && example != "" {
							button["example"] = []string{example}
						}
					}
				case "PHONE_NUMBER":
					phoneNum, _ := btnMap["phone_number"].(string)
					if phoneNum == "" {
						continue // Skip phone buttons without number
					}
					button["type"] = "PHONE_NUMBER"
					button["text"] = btnText
					button["phone_number"] = phoneNum
				case "COPY_CODE":
					button["type"] = "COPY_CODE"
					button["text"] = btnText
					if example, ok := btnMap["example"].(string); ok && example != "" {
						button["example"] = example
					}
				default:
					// Default to QUICK_REPLY if type is not recognized
					button["type"] = "QUICK_REPLY"
					button["text"] = btnText
				}

				if len(button) > 0 {
					buttons = append(buttons, button)
				}
			}
		}
		if len(buttons) > 0 {
			components = append(components, map[string]interface{}{
				"type":    "BUTTONS",
				"buttons": buttons,
			})
		}
	}

	// Build request payload
	payload := map[string]interface{}{
		"name":       template.Name,
		"language":   template.Language,
		"category":   template.Category,
		"components": components,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	a.Log.Info("Submitting template to Meta", "url", url, "name", template.Name, "payload", string(payloadBytes))

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+account.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		a.Log.Error("Meta API error response", "status", resp.StatusCode, "body", string(respBody))
		// Parse error response
		var errorResp struct {
			Error struct {
				Message      string `json:"message"`
				Type         string `json:"type"`
				Code         int    `json:"code"`
				ErrorUserMsg string `json:"error_user_msg"`
				ErrorUserTitle string `json:"error_user_title"`
				ErrorData    struct {
					MessagingProduct string `json:"messaging_product"`
					Details          string `json:"details"`
				} `json:"error_data"`
				ErrorSubcode int    `json:"error_subcode"`
				FBTraceID    string `json:"fbtrace_id"`
			} `json:"error"`
		}
		json.Unmarshal(respBody, &errorResp)

		// Prefer error_user_msg (more descriptive), then error_data.details, then message
		errMsg := errorResp.Error.ErrorUserMsg
		if errMsg == "" && errorResp.Error.ErrorData.Details != "" {
			errMsg = errorResp.Error.ErrorData.Details
		}
		if errMsg == "" {
			errMsg = errorResp.Error.Message
		}
		if errMsg != "" {
			return "", fmt.Errorf("%s", errMsg)
		}
		return "", fmt.Errorf("Meta API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse success response
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.ID, nil
}

// extractExamplesForComponent extracts example values for a specific component from sample_values
// Supports format: [{"component": "body", "index": 1, "value": "John"}, ...]
func extractExamplesForComponent(sampleValues models.JSONBArray, componentType string) []string {
	// Collect samples with their indices
	type indexedSample struct {
		index int
		value string
	}
	samples := []indexedSample{}

	for _, sv := range sampleValues {
		if svMap, ok := sv.(map[string]interface{}); ok {
			comp, _ := svMap["component"].(string)
			if comp == componentType {
				value, _ := svMap["value"].(string)
				if value != "" {
					// Get index (can be float64 from JSON)
					idx := 1
					if idxFloat, ok := svMap["index"].(float64); ok {
						idx = int(idxFloat)
					} else if idxInt, ok := svMap["index"].(int); ok {
						idx = idxInt
					}
					samples = append(samples, indexedSample{index: idx, value: value})
				}
			}
			// Also support legacy format with "values" array
			if svMap["component"] == componentType {
				if values, ok := svMap["values"].([]interface{}); ok {
					for i, v := range values {
						if str, ok := v.(string); ok {
							samples = append(samples, indexedSample{index: i + 1, value: str})
						}
					}
				}
			}
		}
	}

	// Sort by index and extract values
	if len(samples) > 0 {
		// Sort samples by index
		for i := 0; i < len(samples)-1; i++ {
			for j := i + 1; j < len(samples); j++ {
				if samples[i].index > samples[j].index {
					samples[i], samples[j] = samples[j], samples[i]
				}
			}
		}
		examples := make([]string, len(samples))
		for i, s := range samples {
			examples[i] = s.value
		}
		return examples
	}

	// Fallback: if no component-specific samples, try to get all string values
	examples := []string{}
	for _, sv := range sampleValues {
		if str, ok := sv.(string); ok {
			examples = append(examples, str)
		}
	}
	return examples
}

// SyncTemplates syncs templates from Meta API
func (a *App) SyncTemplates(r *fastglue.Request) error {
	orgID, err := getOrganizationID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	// Get account name from query or body
	accountName := string(r.RequestCtx.QueryArgs().Peek("account"))
	if accountName == "" {
		var body struct {
			WhatsAppAccount string `json:"whatsapp_account"`
		}
		r.Decode(&body, "json")
		accountName = body.WhatsAppAccount
	}

	if accountName == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "whatsapp_account is required", nil, "")
	}

	var account models.WhatsAppAccount
	if err := a.DB.Where("name = ? AND organization_id = ?", accountName, orgID).First(&account).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "WhatsApp account not found", nil, "")
	}

	// Fetch templates from Meta API
	templates, err := a.fetchTemplatesFromMeta(&account)
	if err != nil {
		a.Log.Error("Failed to fetch templates from Meta", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusBadGateway, "Failed to fetch templates from Meta: "+err.Error(), nil, "")
	}

	// Sync to database
	synced := 0
	for _, metaTemplate := range templates {
		template := models.Template{
			OrganizationID:  orgID,
			WhatsAppAccount: account.Name,
			MetaTemplateID:  metaTemplate.ID,
			Name:            metaTemplate.Name,
			DisplayName:     metaTemplate.Name,
			Language:        metaTemplate.Language,
			Category:        metaTemplate.Category,
			Status:          metaTemplate.Status,
		}

		// Parse components
		for _, comp := range metaTemplate.Components {
			switch comp.Type {
			case "HEADER":
				template.HeaderType = comp.Format
				if comp.Text != "" {
					template.HeaderContent = comp.Text
				}
			case "BODY":
				template.BodyContent = comp.Text
			case "FOOTER":
				template.FooterContent = comp.Text
			case "BUTTONS":
				template.Buttons = convertToJSONBArray(comp.Buttons)
			}
		}

		// Upsert
		existing := models.Template{}
		if err := a.DB.Where("organization_id = ? AND whats_app_account = ? AND name = ? AND language = ?",
			orgID, account.Name, template.Name, template.Language).First(&existing).Error; err == nil {
			// Update existing
			template.ID = existing.ID
			a.DB.Save(&template)
		} else {
			// Create new
			a.DB.Create(&template)
		}
		synced++
	}

	return r.SendEnvelope(map[string]interface{}{
		"message": fmt.Sprintf("Synced %d templates", synced),
		"count":   synced,
	})
}

// Meta API types
type MetaTemplate struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Language   string          `json:"language"`
	Category   string          `json:"category"`
	Status     string          `json:"status"`
	Components []MetaComponent `json:"components"`
}

type MetaComponent struct {
	Type    string        `json:"type"`
	Format  string        `json:"format,omitempty"`
	Text    string        `json:"text,omitempty"`
	Buttons []interface{} `json:"buttons,omitempty"`
}

func (a *App) fetchTemplatesFromMeta(account *models.WhatsAppAccount) ([]MetaTemplate, error) {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates?limit=100",
		account.APIVersion, account.BusinessID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+account.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Meta API error: %s", string(body))
	}

	var result struct {
		Data []MetaTemplate `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (a *App) deleteTemplateFromMeta(account *models.WhatsAppAccount, templateName string) {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates?name=%s",
		account.APIVersion, account.BusinessID, templateName)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+account.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.Log.Error("Failed to delete template from Meta", "error", err, "template", templateName)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		a.Log.Error("Meta API error deleting template", "status", resp.StatusCode, "body", string(body))
	}
}

// Helper functions

func templateToResponse(t models.Template) TemplateResponse {
	return TemplateResponse{
		ID:              t.ID,
		WhatsAppAccount: t.WhatsAppAccount,
		MetaTemplateID:  t.MetaTemplateID,
		Name:            t.Name,
		DisplayName:     t.DisplayName,
		Language:        t.Language,
		Category:        t.Category,
		Status:          t.Status,
		HeaderType:      t.HeaderType,
		HeaderContent:   t.HeaderContent,
		BodyContent:     t.BodyContent,
		FooterContent:   t.FooterContent,
		Buttons:         convertFromJSONBArray(t.Buttons),
		SampleValues:    convertFromJSONBArray(t.SampleValues),
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func normalizeTemplateName(name string) string {
	// Convert to lowercase and replace spaces with underscores
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	// Remove any non-alphanumeric characters except underscores
	var result strings.Builder
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			result.WriteRune(c)
		}
	}
	return result.String()
}

func convertToJSONBArray(arr []interface{}) models.JSONBArray {
	if arr == nil {
		return models.JSONBArray{}
	}
	return models.JSONBArray(arr)
}

func convertFromJSONBArray(arr models.JSONBArray) []interface{} {
	if arr == nil {
		return []interface{}{}
	}
	return []interface{}(arr)
}
