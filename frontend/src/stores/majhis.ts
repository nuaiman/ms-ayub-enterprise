// src/stores/majhis.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Majhi,
  CreateMajhiPayload,
  UpdateMajhiPayload,
  MajhiSortField,
  SortDirection,
} from "@/types/majhi";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useMajhisStore = defineStore("majhis", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const majhis = ref<Majhi[]>([]);
  const searchQuery = ref("");
  const sortField = ref<MajhiSortField>("name");
  const sortDirection = ref<SortDirection>("asc");

  // ============= COMPUTED =============
  const filteredMajhis = computed(() => {
    let result = [...majhis.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(
        (m) =>
          m.name.toLowerCase().includes(query) ||
          (m.phone && m.phone.toLowerCase().includes(query)) ||
          (m.notes && m.notes.toLowerCase().includes(query))
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

  const totalMajhis = computed(() => majhis.value.length);

  // ============= ACTIONS =============

  // GET ALL MAJHIS
  const fetchMajhis = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Majhi[]>>("/majhis");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      majhis.value = res.data.data;
      return majhis.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch majhis");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH MAJHIS
  const searchMajhis = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Majhi[]>>("/majhis", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      majhis.value = res.data.data;
      return majhis.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search majhis");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE MAJHI
  const createMajhi = async (payload: CreateMajhiPayload): Promise<Majhi | null> => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<Majhi>>("/majhis", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      majhis.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create majhi");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE MAJHI
  const updateMajhi = async (id: number, payload: UpdateMajhiPayload): Promise<Majhi | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Majhi>>(`/majhis/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = majhis.value.findIndex((m) => m.id === id);
      if (index !== -1) {
        majhis.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update majhi");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE MAJHI
  const deleteMajhi = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/majhis/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      majhis.value = majhis.value.filter((m) => m.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete majhi");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: MajhiSortField) => {
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

  // Get majhi name by ID
  const getMajhiName = (id: number): string => {
    const majhi = majhis.value.find((m) => m.id === id);
    return majhi ? majhi.name : `Majhi #${id}`;
  };

  return {
    // State
    majhis,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredMajhis,
    totalMajhis,

    // Fetch
    fetchMajhis,
    searchMajhis,

    // CRUD
    createMajhi,
    updateMajhi,
    deleteMajhi,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getMajhiName,
  };
});