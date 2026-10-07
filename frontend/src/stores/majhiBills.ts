// src/stores/majhiBills.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
    MajhiBill,
    CreateMajhiBillPayload,
    UpdateMajhiBillPayload,
    CreateMajhiBillPaymentPayload,
    MajhiBillSortField,
    SortDirection,
} from "@/types/majhiBill";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useMajhiBillsStore = defineStore("majhiBills", () => {
    const { displayLoader, destroyLoader } = useGlobalLoader();

    const bills = ref<MajhiBill[]>([]);
    const searchQuery = ref("");
    const sortField = ref<MajhiBillSortField>("created_at");
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

    const fetchMajhiBills = async (params?: {
        store_id?: number;
        majhi_id?: number;
        delivery_item_id?: number;
    }) => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<MajhiBill[]>>("/majhi-bills", {
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
            push.error(err.response?.data?.message || "Failed to fetch majhi bills");
            return [];
        } finally {
            destroyLoader();
        }
    };

    const fetchMajhiBillsByStoreId = async (
        storeId: number
    ): Promise<MajhiBill[]> => {
        try {
            const res = await api.get<ApiResponse<MajhiBill[]>>("/majhi-bills", {
                params: { store_id: storeId },
            });
            if (!res.data.success) return [];
            return res.data.data;
        } catch {
            return [];
        }
    };

    const fetchMajhiBillsByDeliveryItemId = async (
        deliveryItemId: number
    ): Promise<MajhiBill[]> => {
        try {
            const res = await api.get<ApiResponse<MajhiBill[]>>("/majhi-bills", {
                params: { delivery_item_id: deliveryItemId },
            });
            if (!res.data.success) return [];
            return res.data.data;
        } catch {
            return [];
        }
    };

    const createMajhiBill = async (
        payload: CreateMajhiBillPayload
    ): Promise<MajhiBill | null> => {
        displayLoader();
        try {
            if (!payload.majhi_id) {
                push.error("Majhi is required");
                return null;
            }
            const hasStore = !!payload.store_id;
            const hasDeliveryItem = !!payload.delivery_item_id;
            if (hasStore === hasDeliveryItem) {
                push.error("Exactly one of store or delivery item must be provided");
                return null;
            }

            const res = await api.post<ApiResponse<MajhiBill>>(
                "/majhi-bills",
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
            push.error(err.response?.data?.message || "Failed to create majhi bill");
            return null;
        } finally {
            destroyLoader();
        }
    };

    const updateMajhiBill = async (
        id: number,
        payload: UpdateMajhiBillPayload
    ): Promise<MajhiBill | null> => {
        displayLoader();
        try {
            const res = await api.patch<ApiResponse<MajhiBill>>(
                `/majhi-bills/${id}`,
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
            push.error(err.response?.data?.message || "Failed to update majhi bill");
            return null;
        } finally {
            destroyLoader();
        }
    };

    const createMajhiBillPayment = async (
        billId: number,
        payload: CreateMajhiBillPaymentPayload
    ): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.post<ApiResponse<unknown>>(
                `/majhi-bills/${billId}/payments`,
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

    const deleteMajhiBill = async (id: number): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.delete<ApiResponse<null>>(`/majhi-bills/${id}`);
            if (!res.data.success) {
                push.error(res.data.message);
                return false;
            }
            bills.value = bills.value.filter((b) => b.id !== id);
            push.success(res.data.message);
            return true;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to delete majhi bill");
            return false;
        } finally {
            destroyLoader();
        }
    };

    const setSort = (field: MajhiBillSortField) => {
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

    const getBillById = (id: number): MajhiBill | undefined => {
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

        fetchMajhiBills,
        fetchMajhiBillsByStoreId,
        fetchMajhiBillsByDeliveryItemId,
        createMajhiBill,
        updateMajhiBill,
        createMajhiBillPayment,
        deleteMajhiBill,

        setSort,
        setSearchQuery,
        clearSearch,

        getBillById,
    };
});