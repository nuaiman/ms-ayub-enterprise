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

// =====================================================
// Prevent mouse wheel from mutating input values
// =====================================================
// Browser default: hovering/focusing an <input type="number">,
// <input type="date">, <input type="month">, <input type="time">,
// or <input type="datetime-local"> and scrolling the wheel
// increments/decrements the value. Users hate this. We suppress
// it by preventing the wheel event's default when the pointer
// is over one of these inputs.
//
// Using capture phase + passive:false so preventDefault() sticks.
const WHEEL_BLOCKED_INPUT_TYPES = new Set([
    'number',
    'date',
    'month',
    'time',
    'week',
    'datetime-local',
])

document.addEventListener(
    'wheel',
    (event) => {
        const target = event.target as HTMLElement | null
        if (!target) return
        if (target.tagName !== 'INPUT') return

        const input = target as HTMLInputElement
        const type = (input.getAttribute('type') || '').toLowerCase()
        if (!WHEEL_BLOCKED_INPUT_TYPES.has(type)) return

        // Only block if the input is actually focused — otherwise the
        // user is just scrolling past it and browser wouldn't step
        // the value anyway. This keeps tabbing/keyboard flow intact.
        if (document.activeElement !== input) return

        event.preventDefault()
    },
    { passive: false, capture: true },
)

// =====================================================
// App bootstrap
// =====================================================

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