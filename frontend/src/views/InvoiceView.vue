<!-- src/views/InvoicesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <div class="flex-1 min-h-0">
                <div class="max-w-6xl mx-auto">
                    <InvoiceGenerator />
                </div>
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import AppLayout from '@/components/layouts/AppLayout.vue'
import InvoiceGenerator from '@/components/features/invoices/InvoiceGenerator.vue'
</script>