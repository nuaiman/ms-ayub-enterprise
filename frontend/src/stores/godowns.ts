// src/stores/godowns.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Godown,
  CreateGodownPayload,
  UpdateGodownPayload,
  GodownSortField,
  SortDirection,
} from "@/types/godown";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useGodownsStore = defineStore("godowns", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const godowns = ref<Godown[]>([]);
  const searchQuery = ref("");
  const sortField = ref<GodownSortField>("name");
  const sortDirection = ref<SortDirection>("asc");

  // ============= COMPUTED =============
  const filteredGodowns = computed(() => {
    let result = [...godowns.value];

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(
        (g) =>
          g.name.toLowerCase().includes(query) ||
          (g.phone && g.phone.toLowerCase().includes(query)) ||
          (g.notes && g.notes.toLowerCase().includes(query))
      );
    }

    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "name":
          comparison = a.name.localeCompare(b.name);
          break;
        case "phone":
          comparison = (a.phone || "").localeCompare(b.phone || "");
          break;
        case "is_active":
          comparison = (a.is_active === b.is_active) ? 0 : a.is_active ? -1 : 1;
          break;
        case "monthly_rent":
          comparison = a.monthly_rent - b.monthly_rent;
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

  const activeGodowns = computed(() => {
    return godowns.value.filter((g) => g.is_active);
  });

  const totalGodowns = computed(() => godowns.value.length);

  // ============= ACTIONS =============

  const fetchGodowns = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Godown[]>>("/godowns");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      godowns.value = res.data.data;
      return godowns.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch godowns");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchActiveGodowns = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Godown[]>>("/godowns", {
        params: { active: true },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      godowns.value = res.data.data;
      return godowns.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch active godowns");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const searchGodowns = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Godown[]>>("/godowns", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      godowns.value = res.data.data;
      return godowns.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search godowns");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const createGodown = async (payload: CreateGodownPayload): Promise<Godown | null> => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<Godown>>("/godowns", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      godowns.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create godown");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateGodown = async (id: number, payload: UpdateGodownPayload): Promise<Godown | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Godown>>(`/godowns/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = godowns.value.findIndex((g) => g.id === id);
      if (index !== -1) {
        godowns.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update godown");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const toggleGodownActive = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Godown>>(`/godowns/${id}/toggle-active`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = godowns.value.findIndex((g) => g.id === id);
      if (index !== -1) {
        godowns.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to toggle godown status");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const deleteGodown = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/godowns/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      godowns.value = godowns.value.filter((g) => g.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete godown");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: GodownSortField) => {
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
  const getGodownName = (id: number): string => {
    const godown = godowns.value.find((g) => g.id === id);
    return godown ? godown.name : `Godown #${id}`;
  };

  const getGodownById = (id: number): Godown | undefined => {
    return godowns.value.find((g) => g.id === id);
  };

  const getTotalMonthlyRent = (): number => {
    return godowns.value
      .filter((g) => g.is_active)
      .reduce((sum, g) => sum + g.monthly_rent, 0);
  };

  return {
    godowns,
    searchQuery,
    sortField,
    sortDirection,

    filteredGodowns,
    activeGodowns,
    totalGodowns,

    fetchGodowns,
    fetchActiveGodowns,
    searchGodowns,

    createGodown,
    updateGodown,
    toggleGodownActive,
    deleteGodown,

    setSort,
    setSearchQuery,
    clearSearch,

    getGodownName,
    getGodownById,
    getTotalMonthlyRent,
  };
});