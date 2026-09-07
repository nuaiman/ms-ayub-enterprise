// src/stores/lots.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Lot,
  CreateLotPayload,
  UpdateLotPayload,
  CustomerChargeType,
  MajhiBillType,
  LotSortField,
  SortDirection,
} from "@/types/lot";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useItemsStore } from "./items";

export const useLotsStore = defineStore("lots", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const lots = ref<Lot[]>([]);
  const searchQuery = ref("");
  const sortField = ref<LotSortField>("lot_number");
  const sortDirection = ref<SortDirection>("asc");

  // ============= COMPUTED =============
  const filteredLots = computed(() => {
    let result = [...lots.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const itemsStore = useItemsStore();
      result = result.filter(
        (lot) =>
          String(lot.lot_number).includes(query) ||
          lot.customer_charge_type.toLowerCase().includes(query) ||
          lot.majhi_bill_type.toLowerCase().includes(query) ||
          (lot.notes && lot.notes.toLowerCase().includes(query)) ||
          itemsStore.getItemName(lot.item_id).toLowerCase().includes(query)
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "item_id":
          comparison = a.item_id - b.item_id;
          break;
        case "lot_number":
          comparison = a.lot_number - b.lot_number;
          break;
        case "customer_charge_type":
          comparison = a.customer_charge_type.localeCompare(b.customer_charge_type);
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

  const activeLots = computed(() => {
    return lots.value.filter((lot) => lot.is_active);
  });

  const totalLots = computed(() => lots.value.length);

  // ============= ACTIONS =============

  // GET ALL LOTS
  const fetchLots = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Lot[]>>("/lots");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      lots.value = res.data.data;
      return lots.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch lots");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // SEARCH LOTS
  const searchLots = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Lot[]>>("/lots", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      lots.value = res.data.data;
      return lots.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search lots");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET LOTS BY ITEM
  const fetchLotsByItem = async (itemId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Lot[]>>("/lots", {
        params: { item_id: itemId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      lots.value = res.data.data;
      return lots.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch item lots");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET ACTIVE LOTS
  const fetchActiveLots = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Lot[]>>("/lots", {
        params: { active: true },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      lots.value = res.data.data;
      return lots.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch active lots");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE LOT
  const createLot = async (payload: CreateLotPayload): Promise<Lot | null> => {
    displayLoader();
    try {
      // Validate: item_id is required
      if (!payload.item_id) {
        push.error("Item is required");
        return null;
      }

      // Validate: lot_number is required
      if (!payload.lot_number) {
        push.error("Lot number is required");
        return null;
      }

      // Validate: customer_charge_type
      if (payload.customer_charge_type !== 'weight' && payload.customer_charge_type !== 'quantity') {
        push.error("customer_charge_type must be 'weight' or 'quantity'");
        return null;
      }

      // Validate: majhi_bill_type
      if (payload.majhi_bill_type !== 'weight' && payload.majhi_bill_type !== 'quantity' && payload.majhi_bill_type !== 'job') {
        push.error("majhi_bill_type must be 'weight', 'quantity', or 'job'");
        return null;
      }

      // Set defaults for new fields if not provided
      const requestPayload = {
        ...payload,
        customer_last_paid_through: payload.customer_last_paid_through || null,
        customer_last_paid_amount: payload.customer_last_paid_amount || 0,
        customer_paid_unload_amount: payload.customer_paid_unload_amount || 0,
        majhi_total_paid: payload.majhi_total_paid || 0,
      };

      const res = await api.post<ApiResponse<Lot>>("/lots", requestPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      lots.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create lot");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE LOT
  const updateLot = async (id: number, payload: UpdateLotPayload): Promise<Lot | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Lot>>(`/lots/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = lots.value.findIndex((lot) => lot.id === id);
      if (index !== -1) {
        lots.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update lot");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE LOT CUSTOMER PAYMENT (storage bills)
  const updateLotCustomerPayment = async (id: number, paidThrough: string | null, paidAmount: number): Promise<Lot | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Lot>>(`/lots/${id}/customer-payment`, {
        customer_last_paid_through: paidThrough,
        customer_last_paid_amount: paidAmount,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = lots.value.findIndex((lot) => lot.id === id);
      if (index !== -1) {
        lots.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update customer payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE LOT CUSTOMER UNLOAD PAYMENT (one-time unload bills)
  const updateLotCustomerUnloadPayment = async (id: number, paidAmount: number, paidThrough?: string | null): Promise<Lot | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Lot>>(`/lots/${id}/customer-unload-payment`, {
        customer_paid_unload_amount: paidAmount,
        paid_through: paidThrough || null,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = lots.value.findIndex((lot) => lot.id === id);
      if (index !== -1) {
        lots.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update unload payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE LOT MAJHI PAYMENT
  const updateLotMajhiPayment = async (id: number, totalPaid: number): Promise<Lot | null> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Lot>>(`/lots/${id}/majhi-payment`, {
        majhi_total_paid: totalPaid,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      const index = lots.value.findIndex((lot) => lot.id === id);
      if (index !== -1) {
        lots.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update majhi payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // TOGGLE LOT ACTIVE
  const toggleLotActive = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Lot>>(`/lots/${id}/toggle-active`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = lots.value.findIndex((lot) => lot.id === id);
      if (index !== -1) {
        lots.value[index] = res.data.data;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to toggle lot status");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // DELETE LOT
  const deleteLot = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/lots/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      lots.value = lots.value.filter((lot) => lot.id !== id);
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete lot");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: LotSortField) => {
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

  // Get lot display name
  const getLotDisplayName = (lot: Lot): string => {
    const itemsStore = useItemsStore();
    const itemName = itemsStore.getItemName(lot.item_id);
    return `${itemName} - Lot #${lot.lot_number}`;
  };

  // Get lot name by ID
  const getLotName = (id: number): string => {
    const lot = lots.value.find((l) => l.id === id);
    if (!lot) return `Lot #${id}`;
    return getLotDisplayName(lot);
  };

  // Get lot by ID
  const getLotById = (id: number): Lot | undefined => {
    return lots.value.find((l) => l.id === id);
  };

  // Get lots by item ID
  const getLotsByItemId = (itemId: number): Lot[] => {
    return lots.value.filter((lot) => lot.item_id === itemId);
  };

  // Get unload bill amount for a lot
  const getUnloadBillAmount = (lot: Lot): number => {
    // This will be calculated in the customerLotBills store
    // Returns 0 if no unload_rate
    return lot.unload_rate || 0;
  };

  // Check if unload bill is paid
  const isUnloadBillPaid = (lot: Lot): boolean => {
    const billAmount = getUnloadBillAmount(lot);
    if (billAmount === 0) return true;
    return (lot.customer_paid_unload_amount || 0) >= billAmount;
  };

  // Format types for display
  const formatCustomerChargeType = (type: CustomerChargeType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  const formatMajhiBillType = (type: MajhiBillType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  // Get majhi name for a lot
  const getMajhiNameForLot = (lot: Lot): string => {
    if (!lot.majhi_id) return "No Majhi Assigned";
    // This will be resolved by the majhis store
    return `Majhi #${lot.majhi_id}`;
  };

  return {
    // State
    lots,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredLots,
    activeLots,
    totalLots,

    // Fetch
    fetchLots,
    searchLots,
    fetchLotsByItem,
    fetchActiveLots,

    // CRUD
    createLot,
    updateLot,
    updateLotCustomerPayment,
    updateLotCustomerUnloadPayment,
    updateLotMajhiPayment,
    toggleLotActive,
    deleteLot,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getLotDisplayName,
    getLotName,
    getLotById,
    getLotsByItemId,
    getUnloadBillAmount,
    isUnloadBillPaid,
    formatCustomerChargeType,
    formatMajhiBillType,
    getMajhiNameForLot,
  };
});