// src/stores/transports.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Transport,
  CreateTransportPayload,
  UpdateTransportPayload,
  UpdateCustomerPaymentPayload,  // NEW
  DeliveryType,
  TransportSortField,
  SortDirection,
} from "@/types/transport";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useCustomersStore } from "./customers";
import { useUsersStore } from "./users";
import { useVehiclesStore } from "./vehicles";
import { useExpensesStore } from "./expenses";

export const useTransportsStore = defineStore("transports", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const transports = ref<Transport[]>([]);
  const searchQuery = ref("");
  const sortField = ref<TransportSortField>("transport_date");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredTransports = computed(() => {
    let result = [...transports.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const customersStore = useCustomersStore();
      const usersStore = useUsersStore();
      result = result.filter(
        (transport) =>
          transport.from_location.toLowerCase().includes(query) ||
          (transport.to_location && transport.to_location.toLowerCase().includes(query)) ||
          (transport.notes && transport.notes.toLowerCase().includes(query)) ||
          String(transport.vehicle_quantity).includes(query) ||
          (transport.delivery_type && transport.delivery_type.toLowerCase().includes(query)) ||
          String(transport.office_commission_amount).includes(query) ||
          String(transport.customer_total_paid).includes(query) ||  // NEW
          (transport.customer_id && customersStore.getCustomerName(transport.customer_id).toLowerCase().includes(query)) ||
          usersStore.getUserName(transport.user_id).toLowerCase().includes(query)
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "customer_id":
          comparison = (a.customer_id || 0) - (b.customer_id || 0);
          break;
        case "from_location":
          comparison = a.from_location.localeCompare(b.from_location);
          break;
        case "to_location":
          comparison = (a.to_location || "").localeCompare(b.to_location || "");
          break;
        case "vehicle_quantity":
          comparison = a.vehicle_quantity - b.vehicle_quantity;
          break;
        case "delivery_type":
          comparison = (a.delivery_type || "").localeCompare(b.delivery_type || "");
          break;
        case "transport_date":
          comparison = new Date(a.transport_date).getTime() - new Date(b.transport_date).getTime();
          break;
        case "office_commission_amount":
          comparison = a.office_commission_amount - b.office_commission_amount;
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

  const totalTransports = computed(() => transports.value.length);

  const totalCommission = computed(() => {
    return transports.value.reduce((sum, transport) => sum + transport.office_commission_amount, 0);
  });

  const totalCustomerPaid = computed(() => {  // NEW
    return transports.value.reduce((sum, transport) => sum + transport.customer_total_paid, 0);
  });

  // ============= EXPENSE HELPERS =============

  const generateExpenseTitle = (transport: Transport): string => {
    let title = `Transport Commission: #${transport.id}`;
    if (transport.from_location) {
      title = `${title} - From: ${transport.from_location}`;
    }
    if (transport.to_location) {
      title = `${title} To: ${transport.to_location}`;
    }
    if (transport.customer_id) {
      const customersStore = useCustomersStore();
      const customerName = customersStore.getCustomerName(transport.customer_id);
      title = `${title} (${customerName})`;
    }
    return title.slice(0, 255);
  };

  const generateExpenseNotes = (transport: Transport): string => {
    let notes = `Transport record #${transport.id}`;
    if (transport.notes) {
      notes = `${notes}: ${transport.notes}`;
    }
    return notes;
  };

  const findAssociatedExpense = async (transportId: number): Promise<number | null> => {
    const expensesStore = useExpensesStore();

    if (expensesStore.expenses.length === 0) {
      await expensesStore.fetchExpenses();
    }

    const matchingExpense = expensesStore.expenses.find((e) =>
      e.notes && e.notes.includes(`Transport record #${transportId}`)
    );

    return matchingExpense ? matchingExpense.id : null;
  };

  // ============= ACTIONS =============

  const fetchTransports = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Transport[]>>("/transports");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      transports.value = res.data.data;
      return transports.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch transports");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const searchTransports = async (query: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Transport[]>>("/transports", {
        params: { search: query },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      transports.value = res.data.data;
      return transports.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to search transports");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchTransportsByCustomer = async (customerId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Transport[]>>("/transports", {
        params: { customer_id: customerId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      transports.value = res.data.data;
      return transports.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch customer transports");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchTransportsByUser = async (userId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Transport[]>>("/transports", {
        params: { user_id: userId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      transports.value = res.data.data;
      return transports.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user transports");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchTransportsByDateRange = async (startDate: string, endDate: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Transport[]>>("/transports", {
        params: { start_date: startDate, end_date: endDate },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      transports.value = res.data.data;
      return transports.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch transports by date");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const createTransport = async (payload: CreateTransportPayload): Promise<Transport | null> => {
    displayLoader();
    try {
      if (!payload.from_location) {
        push.error("From location is required");
        return null;
      }

      if (payload.vehicle_quantity < 0) {
        push.error("Vehicle quantity cannot be negative");
        return null;
      }

      const res = await api.post<ApiResponse<Transport>>("/transports", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const newTransport = res.data.data;
      transports.value.push(newTransport);

      // Create expense for this transport if commission > 0
      const expensesStore = useExpensesStore();

      if (newTransport.office_commission_amount > 0) {
        const expensePayload = {
          title: generateExpenseTitle(newTransport),
          amount: newTransport.office_commission_amount,
          expense_date: newTransport.transport_date,
          notes: generateExpenseNotes(newTransport),
        };

        const expense = await expensesStore.createExpense(expensePayload);

        if (expense) {
          push.success(`${res.data.message} (Expense record created for commission)`);
        } else {
          push.warning(`${res.data.message} but expense creation failed`);
        }
      } else {
        push.success(res.data.message);
      }

      return newTransport;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create transport");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateTransport = async (id: number, payload: UpdateTransportPayload): Promise<Transport | null> => {
    displayLoader();
    try {
      const oldTransport = transports.value.find((t) => t.id === id);

      const res = await api.patch<ApiResponse<Transport>>(`/transports/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const updatedTransport = res.data.data;
      const index = transports.value.findIndex((transport) => transport.id === id);
      if (index !== -1) {
        transports.value[index] = updatedTransport;
      }

      // Find associated expense
      const expenseId = await findAssociatedExpense(id);
      const expensesStore = useExpensesStore();

      if (expenseId) {
        const expenseUpdatePayload: {
          title?: string;
          amount?: number;
          expense_date?: string;
          notes?: string;
        } = {};

        if (oldTransport && oldTransport.office_commission_amount !== updatedTransport.office_commission_amount) {
          if (updatedTransport.office_commission_amount > 0) {
            expenseUpdatePayload.amount = updatedTransport.office_commission_amount;
          } else {
            const deleted = await expensesStore.deleteExpense(expenseId);
            if (deleted) {
              push.success(`${res.data.message} (Associated expense deleted - commission set to 0)`);
            } else {
              push.warning(`${res.data.message} but expense deletion failed`);
            }
            push.success(res.data.message);
            return updatedTransport;
          }
        }

        if (oldTransport &&
          (oldTransport.from_location !== updatedTransport.from_location ||
            oldTransport.to_location !== updatedTransport.to_location)) {
          expenseUpdatePayload.title = generateExpenseTitle(updatedTransport);
        }

        if (oldTransport && oldTransport.transport_date !== updatedTransport.transport_date) {
          expenseUpdatePayload.expense_date = updatedTransport.transport_date;
        }

        if (oldTransport && oldTransport.notes !== updatedTransport.notes) {
          expenseUpdatePayload.notes = generateExpenseNotes(updatedTransport);
        }

        if (Object.keys(expenseUpdatePayload).length > 0) {
          const updatedExpense = await expensesStore.updateExpense(expenseId, expenseUpdatePayload);
          if (updatedExpense) {
            push.success(`${res.data.message} (Associated expense updated)`);
          } else {
            push.warning(`${res.data.message} but expense update failed`);
          }
        } else {
          push.success(res.data.message);
        }
      } else {
        if (updatedTransport.office_commission_amount > 0) {
          const expensePayload = {
            title: generateExpenseTitle(updatedTransport),
            amount: updatedTransport.office_commission_amount,
            expense_date: updatedTransport.transport_date,
            notes: generateExpenseNotes(updatedTransport),
          };

          const expense = await expensesStore.createExpense(expensePayload);
          if (expense) {
            push.success(`${res.data.message} (Expense record created for commission)`);
          } else {
            push.success(res.data.message);
          }
        } else {
          push.success(res.data.message);
        }
      }

      return updatedTransport;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update transport");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // NEW: Update customer payment for a transport
  const updateCustomerPayment = async (id: number, payload: UpdateCustomerPaymentPayload): Promise<Transport | null> => {
    displayLoader();
    try {
      if (payload.customer_total_paid < 0) {
        push.error("Customer total paid cannot be negative");
        return null;
      }

      const res = await api.patch<ApiResponse<Transport>>(`/transports/${id}/customer-payment`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const index = transports.value.findIndex((transport) => transport.id === id);
      if (index !== -1) {
        transports.value[index] = res.data.data;
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

  const deleteTransport = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const expenseId = await findAssociatedExpense(id);

      const res = await api.delete<ApiResponse<null>>(`/transports/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      transports.value = transports.value.filter((transport) => transport.id !== id);

      if (expenseId) {
        const expensesStore = useExpensesStore();
        const deleted = await expensesStore.deleteExpense(expenseId);
        if (deleted) {
          push.success(`${res.data.message} (Associated expense deleted)`);
        } else {
          push.warning(`${res.data.message} but expense deletion failed`);
        }
      } else {
        push.success(res.data.message);
      }

      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete transport");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const deleteTransportsByCustomer = async (customerId: number): Promise<boolean> => {
    displayLoader();
    try {
      const customerTransports = transports.value.filter((t) => t.customer_id === customerId);

      for (const transport of customerTransports) {
        const expenseId = await findAssociatedExpense(transport.id);

        const res = await api.delete<ApiResponse<null>>(`/transports/${transport.id}`);
        if (!res.data.success) {
          push.error(`Failed to delete transport ${transport.id}`);
          return false;
        }

        if (expenseId) {
          const expensesStore = useExpensesStore();
          await expensesStore.deleteExpense(expenseId);
        }
      }

      transports.value = transports.value.filter((t) => t.customer_id !== customerId);
      push.success("All customer transports deleted");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete customer transports");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: TransportSortField) => {
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

  const getTransportById = (id: number): Transport | undefined => {
    return transports.value.find((t) => t.id === id);
  };

  const getTransportsByCustomerId = (customerId: number): Transport[] => {
    return transports.value.filter((transport) => transport.customer_id === customerId);
  };

  const getCustomerNameForTransport = (transport: Transport): string => {
    if (!transport.customer_id) return "N/A";
    const customersStore = useCustomersStore();
    return customersStore.getCustomerName(transport.customer_id);
  };

  const formatDeliveryType = (deliveryType: DeliveryType | null): string => {
    if (!deliveryType) return "N/A";
    return deliveryType.charAt(0).toUpperCase() + deliveryType.slice(1);
  };

  const formatTransportDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const hasVehicles = (transportId: number): boolean => {
    const vehiclesStore = useVehiclesStore();
    const vehicles = vehiclesStore.getVehiclesByTransportId(transportId);
    return vehicles.length > 0;
  };

  const getTransportExpenseId = async (transportId: number): Promise<number | null> => {
    return await findAssociatedExpense(transportId);
  };

  return {
    transports,
    searchQuery,
    sortField,
    sortDirection,

    filteredTransports,
    totalTransports,
    totalCommission,
    totalCustomerPaid,  // NEW

    fetchTransports,
    searchTransports,
    fetchTransportsByCustomer,
    fetchTransportsByUser,
    fetchTransportsByDateRange,

    createTransport,
    updateTransport,
    updateCustomerPayment,  // NEW
    deleteTransport,
    deleteTransportsByCustomer,

    setSort,
    setSearchQuery,
    clearSearch,

    getTransportById,
    getTransportsByCustomerId,
    getCustomerNameForTransport,
    formatDeliveryType,
    formatTransportDate,
    hasVehicles,
    getTransportExpenseId,
  };
});