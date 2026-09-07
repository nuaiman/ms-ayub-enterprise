// src/stores/stores.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Store,
  CreateStorePayload,
  UpdateStorePayload,
  StoreBillType,
  StoreSortField,
  SortDirection,
} from "@/types/store";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useLotsStore } from "./lots";
import { useGodownsStore } from "./godowns";

export const useStoresStore = defineStore("stores", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const stores = ref<Store[]>([]);
  const searchQuery = ref("");
  const sortField = ref<StoreSortField>("created_at");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredStores = computed(() => {
    let result = [...stores.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const lotsStore = useLotsStore();
      const godownsStore = useGodownsStore();
      result = result.filter(
        (store) =>
          lotsStore.getLotName(store.lot_id).toLowerCase().includes(query) ||
          godownsStore.getGodownName(store.godown_id).toLowerCase().includes(query) ||
          store.store_bill_type.toLowerCase().includes(query) ||
          String(store.quantity).includes(query) ||
          String(store.weight).includes(query) ||
          (store.notes && store.notes.toLowerCase().includes(query))
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "lot_id":
          comparison = a.lot_id - b.lot_id;
          break;
        case "godown_id":
          comparison = a.godown_id - b.godown_id;
          break;
        case "quantity":
          comparison = a.quantity - b.quantity;
          break;
        case "weight":
          comparison = a.weight - b.weight;
          break;
        case "is_active":
          comparison = (a.is_active === b.is_active) ? 0 : a.is_active ? -1 : 1;
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

  const activeStores = computed(() => {
    return stores.value.filter((store) => store.is_active);
  });

  const storesWithInventory = computed(() => {
    return stores.value.filter((store) => store.quantity > 0 || store.weight > 0);
  });

  const totalStores = computed(() => stores.value.length);

  // ============= ACTIONS =============

  // GET ALL STORES
  const fetchStores = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Store[]>>("/stores");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      stores.value = res.data.data;
      return stores.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch stores");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET STORES BY LOT
  const fetchStoresByLot = async (lotId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Store[]>>("/stores", {
        params: { lot_id: lotId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      stores.value = res.data.data;
      return stores.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch lot stores");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET STORES BY GODOWN
  const fetchStoresByGodown = async (godownId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Store[]>>("/stores", {
        params: { godown_id: godownId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      stores.value = res.data.data;
      return stores.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch godown stores");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET ACTIVE STORES
  const fetchActiveStores = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Store[]>>("/stores", {
        params: { active: true },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      stores.value = res.data.data;
      return stores.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch active stores");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET STORES WITH INVENTORY
  const fetchStoresWithInventory = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Store[]>>("/stores", {
        params: { with_inventory: true },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      stores.value = res.data.data;
      return stores.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch stores with inventory");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE STORE
  const createStore = async (payload: CreateStorePayload): Promise<Store | null> => {
    displayLoader();
    try {
      // Validate: lot_id is required
      if (!payload.lot_id) {
        push.error("Lot is required");
        return null;
      }

      // Validate: godown_id is required
      if (!payload.godown_id) {
        push.error("Godown is required");
        return null;
      }

      // Validate: store_bill_type if provided
      if (payload.store_bill_type && payload.store_bill_type !== 'weight' && payload.store_bill_type !== 'quantity') {
        push.error("store_bill_type must be 'weight' or 'quantity'");
        return null;
      }

      const res = await api.post<ApiResponse<Store>>("/stores", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      stores.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      // Check if it's a 409 conflict error
      if (err.response?.status === 409) {
        const lotsStore = useLotsStore();
        const godownsStore = useGodownsStore();
        const lotName = lotsStore.getLotName(payload.lot_id);
        const godownName = godownsStore.getGodownName(payload.godown_id);
        push.error(`A store already exists for "${lotName}" in "${godownName}"`);
      } else {
        push.error(err.response?.data?.message || "Failed to create store");
      }
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE STORE
  const updateStore = async (id: number, payload: UpdateStorePayload): Promise<Store | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Store>>(`/stores/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = stores.value.findIndex((store) => store.id === id);
      if (index !== -1) {
        stores.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update store");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // TOGGLE STORE ACTIVE
  const toggleStoreActive = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Store>>(`/stores/${id}/toggle-active`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = stores.value.findIndex((store) => store.id === id);
      if (index !== -1) {
        stores.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to toggle store status");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // DELETE STORE
  const deleteStore = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/stores/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      stores.value = stores.value.filter((store) => store.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete store");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: StoreSortField) => {
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

  // Get store display name
  const getStoreDisplayName = (store: Store): string => {
    const lotsStore = useLotsStore();
    const godownsStore = useGodownsStore();
    const lotName = lotsStore.getLotName(store.lot_id);
    const godownName = godownsStore.getGodownName(store.godown_id);
    return `${lotName} @ ${godownName}`;
  };

  // Get store by ID
  const getStoreById = (id: number): Store | undefined => {
    return stores.value.find((s) => s.id === id);
  };

  // Get stores by lot ID
  const getStoresByLotId = (lotId: number): Store[] => {
    return stores.value.filter((store) => store.lot_id === lotId);
  };

  // Get stores by godown ID
  const getStoresByGodownId = (godownId: number): Store[] => {
    return stores.value.filter((store) => store.godown_id === godownId);
  };

  // Get total quantity by lot
  const getTotalQuantityByLot = (lotId: number): number => {
    return stores.value
      .filter((store) => store.lot_id === lotId)
      .reduce((sum, store) => sum + store.quantity, 0);
  };

  // Get total weight by lot
  const getTotalWeightByLot = (lotId: number): number => {
    return stores.value
      .filter((store) => store.lot_id === lotId)
      .reduce((sum, store) => sum + store.weight, 0);
  };

  // Check if store has inventory
  const hasInventory = (store: Store): boolean => {
    return store.quantity > 0 || store.weight > 0;
  };

  // Format store bill type
  const formatStoreBillType = (type: StoreBillType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  return {
    // State
    stores,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredStores,
    activeStores,
    storesWithInventory,
    totalStores,

    // Fetch
    fetchStores,
    fetchStoresByLot,
    fetchStoresByGodown,
    fetchActiveStores,
    fetchStoresWithInventory,

    // CRUD
    createStore,
    updateStore,
    toggleStoreActive,
    deleteStore,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getStoreDisplayName,
    getStoreById,
    getStoresByLotId,
    getStoresByGodownId,
    getTotalQuantityByLot,
    getTotalWeightByLot,
    hasInventory,
    formatStoreBillType,
  };
});