<!-- src/components/features/majhiBills/MajhiBillFields.vue -->
<template>
    <div>
        <!-- Derived mode: 3-column inline grid (Majhi | Bill Type | Rate) -->
        <div v-if="derived" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <!-- Majhi picker -->
            <div v-if="!locked">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Majhi <span class="text-(--color-red)">*</span>
                </label>
                <select :value="majhiId ?? ''"
                    @change="$emit('update:majhiId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a majhi</option>
                    <option v-for="m in majhiOptions" :key="m.id" :value="m.id">
                        {{ m.name }}
                    </option>
                </select>
            </div>
            <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedMajhiLabel }}</p>
            </div>

            <!-- Bill Type -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Bill Type <span class="text-(--color-red)">*</span>
                </label>
                <select :value="billType"
                    @change="$emit('update:billType', ($event.target as HTMLSelectElement).value as MajhiBillType)"
                    :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                    <option value="job">Job (Fixed)</option>
                </select>
            </div>

            <!-- Rate -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Rate <span class="text-(--color-red)">*</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="rate"
                        @input="$emit('update:rate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Default mode: two-column grid -->
        <div v-else class="space-y-4">
            <!-- Majhi + Store row -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Majhi picker (create) -->
                <div v-if="!locked">
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Majhi <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="majhiId ?? ''"
                        @change="$emit('update:majhiId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                        :disabled="disabled"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a majhi</option>
                        <option v-for="m in majhiOptions" :key="m.id" :value="m.id">
                            {{ m.name }}
                        </option>
                    </select>
                </div>
                <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                    <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedMajhiLabel }}</p>
                </div>

                <!-- Store picker (create) — hidden when store is implicit (derived, but that branch is above) -->
                <div v-if="!locked">
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="storeId ?? ''"
                        @change="$emit('update:storeId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                        :disabled="disabled"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a store</option>
                        <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                            {{ storeLabel(store) }}
                        </option>
                    </select>
                </div>
                <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                    <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
                </div>
            </div>

            <!-- Bill Type + Rate row -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Bill Type -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Bill Type <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="billType"
                        @change="$emit('update:billType', ($event.target as HTMLSelectElement).value as MajhiBillType)"
                        :disabled="disabled" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="quantity">Quantity</option>
                        <option value="weight">Weight</option>
                        <option value="job">Job (Fixed)</option>
                    </select>
                </div>

                <!-- Rate -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Rate <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input :value="rate"
                            @input="$emit('update:rate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                            type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import type { MajhiBillType } from '@/types/majhiBill'
import { useLotsStore } from '@/stores/lots'

withDefaults(defineProps<{
    majhiId?: number | null
    storeId?: number | null
    billType: MajhiBillType
    rate: number
    majhiOptions?: Majhi[]
    storeOptions?: Store[]
    disabled?: boolean
    locked?: boolean
    lockedMajhiLabel?: string
    lockedStoreLabel?: string
    derived?: boolean
}>(), {
    majhiId: null,
    storeId: null,
    majhiOptions: () => [],
    storeOptions: () => [],
    disabled: false,
    locked: false,
    lockedMajhiLabel: '',
    lockedStoreLabel: '',
    derived: false,
})

defineEmits<{
    (e: 'update:majhiId', value: number | null): void
    (e: 'update:storeId', value: number | null): void
    (e: 'update:billType', value: MajhiBillType): void
    (e: 'update:rate', value: number): void
}>()

const lotsStore = useLotsStore()

const storeLabel = (store: Store): string => {
    const lot = lotsStore.getLotById(store.lot_id)
    const lotName = lot ? lot.product_name : '—'
    const lotNum = lot ? lot.lot_number : '—'
    return `Store #${store.id} — ${lotName} (Lot ${lotNum})`
}
</script>