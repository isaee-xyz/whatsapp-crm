<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { chatbotService } from '@/services/api'
import { toast } from 'vue-sonner'
import { Plus, Pencil, Trash2, Sparkles, ArrowLeft, FileText, Globe } from 'lucide-vue-next'

interface ApiConfig {
  url: string
  method: string
  headers: Record<string, string>
  body: string
  response_path: string
}

interface AIContext {
  id: string
  name: string
  context_type: string
  trigger_keywords: string[]
  static_content: string
  api_config: ApiConfig
  priority: number
  enabled: boolean
  created_at: string
}

const contexts = ref<AIContext[]>([])
const isLoading = ref(true)
const isDialogOpen = ref(false)
const isSubmitting = ref(false)
const editingContext = ref<AIContext | null>(null)

const formData = ref({
  name: '',
  context_type: 'static',
  trigger_keywords: '',
  static_content: '',
  api_url: '',
  api_method: 'GET',
  api_headers: '',
  api_response_path: '',
  priority: 10,
  enabled: true
})

// Helper to display variable placeholders without Vue parsing issues
const variableExample = (name: string) => `{{${name}}}`

onMounted(async () => {
  await fetchContexts()
})

async function fetchContexts() {
  isLoading.value = true
  try {
    const response = await chatbotService.listAIContexts()
    // API response is wrapped in { status: "success", data: { contexts: [...] } }
    const data = response.data.data || response.data
    contexts.value = data.contexts || []
  } catch (error) {
    console.error('Failed to load AI contexts:', error)
    contexts.value = []
  } finally {
    isLoading.value = false
  }
}

function openCreateDialog() {
  editingContext.value = null
  formData.value = {
    name: '',
    context_type: 'static',
    trigger_keywords: '',
    static_content: '',
    api_url: '',
    api_method: 'GET',
    api_headers: '',
    api_response_path: '',
    priority: 10,
    enabled: true
  }
  isDialogOpen.value = true
}

function openEditDialog(context: AIContext) {
  editingContext.value = context
  const apiConfig = context.api_config || {} as ApiConfig
  formData.value = {
    name: context.name,
    context_type: context.context_type || 'static',
    trigger_keywords: (context.trigger_keywords || []).join(', '),
    static_content: context.static_content || '',
    api_url: apiConfig.url || '',
    api_method: apiConfig.method || 'GET',
    api_headers: apiConfig.headers ? JSON.stringify(apiConfig.headers, null, 2) : '',
    api_response_path: apiConfig.response_path || '',
    priority: context.priority || 10,
    enabled: context.enabled
  }
  isDialogOpen.value = true
}

async function saveContext() {
  if (!formData.value.name.trim()) {
    toast.error('Please enter a name')
    return
  }

  if (formData.value.context_type === 'api' && !formData.value.api_url.trim()) {
    toast.error('Please enter an API URL')
    return
  }

  isSubmitting.value = true
  try {
    // Parse headers JSON if provided
    let headers = {}
    if (formData.value.api_headers.trim()) {
      try {
        headers = JSON.parse(formData.value.api_headers)
      } catch (e) {
        toast.error('Invalid JSON format for headers')
        isSubmitting.value = false
        return
      }
    }

    const data: any = {
      name: formData.value.name,
      context_type: formData.value.context_type,
      trigger_keywords: formData.value.trigger_keywords.split(',').map(k => k.trim()).filter(Boolean),
      static_content: formData.value.static_content,
      api_config: formData.value.context_type === 'api' ? {
        url: formData.value.api_url,
        method: formData.value.api_method,
        headers: headers,
        response_path: formData.value.api_response_path
      } : null,
      priority: formData.value.priority,
      enabled: formData.value.enabled
    }

    if (editingContext.value) {
      await chatbotService.updateAIContext(editingContext.value.id, data)
      toast.success('AI context updated')
    } else {
      await chatbotService.createAIContext(data)
      toast.success('AI context created')
    }

    isDialogOpen.value = false
    await fetchContexts()
  } catch (error) {
    toast.error('Failed to save AI context')
  } finally {
    isSubmitting.value = false
  }
}

