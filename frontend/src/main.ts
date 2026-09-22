// src/main.ts
import './assets/main.css'

import { getInitialTheme, applyThemeToDOM } from '@/stores/theme'
applyThemeToDOM(getInitialTheme())

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createNotivue } from 'notivue'
import 'notivue/notification.css'
import 'notivue/animations.css'  // Make sure this is imported
import { globalLoader } from 'vue-global-loader'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(createNotivue({
    position: 'top-right',
    limit: 3,
    enqueue: true,
    pauseOnHover: true,
    notifications: {
        global: {
            duration: 5000,
        },
        error: {
            duration: 6000,
        },
        success: {
            duration: 4000,
        }
    }
}))
app.use(globalLoader)

const auth = useAuthStore()
auth.initAuth().finally(() => {
    app.mount('#app')
})




// That's Batch 4 done.

// Next up is Batch 5 — components. This is the biggest one:

// DeliveryItemFields.vue, DeliveryItemForm.vue, DeliveryItemDetail.vue, DeliveryItemRow.vue, DeliveryItemList.vue

// DeliveryForm.vue (inline item handling)

// DamageForm.vue, DamageDetail.vue, DamageRow.vue

// CustomerLedger.vue, CustomerLedgerPrintView.vue

// Say "continue" and I'll start with the Lot* components.