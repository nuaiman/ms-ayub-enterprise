// src/stores/stores.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Store,
  CreateStorePayload,
  UpdateStorePayload,
  StoreSortField,
  SortDirection,
} from "@/types/store";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useLotsStore } from "./lots";

export const useStoresStore = defineStore("stores", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  const stores = ref<Store[]>([]);
  const searchQuery = ref("");
  const sortField = ref<StoreSortField>("created_at");
  const sortDirection = ref<SortDirection>("desc");

  const filteredStores = computed(() => {
    let result = [...stores.value];

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const lotsStore = useLotsStore();
      result = result.filter(
        (store) =>
          lotsStore.getLotName(store.lot_id).toLowerCase().includes(query) ||
          String(store.quantity).includes(query) ||
          String(store.weight).includes(query)
      );
    }

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
        case "start_date":
          comparison = new Date(a.start_date).getTime() - new Date(b.start_date).getTime();
          break;
        case "is_active":
          comparison = a.is_active === b.is_active ? 0 : a.is_active ? -1 : 1;
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

  const activeStores = computed(() => stores.value.filter((store) => store.is_active));
  const storesWithInventory = computed(() =>
    stores.value.filter((store) => store.quantity > 0 || store.weight > 0)
  );
  const totalStores = computed(() => stores.value.length);

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

  const createStore = async (payload: CreateStorePayload): Promise<Store | null> => {
    displayLoader();
    try {
      if (!payload.lot_id) {
        push.error("Lot is required");
        return null;
      }
      if (!payload.godown_id) {
        push.error("Godown is required");
        return null;
      }
      if (payload.weight < 0) {
        push.error("Weight cannot be negative");
        return null;
      }
      if (payload.quantity < 0) {
        push.error("Quantity cannot be negative");
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
      push.error(err.response?.data?.message || "Failed to create store");
      return null;
    } finally {
      destroyLoader();
    }
  };

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

  const setSort = (field: StoreSortField) => {
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

  const getStoreDisplayName = (store: Store): string => {
    const lotsStore = useLotsStore();
    return lotsStore.getLotName(store.lot_id);
  };

  const getStoreById = (id: number): Store | undefined => {
    return stores.value.find((s) => s.id === id);
  };

  const getStoresByLotId = (lotId: number): Store[] => {
    return stores.value.filter((store) => store.lot_id === lotId);
  };

  const getTotalQuantityByLot = (lotId: number): number => {
    return stores.value
      .filter((store) => store.lot_id === lotId)
      .reduce((sum, store) => sum + store.quantity, 0);
  };

  const getTotalWeightByLot = (lotId: number): number => {
    return stores.value
      .filter((store) => store.lot_id === lotId)
      .reduce((sum, store) => sum + store.weight, 0);
  };

  const hasInventory = (store: Store): boolean => {
    return store.quantity > 0 || store.weight > 0;
  };

  const refreshStoreBills = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<unknown>>(`/stores/${id}/refresh-bills`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      push.success(res.data.message || "Bills refreshed");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to refresh bills");
      return false;
    } finally {
      destroyLoader();
    }
  };

  return {
    stores,
    searchQuery,
    sortField,
    sortDirection,

    filteredStores,
    activeStores,
    storesWithInventory,
    totalStores,

    fetchStores,
    fetchStoresByLot,
    fetchActiveStores,
    fetchStoresWithInventory,

    createStore,
    updateStore,
    deleteStore,

    setSort,
    setSearchQuery,
    clearSearch,

    getStoreDisplayName,
    getStoreById,
    getStoresByLotId,
    getTotalQuantityByLot,
    getTotalWeightByLot,
    hasInventory,

    refreshStoreBills,
  };
});