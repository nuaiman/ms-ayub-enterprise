// src/stores/godownBills.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
    GodownBill,
    CreateGodownBillPayload,
    UpdateGodownBillPayload,
    CreateBillPaymentPayload,
    GodownBillSortField,
    SortDirection,
} from "@/types/godownBill";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useGodownBillsStore = defineStore("godownBills", () => {
    const { displayLoader, destroyLoader } = useGlobalLoader();

    const bills = ref<GodownBill[]>([]);
    const searchQuery = ref("");
    const sortField = ref<GodownBillSortField>("created_at");
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
                case "godown_id":
                    comparison = a.godown_id - b.godown_id;
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
                        new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
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

    const fetchGodownBills = async (params?: {
        store_id?: number;
        godown_id?: number;
        month_year?: string;
    }) => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<GodownBill[]>>("/godown-bills", {
                params,
            });
            if (!res.data.success) {
                push.error(res.data.message);
                return [];
            }
            bills.value = res.data.data;
            return bills.value;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to fetch godown bills");
            return [];
        } finally {
            destroyLoader();
        }
    };

    const fetchGodownBillsByStoreId = async (
        storeId: number
    ): Promise<GodownBill[]> => {
        try {
            const res = await api.get<ApiResponse<GodownBill[]>>("/godown-bills", {
                params: { store_id: storeId },
            });
            if (!res.data.success) return [];
            return res.data.data;
        } catch {
            return [];
        }
    };

    const createGodownBill = async (
        payload: CreateGodownBillPayload
    ): Promise<GodownBill | null> => {
        displayLoader();
        try {
            if (!payload.store_id) {
                push.error("Store is required");
                return null;
            }
            if (!payload.godown_id) {
                push.error("Godown is required");
                return null;
            }

            const res = await api.post<ApiResponse<GodownBill>>(
                "/godown-bills",
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
            push.error(err.response?.data?.message || "Failed to create godown bill");
            return null;
        } finally {
            destroyLoader();
        }
    };

    const updateGodownBill = async (
        id: number,
        payload: UpdateGodownBillPayload
    ): Promise<GodownBill | null> => {
        displayLoader();
        try {
            const res = await api.patch<ApiResponse<GodownBill>>(
                `/godown-bills/${id}`,
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
            push.error(err.response?.data?.message || "Failed to update godown bill");
            return null;
        } finally {
            destroyLoader();
        }
    };

    const createGodownBillPayment = async (
        billId: number,
        payload: CreateBillPaymentPayload
    ): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.post<ApiResponse<unknown>>(
                `/godown-bills/${billId}/payments`,
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

    const deleteGodownBill = async (id: number): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.delete<ApiResponse<null>>(`/godown-bills/${id}`);
            if (!res.data.success) {
                push.error(res.data.message);
                return false;
            }
            bills.value = bills.value.filter((b) => b.id !== id);
            push.success(res.data.message);
            return true;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to delete godown bill");
            return false;
        } finally {
            destroyLoader();
        }
    };

    const setSort = (field: GodownBillSortField) => {
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

    const getBillById = (id: number): GodownBill | undefined => {
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

        fetchGodownBills,
        fetchGodownBillsByStoreId,
        createGodownBill,
        updateGodownBill,
        createGodownBillPayment,
        deleteGodownBill,

        setSort,
        setSearchQuery,
        clearSearch,

        getBillById,
    };
});