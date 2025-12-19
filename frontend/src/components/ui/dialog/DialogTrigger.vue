<script setup lang="ts">
import { DialogTrigger as DialogTriggerPrimitive } from 'radix-vue'
</script>

<template>
  <DialogTriggerPrimitive as-child>
    <slot />
  </DialogTriggerPrimitive>
</template>
