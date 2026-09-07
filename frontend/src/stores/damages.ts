// src/stores/damages.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Damage,
  CreateDamagePayload,
  UpdateDamagePayload,
  DamageSortField,
  SortDirection,
} from "@/types/damage";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useStoresStore } from "./stores";
import { useUsersStore } from "./users";
import { useExpensesStore } from "./expenses";
import { useLotsStore } from "./lots";
import { useGodownsStore } from "./godowns";

export const useDamagesStore = defineStore("damages", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const damages = ref<Damage[]>([]);
  const searchQuery = ref("");
  const sortField = ref<DamageSortField>("damage_date");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredDamages = computed(() => {
    let result = [...damages.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const storesStore = useStoresStore();
      const usersStore = useUsersStore();
      result = result.filter((damage) => {
        const store = storesStore.getStoreById(damage.store_id);
        const storeName = store ? storesStore.getStoreDisplayName(store) : `Store #${damage.store_id}`;

        return (
          damage.reason.toLowerCase().includes(query) ||
          String(damage.quantity).includes(query) ||
          String(damage.weight).includes(query) ||
          String(damage.amount).includes(query) ||
          (damage.notes && damage.notes.toLowerCase().includes(query)) ||
          storeName.toLowerCase().includes(query) ||
          usersStore.getUserName(damage.user_id).toLowerCase().includes(query)
        );
      });
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "store_id":
          comparison = a.store_id - b.store_id;
          break;
        case "user_id":
          comparison = a.user_id - b.user_id;
          break;
        case "quantity":
          comparison = a.quantity - b.quantity;
          break;
        case "weight":
          comparison = a.weight - b.weight;
          break;
        case "damage_date":
          comparison = new Date(a.damage_date).getTime() - new Date(b.damage_date).getTime();
          break;
        case "amount":
          comparison = a.amount - b.amount;
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

  const totalDamages = computed(() => damages.value.length);

  const totalDamageAmount = computed(() => {
    return damages.value.reduce((sum, damage) => sum + damage.amount, 0);
  });

  // ============= HELPERS =============

  // Generate expense title from damage
  const generateExpenseTitle = (damage: Damage): string => {
    const storesStore = useStoresStore();
    const lotsStore = useLotsStore();
    const godownsStore = useGodownsStore();

    const store = storesStore.getStoreById(damage.store_id);
    const lotName = store ? lotsStore.getLotName(store.lot_id) : `Store #${damage.store_id}`;
    const godownName = store ? godownsStore.getGodownName(store.godown_id) : "";

    let title = `Damage: ${damage.reason}`;
    if (lotName) {
      title = `${title} - ${lotName}`;
    }
    if (godownName) {
      title = `${title} (${godownName})`;
    }

    if (damage.quantity > 0) {
      title = `${title} - ${damage.quantity} ${damage.quantity_unit}`;
    } else if (damage.weight > 0) {
      title = `${title} - ${damage.weight} ${damage.weight_unit}`;
    }

    return title.slice(0, 255);
  };

  // Generate expense notes from damage with damage_id stored
  const generateExpenseNotes = (damage: Damage): string => {
    return `Damage record #${damage.id}: ${damage.reason}${damage.notes ? " - " + damage.notes : ""}`;
  };

  // Find associated expense for a damage by searching notes for "Damage record #{damageId}"
  const findAssociatedExpense = async (damageId: number): Promise<number | null> => {
    const expensesStore = useExpensesStore();

    // Ensure expenses are loaded
    if (expensesStore.expenses.length === 0) {
      await expensesStore.fetchExpenses();
    }

    // Search for expense with "Damage record #{damageId}" in notes
    const matchingExpense = expensesStore.expenses.find((e) =>
      e.notes && e.notes.includes(`Damage record #${damageId}`)
    );

    return matchingExpense ? matchingExpense.id : null;
  };

  // ============= ACTIONS =============

  // GET ALL DAMAGES
  const fetchDamages = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Damage[]>>("/damages");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      damages.value = res.data.data;
      return damages.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch damages");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DAMAGES BY STORE
  const fetchDamagesByStore = async (storeId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Damage[]>>("/damages", {
        params: { store_id: storeId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      damages.value = res.data.data;
      return damages.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch store damages");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DAMAGES BY USER
  const fetchDamagesByUser = async (userId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Damage[]>>("/damages", {
        params: { user_id: userId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      damages.value = res.data.data;
      return damages.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user damages");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DAMAGES BY DATE RANGE
  const fetchDamagesByDateRange = async (startDate: string, endDate: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Damage[]>>("/damages", {
        params: { start_date: startDate, end_date: endDate },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      damages.value = res.data.data;
      return damages.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch damages by date");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE DAMAGE (with expense creation)
  const createDamage = async (payload: CreateDamagePayload): Promise<Damage | null> => {
    displayLoader();
    try {
      // Validate: either quantity or weight must be greater than 0
      if (payload.quantity <= 0 && payload.weight <= 0) {
        push.error("Either quantity or weight must be greater than 0");
        return null;
      }

      const res = await api.post<ApiResponse<Damage>>("/damages", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const newDamage = res.data.data;
      damages.value.push(newDamage);

      // Create expense for this damage with damage_id in notes
      const expensesStore = useExpensesStore();

      const expensePayload = {
        title: generateExpenseTitle(newDamage),
        amount: newDamage.amount,
        expense_date: newDamage.damage_date,
        notes: generateExpenseNotes(newDamage),
      };

      const expense = await expensesStore.createExpense(expensePayload);

      if (expense) {
        push.success(`${res.data.message} (Expense record created)`);
      } else {
        push.warning(`${res.data.message} but expense creation failed`);
      }

      return newDamage;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create damage record");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE DAMAGE (with expense update)
  const updateDamage = async (id: number, payload: UpdateDamagePayload): Promise<Damage | null> => {
    displayLoader();
    try {
      // IMPORTANT: Get the old damage BEFORE updating
      const oldDamage = damages.value.find((d) => d.id === id);

      const res = await api.patch<ApiResponse<Damage>>(`/damages/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const updatedDamage = res.data.data;
      const index = damages.value.findIndex((damage) => damage.id === id);
      if (index !== -1) {
        damages.value[index] = updatedDamage;
      }

      // Find associated expense
      const expenseId = await findAssociatedExpense(id);
      const expensesStore = useExpensesStore();

      if (expenseId) {
        // Update expense with new damage details
        const expenseUpdatePayload: {
          title?: string;
          amount?: number;
          expense_date?: string;
          notes?: string;
        } = {};

        // Check if amount changed - use oldDamage for comparison
        if (oldDamage && oldDamage.amount !== updatedDamage.amount) {
          expenseUpdatePayload.amount = updatedDamage.amount;
          console.log(`[DAMAGES] Amount changed: ${oldDamage.amount} -> ${updatedDamage.amount}`);
        }

        // Check if reason changed
        if (oldDamage && oldDamage.reason !== updatedDamage.reason) {
          expenseUpdatePayload.title = generateExpenseTitle(updatedDamage);
        }

        // Check if damage date changed
        if (oldDamage && oldDamage.damage_date !== updatedDamage.damage_date) {
          expenseUpdatePayload.expense_date = updatedDamage.damage_date;
        }

        // Check if notes changed
        if (oldDamage && oldDamage.notes !== updatedDamage.notes) {
          expenseUpdatePayload.notes = generateExpenseNotes(updatedDamage);
        }

        // Only update if there are changes
        if (Object.keys(expenseUpdatePayload).length > 0) {
          const updatedExpense = await expensesStore.updateExpense(expenseId, expenseUpdatePayload);
          if (updatedExpense) {
            push.success(`${res.data.message} (Associated expense updated)`);
          } else {
            push.warning(`${res.data.message} but expense update failed`);
          }
        } else {
          push.success(res.data.message);
        }
      } else {
        // No associated expense found - create one
        const expensePayload = {
          title: generateExpenseTitle(updatedDamage),
          amount: updatedDamage.amount,
          expense_date: updatedDamage.damage_date,
          notes: generateExpenseNotes(updatedDamage),
        };

        const expense = await expensesStore.createExpense(expensePayload);
        if (expense) {
          push.success(`${res.data.message} (Expense record created)`);
        } else {
          push.success(res.data.message);
        }
      }

      return updatedDamage;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update damage record");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE DAMAGE
  const deleteDamage = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      // Find associated expense first
      const expenseId = await findAssociatedExpense(id);

      const res = await api.delete<ApiResponse<null>>(`/damages/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      damages.value = damages.value.filter((damage) => damage.id !== id);

      // Delete associated expense if found
      if (expenseId) {
        const expensesStore = useExpensesStore();
        const deleted = await expensesStore.deleteExpense(expenseId);
        if (deleted) {
          push.success(`${res.data.message} (Associated expense deleted)`);
        } else {
          push.warning(`${res.data.message} but expense deletion failed`);
        }
      } else {
        push.success(res.data.message);
      }

      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete damage record");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: DamageSortField) => {
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

  // Get damage by ID
  const getDamageById = (id: number): Damage | undefined => {
    return damages.value.find((d) => d.id === id);
  };

  // Get damages by store ID
  const getDamagesByStoreId = (storeId: number): Damage[] => {
    return damages.value.filter((damage) => damage.store_id === storeId);
  };

  // Get total damages by store
  const getTotalDamagesByStore = (storeId: number): number => {
    return damages.value
      .filter((damage) => damage.store_id === storeId)
      .reduce((sum, damage) => sum + damage.amount, 0);
  };

  // Format damage date
  const formatDamageDate = (dateStr: string): string => {
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
    damages,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredDamages,
    totalDamages,
    totalDamageAmount,

    // Fetch
    fetchDamages,
    fetchDamagesByStore,
    fetchDamagesByUser,
    fetchDamagesByDateRange,

    // CRUD
    createDamage,
    updateDamage,
    deleteDamage,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getDamageById,
    getDamagesByStoreId,
    getTotalDamagesByStore,
    formatDamageDate,
  };
});