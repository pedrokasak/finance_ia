import { create } from "zustand";
import { persist, createJSONStorage, devtools } from "zustand/middleware";

interface AuthenticationResponse {
  email: string;
  token: string;
  password?: string;
  success: boolean;
}

interface AuthState {
  user: AuthenticationResponse | null;
  token: string | null;
  isAuthenticated: boolean;
}

interface AuthActions {
  setUser: (user: AuthenticationResponse) => void;
  clearAuth: () => void;
  setToken: (token: string) => void;
  hydrate: () => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      (set, get) => ({
        user: null,
        token: null,
        isAuthenticated: false,

        setUser: (user) => {
          set(
            {
              user,
              token: user.token,
              isAuthenticated: true,
            },
            false,
            "auth/setUser",
          );
        },

        clearAuth: () => {
          localStorage.removeItem("authToken");
          set(
            {
              user: null,
              token: null,
              isAuthenticated: false,
            },
            false,
            "auth/clearAuth",
          );
        },

        setToken: (token) => {
          localStorage.setItem("authToken", token);
          set({ token }, false, "auth/setToken");
        },

        hydrate: () => {
          const token = localStorage.getItem("authToken");
          const storedUser = localStorage.getItem("auth-storage");

          if (token && storedUser) {
            try {
              const parsed = JSON.parse(storedUser);
              if (parsed.state?.user) {
                set(
                  {
                    user: parsed.state.user,
                    token,
                    isAuthenticated: true,
                  },
                  false,
                  "auth/hydrate",
                );
              }
            } catch (error) {
              console.error("Error hydrating auth state:", error);
              get().clearAuth();
            }
          }
        },
      }),
      {
        name: "auth-storage",
        storage: createJSONStorage(() => localStorage),
        partialize: (state) => ({
          user: state.user,
          token: state.token,
          isAuthenticated: state.isAuthenticated,
        }),
      },
    ),
    {
      name: "AuthStore",
      enabled: import.meta.env.VITE_ENV === "development" ? true : false,
    },
  ),
);
