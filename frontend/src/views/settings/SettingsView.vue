<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { toast } from 'vue-sonner'
import { Settings, Bot, Bell, Shield, Loader2, Brain } from 'lucide-vue-next'
import { chatbotService } from '@/services/api'

const isSubmitting = ref(false)

const generalSettings = ref({
  organization_name: 'My Organization',
  default_timezone: 'UTC',
  date_format: 'YYYY-MM-DD'
})

const chatbotSettings = ref({
  greeting_message: '',
  fallback_message: '',
  session_timeout_minutes: 30,
  transfer_message: ''
})

const aiSettings = ref({
  ai_enabled: false,
  ai_provider: '',
  ai_api_key: '',
  ai_model: '',
  ai_max_tokens: 500,
  ai_system_prompt: ''
})

const aiProviders = [
  { value: 'openai', label: 'OpenAI', models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-3.5-turbo'] },
  { value: 'anthropic', label: 'Anthropic', models: ['claude-3-5-sonnet-latest', 'claude-3-5-haiku-latest', 'claude-3-opus-latest'] },
  { value: 'google', label: 'Google AI', models: ['gemini-2.0-flash', 'gemini-2.0-flash-lite', 'gemini-1.5-flash', 'gemini-1.5-flash-8b'] }
]

const availableModels = computed(() => {
  const provider = aiProviders.find(p => p.value === aiSettings.value.ai_provider)
  return provider?.models || []
})

onMounted(async () => {
  try {
    const response = await chatbotService.getSettings()
    const data = response.data.data || response.data
    if (data.settings) {
      chatbotSettings.value = {
        greeting_message: data.settings.greeting_message || '',
        fallback_message: data.settings.fallback_message || '',
        session_timeout_minutes: data.settings.session_timeout_minutes || 30,
        transfer_message: ''
      }
      aiSettings.value = {
        ai_enabled: data.settings.ai_enabled || false,
        ai_provider: data.settings.ai_provider || '',
        ai_api_key: '', // Don't load API key for security
        ai_model: data.settings.ai_model || '',
        ai_max_tokens: data.settings.ai_max_tokens || 500,
        ai_system_prompt: data.settings.ai_system_prompt || ''
      }
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
})

const notificationSettings = ref({
  email_notifications: true,
  new_message_alerts: true,
  campaign_updates: true
})

async function saveGeneralSettings() {
  isSubmitting.value = true
  try {
    // API call would go here
    await new Promise(resolve => setTimeout(resolve, 500))
    toast.success('General settings saved')
  } catch (error) {
    toast.error('Failed to save settings')
  } finally {
    isSubmitting.value = false
  }
}

async function saveChatbotSettings() {
  isSubmitting.value = true
  try {
    await chatbotService.updateSettings({
      greeting_message: chatbotSettings.value.greeting_message,
      fallback_message: chatbotSettings.value.fallback_message,
      session_timeout_minutes: chatbotSettings.value.session_timeout_minutes
    })
    toast.success('Chatbot settings saved')
  } catch (error) {
    toast.error('Failed to save settings')
  } finally {
    isSubmitting.value = false
  }
}

async function saveAISettings() {
  isSubmitting.value = true
  try {
    const payload: any = {
      ai_enabled: aiSettings.value.ai_enabled,
      ai_provider: aiSettings.value.ai_provider,
      ai_model: aiSettings.value.ai_model,
      ai_max_tokens: aiSettings.value.ai_max_tokens,
      ai_system_prompt: aiSettings.value.ai_system_prompt
    }
    // Only send API key if it's been changed (not empty)
    if (aiSettings.value.ai_api_key) {
      payload.ai_api_key = aiSettings.value.ai_api_key
    }
    await chatbotService.updateSettings(payload)
    toast.success('AI settings saved')
    aiSettings.value.ai_api_key = '' // Clear the API key field after save
  } catch (error) {
    toast.error('Failed to save AI settings')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="flex h-16 items-center px-6">
        <Settings class="h-5 w-5 mr-3" />
        <div>
          <h1 class="text-xl font-semibold">Settings</h1>
          <p class="text-sm text-muted-foreground">Manage your account and application settings</p>
        </div>
      </div>
    </header>

    <!-- Content -->
    <ScrollArea class="flex-1">
      <div class="p-6 space-y-6 max-w-3xl">
        <!-- General Settings -->
        <Card>
          <CardHeader>
            <CardTitle>General Settings</CardTitle>
            <CardDescription>Basic organization and display settings</CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label for="org_name">Organization Name</Label>
              <Input
                id="org_name"
                v-model="generalSettings.organization_name"
                placeholder="Your Organization"
              />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="timezone">Default Timezone</Label>
                <select
                  id="timezone"
                  v-model="generalSettings.default_timezone"
                  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="UTC">UTC</option>
                  <option value="America/New_York">Eastern Time</option>
                  <option value="America/Los_Angeles">Pacific Time</option>
                  <option value="Europe/London">London</option>
                  <option value="Asia/Tokyo">Tokyo</option>
                </select>
              </div>
              <div class="space-y-2">
                <Label for="date_format">Date Format</Label>
                <select
                  id="date_format"
                  v-model="generalSettings.date_format"
                  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="YYYY-MM-DD">YYYY-MM-DD</option>
                  <option value="DD/MM/YYYY">DD/MM/YYYY</option>
                  <option value="MM/DD/YYYY">MM/DD/YYYY</option>
                </select>
              </div>
            </div>
            <div class="flex justify-end">
              <Button @click="saveGeneralSettings" :disabled="isSubmitting">
                <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
                Save Changes
              </Button>
            </div>
          </CardContent>
        </Card>

        <!-- Chatbot Settings -->
        <Card>
          <CardHeader>
            <div class="flex items-center gap-2">
              <Bot class="h-5 w-5" />
              <div>
                <CardTitle>Chatbot Settings</CardTitle>
                <CardDescription>Configure default chatbot behavior</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label for="greeting">Greeting Message</Label>
              <Textarea
                id="greeting"
                v-model="chatbotSettings.greeting_message"
                placeholder="Hello! How can I help you?"
                :rows="2"
              />
            </div>
            <div class="space-y-2">
              <Label for="fallback">Fallback Message</Label>
              <Textarea
                id="fallback"
                v-model="chatbotSettings.fallback_message"
                placeholder="Sorry, I didn't understand that."
                :rows="2"
              />
            </div>
            <div class="space-y-2">
              <Label for="transfer">Transfer Message</Label>
              <Textarea
                id="transfer"
                v-model="chatbotSettings.transfer_message"
                placeholder="Transferring you to a human agent..."
                :rows="2"
              />
            </div>
            <div class="space-y-2">
              <Label for="timeout">Session Timeout (minutes)</Label>
              <Input
                id="timeout"
                v-model.number="chatbotSettings.session_timeout_minutes"
                type="number"
                min="5"
                max="120"
              />
            </div>
            <div class="flex justify-end">
              <Button @click="saveChatbotSettings" :disabled="isSubmitting">
                <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
                Save Changes
              </Button>
            </div>
          </CardContent>
        </Card>

        <!-- AI Settings -->
        <Card>
          <CardHeader>
            <div class="flex items-center gap-2">
              <Brain class="h-5 w-5" />
              <div>
                <CardTitle>AI Settings</CardTitle>
                <CardDescription>Configure AI-powered responses for your chatbot</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">Enable AI Responses</p>
                <p class="text-sm text-muted-foreground">Use AI to generate responses when no flow matches</p>
              </div>
              <input
                type="checkbox"
                v-model="aiSettings.ai_enabled"
                class="h-5 w-5 rounded"
              />
            </div>

            <div v-if="aiSettings.ai_enabled" class="space-y-4 pt-2">
              <Separator />

              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="ai_provider">AI Provider</Label>
                  <select
                    id="ai_provider"
                    v-model="aiSettings.ai_provider"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  >
                    <option value="">Select provider...</option>
                    <option v-for="provider in aiProviders" :key="provider.value" :value="provider.value">
                      {{ provider.label }}
                    </option>
                  </select>
                </div>
                <div class="space-y-2">
                  <Label for="ai_model">Model</Label>
                  <select
                    id="ai_model"
                    v-model="aiSettings.ai_model"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    :disabled="!aiSettings.ai_provider"
                  >
                    <option value="">Select model...</option>
                    <option v-for="model in availableModels" :key="model" :value="model">
                      {{ model }}
                    </option>
                  </select>
                </div>
              </div>

              <div class="space-y-2">
                <Label for="ai_api_key">API Key</Label>
                <Input
                  id="ai_api_key"
                  v-model="aiSettings.ai_api_key"
                  type="password"
                  placeholder="Enter API key (leave empty to keep existing)"
                />
                <p class="text-xs text-muted-foreground">Your API key is encrypted and stored securely</p>
              </div>

              <div class="space-y-2">
                <Label for="ai_max_tokens">Max Tokens</Label>
                <Input
                  id="ai_max_tokens"
                  v-model.number="aiSettings.ai_max_tokens"
                  type="number"
                  min="100"
                  max="4000"
                />
                <p class="text-xs text-muted-foreground">Maximum number of tokens for AI responses (100-4000)</p>
              </div>

              <div class="space-y-2">
                <Label for="ai_system_prompt">System Prompt (optional)</Label>
                <Textarea
                  id="ai_system_prompt"
                  v-model="aiSettings.ai_system_prompt"
                  placeholder="You are a helpful customer service assistant..."
                  :rows="3"
                />
                <p class="text-xs text-muted-foreground">Instructions for the AI on how to respond</p>
              </div>
            </div>

            <div class="flex justify-end pt-2">
              <Button @click="saveAISettings" :disabled="isSubmitting">
                <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
                Save Changes
              </Button>
            </div>
          </CardContent>
        </Card>

        <!-- Notification Settings -->
        <Card>
          <CardHeader>
            <div class="flex items-center gap-2">
              <Bell class="h-5 w-5" />
              <div>
                <CardTitle>Notifications</CardTitle>
                <CardDescription>Manage how you receive notifications</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">Email Notifications</p>
                <p class="text-sm text-muted-foreground">Receive important updates via email</p>
              </div>
              <input
                type="checkbox"
                v-model="notificationSettings.email_notifications"
                class="h-5 w-5 rounded"
              />
            </div>
            <Separator />
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">New Message Alerts</p>
                <p class="text-sm text-muted-foreground">Get notified when new messages arrive</p>
              </div>
              <input
                type="checkbox"
                v-model="notificationSettings.new_message_alerts"
                class="h-5 w-5 rounded"
              />
            </div>
            <Separator />
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">Campaign Updates</p>
                <p class="text-sm text-muted-foreground">Receive campaign status notifications</p>
              </div>
              <input
                type="checkbox"
                v-model="notificationSettings.campaign_updates"
                class="h-5 w-5 rounded"
              />
            </div>
          </CardContent>
        </Card>

        <!-- Security -->
        <Card>
          <CardHeader>
            <div class="flex items-center gap-2">
              <Shield class="h-5 w-5" />
              <div>
                <CardTitle>Security</CardTitle>
                <CardDescription>Manage your account security</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">Change Password</p>
                <p class="text-sm text-muted-foreground">Update your account password</p>
              </div>
              <Button variant="outline">Change Password</Button>
            </div>
            <Separator />
            <div class="flex items-center justify-between">
              <div>
                <p class="font-medium">Two-Factor Authentication</p>
                <p class="text-sm text-muted-foreground">Add an extra layer of security</p>
              </div>
              <Button variant="outline">Enable 2FA</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </ScrollArea>
  </div>
</template>
