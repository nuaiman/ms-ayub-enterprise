// src/stores/vehicles.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Vehicle,
  CreateVehiclePayload,
  UpdateVehiclePayload,
  UpdateBrokerPaymentPayload,
  VehicleSortField,
  SortDirection,
} from "@/types/vehicle";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useBrokersStore } from "./brokers";
import { useUsersStore } from "./users";
import { useExpensesStore } from "./expenses";

export const useVehiclesStore = defineStore("vehicles", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  const vehicles = ref<Vehicle[]>([]);
  const searchQuery = ref("");
  const sortField = ref<VehicleSortField>("vehicle_number");
  const sortDirection = ref<SortDirection>("asc");

  const filteredVehicles = computed(() => {
    let result = [...vehicles.value];

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const brokersStore = useBrokersStore();
      const usersStore = useUsersStore();
      result = result.filter(
        (vehicle) =>
          vehicle.vehicle_number.toLowerCase().includes(query) ||
          (vehicle.notes && vehicle.notes.toLowerCase().includes(query)) ||
          String(vehicle.joma_cost).includes(query) ||
          String(vehicle.vehicle_cost).includes(query) ||
          String(vehicle.total_paid_to_broker).includes(query) ||
          String(vehicle.other_cost).includes(query) ||
          String(vehicle.labour_cost).includes(query) ||
          String(vehicle.demarage_cost).includes(query) ||
          (vehicle.broker_id && brokersStore.getBrokerName(vehicle.broker_id).toLowerCase().includes(query)) ||
          usersStore.getUserName(vehicle.user_id).toLowerCase().includes(query)
      );
    }

    result.sort((a, b) => {
      let comparison = 0;
      switch (sortField.value) {
        case "transport_id":
          comparison = a.transport_id - b.transport_id;
          break;
        case "vehicle_number":
          comparison = a.vehicle_number.localeCompare(b.vehicle_number);
          break;
        case "broker_id":
          comparison = (a.broker_id || 0) - (b.broker_id || 0);
          break;
        case "joma_cost":
          comparison = a.joma_cost - b.joma_cost;
          break;
        case "vehicle_cost":
          comparison = a.vehicle_cost - b.vehicle_cost;
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

  const totalVehicles = computed(() => vehicles.value.length);

  const totalJomaCost = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.joma_cost, 0)
  );
  const totalVehicleCost = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.vehicle_cost, 0)
  );
  const totalOtherCost = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.other_cost, 0)
  );
  const totalLabourCost = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.labour_cost, 0)
  );
  const totalDemarageCost = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.demarage_cost, 0)
  );
  const totalBrokerPaid = computed(() =>
    vehicles.value.reduce((sum, v) => sum + v.total_paid_to_broker, 0)
  );

  // ============= EXPENSE HELPERS =============

  const generateExpenseTitle = (vehicle: Vehicle, costType: string): string => {
    let title = `${costType}: Vehicle ${vehicle.vehicle_number}`;
    if (vehicle.transport_id) {
      title = `${title} (Transport #${vehicle.transport_id})`;
    }
    return title.slice(0, 255);
  };

  const generateExpenseNotes = (vehicle: Vehicle, costType: string): string => {
    let notes = `Vehicle record #${vehicle.id} - ${costType}`;
    if (vehicle.notes) {
      notes = `${notes} - ${vehicle.notes}`;
    }
    return notes;
  };

  const findAssociatedExpense = async (vehicleId: number, costType: string): Promise<number | null> => {
    const expensesStore = useExpensesStore();

    if (expensesStore.expenses.length === 0) {
      await expensesStore.fetchExpenses();
    }

    const matchingExpense = expensesStore.expenses.find((e) =>
      e.notes && e.notes.includes(`Vehicle record #${vehicleId} - ${costType}`)
    );

    return matchingExpense ? matchingExpense.id : null;
  };

  const findAssociatedExpenses = async (vehicleId: number): Promise<number[]> => {
    const expensesStore = useExpensesStore();

    if (expensesStore.expenses.length === 0) {
      await expensesStore.fetchExpenses();
    }

    const matchingExpenses = expensesStore.expenses.filter((e) =>
      e.notes && e.notes.includes(`Vehicle record #${vehicleId}`)
    );

    return matchingExpenses.map((e) => e.id);
  };

  // ============= ACTIONS =============

  const fetchVehicles = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Vehicle[]>>("/vehicles");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      vehicles.value = res.data.data;
      return vehicles.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch vehicles");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchVehiclesByTransport = async (transportId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Vehicle[]>>("/vehicles", {
        params: { transport_id: transportId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      vehicles.value = res.data.data;
      return vehicles.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch transport vehicles");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchVehiclesByBroker = async (brokerId: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Vehicle[]>>("/vehicles", {
        params: { broker_id: brokerId },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      vehicles.value = res.data.data;
      return vehicles.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch broker vehicles");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchVehiclesByNumber = async (vehicleNumber: string) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Vehicle[]>>("/vehicles", {
        params: { vehicle_number: vehicleNumber },
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      vehicles.value = res.data.data;
      return vehicles.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch vehicles by number");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const createVehicle = async (payload: CreateVehiclePayload): Promise<Vehicle | null> => {
    displayLoader();
    try {
      if (!payload.vehicle_number) {
        push.error("Vehicle number is required");
        return null;
      }
      if (!payload.broker_id) {
        push.error("Broker is required");
        return null;
      }
      if (payload.joma_cost < 0 || payload.vehicle_cost < 0 ||
        payload.other_cost < 0 || payload.labour_cost < 0 ||
        payload.demarage_cost < 0) {
        push.error("Costs cannot be negative");
        return null;
      }

      const res = await api.post<ApiResponse<Vehicle>>("/vehicles", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const newVehicle = res.data.data;
      vehicles.value.push(newVehicle);

      const expensesStore = useExpensesStore();

      if (newVehicle.other_cost > 0) {
        await expensesStore.createExpense({
          title: generateExpenseTitle(newVehicle, 'Other Cost'),
          amount: newVehicle.other_cost,
          expense_date: new Date().toISOString(),
          notes: generateExpenseNotes(newVehicle, 'Other Cost'),
        });
      }

      if (newVehicle.labour_cost > 0) {
        await expensesStore.createExpense({
          title: generateExpenseTitle(newVehicle, 'Labour Cost'),
          amount: newVehicle.labour_cost,
          expense_date: new Date().toISOString(),
          notes: generateExpenseNotes(newVehicle, 'Labour Cost'),
        });
      }

      if (newVehicle.demarage_cost > 0) {
        await expensesStore.createExpense({
          title: generateExpenseTitle(newVehicle, 'Demarage Cost'),
          amount: newVehicle.demarage_cost,
          expense_date: new Date().toISOString(),
          notes: generateExpenseNotes(newVehicle, 'Demarage Cost'),
        });
      }

      push.success(res.data.message);
      return newVehicle;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create vehicle");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateVehicle = async (id: number, payload: UpdateVehiclePayload): Promise<Vehicle | null> => {
    displayLoader();
    try {
      const oldVehicle = vehicles.value.find((v) => v.id === id);

      const res = await api.patch<ApiResponse<Vehicle>>(`/vehicles/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const updatedVehicle = res.data.data;
      const index = vehicles.value.findIndex((vehicle) => vehicle.id === id);
      if (index !== -1) {
        vehicles.value[index] = updatedVehicle;
      }

      const expensesStore = useExpensesStore();

      // Other Cost
      if (oldVehicle && oldVehicle.other_cost !== updatedVehicle.other_cost) {
        const expenseId = await findAssociatedExpense(id, 'Other Cost');
        if (expenseId) {
          if (updatedVehicle.other_cost > 0) {
            await expensesStore.updateExpense(expenseId, {
              amount: updatedVehicle.other_cost,
              title: generateExpenseTitle(updatedVehicle, 'Other Cost'),
              notes: generateExpenseNotes(updatedVehicle, 'Other Cost'),
            });
          } else {
            await expensesStore.deleteExpense(expenseId);
          }
        } else if (updatedVehicle.other_cost > 0) {
          await expensesStore.createExpense({
            title: generateExpenseTitle(updatedVehicle, 'Other Cost'),
            amount: updatedVehicle.other_cost,
            expense_date: new Date().toISOString(),
            notes: generateExpenseNotes(updatedVehicle, 'Other Cost'),
          });
        }
      }

      // Labour Cost
      if (oldVehicle && oldVehicle.labour_cost !== updatedVehicle.labour_cost) {
        const expenseId = await findAssociatedExpense(id, 'Labour Cost');
        if (expenseId) {
          if (updatedVehicle.labour_cost > 0) {
            await expensesStore.updateExpense(expenseId, {
              amount: updatedVehicle.labour_cost,
              title: generateExpenseTitle(updatedVehicle, 'Labour Cost'),
              notes: generateExpenseNotes(updatedVehicle, 'Labour Cost'),
            });
          } else {
            await expensesStore.deleteExpense(expenseId);
          }
        } else if (updatedVehicle.labour_cost > 0) {
          await expensesStore.createExpense({
            title: generateExpenseTitle(updatedVehicle, 'Labour Cost'),
            amount: updatedVehicle.labour_cost,
            expense_date: new Date().toISOString(),
            notes: generateExpenseNotes(updatedVehicle, 'Labour Cost'),
          });
        }
      }

      // Demarage Cost
      if (oldVehicle && oldVehicle.demarage_cost !== updatedVehicle.demarage_cost) {
        const expenseId = await findAssociatedExpense(id, 'Demarage Cost');
        if (expenseId) {
          if (updatedVehicle.demarage_cost > 0) {
            await expensesStore.updateExpense(expenseId, {
              amount: updatedVehicle.demarage_cost,
              title: generateExpenseTitle(updatedVehicle, 'Demarage Cost'),
              notes: generateExpenseNotes(updatedVehicle, 'Demarage Cost'),
            });
          } else {
            await expensesStore.deleteExpense(expenseId);
          }
        } else if (updatedVehicle.demarage_cost > 0) {
          await expensesStore.createExpense({
            title: generateExpenseTitle(updatedVehicle, 'Demarage Cost'),
            amount: updatedVehicle.demarage_cost,
            expense_date: new Date().toISOString(),
            notes: generateExpenseNotes(updatedVehicle, 'Demarage Cost'),
          });
        }
      }

      push.success(res.data.message);
      return updatedVehicle;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update vehicle");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateVehicleBrokerPayment = async (id: number, payload: UpdateBrokerPaymentPayload): Promise<Vehicle | null> => {
    displayLoader();
    try {
      if (payload.total_paid_to_broker < 0) {
        push.error("Broker total paid cannot be negative");
        return null;
      }

      const res = await api.patch<ApiResponse<Vehicle>>(`/vehicles/${id}/broker-payment`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      const index = vehicles.value.findIndex((vehicle) => vehicle.id === id);
      if (index !== -1) {
        vehicles.value[index] = res.data.data;
      }

      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update broker payment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const deleteVehicle = async (id: number): Promise<boolean> => {
    displayLoader();
    try {
      const expenseIds = await findAssociatedExpenses(id);

      const res = await api.delete<ApiResponse<null>>(`/vehicles/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      vehicles.value = vehicles.value.filter((vehicle) => vehicle.id !== id);

      if (expenseIds.length > 0) {
        const expensesStore = useExpensesStore();
        let deletedCount = 0;
        for (const expenseId of expenseIds) {
          const deleted = await expensesStore.deleteExpense(expenseId);
          if (deleted) deletedCount++;
        }
        if (deletedCount > 0) {
          push.success(`${res.data.message} (${deletedCount} associated expense(s) deleted)`);
        } else {
          push.success(res.data.message);
        }
      } else {
        push.success(res.data.message);
      }

      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete vehicle");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const deleteVehiclesByTransport = async (transportId: number): Promise<boolean> => {
    displayLoader();
    try {
      const transportVehicles = vehicles.value.filter((v) => v.transport_id === transportId);

      for (const vehicle of transportVehicles) {
        const expenseIds = await findAssociatedExpenses(vehicle.id);
        const res = await api.delete<ApiResponse<null>>(`/vehicles/${vehicle.id}`);
        if (!res.data.success) {
          push.error(`Failed to delete vehicle ${vehicle.id}`);
          return false;
        }

        if (expenseIds.length > 0) {
          const expensesStore = useExpensesStore();
          for (const expenseId of expenseIds) {
            await expensesStore.deleteExpense(expenseId);
          }
        }
      }

      vehicles.value = vehicles.value.filter((vehicle) => vehicle.transport_id !== transportId);
      push.success("All transport vehicles deleted");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete transport vehicles");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const setSort = (field: VehicleSortField) => {
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

  const getVehicleById = (id: number): Vehicle | undefined => {
    return vehicles.value.find((v) => v.id === id);
  };

  const getVehiclesByTransportId = (transportId: number): Vehicle[] => {
    return vehicles.value.filter((vehicle) => vehicle.transport_id === transportId);
  };

  const getVehiclesByBrokerId = (brokerId: number): Vehicle[] => {
    return vehicles.value.filter((vehicle) => vehicle.broker_id === brokerId);
  };

  const getBrokerNameForVehicle = (vehicle: Vehicle): string => {
    if (!vehicle.broker_id) return "N/A";
    const brokersStore = useBrokersStore();
    return brokersStore.getBrokerName(vehicle.broker_id);
  };

  const getTotalCostByTransport = (transportId: number): number => {
    return vehicles.value
      .filter((vehicle) => vehicle.transport_id === transportId)
      .reduce((sum, vehicle) => sum + vehicle.joma_cost + vehicle.vehicle_cost +
        vehicle.other_cost + vehicle.labour_cost + vehicle.demarage_cost, 0);
  };

  const getTotalJomaPlusVehicleByTransport = (transportId: number): number => {
    return vehicles.value
      .filter((vehicle) => vehicle.transport_id === transportId)
      .reduce((sum, vehicle) => sum + vehicle.joma_cost + vehicle.vehicle_cost, 0);
  };

  const getTotalBrokerPaidByTransport = (transportId: number): number => {
    return vehicles.value
      .filter((vehicle) => vehicle.transport_id === transportId)
      .reduce((sum, vehicle) => sum + vehicle.total_paid_to_broker, 0);
  };

  const getVehicleExpenses = async (vehicleId: number): Promise<number[]> => {
    return await findAssociatedExpenses(vehicleId);
  };

  return {
    vehicles,
    searchQuery,
    sortField,
    sortDirection,

    filteredVehicles,
    totalVehicles,
    totalJomaCost,
    totalVehicleCost,
    totalOtherCost,
    totalLabourCost,
    totalDemarageCost,
    totalBrokerPaid,

    fetchVehicles,
    fetchVehiclesByTransport,
    fetchVehiclesByBroker,
    fetchVehiclesByNumber,

    createVehicle,
    updateVehicle,
    updateVehicleBrokerPayment,
    deleteVehicle,
    deleteVehiclesByTransport,

    setSort,
    setSearchQuery,
    clearSearch,

    getVehicleById,
    getVehiclesByTransportId,
    getVehiclesByBrokerId,
    getBrokerNameForVehicle,
    getTotalCostByTransport,
    getTotalJomaPlusVehicleByTransport,
    getTotalBrokerPaidByTransport,
    getVehicleExpenses,
  };
});