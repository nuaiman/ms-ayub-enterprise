// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/auth'
        },
        {
            path: '/auth',
            name: 'auth',
            component: () => import('@/views/AuthView.vue'),
            meta: { layout: 'auth' } // No header layout
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: () => import('@/views/DashboardView.vue'),
            meta: { requiresAuth: true, layout: 'app' }
        },
        {
            path: '/support',
            name: 'support',
            component: () => import('@/views/SupportView.vue'),
            meta: { layout: 'app' }
        },
        {
            path: '/users',
            name: 'users',
            component: () => import('@/views/UsersView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/salaries',
            name: 'salaries',
            component: () => import('@/views/SalariesView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/expenses',
            name: 'expenses',
            component: () => import('@/views/ExpensesView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/logs',
            name: 'logs',
            component: () => import('@/views/LogsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/customers',
            name: 'customers',
            component: () => import('@/views/CustomersView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/brokers',
            name: 'brokers',
            component: () => import('@/views/BrokersView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/majhis',
            name: 'majhis',
            component: () => import('@/views/MajhisView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/godowns',
            name: 'godowns',
            component: () => import('@/views/GodownsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/rents',
            name: 'rents',
            component: () => import('@/views/RentsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/lots',
            name: 'lots',
            component: () => import('@/views/LotsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/stores',
            name: 'stores',
            component: () => import('@/views/StoresView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/godown-store-bills',
            name: 'godown-store-bills',
            component: () => import('@/views/GodownStoreBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/customer-storage-bills',
            name: 'customer-storage-bills',
            component: () => import('@/views/CustomerStorageBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/customer-lot-bills',
            name: 'customer-lot-bills',
            component: () => import('@/views/CustomerLotBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/majhi-lot-bills',
            name: 'majhi-lot-bills',
            component: () => import('@/views/MajhiLotBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/customer-delivery-bills',
            name: 'customer-delivery-bills',
            component: () => import('@/views/CustomerDeliveryBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/majhi-loading-bills',
            name: 'majhi-loading-bills',
            component: () => import('@/views/MajhiLoadingBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/damages',
            name: 'damages',
            component: () => import('@/views/DamagesView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/deliveries',
            name: 'deliveries',
            component: () => import('@/views/DeliveriesView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/delivery-items',
            name: 'delivery-items',
            component: () => import('@/views/DeliveryItemsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/transports',
            name: 'transports',
            component: () => import('@/views/TransportsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/vehicles',
            name: 'vehicles',
            component: () => import('@/views/VehiclesView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/broker-vehicle-bills',
            name: 'broker-vehicle-bills',
            component: () => import('@/views/BrokerVehicleBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/customer-transport-bills',
            name: 'customer-transport-bills',
            component: () => import('@/views/CustomerTransportBillsView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/invoices',
            name: 'invoices',
            component: () => import('@/views/InvoiceView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'not-found',
            component: () => import('@/views/NotFoundView.vue'),
        }
    ]
})

let authInitialized = false

router.beforeEach(async (to) => {
    const auth = useAuthStore()

    // ✅ If auth hasn't been initialized yet, wait for it
    if (!authInitialized) {
        // Check if we have a token in localStorage
        const token = localStorage.getItem('access_token')
        if (token) {
            // Try to initialize auth
            await auth.initAuth()
        }
        authInitialized = true
    }

    const isAuthPage = to.path === '/auth'

    // If route requires auth and user is not authenticated
    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        return '/auth'
    }

    // If user is authenticated and trying to go to auth page
    if (isAuthPage && auth.isAuthenticated) {
        return '/dashboard'
    }

    return true
})

export default router