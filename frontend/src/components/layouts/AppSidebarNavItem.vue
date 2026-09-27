<!-- src/components/layouts/AppSidebarNavItem.vue -->
<template>
    <router-link :to="to"
        class="group flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium w-full transition-colors duration-150"
        :class="active
            ? 'bg-(--color-blue)/10 text-(--color-blue)'
            : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)'"
        @click="$emit('click')">
        <span class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
            :class="active
                ? 'bg-(--color-blue)/15 border-(--color-blue)/20'
                : 'bg-(--color-muted-bg)/50 border-(--color-border)/30 group-hover:border-(--color-border)/50'">
            <component :is="iconComponent" class="w-3.5 h-3.5" :stroke-width="2" />
        </span>
        <span class="text-[13px] leading-tight">{{ label }}</span>
    </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Component } from 'vue'
import {
    LayoutDashboard,
    Users,
    Handshake,
    HardHat,
    Warehouse,
    Package,
    ClipboardList,
    AlertCircle,
    Truck,
    PackageCheck,
    Car,
    FileText,
    Receipt,
    Wallet,
    CreditCard,
    ScrollText,
    HelpCircle,
} from 'lucide-vue-next'

const props = defineProps<{
    to: string
    icon: string
    label: string
    active: boolean
}>()

defineEmits<{
    (e: 'click'): void
}>()

const iconRegistry: Record<string, Component> = {
    dashboard: LayoutDashboard,
    users: Users,
    brokers: Handshake,
    majhis: HardHat,
    godowns: Warehouse,
    lots: Package,
    stores: ClipboardList,
    damages: AlertCircle,
    deliveries: Truck,
    'delivery-items': PackageCheck,
    transports: Truck,
    vehicles: Car,
    'godown-store-bills': FileText,
    'customer-storage-bills': Receipt,
    'customer-lot-bills': Receipt,
    'majhi-lot-bills': HardHat,
    'customer-delivery-bills': Truck,
    'majhi-loading-bills': HardHat,
    'broker-vehicle-bills': Car,
    'customer-transport-bills': Truck,
    'customer-additional-bills': FileText,
    invoices: FileText,
    salaries: Wallet,
    expenses: CreditCard,
    logs: ScrollText,
    support: HelpCircle,
}

const iconComponent = computed<Component>(() => {
    return iconRegistry[props.icon] ?? LayoutDashboard
})
</script>