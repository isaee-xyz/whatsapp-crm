<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
  LayoutDashboard,
  MessageSquare,
  Bot,
  FileText,
  Megaphone,
  Settings,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Users,
  Workflow,
  Sparkles,
  Key
} from 'lucide-vue-next'
import { getInitials } from '@/lib/utils'

const route = useRoute()
const authStore = useAuthStore()
const isCollapsed = ref(false)

const navigation = computed(() => [
  {
    name: 'Dashboard',
    path: '/',
    icon: LayoutDashboard,
    active: route.name === 'dashboard'
  },
  {
    name: 'Chat',
    path: '/chat',
    icon: MessageSquare,
    active: route.name === 'chat' || route.name === 'chat-conversation'
  },
  {
    name: 'Chatbot',
    path: '/chatbot',
    icon: Bot,
    active: route.path.startsWith('/chatbot'),
    children: [
      { name: 'Overview', path: '/chatbot', icon: Bot },
      { name: 'Keywords', path: '/chatbot/keywords', icon: Key },
      { name: 'Flows', path: '/chatbot/flows', icon: Workflow },
      { name: 'AI Contexts', path: '/chatbot/ai', icon: Sparkles }
    ]
  },
  {
    name: 'Templates',
    path: '/templates',
    icon: FileText,
    active: route.name === 'templates'
  },
  {
    name: 'Flows',
    path: '/flows',
    icon: Workflow,
    active: route.name === 'flows'
  },
  {
    name: 'Campaigns',
    path: '/campaigns',
    icon: Megaphone,
    active: route.name === 'campaigns'
  },
  {
    name: 'Settings',
    path: '/settings',
    icon: Settings,
    active: route.path.startsWith('/settings'),
    children: [
      { name: 'General', path: '/settings', icon: Settings },
      { name: 'Accounts', path: '/settings/accounts', icon: Users }
    ]
  }
])

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const handleLogout = async () => {
  await authStore.logout()
}
</script>

<template>
  <div class="flex h-screen bg-background">
    <!-- Sidebar -->
    <aside
      :class="[
        'flex flex-col border-r bg-card transition-all duration-300',
        isCollapsed ? 'w-16' : 'w-64'
      ]"
    >
      <!-- Logo -->
      <div class="flex h-16 items-center justify-between px-4 border-b">
        <RouterLink to="/" class="flex items-center gap-2">
          <div class="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
            <MessageSquare class="h-5 w-5 text-primary-foreground" />
          </div>
          <span
            v-if="!isCollapsed"
            class="font-bold text-lg text-foreground"
          >
            Whatomate
          </span>
        </RouterLink>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          @click="toggleSidebar"
        >
          <ChevronLeft v-if="!isCollapsed" class="h-4 w-4" />
          <ChevronRight v-else class="h-4 w-4" />
        </Button>
      </div>

      <!-- Navigation -->
      <ScrollArea class="flex-1 py-4">
        <nav class="space-y-1 px-2">
          <template v-for="item in navigation" :key="item.path">
            <RouterLink
              :to="item.path"
              :class="[
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                item.active
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                isCollapsed && 'justify-center px-2'
              ]"
            >
              <component :is="item.icon" class="h-5 w-5 shrink-0" />
              <span v-if="!isCollapsed">{{ item.name }}</span>
            </RouterLink>

            <!-- Submenu items -->
            <template v-if="item.children && item.active && !isCollapsed">
              <RouterLink
                v-for="child in item.children"
                :key="child.path"
                :to="child.path"
                :class="[
                  'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ml-4',
                  route.path === child.path
                    ? 'bg-accent text-accent-foreground'
                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                ]"
              >
                <component :is="child.icon" class="h-4 w-4 shrink-0" />
                <span>{{ child.name }}</span>
              </RouterLink>
            </template>
          </template>
        </nav>
      </ScrollArea>

      <!-- User section -->
      <div class="border-t p-4">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button
              variant="ghost"
              :class="[
                'w-full justify-start gap-3',
                isCollapsed && 'justify-center px-2'
              ]"
            >
              <Avatar class="h-8 w-8">
                <AvatarImage :src="undefined" />
                <AvatarFallback>
                  {{ getInitials(authStore.user?.full_name || 'U') }}
                </AvatarFallback>
              </Avatar>
              <div v-if="!isCollapsed" class="flex flex-col items-start text-left">
                <span class="text-sm font-medium truncate max-w-[140px]">
                  {{ authStore.user?.full_name }}
                </span>
                <span class="text-xs text-muted-foreground truncate max-w-[140px]">
                  {{ authStore.user?.email }}
                </span>
              </div>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-56">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleLogout">
              <LogOut class="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-hidden">
      <RouterView />
    </main>
  </div>
</template>
