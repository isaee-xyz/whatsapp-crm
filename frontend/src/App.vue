<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Toaster } from 'vue-sonner'
</script>

<template>
  <div class="min-h-screen bg-background">
    <RouterView />
    <Toaster position="top-right" richColors />
  </div>
</template>
