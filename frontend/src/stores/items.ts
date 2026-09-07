// src/stores/items.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Item,
  CreateItemPayload,
  UpdateItemPayload,
  ItemSortField,
  SortDirection,
} from "@/types/item";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useCustomersStore } from "./customers";
import { useAuthStore } from "./auth";

export const useItemsStore = defineStore("items", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const items = ref<Item[]>([]);
  const searchQuery = ref("");
  const sortField = ref<ItemSortField>("product_name");
  const sortDirection = ref<SortDirection>("asc");

  // ============= COMPUTED =============
  const filteredItems = computed(() => {
    let result = [...items.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(
        (item) =>
          (item.product_name && item.product_name.toLowerCase().includes(query)) ||
          (item.category && item.category.toLowerCase().includes(query)) ||
          (item.notes && item.notes.toLowerCase().includes(query))
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "product_name":
          comparison = (a.product_name || "").localeCompare(b.product_name || "");
          break;
        case "category":
          comparison = (a.category || "").localeCompare(b.category || "");
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

  const activeItems = computed(() => {
    return items.value.filter((item) => item.is_active);
  });

  const totalItems = computed(() => items.value.length);

  // ============= ACTIONS =============

  // GET ALL ITEMS
  const fetchItems = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH ITEMS
  const searchItems = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET ITEMS BY USER
  const fetchItemsByUser = async (userId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items", {
        params: { user_id: userId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET ITEMS BY CUSTOMER
  const fetchItemsByCustomer = async (customerId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items", {
        params: { customer_id: customerId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch customer items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET ACTIVE ITEMS
  const fetchActiveItems = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items", {
        params: { active: true },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch active items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE ITEM
  const createItem = async (payload: CreateItemPayload): Promise<Item | null> => {
    displayLoader();
    try {
      // Validate: at least product_name or category is required
      if (!payload.product_name && !payload.category) {
        push.error("Either product name or category is required");
        return null;
      }

      const res = await api.post<ApiResponse<Item>>("/items", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      items.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE ITEM
  const updateItem = async (id: number, payload: UpdateItemPayload): Promise<Item | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Item>>(`/items/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = items.value.findIndex((item) => item.id === id);
      if (index !== -1) {
        items.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // TOGGLE ITEM ACTIVE
  const toggleItemActive = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Item>>(`/items/${id}/toggle-active`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = items.value.findIndex((item) => item.id === id);
      if (index !== -1) {
        items.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to toggle item status");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // DELETE ITEM
  const deleteItem = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/items/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      items.value = items.value.filter((item) => item.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete item");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: ItemSortField) => {
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

  // Get item display name (product name or category)
  const getItemDisplayName = (item: Item): string => {
    if (item.product_name) return item.product_name;
    if (item.category) return item.category;
    return `Item #${item.id}`;
  };

  // Get item name by ID
  const getItemName = (id: number): string => {
    const item = items.value.find((i) => i.id === id);
    if (!item) return `Item #${id}`;
    return getItemDisplayName(item);
  };

  // Get item by ID
  const getItemById = (id: number): Item | undefined => {
    return items.value.find((i) => i.id === id);
  };

  // Get customer name for item
  const getCustomerNameForItem = (item: Item): string => {
    if (!item.customer_id) return "N/A";
    const customersStore = useCustomersStore();
    return customersStore.getCustomerName(item.customer_id);
  };

  // Get current user's items (filter by logged in user)
  const getMyItems = computed(() => {
    const auth = useAuthStore();
    if (!auth.user) return [];
    return items.value.filter((item) => item.user_id === auth.user!.id);
  });

  return {
    // State
    items,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredItems,
    activeItems,
    totalItems,
    getMyItems,

    // Fetch
    fetchItems,
    searchItems,
    fetchItemsByUser,
    fetchItemsByCustomer,
    fetchActiveItems,

    // CRUD
    createItem,
    updateItem,
    toggleItemActive,
    deleteItem,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getItemDisplayName,
    getItemName,
    getItemById,
    getCustomerNameForItem,
  };
});