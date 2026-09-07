// src/stores/brokers.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Broker,
  CreateBrokerPayload,
  UpdateBrokerPayload,
  BrokerSortField,
  SortDirection,
} from "@/types/broker";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useBrokersStore = defineStore("brokers", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const brokers = ref<Broker[]>([]);
  const searchQuery = ref("");
  const sortField = ref<BrokerSortField>("name");
  const sortDirection = ref<SortDirection>("asc");

  // ============= COMPUTED =============
  const filteredBrokers = computed(() => {
    let result = [...brokers.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(
        (b) =>
          b.name.toLowerCase().includes(query) ||
          (b.phone && b.phone.toLowerCase().includes(query)) ||
          (b.notes && b.notes.toLowerCase().includes(query))
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "name":
          comparison = a.name.localeCompare(b.name);
          break;
        case "phone":
          comparison = (a.phone || "").localeCompare(b.phone || "");
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

  const totalBrokers = computed(() => brokers.value.length);

  // ============= ACTIONS =============

  // GET ALL BROKERS
  const fetchBrokers = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Broker[]>>("/brokers");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      brokers.value = res.data.data;
      return brokers.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch brokers");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH BROKERS
  const searchBrokers = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Broker[]>>("/brokers", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      brokers.value = res.data.data;
      return brokers.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search brokers");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE BROKER
  const createBroker = async (payload: CreateBrokerPayload): Promise<Broker | null> => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<Broker>>("/brokers", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      brokers.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create broker");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE BROKER
  const updateBroker = async (id: number, payload: UpdateBrokerPayload): Promise<Broker | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Broker>>(`/brokers/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = brokers.value.findIndex((b) => b.id === id);
      if (index !== -1) {
        brokers.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update broker");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE BROKER
  const deleteBroker = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/brokers/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      brokers.value = brokers.value.filter((b) => b.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete broker");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: BrokerSortField) => {
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

  // Get broker name by ID
  const getBrokerName = (id: number): string => {
    const broker = brokers.value.find((b) => b.id === id);
    return broker ? broker.name : `Broker #${id}`;
  };

  return {
    // State
    brokers,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredBrokers,
    totalBrokers,

    // Fetch
    fetchBrokers,
    searchBrokers,

    // CRUD
    createBroker,
    updateBroker,
    deleteBroker,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getBrokerName,
  };
});