// src/stores/deliveries.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Delivery,
  CreateDeliveryPayload,
  UpdateDeliveryPayload,
  DeliverySortField,
  SortDirection,
} from "@/types/delivery";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useCustomersStore } from "./customers";
import { useUsersStore } from "./users";
import { useDeliveryItemsStore } from "./deliveryItems";

export const useDeliveriesStore = defineStore("deliveries", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const deliveries = ref<Delivery[]>([]);
  const searchQuery = ref("");
  const sortField = ref<DeliverySortField>("delivery_date");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredDeliveries = computed(() => {
    let result = [...deliveries.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const customersStore = useCustomersStore();
      const usersStore = useUsersStore();
      result = result.filter(
        (delivery) =>
          (delivery.receiver_name && delivery.receiver_name.toLowerCase().includes(query)) ||
          (delivery.receiver_phone && delivery.receiver_phone.toLowerCase().includes(query)) ||
          (delivery.from_location && delivery.from_location.toLowerCase().includes(query)) ||
          (delivery.to_location && delivery.to_location.toLowerCase().includes(query)) ||
          (delivery.notes && delivery.notes.toLowerCase().includes(query)) ||
          (delivery.customer_id && customersStore.getCustomerName(delivery.customer_id).toLowerCase().includes(query)) ||
          usersStore.getUserName(delivery.user_id).toLowerCase().includes(query)
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "customer_id":
          comparison = (a.customer_id || 0) - (b.customer_id || 0);
          break;
        case "delivery_date":
          comparison = new Date(a.delivery_date).getTime() - new Date(b.delivery_date).getTime();
          break;
        case "receiver_name":
          comparison = (a.receiver_name || "").localeCompare(b.receiver_name || "");
          break;
        case "from_location":
          comparison = (a.from_location || "").localeCompare(b.from_location || "");
          break;
        case "to_location":
          comparison = (a.to_location || "").localeCompare(b.to_location || "");
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

  const totalDeliveries = computed(() => deliveries.value.length);

  // ============= ACTIONS =============

  // GET ALL DELIVERIES
  const fetchDeliveries = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Delivery[]>>("/deliveries");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveries.value = res.data.data;
      return deliveries.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch deliveries");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH DELIVERIES
  const searchDeliveries = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Delivery[]>>("/deliveries", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveries.value = res.data.data;
      return deliveries.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search deliveries");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERIES BY CUSTOMER
  const fetchDeliveriesByCustomer = async (customerId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Delivery[]>>("/deliveries", {
        params: { customer_id: customerId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveries.value = res.data.data;
      return deliveries.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch customer deliveries");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERIES BY USER
  const fetchDeliveriesByUser = async (userId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Delivery[]>>("/deliveries", {
        params: { user_id: userId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveries.value = res.data.data;
      return deliveries.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user deliveries");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERIES BY DATE RANGE
  const fetchDeliveriesByDateRange = async (startDate: string, endDate: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Delivery[]>>("/deliveries", {
        params: { start_date: startDate, end_date: endDate },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveries.value = res.data.data;
      return deliveries.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch deliveries by date");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE DELIVERY
  const createDelivery = async (payload: CreateDeliveryPayload): Promise<Delivery | null> => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<Delivery>>("/deliveries", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      deliveries.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create delivery");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE DELIVERY
  const updateDelivery = async (id: number, payload: UpdateDeliveryPayload): Promise<Delivery | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Delivery>>(`/deliveries/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = deliveries.value.findIndex((delivery) => delivery.id === id);
      if (index !== -1) {
        deliveries.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update delivery");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE DELIVERY
  const deleteDelivery = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/deliveries/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      deliveries.value = deliveries.value.filter((delivery) => delivery.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete delivery");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: DeliverySortField) => {
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

  // Get delivery by ID
  const getDeliveryById = (id: number): Delivery | undefined => {
    return deliveries.value.find((d) => d.id === id);
  };

  // Get deliveries by customer ID
  const getDeliveriesByCustomerId = (customerId: number): Delivery[] => {
    return deliveries.value.filter((delivery) => delivery.customer_id === customerId);
  };

  // Get customer name for delivery
  const getCustomerNameForDelivery = (delivery: Delivery): string => {
    if (!delivery.customer_id) return "N/A";
    const customersStore = useCustomersStore();
    return customersStore.getCustomerName(delivery.customer_id);
  };

  // Format delivery date
  const formatDeliveryDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // Check if delivery has items
  const hasItems = (deliveryId: number): boolean => {
    const deliveryItemsStore = useDeliveryItemsStore();
    const items = deliveryItemsStore.getDeliveryItemsByDeliveryId(deliveryId);
    return items.length > 0;
  };

  return {
    // State
    deliveries,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredDeliveries,
    totalDeliveries,

    // Fetch
    fetchDeliveries,
    searchDeliveries,
    fetchDeliveriesByCustomer,
    fetchDeliveriesByUser,
    fetchDeliveriesByDateRange,

    // CRUD
    createDelivery,
    updateDelivery,
    deleteDelivery,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getDeliveryById,
    getDeliveriesByCustomerId,
    getCustomerNameForDelivery,
    formatDeliveryDate,
    hasItems,
  };
});