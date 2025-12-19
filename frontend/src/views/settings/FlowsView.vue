<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { flowsService } from '@/services/api'
import { toast } from 'vue-sonner'
import { Plus, Pencil, Trash2, Workflow, Play, ExternalLink } from 'lucide-vue-next'

interface WhatsAppFlow {
  id: string
  name: string
  flow_id: string
  status: 'DRAFT' | 'PUBLISHED' | 'DEPRECATED'
  categories: string[]
  preview_url?: string
  created_at: string
}

const flows = ref<WhatsAppFlow[]>([])
const isLoading = ref(true)

onMounted(async () => {
  await fetchFlows()
})

async function fetchFlows() {
  isLoading.value = true
  try {
    const response = await flowsService.list()
    flows.value = response.data.flows || []
  } catch (error) {
    // Demo data
    flows.value = [
      {
        id: '1',
        name: 'Customer Feedback Survey',
        flow_id: 'FLOW_123456',
        status: 'PUBLISHED',
        categories: ['SURVEY'],
        preview_url: 'https://wa.me/...',
        created_at: new Date().toISOString()
      },
      {
        id: '2',
        name: 'Product Catalog',
        flow_id: 'FLOW_789012',
        status: 'DRAFT',
        categories: ['CATALOG'],
        created_at: new Date().toISOString()
      },
      {
        id: '3',
        name: 'Appointment Booking',
        flow_id: 'FLOW_345678',
        status: 'PUBLISHED',
        categories: ['BOOKING'],
        preview_url: 'https://wa.me/...',
        created_at: new Date().toISOString()
      }
    ]
  } finally {
    isLoading.value = false
  }
}

async function publishFlow(flow: WhatsAppFlow) {
  try {
    await flowsService.publish(flow.id)
    toast.success('Flow published successfully')
    await fetchFlows()
  } catch (error) {
    toast.error('Failed to publish flow')
  }
}

async function deleteFlow(flow: WhatsAppFlow) {
  if (!confirm(`Are you sure you want to delete "${flow.name}"?`)) return

  try {
    await flowsService.delete(flow.id)
    toast.success('Flow deleted')
    await fetchFlows()
  } catch (error) {
    toast.error('Failed to delete flow')
  }
}

function getStatusVariant(status: string) {
  switch (status) {
    case 'PUBLISHED':
      return 'default'
    case 'DRAFT':
      return 'secondary'
    case 'DEPRECATED':
      return 'destructive'
    default:
      return 'outline'
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="flex h-16 items-center justify-between px-6">
        <div>
          <h1 class="text-xl font-semibold">WhatsApp Flows</h1>
          <p class="text-sm text-muted-foreground">Create interactive flows for your customers</p>
        </div>
        <Button>
          <Plus class="h-4 w-4 mr-2" />
          Create Flow
        </Button>
      </div>
    </header>

    <!-- Flows List -->
    <ScrollArea class="flex-1">
      <div class="p-6 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card v-for="flow in flows" :key="flow.id">
          <CardHeader>
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-3">
                <div class="h-10 w-10 rounded-lg bg-indigo-100 dark:bg-indigo-900 flex items-center justify-center">
                  <Workflow class="h-5 w-5 text-indigo-600 dark:text-indigo-400" />
                </div>
                <div>
                  <CardTitle class="text-base">{{ flow.name }}</CardTitle>
                  <p class="text-xs text-muted-foreground font-mono">{{ flow.flow_id }}</p>
                </div>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div class="flex flex-wrap gap-2 mb-3">
              <Badge :variant="getStatusVariant(flow.status)">
                {{ flow.status }}
              </Badge>
              <Badge v-for="category in flow.categories" :key="category" variant="outline">
                {{ category }}
              </Badge>
            </div>
          </CardContent>
          <div class="px-6 pb-4 flex items-center justify-between border-t pt-4">
            <div class="flex gap-2">
              <Button variant="ghost" size="icon">
                <Pencil class="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" @click="deleteFlow(flow)">
                <Trash2 class="h-4 w-4 text-destructive" />
              </Button>
            </div>
            <div class="flex gap-2">
              <Button
                v-if="flow.preview_url"
                variant="outline"
                size="sm"
                as="a"
                :href="flow.preview_url"
                target="_blank"
              >
                <ExternalLink class="h-4 w-4 mr-1" />
                Preview
              </Button>
              <Button
                v-if="flow.status === 'DRAFT'"
                size="sm"
                @click="publishFlow(flow)"
              >
                <Play class="h-4 w-4 mr-1" />
                Publish
              </Button>
            </div>
          </div>
        </Card>

        <Card v-if="flows.length === 0 && !isLoading" class="col-span-full">
          <CardContent class="py-12 text-center text-muted-foreground">
            <Workflow class="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p class="text-lg font-medium">No WhatsApp Flows yet</p>
            <p class="text-sm mb-4">Create interactive flows to engage your customers.</p>
            <Button>
              <Plus class="h-4 w-4 mr-2" />
              Create Flow
            </Button>
          </CardContent>
        </Card>
      </div>
    </ScrollArea>
  </div>
</template>
