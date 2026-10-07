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

// Marker used to link an expense row back to its damage.
// Stored on its own line at the end of the expense notes.
// Anchored by endsWith() during lookup so [D:4] won't match [D:42].
const damageMarker = (damageId: number): string => `[D:${damageId}]`;

const generateExpenseTitle = (damage: Damage): string => {
    const storesStore = useStoresStore();
    const lotsStore = useLotsStore();
    const godownsStore = useGodownsStore();

    const store = storesStore.getStoreById(damage.store_id);
    const lotName = store ? lotsStore.getLotName(store.lot_id) : "";
    const godownName = store ? godownsStore.getGodownName(store.godown_id) : "";

    const parts: string[] = [`Damage: ${damage.reason}`];

    if (lotName) {
        parts.push(godownName ? `${lotName} (${godownName})` : lotName);
    } else {
        parts.push(`Store #${damage.store_id}`);
    }

    if (damage.quantity > 0) {
        parts.push(`${damage.quantity} ${damage.quantity_unit}`);
    } else if (damage.weight > 0) {
        parts.push(`${damage.weight} ${damage.weight_unit}`);
    }

    return parts.join(" — ");
};

const generateExpenseNotes = (damage: Damage): string => {
    const marker = damageMarker(damage.id);
    return damage.notes ? `${damage.notes}\n${marker}` : marker;
};

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

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase();
            const storesStore = useStoresStore();
            const usersStore = useUsersStore();
            result = result.filter((damage) => {
                const store = storesStore.getStoreById(damage.store_id);
                const storeName = store
                    ? storesStore.getStoreDisplayName(store)
                    : `Store #${damage.store_id}`;

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
                    comparison =
                        new Date(a.damage_date).getTime() - new Date(b.damage_date).getTime();
                    break;
                case "amount":
                    comparison = a.amount - b.amount;
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

    const totalDamages = computed(() => damages.value.length);

    const totalDamageAmount = computed(() => {
        return damages.value.reduce((sum, damage) => sum + damage.amount, 0);
    });

    // ============= EXPENSE LINK HELPERS =============

    // Find the expense row linked to a damage by scanning its notes for
    // the trailing [D:<id>] marker. Loads the expenses list lazily if
    // it's currently empty; returns null on any failure.
    const findAssociatedExpense = async (damageId: number): Promise<number | null> => {
        const expensesStore = useExpensesStore();

        if (expensesStore.expenses.length === 0) {
            await expensesStore.fetchExpenses();
        }

        const marker = damageMarker(damageId);
        const match = expensesStore.expenses.find(
            (e) => typeof e.notes === "string" && e.notes.endsWith(marker)
        );

        return match ? match.id : null;
    };

    // Create the expense row for a freshly-created damage.
    const createExpenseForDamage = async (damage: Damage) => {
        const expensesStore = useExpensesStore();
        const expense = await expensesStore.createExpense({
            title: generateExpenseTitle(damage),
            amount: damage.amount,
            expense_date: damage.damage_date,
            notes: generateExpenseNotes(damage),
        });
        if (!expense) {
            push.warning("Damage saved, but the linked expense could not be created.");
        }
    };

    // Update the expense row for a damage that already has one.
    // Silent no-op if no linked expense is found.
    const updateExpenseForDamage = async (damage: Damage) => {
        const expenseId = await findAssociatedExpense(damage.id);
        if (expenseId === null) return;

        const expensesStore = useExpensesStore();
        const updated = await expensesStore.updateExpense(expenseId, {
            title: generateExpenseTitle(damage),
            amount: damage.amount,
            expense_date: damage.damage_date,
            notes: generateExpenseNotes(damage),
        });
        if (!updated) {
            push.warning("Damage updated, but the linked expense could not be updated.");
        }
    };

    // Delete the expense row for a damage that is being removed.
    // Silent no-op if no linked expense is found.
    const deleteExpenseForDamage = async (damageId: number) => {
        const expenseId = await findAssociatedExpense(damageId);
        if (expenseId === null) return;

        const expensesStore = useExpensesStore();
        const ok = await expensesStore.deleteExpense(expenseId);
        if (!ok) {
            push.warning("Damage deleted, but the linked expense could not be deleted.");
        }
    };

    // ============= FETCH =============

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

    // ============= CREATE =============

    const createDamage = async (payload: CreateDamagePayload): Promise<Damage | null> => {
        displayLoader();
        try {
            const qty = payload.quantity ?? 0;
            const wt = payload.weight ?? 0;
            if (qty <= 0 && wt <= 0) {
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
            push.success(res.data.message);

            // Mirror into expenses.
            await createExpenseForDamage(newDamage);

            return newDamage;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to create damage record");
            return null;
        } finally {
            destroyLoader();
        }
    };

    // ============= UPDATE =============

    const updateDamage = async (
        id: number,
        payload: UpdateDamagePayload
    ): Promise<Damage | null> => {
        displayLoader();
        try {
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
            push.success(res.data.message);

            // Keep the linked expense in sync.
            await updateExpenseForDamage(updatedDamage);

            return updatedDamage;
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>;
            push.error(err.response?.data?.message || "Failed to update damage record");
            return null;
        } finally {
            destroyLoader();
        }
    };

    // ============= DELETE =============

    const deleteDamage = async (id: number): Promise<boolean> => {
        displayLoader();
        try {
            const res = await api.delete<ApiResponse<null>>(`/damages/${id}`);
            if (!res.data.success) {
                push.error(res.data.message);
                return false;
            }

            damages.value = damages.value.filter((damage) => damage.id !== id);
            push.success(res.data.message);

            // Remove the linked expense.
            await deleteExpenseForDamage(id);

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

    const getDamageById = (id: number): Damage | undefined => {
        return damages.value.find((d) => d.id === id);
    };

    const getDamagesByStoreId = (storeId: number): Damage[] => {
        return damages.value.filter((damage) => damage.store_id === storeId);
    };

    const getTotalDamagesByStore = (storeId: number): number => {
        return damages.value
            .filter((damage) => damage.store_id === storeId)
            .reduce((sum, damage) => sum + damage.amount, 0);
    };

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