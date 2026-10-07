<!-- src/components/features/deliveries/DeliveryItemAddForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Item Details
            </h3>

            <DeliveryItemFields v-model:store-id="form.store_id" v-model:majhi-id="form.majhi_id"
                v-model:vehicle-number="form.vehicle_number" v-model:driver-number="form.driver_number"
                v-model:quantity="form.quantity" v-model:weight="form.weight" v-model:cd-bill-type="form.cd_bill_type"
                v-model:cd-bill-rate="form.cd_bill_rate" v-model:majhi-bill-type="form.majhi_bill_type"
                v-model:majhi-bill-rate="form.majhi_bill_rate" :store-options="storeOptions"
                :majhi-options="majhiOptions" :disabled="submitting" />
        </div>

        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Creating...
                </span>
                <span v-else>Create Item</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { CustomerDeliveryBillType } from '@/types/customerDeliveryBill'
import type { MajhiBillType } from '@/types/majhiBill'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useStoresStore } from '@/stores/stores'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { push } from 'notivue'
import DeliveryItemFields from './DeliveryItemFields.vue'

const props = defineProps<{
    deliveryId: number
    customerId: number | null
}>()

const emit = defineEmits<{
    'item-created': []
    'cancel': []
}>()

const deliveriesStore = useDeliveriesStore()
const storesStore = useStoresStore()
const majhisStore = useMajhisStore()
const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const majhiBillsStore = useMajhiBillsStore()

const submitting = ref(false)
const storeOptions = computed(() => storesStore.stores)
const majhiOptions = computed(() => majhisStore.majhis)

const form = ref({
    store_id: null as number | null,
    majhi_id: null as number | null,
    vehicle_number: '',
    driver_number: '',
    quantity: 0,
    weight: 0,
    cd_bill_type: 'quantity' as CustomerDeliveryBillType,
    cd_bill_rate: 0,
    majhi_bill_type: 'quantity' as MajhiBillType,
    majhi_bill_rate: 0,
})

const submit = async () => {
    if (!form.value.store_id) {
        push.error('Store is required')
        return
    }
    if (!form.value.majhi_id) {
        push.error('Majhi is required')
        return
    }
    if (form.value.quantity < 0 || form.value.weight < 0) {
        push.error('Quantity/weight cannot be negative')
        return
    }
    if (form.value.quantity === 0 && form.value.weight === 0) {
        push.error('Quantity or weight must be greater than 0')
        return
    }
    if (!form.value.cd_bill_rate || form.value.cd_bill_rate <= 0) {
        push.error('Customer delivery bill rate must be greater than 0')
        return
    }
    if (!form.value.majhi_bill_rate || form.value.majhi_bill_rate <= 0) {
        push.error('Majhi bill rate must be greater than 0')
        return
    }

    submitting.value = true
    try {
        const created = await deliveriesStore.createDeliveryItem(props.deliveryId, {
            store_id: form.value.store_id,
            majhi_id: form.value.majhi_id,
            vehicle_number: form.value.vehicle_number.trim() || null,
            driver_number: form.value.driver_number.trim() || null,
            quantity: form.value.quantity,
            weight: form.value.weight,
        })
        if (!created) return

        if (props.customerId) {
            await customerDeliveryBillsStore.createCustomerDeliveryBill({
                customer_id: props.customerId,
                delivery_item_id: created.id,
                bill_type: form.value.cd_bill_type,
                rate: form.value.cd_bill_rate,
            })
        }

        await majhiBillsStore.createMajhiBill({
            majhi_id: form.value.majhi_id,
            delivery_item_id: created.id,
            bill_type: form.value.majhi_bill_type,
            rate: form.value.majhi_bill_rate,
        })

        emit('item-created')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
})
</script>