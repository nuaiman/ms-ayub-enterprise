// src/stores/incomes.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
    Income,
    CreateIncomePayload,
    UpdateIncomePayload,
    IncomeSortField,
    SortDirection,
} from "@/types/income";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useUsersStore } from "./users";

export const useIncomesStore = defineStore("incomes", () => {
    const { displayLoader, destroyLoader } = useGlobalLoader();

    // ============= STATE =============
    const incomes = ref<Income[]>([]);
    const searchQuery = ref("");
    const sortField = ref<IncomeSortField>("income_date");
    const sortDirection = ref<SortDirection>("desc");

    // ============= COMPUTED =============
    const filteredIncomes = computed(() => {
        let result = [...incomes.value];

        // Filter by search query
        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase();
            const usersStore = useUsersStore();
            result = result.filter(
                (income) =>
                    income.title.toLowerCase().includes(query) ||
                    String(income.amount).includes(query) ||
                    (income.notes && income.notes.toLowerCase().includes(query)) ||
                    usersStore.getUserName(income.user_id).toLowerCase().includes(query)
            );
        }

        // Sort
        result.sort((a, b) => {
            let comparison = 0;
            switch (sortField.value) {
                case "title":
                    comparison = a.title.localeCompare(b.title);
                    break;
                case "amount":
                    comparison = a.amount - b.amount;
                    break;
                case "income_date":
                    comparison = new Date(a.income_date).getTime() - new Date(b.income_date).getTime();
                    break;
                case "created_at":
                    comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
                    break;
                default:
                    comparison = 0;
            }
            return sortDirection.value === "desc" ? -comparison : comparison;
        });

        return result;
    });

    const totalIncomes = computed(() => incomes.value.length);

    const totalAmount = computed(() => {
        return incomes.value.reduce((sum, income) => sum + income.amount, 0);
    });

    // ============= ACTIONS =============

    // GET ALL INCOMES
    const fetchIncomes = async () => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<Income[]>>("/incomes");
            if (!res.data.success) {
                push.error(res.data.message);
                return [];
            }
            incomes.value = res.data.data;
            return incomes.value;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to fetch incomes");
            return [];
        } finally {
            destroyLoader();
        }
    };

    // SEARCH INCOMES
    const searchIncomes = async (query: string) => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<Income[]>>("/incomes", {
                params: { search: query },
            });
            if (!res.data.success) {
                push.error(res.data.message);
                return [];
            }
            incomes.value = res.data.data;
            return incomes.value;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to search incomes");
            return [];
        } finally {
            destroyLoader();
        }
    };

    // GET INCOMES BY USER
    const fetchIncomesByUser = async (userId: number) => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<Income[]>>("/incomes", {
                params: { user_id: userId },
            });
            if (!res.data.success) {
                push.error(res.data.message);
                return [];
            }
            incomes.value = res.data.data;
            return incomes.value;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to fetch user incomes");
            return [];
        } finally {
            destroyLoader();
        }
    };

    // GET INCOMES BY DATE RANGE
    const fetchIncomesByDateRange = async (startDate: string, endDate: string) => {
        displayLoader();
        try {
            const res = await api.get<ApiResponse<Income[]>>("/incomes", {
                params: { start_date: startDate, end_date: endDate },
            });
            if (!res.data.success) {
                push.error(res.data.message);
                return [];
            }
            incomes.value = res.data.data;
            return incomes.value;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to fetch incomes by date");
            return [];
        } finally {
            destroyLoader();
        }
    };

    // CREATE INCOME
    const createIncome = async (payload: CreateIncomePayload): Promise<Income | null> => {
        displayLoader();
        try {
            // Validate: title is required
            if (!payload.title) {
                push.error("Title is required");
                return null;
            }

            // Validate: amount must be greater than 0
            if (payload.amount <= 0) {
                push.error("Amount must be greater than 0");
                return null;
            }

            const res = await api.post<ApiResponse<Income>>("/incomes", payload);
            if (!res.data.success) {
                push.error(res.data.message);
                return null;
            }
            incomes.value.push(res.data.data);
            push.success(res.data.message);
            return res.data.data;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to create income");
            return null;
        } finally {
            destroyLoader();
        }
    };

    // UPDATE INCOME
    const updateIncome = async (id: number, payload: UpdateIncomePayload): Promise<Income | null> => {
        displayLoader();
        try {
            const res = await api.patch<ApiResponse<Income>>(`/incomes/${id}`, payload);
            if (!res.data.success) {
                push.error(res.data.message);
                return null;
            }
            const index = incomes.value.findIndex((income) => income.id === id);
            if (index !== -1) {
                incomes.value[index] = res.data.data;
            }
            push.success(res.data.message);
            return res.data.data;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to update income");
            return null;
        } finally {
            destroyLoader();
        }
    };

    // DELETE INCOME
    const deleteIncome = async (id: number): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.delete<ApiResponse<null>>(`/incomes/${id}`);
            if (!res.data.success) {
                push.error(res.data.message);
                return false;
            }
            incomes.value = incomes.value.filter((income) => income.id !== id);
            push.success(res.data.message);
            return true;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to delete income");
            return false;
        } finally {
            destroyLoader();
        }
    };

    // ============= SORT =============
    const setSort = (field: IncomeSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
        } else {
            sortField.value = field;
            sortDirection.value = "asc";
        }
    };

    // ============= SEARCH =============
    const setSearchQuery = (query: string) => {
        searchQuery.value = query;
    };

    const clearSearch = () => {
        searchQuery.value = "";
    };

    // ============= UTILITIES =============

    const getIncomeById = (id: number): Income | undefined => {
        return incomes.value.find((i) => i.id === id);
    };

    const getIncomesByUserId = (userId: number): Income[] => {
        return incomes.value.filter((income) => income.user_id === userId);
    };

    const getTotalAmountByUser = (userId: number): number => {
        return incomes.value
            .filter((income) => income.user_id === userId)
            .reduce((sum, income) => sum + income.amount, 0);
    };

    const getTotalAmountByDateRange = (startDate: string, endDate: string): number => {
        return incomes.value
            .filter((income) => {
                const incomeDate = new Date(income.income_date);
                const start = new Date(startDate);
                const end = new Date(endDate);
                return incomeDate >= start && incomeDate <= end;
            })
            .reduce((sum, income) => sum + income.amount, 0);
    };

    const formatIncomeDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString("en-US", {
            month: "short",
            day: "numeric",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    return {
        // State
        incomes,
        searchQuery,
        sortField,
        sortDirection,

        // Computed
        filteredIncomes,
        totalIncomes,
        totalAmount,

        // Fetch
        fetchIncomes,
        searchIncomes,
        fetchIncomesByUser,
        fetchIncomesByDateRange,

        // CRUD
        createIncome,
        updateIncome,
        deleteIncome,

        // Sort
        setSort,

        // Search
        setSearchQuery,
        clearSearch,

        // Utilities
        getIncomeById,
        getIncomesByUserId,
        getTotalAmountByUser,
        getTotalAmountByDateRange,
        formatIncomeDate,
    };
});