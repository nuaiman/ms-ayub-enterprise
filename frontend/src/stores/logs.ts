// src/stores/logs.ts

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  Log,
  LogSortField,
  SortDirection,
} from "@/types/log";
import type { ApiResponse } from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useUsersStore } from "./users";

export const useLogsStore = defineStore("logs", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // ============= STATE =============
  const logs = ref<Log[]>([]);
  const searchQuery = ref("");
  const sortField = ref<LogSortField>("created_at");
  const sortDirection = ref<SortDirection>("desc");

  // ============= COMPUTED =============
  const filteredLogs = computed(() => {
    let result = [...logs.value];

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const usersStore = useUsersStore();
      result = result.filter(
        (log: Log) =>
          log.action.toLowerCase().includes(query) ||
          log.description.toLowerCase().includes(query) ||
          log.entity_type.toLowerCase().includes(query) ||
          String(log.entity_id).includes(query) ||
          (log.ip_address && log.ip_address.toLowerCase().includes(query)) ||
          usersStore.getUserName(log.user_id).toLowerCase().includes(query)
      );
    }

    // Sort
    result.sort((a: Log, b: Log) => {
      let comparison = 0;
      switch (sortField.value) {
        case "user_id":
          comparison = a.user_id - b.user_id;
          break;
        case "action":
          comparison = a.action.localeCompare(b.action);
          break;
        case "entity_type":
          comparison = a.entity_type.localeCompare(b.entity_type);
          break;
        case "entity_id":
          comparison = a.entity_id - b.entity_id;
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

  const totalLogs = computed(() => logs.value.length);

  // ============= ACTIONS =============

  // GET ALL LOGS
  const fetchLogs = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Log[]>>("/logs");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      logs.value = res.data.data;
      return logs.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch logs");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET LOGS BY USER
  const fetchLogsByUser = async (userId: number) => {
    displayLoader();
    try {
      // Since the backend doesn't have a specific endpoint for user logs,
      // we fetch all and filter client-side
      const res = await api.get<ApiResponse<Log[]>>("/logs");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      logs.value = res.data.data.filter((log: Log) => log.user_id === userId);
      return logs.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch user logs");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // GET LOGS BY ENTITY
  const fetchLogsByEntity = async (entityType: string, entityId: number) => {
    displayLoader();
    try {
      // Since the backend doesn't have a specific endpoint for entity logs,
      // we fetch all and filter client-side
      const res = await api.get<ApiResponse<Log[]>>("/logs");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      logs.value = res.data.data.filter(
        (log: Log) => log.entity_type === entityType && log.entity_id === entityId
      );
      return logs.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch entity logs");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // ============= SORT =============
  const setSort = (field: LogSortField) => {
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

  // Get log by ID
  const getLogById = (id: number): Log | undefined => {
    return logs.value.find((log: Log) => log.id === id);
  };

  // Get logs by user ID
  const getLogsByUserId = (userId: number): Log[] => {
    return logs.value.filter((log: Log) => log.user_id === userId);
  };

  // Get logs by entity
  const getLogsByEntity = (entityType: string, entityId: number): Log[] => {
    return logs.value.filter(
      (log: Log) => log.entity_type === entityType && log.entity_id === entityId
    );
  };

  // Get action label
  const getActionLabel = (action: string): string => {
    const labels: Record<string, string> = {
      'create': 'Created',
      'update': 'Updated',
      'delete': 'Deleted',
      'login': 'Logged In',
      'logout': 'Logged Out',
      'approve': 'Approved',
      'reject': 'Rejected',
      'pay': 'Paid',
      'status_change': 'Status Changed',
    };
    return labels[action] || action.charAt(0).toUpperCase() + action.slice(1);
  };

  // Get action badge class
  const getActionBadgeClass = (action: string): string => {
    switch (action) {
      case 'create':
        return 'bg-success-bg text-success-text';
      case 'update':
        return 'bg-info-bg text-info-text';
      case 'delete':
        return 'bg-danger-bg text-danger-text';
      case 'login':
        return 'bg-info-bg text-info-text';
      case 'logout':
        return 'bg-warning-bg text-warning-text';
      case 'approve':
        return 'bg-success-bg text-success-text';
      case 'reject':
        return 'bg-danger-bg text-danger-text';
      case 'pay':
        return 'bg-success-bg text-success-text';
      case 'status_change':
        return 'bg-warning-bg text-warning-text';
      default:
        return 'bg-surface-alt text-secondary';
    }
  };

  // Format log date
  const formatLogDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  };

  // Truncate description for display
  const truncateDescription = (description: string, maxLength: number = 100): string => {
    if (description.length <= maxLength) return description;
    return description.substring(0, maxLength) + '...';
  };

  // Check if log has old_data
  const hasOldData = (log: Log): boolean => {
    return log.old_data !== null && log.old_data !== '';
  };

  // Check if log has new_data
  const hasNewData = (log: Log): boolean => {
    return log.new_data !== null && log.new_data !== '';
  };

  return {
    // State
    logs,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredLogs,
    totalLogs,

    // Fetch
    fetchLogs,
    fetchLogsByUser,
    fetchLogsByEntity,

    // Sort
    setSort,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getLogById,
    getLogsByUserId,
    getLogsByEntity,
    getActionLabel,
    getActionBadgeClass,
    formatLogDate,
    truncateDescription,
    hasOldData,
    hasNewData,
  };
});