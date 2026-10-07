// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layouts/AppLayout.vue'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/customers'
        },
        {
            path: '/auth',
            name: 'auth',
            component: () => import('@/views/AuthView.vue'),
            meta: { layout: 'auth' }
        },
        {
            path: '/',
            component: AppLayout,
            children: [
                {
                    path: 'support',
                    name: 'support',
                    component: () => import('@/views/SupportView.vue'),
                },
                {
                    path: 'users',
                    name: 'users',
                    component: () => import('@/views/UsersView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'salaries',
                    name: 'salaries',
                    component: () => import('@/views/SalariesView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'expenses',
                    name: 'expenses',
                    component: () => import('@/views/ExpensesView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'incomes',
                    name: 'incomes',
                    component: () => import('@/views/IncomesView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'logs',
                    name: 'logs',
                    component: () => import('@/views/LogsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'customers',
                    name: 'customers',
                    component: () => import('@/views/CustomersView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'brokers',
                    name: 'brokers',
                    component: () => import('@/views/BrokersView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'majhis',
                    name: 'majhis',
                    component: () => import('@/views/MajhisView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'godowns',
                    name: 'godowns',
                    component: () => import('@/views/GodownsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'lots',
                    name: 'lots',
                    component: () => import('@/views/LotsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'deliveries',
                    name: 'deliveries',
                    component: () => import('@/views/DeliveriesView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'customer-delivery-bills',
                    name: 'customer-delivery-bills',
                    component: () => import('@/views/CustomerDeliveryBillsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'damages',
                    name: 'damages',
                    component: () => import('@/views/DamagesView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'customer-additional-bills',
                    name: 'customer-additional-bills',
                    component: () => import('@/views/CustomerAdditionalBillsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'store-adjustments',
                    name: 'store-adjustments',
                    component: () => import('@/views/StoreAdjustmentsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'store-transfers',
                    name: 'store-transfers',
                    component: () => import('@/views/StoreTransfersView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'lot-transfers',
                    name: 'lot-transfers',
                    component: () => import('@/views/LotTransfersView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'majhi-bills',
                    name: 'majhi-bills',
                    component: () => import('@/views/MajhiBillsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'godown-bills',
                    name: 'godown-bills',
                    component: () => import('@/views/GodownBillsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'customer-store-bills',
                    name: 'customer-store-bills',
                    component: () => import('@/views/CustomerStoreBillsView.vue'),
                    meta: { requiresAuth: true }
                },
                {
                    path: 'invoices',
                    name: 'invoices',
                    component: () => import('@/views/InvoicesView.vue'),
                    meta: { requiresAuth: true }
                },
            ]
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

    if (!authInitialized) {
        const token = localStorage.getItem('access_token')
        if (token) {
            await auth.initAuth()
        }
        authInitialized = true
    }

    const isAuthPage = to.path === '/auth'

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        return '/auth'
    }

    if (isAuthPage && auth.isAuthenticated) {
        return '/customers'
    }

    return true
})

export default router