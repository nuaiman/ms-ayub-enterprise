// src/stores/lots.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Lot,
  CreateLotPayload,
  UpdateLotPayload,
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
          lot.lot_number.toLowerCase().includes(query) ||
          lot.product_name.toLowerCase().includes(query) ||
          customersStore.getCustomerName(lot.customer_id).toLowerCase().includes(query)
      );
    }

    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "customer_id":
          comparison = a.customer_id - b.customer_id;
          break;
        case "lot_number":
          comparison = a.lot_number.localeCompare(b.lot_number);
          break;
        case "product_name":
          comparison = a.product_name.localeCompare(b.product_name);
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

  const createLot = async (payload: CreateLotPayload): Promise<Lot | null> => {
    displayLoader();
    try {
      if (!payload.customer_id) {
        push.error("Customer is required");
        return null;
      }
      if (!payload.lot_number.trim()) {
        push.error("Lot number is required");
        return null;
      }
      if (!payload.product_name.trim()) {
        push.error("Product name is required");
        return null;
      }

      const requestPayload: CreateLotPayload = {
        customer_id: payload.customer_id,
        lot_number: payload.lot_number.trim(),
        product_name: payload.product_name.trim(),
        weight_unit: payload.weight_unit || "kg",
        quantity_unit: payload.quantity_unit || "units",
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

  // Display name for a lot: product_name -> "Lot <lot_number>"
  const getLotDisplayName = (lot: Lot): string => {
    if (lot.product_name) return lot.product_name;
    return `Lot ${lot.lot_number}`;
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

  return {
    lots,
    searchQuery,
    sortField,
    sortDirection,

    filteredLots,
    totalLots,

    fetchLots,
    searchLots,
    fetchLotsByCustomer,

    createLot,
    updateLot,
    deleteLot,

    setSort,
    setSearchQuery,
    clearSearch,

    getLotDisplayName,
    getLotName,
    getLotById,
    getLotsByCustomerId,
  };
});