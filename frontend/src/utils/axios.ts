// src/utils/ axios.ts

import axios from "axios";
import { useAuthStore } from "@/stores/auth";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api",
  withCredentials: true,
});

let accessToken: string | null = localStorage.getItem("access_token");
let isRefreshing = false;
let failedQueue: Array<{ resolve: (token: string) => void; reject: (error: unknown) => void }> = [];

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else if (token) {
      resolve(token);
    }
  });
  failedQueue = [];
};

export const setAccessToken = (token: string | null) => {
  accessToken = token;
  if (token) {
    localStorage.setItem("access_token", token);
  } else {
    localStorage.removeItem("access_token");
  }
};

api.interceptors.request.use((config) => {
  // Always check localStorage first for the latest token
  const storedToken = localStorage.getItem("access_token");
  if (storedToken) {
    accessToken = storedToken;
  }

  if (accessToken) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${accessToken}`;
  }

  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Only handle 401 and avoid infinite loops
    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error);
    }

    // Don't try to refresh if the request was to refresh itself or login
    if (originalRequest.url === "/users/refresh-session" || originalRequest.url === "/users/login") {
      return Promise.reject(error);
    }

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then((token) => {
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return api(originalRequest);
      });
    }

    originalRequest._retry = true;
    isRefreshing = true;

    try {
      const res = await axios.post(
        `${api.defaults.baseURL}/users/refresh-session`,
        {},
        { withCredentials: true },
      );

      if (res.data?.success) {
        const newToken = res.data.data.access_token;

        // Update token
        setAccessToken(newToken);
        api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;

        // Update auth store
        const authStore = useAuthStore();
        authStore.updateToken(newToken);

        processQueue(null, newToken);

        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      }

      processQueue(new Error("Refresh failed"));
      return Promise.reject(error);
    } catch (refreshError) {
      processQueue(refreshError);
      // Clear auth state on refresh failure
      const authStore = useAuthStore();
      authStore.clearSession();
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  },
);

export default api;