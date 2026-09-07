<!-- src/components/layouts/AppSidebarNavItem.vue -->
<template>
    <div class="relative w-full">
        <router-link :to="to" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium w-full" :class="[
            active
                ? 'bg-(--color-blue)/10 text-(--color-blue)'
                : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)',
            collapsed ? 'justify-center px-0' : ''
        ]" @click="$emit('click')">
            <span
                class="w-7 h-7 rounded-lg bg-(--color-muted-bg)/50 flex items-center justify-center text-sm shrink-0 border border-(--color-border)/30">
                {{ icon }}
            </span>
            <!-- ✅ Show text when NOT collapsed -->
            <span v-if="!collapsed" class="truncate text-[13px]">{{ label }}</span>
        </router-link>

        <!-- Tooltip for collapsed mode (desktop only) -->
        <div v-if="collapsed"
            class="hidden lg:block absolute left-full top-1/2 -translate-y-1/2 ml-3 px-2.5 py-1 rounded-lg bg-(--color-surface) border border-(--color-border) text-[11px] font-medium text-(--color-text-primary) shadow-lg whitespace-nowrap z-50">
            {{ label }}
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    to: string
    icon: string
    label: string
    active: boolean
    collapsed: boolean
}>()

defineEmits<{
    (e: 'click'): void
}>()
</script>