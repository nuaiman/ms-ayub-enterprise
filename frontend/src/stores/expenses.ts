// src/stores/expenses.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Expense,
  CreateExpensePayload,
  UpdateExpensePayload,
  ExpenseSortField,
  SortDirection,
} from "@/types/expense";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useUsersStore } from "./users";

export const useExpensesStore = defineStore("expenses", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const expenses = ref<Expense[]>([]);
  const searchQuery = ref("");
  const sortField = ref<ExpenseSortField>("expense_date");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredExpenses = computed(() => {
    let result = [...expenses.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const usersStore = useUsersStore();
      result = result.filter(
        (expense) =>
          expense.title.toLowerCase().includes(query) ||
          String(expense.amount).includes(query) ||
          (expense.notes && expense.notes.toLowerCase().includes(query)) ||
          usersStore.getUserName(expense.user_id).toLowerCase().includes(query)
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
        case "expense_date":
          comparison = new Date(a.expense_date).getTime() - new Date(b.expense_date).getTime();
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

  const totalExpenses = computed(() => expenses.value.length);

  const totalAmount = computed(() => {
    return expenses.value.reduce((sum, expense) => sum + expense.amount, 0);
  });

  // ============= ACTIONS =============

  // GET ALL EXPENSES
  const fetchExpenses = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense[]>>("/expenses");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      expenses.value = res.data.data;
      return expenses.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch expenses");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH EXPENSES
  const searchExpenses = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense[]>>("/expenses", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      expenses.value = res.data.data;
      return expenses.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search expenses");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET EXPENSES BY USER
  const fetchExpensesByUser = async (userId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense[]>>("/expenses", {
        params: { user_id: userId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      expenses.value = res.data.data;
      return expenses.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user expenses");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET EXPENSES BY DATE RANGE
  const fetchExpensesByDateRange = async (startDate: string, endDate: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense[]>>("/expenses", {
        params: { start_date: startDate, end_date: endDate },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      expenses.value = res.data.data;
      return expenses.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch expenses by date");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE EXPENSE
  const createExpense = async (payload: CreateExpensePayload): Promise<Expense | null> => {
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

      const res = await api.post<ApiResponse<Expense>>("/expenses", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      expenses.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create expense");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE EXPENSE
  const updateExpense = async (id: number, payload: UpdateExpensePayload): Promise<Expense | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Expense>>(`/expenses/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = expenses.value.findIndex((expense) => expense.id === id);
      if (index !== -1) {
        expenses.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update expense");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE EXPENSE
  const deleteExpense = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/expenses/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      expenses.value = expenses.value.filter((expense) => expense.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete expense");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: ExpenseSortField) => {
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

  // Get expense by ID
  const getExpenseById = (id: number): Expense | undefined => {
    return expenses.value.find((e) => e.id === id);
  };

  // Get expenses by user ID
  const getExpensesByUserId = (userId: number): Expense[] => {
    return expenses.value.filter((expense) => expense.user_id === userId);
  };

  // Get total amount by user
  const getTotalAmountByUser = (userId: number): number => {
    return expenses.value
      .filter((expense) => expense.user_id === userId)
      .reduce((sum, expense) => sum + expense.amount, 0);
  };

  // Get total amount by date range
  const getTotalAmountByDateRange = (startDate: string, endDate: string): number => {
    return expenses.value
      .filter((expense) => {
        const expenseDate = new Date(expense.expense_date);
        const start = new Date(startDate);
        const end = new Date(endDate);
        return expenseDate >= start && expenseDate <= end;
      })
      .reduce((sum, expense) => sum + expense.amount, 0);
  };

  // Format expense date
  const formatExpenseDate = (dateStr: string): string => {
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
    expenses,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredExpenses,
    totalExpenses,
    totalAmount,

    // Fetch
    fetchExpenses,
    searchExpenses,
    fetchExpensesByUser,
    fetchExpensesByDateRange,

    // CRUD
    createExpense,
    updateExpense,
    deleteExpense,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getExpenseById,
    getExpensesByUserId,
    getTotalAmountByUser,
    getTotalAmountByDateRange,
    formatExpenseDate,
  };
});