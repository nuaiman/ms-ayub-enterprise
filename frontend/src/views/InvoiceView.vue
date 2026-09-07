<!-- src/views/InvoicesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Invoice Generator -->
            <div class="flex-1 min-h-0">
                <div class="max-w-6xl mx-auto">
                    <InvoiceGenerator :type="invoiceType" />
                </div>
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layouts/AppLayout.vue'
import InvoiceGenerator from '@/components/features/invoices/InvoiceGenerator.vue'
import type { InvoiceType } from '@/types/invoice'

const route = useRoute()
const invoiceType = ref<InvoiceType>('godown')

onMounted(() => {
    const type = route.query.type as string
    if (type === 'transport' || type === 'godown') {
        invoiceType.value = type
    }
})
</script>