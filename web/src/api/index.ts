import {
  ForgotPasswordRequest,
  ForgotPasswordResponse,
  LogoutResponse,
  ResetPasswordRequest,
  ResetPasswordResponse,
  SignupRequest,
} from "../types/auth";

const baseUrl = import.meta.env.VITE_API_BASE_URL;

export const fetchApi = async <T = unknown>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> => {
  const response = await fetch(`${baseUrl}${endpoint}`, options);
  let data: T;
  try {
    data = await response.json();
  } catch (e: unknown) {
    data = {} as T;
    console.error(e);
  }
  if (!response.ok) {
    throw data;
  }
  return data;
};

const authentication = {
  login: async (email: string, password: string, code?: string) => {
    const response = await fetchApi(`auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password, code }),
    });

    return response as unknown;
  },
  signup: async (data: SignupRequest) => {
    return await fetchApi<unknown>(`user/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        first_name: data.firstName,
        last_name: data.lastName,
        email: data.email,
        password: data.password,
      }),
    });
  },
  logout: async (): Promise<LogoutResponse> => {
    return await fetchApi<LogoutResponse>(`auth/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
  },
  forgotPassword: async (
    data: ForgotPasswordRequest,
  ): Promise<ForgotPasswordResponse> => {
    return await fetchApi<ForgotPasswordResponse>(`auth/forgot-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
  },
  resetPassword: async (
    data: ResetPasswordRequest,
  ): Promise<ResetPasswordResponse> => {
    return await fetchApi<ResetPasswordResponse>(`auth/reset-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
  },
};

export { authentication };
