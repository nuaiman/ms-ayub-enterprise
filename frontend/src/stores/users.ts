// src/stores/users.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type {
  User,
  Role,
  CreateUserPayload,
  UpdateUserProfilePayload,
  UpdateUserSalaryPayload,
  UserSortField,
  SortDirection,
  UserSortKey
} from "@/types/auth";
import type {
  ApiResponse
} from "@/types/api";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";
import { useAuthStore } from "./auth";

export const useUsersStore = defineStore("users", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  // STATE
  const users = ref<User[]>([]);
  const searchQuery = ref('');
  const sortField = ref<UserSortField>('created_at');
  const sortDirection = ref<SortDirection>('desc');

  // ✅ Current user info from auth store
  const currentUserId = computed(() => {
    const auth = useAuthStore()
    return auth.user?.id || null
  })

  const currentUserRole = computed(() => {
    const auth = useAuthStore()
    return auth.user?.role || null
  })

  // COMPUTED
  const filteredUsers = computed(() => {
    let result = [...users.value];

    // Filter
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(u =>
        u.name.toLowerCase().includes(query) ||
        u.username.toLowerCase().includes(query) ||
        (u.email && u.email.toLowerCase().includes(query)) ||
        u.role.toLowerCase().includes(query) ||
        u.phone?.toLowerCase().includes(query) ||
        u.address?.toLowerCase().includes(query) ||
        u.id_type?.toLowerCase().includes(query) ||
        u.id_number?.toLowerCase().includes(query)
      );
    }

    // Sort
    result.sort((a, b) => {
      let comparison = 0;
      if (sortField.value === 'name') {
        comparison = a.name.localeCompare(b.name);
      } else if (sortField.value === 'username') {
        comparison = a.username.localeCompare(b.username);
      } else if (sortField.value === 'role') {
        comparison = a.role.localeCompare(b.role);
      } else if (sortField.value === 'created_at') {
        comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
      }
      return sortDirection.value === 'desc' ? -comparison : comparison;
    });

    return result;
  });

  // ACTIONS

  // GET ALL USERS
  const fetchUsers = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<User[]>>("/users");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      users.value = res.data.data;
      return users.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch users");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE USER
  const createUser = async (payload: CreateUserPayload) => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<User>>("/users", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      users.value.push(res.data.data);
      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create user");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // RESET ALL PASSWORDS
  const resetAllPasswords = async (new_password: string) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<null>>("/users/reset-all-passwords", {
        new_password,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to reset passwords");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // CHANGE USER PASSWORD (ADMIN)
  const changeUserPassword = async (id: number, new_password: string) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/change-password`, {
        new_password,
      });
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to change password");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // CHANGE ROLE
  const changeRole = async (id: number, role: Role) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/change-role`, { role });
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const user = users.value.find((u) => u.id === id);
      if (user) {
        user.role = role;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to change role");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // TOGGLE ACTIVE
  const toggleActive = async (id: number) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/toggle-active`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const user = users.value.find((u) => u.id === id);
      if (user) {
        user.is_active = !user.is_active;
      }
      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update user");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ✅ UPDATED: UPDATE USER PROFILE (without monthly_salary)
  const updateUserProfile = async (id: number, payload: UpdateUserProfilePayload) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<User>>(`/users/${id}/profile`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = users.value.findIndex((u) => u.id === id);
      if (index !== -1) {
        users.value[index] = { ...users.value[index], ...res.data.data };
      }
      push.success(res.data.message || "Profile updated");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update profile");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // ✅ NEW: UPDATE USER SALARY (separate endpoint)
  const updateUserSalary = async (id: number, payload: UpdateUserSalaryPayload) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<User>>(`/users/${id}/salary`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = users.value.findIndex((u) => u.id === id);
      if (index !== -1) {
        users.value[index] = { ...users.value[index], ...res.data.data };
      }
      push.success(res.data.message || "Salary updated");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update salary");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // SORT
  const setSort = (field: UserSortField) => {
    if (sortField.value === field) {
      sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
    } else {
      sortField.value = field;
      sortDirection.value = 'asc';
    }
  };

  const applySortByKey = (key: UserSortKey) => {
    switch (key) {
      case 'newest': sortField.value = 'created_at'; sortDirection.value = 'desc'; break;
      case 'oldest': sortField.value = 'created_at'; sortDirection.value = 'asc'; break;
      case 'name_asc': sortField.value = 'name'; sortDirection.value = 'asc'; break;
      case 'name_desc': sortField.value = 'name'; sortDirection.value = 'desc'; break;
    }
  };

  const getSortKey = (): UserSortKey => {
    if (sortField.value === 'created_at' && sortDirection.value === 'desc') return 'newest';
    if (sortField.value === 'created_at' && sortDirection.value === 'asc') return 'oldest';
    if (sortField.value === 'name' && sortDirection.value === 'asc') return 'name_asc';
    return 'name_desc';
  };

  // SEARCH
  const setSearchQuery = (query: string) => {
    searchQuery.value = query;
  };

  const clearSearch = () => {
    searchQuery.value = '';
  };

  // UTILITY - Get user name by ID
  const getUserName = (userId: number): string => {
    const user = users.value.find(u => u.id === userId);
    return user ? user.name : `User #${userId}`;
  };

  // UTILITY - Get initials
  const getInitials = (name: string): string => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  // UTILITY - Format date
  const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const formatTime = (dateStr: string): string => {
    return new Date(dateStr).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // UTILITY - Role badge styles
  const getRoleBadgeClass = (role: string): string => {
    switch (role) {
      case 'admin': return 'bg-info-bg text-info-text';
      case 'manager': return 'bg-warning-bg text-warning-text';
      default: return 'bg-surface-alt text-secondary';
    }
  };

  const getRoleDotClass = (role: string): string => {
    switch (role) {
      case 'admin': return 'bg-info-text';
      case 'manager': return 'bg-warning-text';
      default: return 'bg-muted';
    }
  };

  // UTILITY - Generate TSV for clipboard
  const generateClipboardData = (): string => {
    const data = filteredUsers.value;
    if (data.length === 0) return '';

    const headers = ['ID', 'Name', 'Username', 'Role', 'Status', 'Email', 'Phone', 'Address', 'Salary', 'ID Type', 'ID Number', 'Created'];
    const rows = data.map(u => [
      u.id, u.name, u.username, u.role, u.is_active ? 'Active' : 'Inactive',
      u.email || '', u.phone || '', u.address || '', u.monthly_salary ? u.monthly_salary.toFixed(2) : '',
      u.id_type || '', u.id_number || '', formatDate(u.created_at),
    ]);

    return [headers, ...rows]
      .map(row => row.map(cell => String(cell).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t'))
      .join('\n');
  };

  return {
    // State
    users,
    searchQuery,
    sortField,
    sortDirection,

    // Computed
    filteredUsers,
    currentUserId,
    currentUserRole,

    // Actions
    fetchUsers,
    createUser,
    resetAllPasswords,
    changeUserPassword,
    changeRole,
    toggleActive,
    updateUserProfile,
    updateUserSalary, // ✅ NEW

    // Sort
    setSort,
    applySortByKey,
    getSortKey,

    // Search
    setSearchQuery,
    clearSearch,

    // Utilities
    getUserName,
    getInitials,
    formatDate,
    formatTime,
    getRoleBadgeClass,
    getRoleDotClass,
    generateClipboardData,
  };
});