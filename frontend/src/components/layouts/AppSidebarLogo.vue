<!-- src/components/layouts/AppSidebarLogo.vue -->
<template>
    <div class="h-14 flex items-center px-4 border-b border-(--color-border)/40 shrink-0">
        <div class="flex items-center gap-3 w-full">
            <div
                class="w-8 h-8 shrink-0 overflow-hidden rounded-lg bg-(--color-muted-bg) border border-(--color-border)/50 p-0.5 flex items-center justify-center">
                <img src="@/assets/logo.png" alt="MS Ayub Enterprise"
                    class="w-full h-full block object-cover rounded-md" />
            </div>
            <div class="flex flex-col">
                <span class="text-sm font-semibold tracking-tight text-(--color-text-primary) leading-tight">
                    MS Ayub
                </span>
                <span class="text-[10px] font-medium text-(--color-text-secondary)/60 tracking-[0.15em] uppercase">
                    Enterprise
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
</script>