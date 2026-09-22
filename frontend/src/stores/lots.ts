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
import { useCustomersStore } from "./customers";

export const useLotsStore = defineStore("lots", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  const lots = ref<Lot[]>([]);
  const searchQuery = ref("");
  const sortField = ref<LotSortField>("lot_number");
  const sortDirection = ref<SortDirection>("asc");

  const filteredLots = computed(() => {
    let result = [...lots.value];

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const customersStore = useCustomersStore();
      result = result.filter(
        (lot) =>
          String(lot.lot_number).includes(query) ||
          (lot.product_name && lot.product_name.toLowerCase().includes(query)) ||
          (lot.category && lot.category.toLowerCase().includes(query)) ||
          lot.customer_charge_type.toLowerCase().includes(query) ||
          lot.majhi_bill_type.toLowerCase().includes(query) ||
          (lot.notes && lot.notes.toLowerCase().includes(query)) ||
          (lot.customer_id
            ? customersStore.getCustomerName(lot.customer_id).toLowerCase().includes(query)
            : false)
      );
    }

    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "customer_id":
          comparison = (a.customer_id || 0) - (b.customer_id || 0);
          break;
        case "lot_number":
          comparison = a.lot_number - b.lot_number;
          break;
        case "customer_charge_type":
          comparison = a.customer_charge_type.localeCompare(b.customer_charge_type);
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

  const activeLots = computed(() => lots.value.filter((lot) => lot.is_active));
  const totalLots = computed(() => lots.value.length);

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

  const fetchLotsByCustomer = async (customerId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Lot[]>>("/lots", {
        params: { customer_id: customerId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      lots.value = res.data.data;
      return lots.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch customer lots");
      return [];
    } finally {
      destroyLoader();
    }
  };

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

  const createLot = async (payload: CreateLotPayload): Promise<Lot | null> => {
    displayLoader();
    try {
      if (!payload.lot_number) {
        push.error("Lot number is required");
        return null;
      }

      if (payload.customer_charge_type !== 'weight' && payload.customer_charge_type !== 'quantity') {
        push.error("customer_charge_type must be 'weight' or 'quantity'");
        return null;
      }

      if (payload.majhi_bill_type !== 'weight' && payload.majhi_bill_type !== 'quantity' && payload.majhi_bill_type !== 'job') {
        push.error("majhi_bill_type must be 'weight', 'quantity', or 'job'");
        return null;
      }

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

  const setSort = (field: LotSortField) => {
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

  // Display name for a lot: product_name → category → Lot #<number>
  const getLotDisplayName = (lot: Lot): string => {
    if (lot.product_name) return lot.product_name;
    if (lot.category) return lot.category;
    return `Lot #${lot.lot_number}`;
  };

  const getLotName = (id: number): string => {
    const lot = lots.value.find((l) => l.id === id);
    if (!lot) return `Lot #${id}`;
    return getLotDisplayName(lot);
  };

  const getLotById = (id: number): Lot | undefined => {
    return lots.value.find((l) => l.id === id);
  };

  const getLotsByCustomerId = (customerId: number): Lot[] => {
    return lots.value.filter((lot) => lot.customer_id === customerId);
  };

  const getUnloadBillAmount = (lot: Lot): number => {
    return lot.unload_rate || 0;
  };

  const isUnloadBillPaid = (lot: Lot): boolean => {
    const billAmount = getUnloadBillAmount(lot);
    if (billAmount === 0) return true;
    return (lot.customer_paid_unload_amount || 0) >= billAmount;
  };

  const formatCustomerChargeType = (type: CustomerChargeType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  const formatMajhiBillType = (type: MajhiBillType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  return {
    lots,
    searchQuery,
    sortField,
    sortDirection,

    filteredLots,
    activeLots,
    totalLots,

    fetchLots,
    searchLots,
    fetchLotsByCustomer,
    fetchActiveLots,

    createLot,
    updateLot,
    updateLotCustomerPayment,
    updateLotCustomerUnloadPayment,
    updateLotMajhiPayment,
    toggleLotActive,
    deleteLot,

    setSort,
    setSearchQuery,
    clearSearch,

    getLotDisplayName,
    getLotName,
    getLotById,
    getLotsByCustomerId,
    getUnloadBillAmount,
    isUnloadBillPaid,
    formatCustomerChargeType,
    formatMajhiBillType,
  };
});