async function deleteContext(context: AIContext) {
  if (!confirm(`Are you sure you want to delete "${context.name}"?`)) return

  try {
    await chatbotService.deleteAIContext(context.id)
    toast.success('AI context deleted')
    await fetchContexts()
  } catch (error) {
    toast.error('Failed to delete AI context')
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="flex h-16 items-center justify-between px-6">
        <div class="flex items-center gap-4">
          <RouterLink to="/chatbot">
            <Button variant="ghost" size="icon">
              <ArrowLeft class="h-5 w-5" />
            </Button>
          </RouterLink>
          <div>
            <h1 class="text-xl font-semibold">AI Contexts</h1>
            <p class="text-sm text-muted-foreground">Knowledge bases for AI-powered responses</p>
          </div>
        </div>
        <Dialog v-model:open="isDialogOpen">
          <DialogTrigger as-child>
            <Button @click="openCreateDialog">
              <Plus class="h-4 w-4 mr-2" />
              Add Context
            </Button>
          </DialogTrigger>
          <DialogContent class="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{{ editingContext ? 'Edit' : 'Create' }} AI Context</DialogTitle>
              <DialogDescription>
                Add knowledge context that the AI can use when responding to messages.
              </DialogDescription>
            </DialogHeader>
            <div class="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="name">Name *</Label>
                  <Input
                    id="name"
                    v-model="formData.name"
                    placeholder="Product FAQ"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="context_type">Type</Label>
                  <select
                    id="context_type"
                    v-model="formData.context_type"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  >
                    <option value="static">Static Content</option>
                    <option value="api">API Fetch</option>
                  </select>
                </div>
              </div>

              <div class="space-y-2">
                <Label for="trigger_keywords">Trigger Keywords (comma-separated, optional)</Label>
                <Input
                  id="trigger_keywords"
                  v-model="formData.trigger_keywords"
                  placeholder="faq, help, info"
                />
                <p class="text-xs text-muted-foreground">
                  Leave empty to always include this context, or specify keywords to include only when mentioned.
                </p>
              </div>

              <!-- Content/Prompt Field - always shown -->
              <div class="space-y-2">
                <Label for="static_content">Content / Prompt</Label>
                <Textarea
                  id="static_content"
                  v-model="formData.static_content"
                  placeholder="Enter knowledge content or prompt for the AI..."
                  :rows="6"
                />
                <p class="text-xs text-muted-foreground">
                  This content will be provided to the AI as context for generating responses.
                </p>
              </div>

              <!-- API Configuration Fields - shown only for API type -->
              <div v-if="formData.context_type === 'api'" class="space-y-4 border-t pt-4">
                <p class="text-sm font-medium">API Configuration</p>
                <p class="text-xs text-muted-foreground">
                  Data fetched from this API will be combined with the content above.
                </p>

                <div class="grid grid-cols-4 gap-4">
                  <div class="col-span-1 space-y-2">
                    <Label for="api_method">Method</Label>
                    <select
                      id="api_method"
                      v-model="formData.api_method"
                      class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    >
                      <option value="GET">GET</option>
                      <option value="POST">POST</option>
                    </select>
                  </div>
                  <div class="col-span-3 space-y-2">
                    <Label for="api_url">API URL *</Label>
                    <Input
                      id="api_url"
                      v-model="formData.api_url"
                      placeholder="https://api.example.com/context"
                    />
                  </div>
                </div>
                <p class="text-xs text-muted-foreground">
                  Variables: <code class="bg-muted px-1 rounded">{{ variableExample('phone_number') }}</code>, <code class="bg-muted px-1 rounded">{{ variableExample('user_message') }}</code>
                </p>

                <div class="space-y-2">
                  <Label for="api_headers">Headers (JSON, optional)</Label>
                  <Textarea
                    id="api_headers"
                    v-model="formData.api_headers"
                    placeholder='{"Authorization": "Bearer xxx"}'
                    :rows="2"
                  />
                </div>

                <div class="space-y-2">
                  <Label for="api_response_path">Response Path (optional)</Label>
                  <Input
                    id="api_response_path"
                    v-model="formData.api_response_path"
                    placeholder="data.context"
                  />
                  <p class="text-xs text-muted-foreground">
                    Dot-notation path to extract from JSON response.
                  </p>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="priority">Priority</Label>
                  <Input
                    id="priority"
                    v-model.number="formData.priority"
                    type="number"
                    min="1"
                    max="100"
                  />
                  <p class="text-xs text-muted-foreground">Higher priority contexts are used first</p>
                </div>
                <div class="flex items-center gap-2 pt-8">
                  <input
                    id="enabled"
                    v-model="formData.enabled"
                    type="checkbox"
                    class="h-4 w-4 rounded border-gray-300"
                  />
                  <Label for="enabled">Enabled</Label>
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" @click="isDialogOpen = false">Cancel</Button>
              <Button @click="saveContext" :disabled="isSubmitting">
                {{ editingContext ? 'Update' : 'Create' }}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </header>

    <!-- Contexts List -->
    <ScrollArea class="flex-1">
      <div class="p-6 grid gap-4 md:grid-cols-2">
        <Card v-for="context in contexts" :key="context.id">
          <CardHeader>
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-3">
                <div
                  class="h-10 w-10 rounded-lg flex items-center justify-center"
                  :class="context.context_type === 'api' ? 'bg-blue-100 dark:bg-blue-900' : 'bg-orange-100 dark:bg-orange-900'"
                >
                  <Globe v-if="context.context_type === 'api'" class="h-5 w-5 text-blue-600 dark:text-blue-400" />
                  <FileText v-else class="h-5 w-5 text-orange-600 dark:text-orange-400" />
                </div>
                <div>
                  <CardTitle class="text-base">{{ context.name }}</CardTitle>
                  <CardDescription>{{ context.context_type === 'api' ? 'API Fetch' : 'Static Content' }}</CardDescription>
                </div>
              </div>
              <Badge :variant="context.enabled ? 'default' : 'secondary'">
                {{ context.enabled ? 'Active' : 'Inactive' }}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div class="flex flex-wrap gap-1 mb-3" v-if="context.trigger_keywords?.length">
              <Badge v-for="kw in context.trigger_keywords" :key="kw" variant="outline" class="text-xs">
                {{ kw }}
              </Badge>
            </div>
            <div class="flex flex-wrap gap-2 mb-3">
              <Badge variant="secondary">Priority: {{ context.priority }}</Badge>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex gap-2">
                <Button variant="ghost" size="icon" @click="openEditDialog(context)">
                  <Pencil class="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" @click="deleteContext(context)">
                  <Trash2 class="h-4 w-4 text-destructive" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card v-if="contexts.length === 0" class="col-span-full">
          <CardContent class="py-12 text-center text-muted-foreground">
            <Sparkles class="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p class="text-lg font-medium">No AI contexts yet</p>
            <p class="text-sm mb-4">Create knowledge contexts that the AI can use to answer questions.</p>
            <Button @click="openCreateDialog">
              <Plus class="h-4 w-4 mr-2" />
              Create Context
            </Button>
          </CardContent>
        </Card>
      </div>
    </ScrollArea>
  </div>
</template>
