import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081/api/v1";

export const api = axios.create({
  baseURL: API_URL,
  headers: { "Content-Type": "application/json" },
});

// Attach JWT token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("authToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto-logout on 401 — dispatch event instead of reloading the page
// so that React state updates gracefully without losing navigation context
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem("authToken");
      localStorage.removeItem("onboarding_completed");
      // Emit a custom event — App.tsx listens and clears state without reloading
      window.dispatchEvent(new CustomEvent("auth:logout"));
    }
    return Promise.reject(err);
  },
);
