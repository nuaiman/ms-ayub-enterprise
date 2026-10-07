// src/stores/customerStoreBills.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
    CustomerStoreBill,
    CreateCustomerStoreBillPayload,
    UpdateCustomerStoreBillPayload,
    CreateBillPaymentPayload,
    CustomerStoreBillSortField,
    SortDirection,
} from "@/types/customerStoreBill";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useCustomerStoreBillsStore = defineStore(
    "customerStoreBills",
    () => {
        const { displayLoader, destroyLoader } = useGlobalLoader();

        const bills = ref<CustomerStoreBill[]>([]);
        const searchQuery = ref("");
        const sortField = ref<CustomerStoreBillSortField>("created_at");
        const sortDirection = ref<SortDirection>("desc");

        const filteredBills = computed(() => {
            let result = [...bills.value];

            if (searchQuery.value) {
                const query = searchQuery.value.toLowerCase();
                result = result.filter(
                    (b) =>
                        b.bill_type.toLowerCase().includes(query) ||
                        String(b.rate).includes(query) ||
                        String(b.total_paid).includes(query)
                );
            }

            result.sort((a, b) => {
                let comparison = 0;
                switch (sortField.value) {
                    case "store_id":
                        comparison = a.store_id - b.store_id;
                        break;
                    case "customer_id":
                        comparison = a.customer_id - b.customer_id;
                        break;
                    case "bill_type":
                        comparison = a.bill_type.localeCompare(b.bill_type);
                        break;
                    case "rate":
                        comparison = a.rate - b.rate;
                        break;
                    case "total_paid":
                        comparison = a.total_paid - b.total_paid;
                        break;
                    case "created_at":
                        comparison =
                            new Date(a.created_at).getTime() -
                            new Date(b.created_at).getTime();
                        break;
                    default:
                        comparison = 0;
                }
                return sortDirection.value === "desc" ? -comparison : comparison;
            });

            return result;
        });

        const totalBills = computed(() => bills.value.length);
        const totalBilled = computed(() =>
            bills.value.reduce((sum, b) => sum + b.total_amount, 0)
        );
        const totalPaid = computed(() =>
            bills.value.reduce((sum, b) => sum + b.total_paid, 0)
        );

        const fetchCustomerStoreBills = async (params?: {
            store_id?: number;
            customer_id?: number;
            month_year?: string;
        }) => {
            displayLoader();
            try {
                const res = await api.get<ApiResponse<CustomerStoreBill[]>>(
                    "/customer-store-bills",
                    { params }
                );
                if (!res.data.success) {
                    push.error(res.data.message);
                    return [];
                }
                bills.value = res.data.data;
                return bills.value;
            } catch (error) {
                const err = error as AxiosError<ApiResponse<null>>;
                push.error(
                    err.response?.data?.message ||
                    "Failed to fetch customer store bills"
                );
                return [];
            } finally {
                destroyLoader();
            }
        };

        const fetchCustomerStoreBillsByStoreId = async (
            storeId: number
        ): Promise<CustomerStoreBill[]> => {
            try {
                const res = await api.get<ApiResponse<CustomerStoreBill[]>>(
                    "/customer-store-bills",
                    { params: { store_id: storeId } }
                );
                if (!res.data.success) return [];
                return res.data.data;
            } catch {
                return [];
            }
        };

        const createCustomerStoreBill = async (
            payload: CreateCustomerStoreBillPayload
        ): Promise<CustomerStoreBill | null> => {
            displayLoader();
            try {
                if (!payload.customer_id) {
                    push.error("Customer is required");
                    return null;
                }
                if (!payload.store_id) {
                    push.error("Store is required");
                    return null;
                }

                const res = await api.post<ApiResponse<CustomerStoreBill>>(
                    "/customer-store-bills",
                    payload
                );
                if (!res.data.success) {
                    push.error(res.data.message);
                    return null;
                }
                bills.value.push(res.data.data);
                push.success(res.data.message);
                return res.data.data;
            } catch (error) {
                const err = error as AxiosError<ApiResponse<null>>;
                push.error(
                    err.response?.data?.message ||
                    "Failed to create customer store bill"
                );
                return null;
            } finally {
                destroyLoader();
            }
        };

        const updateCustomerStoreBill = async (
            id: number,
            payload: UpdateCustomerStoreBillPayload
        ): Promise<CustomerStoreBill | null> => {
            displayLoader();
            try {
                const res = await api.patch<ApiResponse<CustomerStoreBill>>(
                    `/customer-store-bills/${id}`,
                    payload
                );
                if (!res.data.success) {
                    push.error(res.data.message);
                    return null;
                }
                const index = bills.value.findIndex((b) => b.id === id);
                if (index !== -1) bills.value[index] = res.data.data;
                push.success(res.data.message);
                return res.data.data;
            } catch (error) {
                const err = error as AxiosError<ApiResponse<null>>;
                push.error(
                    err.response?.data?.message ||
                    "Failed to update customer store bill"
                );
                return null;
            } finally {
                destroyLoader();
            }
        };

        const createCustomerStoreBillPayment = async (
            billId: number,
            payload: CreateBillPaymentPayload
        ): Promise<boolean> => {
            displayLoader();
            try {
                const res = await api.post<ApiResponse<unknown>>(
                    `/customer-store-bills/${billId}/payments`,
                    payload
                );
                if (!res.data.success) {
                    push.error(res.data.message);
                    return false;
                }
                push.success(res.data.message || "Payment recorded");
                return true;
            } catch (error) {
                const err = error as AxiosError<ApiResponse<null>>;
                push.error(err.response?.data?.message || "Failed to record payment");
                return false;
            } finally {
                destroyLoader();
            }
        };

        const deleteCustomerStoreBill = async (id: number): Promise<boolean> => {
            displayLoader();
            try {
                const res = await api.delete<ApiResponse<null>>(
                    `/customer-store-bills/${id}`
                );
                if (!res.data.success) {
                    push.error(res.data.message);
                    return false;
                }
                bills.value = bills.value.filter((b) => b.id !== id);
                push.success(res.data.message);
                return true;
            } catch (error) {
                const err = error as AxiosError<ApiResponse<null>>;
                push.error(
                    err.response?.data?.message || "Failed to delete customer store bill"
                );
                return false;
            } finally {
                destroyLoader();
            }
        };

        const setSort = (field: CustomerStoreBillSortField) => {
            if (sortField.value === field) {
                sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
            } else {
                sortField.value = field;
                sortDirection.value = "asc";
            }
        };

        const setSearchQuery = (query: string) => {
            searchQuery.value = query;
        };

        const clearSearch = () => {
            searchQuery.value = "";
        };

        const getBillById = (id: number): CustomerStoreBill | undefined => {
            return bills.value.find((b) => b.id === id);
        };

        return {
            bills,
            searchQuery,
            sortField,
            sortDirection,

            filteredBills,
            totalBills,
            totalBilled,
            totalPaid,

            fetchCustomerStoreBills,
            fetchCustomerStoreBillsByStoreId,
            createCustomerStoreBill,
            updateCustomerStoreBill,
            createCustomerStoreBillPayment,
            deleteCustomerStoreBill,

            setSort,
            setSearchQuery,
            clearSearch,

            getBillById,
        };
    }
);