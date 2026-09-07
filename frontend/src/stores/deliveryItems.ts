// src/stores/deliveryItems.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  DeliveryItem,
  CreateDeliveryItemPayload,
  UpdateDeliveryItemPayload,
  DeliveryItemSortField,
  SortDirection,
} from "@/types/deliveryItem";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useItemsStore } from "./items";
import { useLotsStore } from "./lots";
import { useMajhisStore } from "./majhis";

export const useDeliveryItemsStore = defineStore("deliveryItems", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const deliveryItems = ref<DeliveryItem[]>([]);
  const searchQuery = ref("");
  const sortField = ref<DeliveryItemSortField>("created_at");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredDeliveryItems = computed(() => {
    let result = [...deliveryItems.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const itemsStore = useItemsStore();
      const lotsStore = useLotsStore();
      const majhisStore = useMajhisStore();
      result = result.filter(
        (item) =>
          String(item.quantity).includes(query) ||
          String(item.weight).includes(query) ||
          String(item.loading_rate).includes(query) ||
          String(item.majhi_cut).includes(query) ||
          (item.vehicle_number && item.vehicle_number.toLowerCase().includes(query)) ||
          (item.driver_number && item.driver_number.toLowerCase().includes(query)) ||
          (item.notes && item.notes.toLowerCase().includes(query)) ||
          itemsStore.getItemName(item.item_id).toLowerCase().includes(query) ||
          lotsStore.getLotName(item.lot_id).toLowerCase().includes(query) ||
          (item.majhi_id && majhisStore.getMajhiName(item.majhi_id).toLowerCase().includes(query))
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "delivery_id":
          comparison = a.delivery_id - b.delivery_id;
          break;
        case "store_id":
          comparison = a.store_id - b.store_id;
          break;
        case "item_id":
          comparison = a.item_id - b.item_id;
          break;
        case "lot_id":
          comparison = a.lot_id - b.lot_id;
          break;
        case "quantity":
          comparison = a.quantity - b.quantity;
          break;
        case "weight":
          comparison = a.weight - b.weight;
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

  const totalDeliveryItems = computed(() => deliveryItems.value.length);

  const totalQuantity = computed(() => {
    return deliveryItems.value.reduce((sum, item) => sum + item.quantity, 0);
  });

  const totalWeight = computed(() => {
    return deliveryItems.value.reduce((sum, item) => sum + item.weight, 0);
  });

  // ============= ACTIONS =============

  // GET ALL DELIVERY ITEMS
  const fetchDeliveryItems = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<DeliveryItem[]>>("/delivery-items");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveryItems.value = res.data.data;
      return deliveryItems.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch delivery items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERY ITEMS BY DELIVERY
  const fetchDeliveryItemsByDelivery = async (deliveryId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<DeliveryItem[]>>("/delivery-items", {
        params: { delivery_id: deliveryId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveryItems.value = res.data.data;
      return deliveryItems.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch delivery items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERY ITEMS BY STORE
  const fetchDeliveryItemsByStore = async (storeId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<DeliveryItem[]>>("/delivery-items", {
        params: { store_id: storeId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveryItems.value = res.data.data;
      return deliveryItems.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch store delivery items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET DELIVERY ITEMS BY LOT
  const fetchDeliveryItemsByLot = async (lotId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<DeliveryItem[]>>("/delivery-items", {
        params: { lot_id: lotId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      deliveryItems.value = res.data.data;
      return deliveryItems.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch lot delivery items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE DELIVERY ITEM
  const createDeliveryItem = async (payload: CreateDeliveryItemPayload): Promise<DeliveryItem | null> => {
    displayLoader();
    try {
      // Validate required fields
      if (!payload.delivery_id) {
        push.error("Delivery ID is required");
        return null;
      }
      if (!payload.store_id) {
        push.error("Store ID is required");
        return null;
      }
      if (!payload.item_id) {
        push.error("Item ID is required");
        return null;
      }
      if (!payload.lot_id) {
        push.error("Lot ID is required");
        return null;
      }
      if (payload.quantity <= 0 && payload.weight <= 0) {
        push.error("Either quantity or weight must be greater than 0");
        return null;
      }

      // Ensure defaults
      const requestPayload: CreateDeliveryItemPayload = {
        delivery_id: payload.delivery_id,
        store_id: payload.store_id,
        majhi_id: payload.majhi_id || null,
        item_id: payload.item_id,
        lot_id: payload.lot_id,
        vehicle_number: payload.vehicle_number || null,
        driver_number: payload.driver_number || null,
        quantity: payload.quantity || 0,
        quantity_unit: payload.quantity_unit || 'units',
        weight: payload.weight || 0,
        weight_unit: payload.weight_unit || 'kg',
        loading_rate: payload.loading_rate || 0,
        majhi_cut: payload.majhi_cut || 0,
        notes: payload.notes || null,
        // New fields - optional, backend will use defaults from lot if not provided
        customer_charge_type: payload.customer_charge_type,
        customer_paid_unload_amount: payload.customer_paid_unload_amount || 0,
        majhi_bill_type: payload.majhi_bill_type,
        majhi_total_paid: payload.majhi_total_paid || 0,
      };

      console.log('[DELIVERY_ITEMS] Creating with payload:', requestPayload);

      const res = await api.post<ApiResponse<DeliveryItem>>("/delivery-items", requestPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      deliveryItems.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      console.error('[DELIVERY_ITEMS] Create error:', err.response?.data);
      push.error(err.response?.data?.message || "Failed to create delivery item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE DELIVERY ITEM
  const updateDeliveryItem = async (id: number, payload: UpdateDeliveryItemPayload): Promise<DeliveryItem | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<DeliveryItem>>(`/delivery-items/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = deliveryItems.value.findIndex((item) => item.id === id);
      if (index !== -1) {
        deliveryItems.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update delivery item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE DELIVERY ITEM CUSTOMER UNLOAD PAYMENT (Hidden - for backend use)
  const updateDeliveryItemCustomerUnloadPayment = async (id: number, paidAmount: number): Promise<DeliveryItem | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<DeliveryItem>>(`/delivery-items/${id}/customer-unload-payment`, {
        customer_paid_unload_amount: paidAmount,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = deliveryItems.value.findIndex((item) => item.id === id);
      if (index !== -1) {
        deliveryItems.value[index] = res.data.data;
      }
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update customer unload payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE DELIVERY ITEM MAJHI PAYMENT (Hidden - for backend use)
  const updateDeliveryItemMajhiPayment = async (id: number, paidAmount: number): Promise<DeliveryItem | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<DeliveryItem>>(`/delivery-items/${id}/majhi-payment`, {
        majhi_total_paid: paidAmount,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = deliveryItems.value.findIndex((item) => item.id === id);
      if (index !== -1) {
        deliveryItems.value[index] = res.data.data;
      }
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update majhi payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // DELETE DELIVERY ITEM
  const deleteDeliveryItem = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/delivery-items/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      deliveryItems.value = deliveryItems.value.filter((item) => item.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete delivery item");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // DELETE DELIVERY ITEMS BY DELIVERY
  const deleteDeliveryItemsByDelivery = async (deliveryId: number): Promise<boolean> => {
    displayLoader();
    try {
      const items = deliveryItems.value.filter((item) => item.delivery_id === deliveryId);

      for (const item of items) {
        const res = await api.delete<ApiResponse<null>>(`/delivery-items/${item.id}`);
        if (!res.data.success) {
          push.error(`Failed to delete delivery item ${item.id}`);
          return false;
        }
      }

      deliveryItems.value = deliveryItems.value.filter((item) => item.delivery_id !== deliveryId);
      push.success("All delivery items deleted");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete delivery items");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: DeliveryItemSortField) => {
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

  const getDeliveryItemById = (id: number): DeliveryItem | undefined => {
    return deliveryItems.value.find((item) => item.id === id);
  };

  const getDeliveryItemsByDeliveryId = (deliveryId: number): DeliveryItem[] => {
    return deliveryItems.value.filter((item) => item.delivery_id === deliveryId);
  };

  const getTotalQuantityByDelivery = (deliveryId: number): number => {
    return deliveryItems.value
      .filter((item) => item.delivery_id === deliveryId)
      .reduce((sum, item) => sum + item.quantity, 0);
  };

  const getTotalWeightByDelivery = (deliveryId: number): number => {
    return deliveryItems.value
      .filter((item) => item.delivery_id === deliveryId)
      .reduce((sum, item) => sum + item.weight, 0);
  };

  const getItemNameForDeliveryItem = (item: DeliveryItem): string => {
    const itemsStore = useItemsStore();
    return itemsStore.getItemName(item.item_id);
  };

  const getLotNameForDeliveryItem = (item: DeliveryItem): string => {
    const lotsStore = useLotsStore();
    return lotsStore.getLotName(item.lot_id);
  };

  const getMajhiNameForDeliveryItem = (item: DeliveryItem): string => {
    if (!item.majhi_id) return "N/A";
    const majhisStore = useMajhisStore();
    return majhisStore.getMajhiName(item.majhi_id);
  };

  // Get charge type label
  const getChargeTypeLabel = (chargeType: 'weight' | 'quantity'): string => {
    return chargeType === 'weight' ? 'Weight' : 'Quantity';
  };

  // Get majhi bill type label
  const getMajhiBillTypeLabel = (billType: 'weight' | 'quantity' | 'job'): string => {
    const labels = {
      weight: 'Weight',
      quantity: 'Quantity',
      job: 'Job (Fixed)',
    };
    return labels[billType] || billType;
  };

  return {
    // State
    deliveryItems,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredDeliveryItems,
    totalDeliveryItems,
    totalQuantity,
    totalWeight,

    // Fetch
    fetchDeliveryItems,
    fetchDeliveryItemsByDelivery,
    fetchDeliveryItemsByStore,
    fetchDeliveryItemsByLot,

    // CRUD
    createDeliveryItem,
    updateDeliveryItem,
    updateDeliveryItemCustomerUnloadPayment,
    updateDeliveryItemMajhiPayment,
    deleteDeliveryItem,
    deleteDeliveryItemsByDelivery,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getDeliveryItemById,
    getDeliveryItemsByDeliveryId,
    getTotalQuantityByDelivery,
    getTotalWeightByDelivery,
    getItemNameForDeliveryItem,
    getLotNameForDeliveryItem,
    getMajhiNameForDeliveryItem,
    getChargeTypeLabel,
    getMajhiBillTypeLabel,
  };
});