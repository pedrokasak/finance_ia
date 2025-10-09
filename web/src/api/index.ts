import {
  ForgotPasswordRequest,
  ForgotPasswordResponse,
  LogoutResponse,
  ResetPasswordRequest,
  ResetPasswordResponse,
  SignupRequest,
} from '../types/auth';

const baseUrl = import.meta.env.VITE_API_BASE_URL;

export const fetchApi = async (endpoint: string, options = {}) => {
  const response = await fetch(`${baseUrl}${endpoint}`, options);
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
};

const authentication = {
  login: async (email: string, password: string) => {
    const response = await fetchApi(`auth/login/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    const data = await response.json();
    return data;
  },
  signup: async (data: SignupRequest) => {
    return await fetchApi(`auth/signup/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  },
  logout: async (): Promise<LogoutResponse> => {
    return await fetchApi(`auth/logout/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });
  },
  forgotPassword: async (
    data: ForgotPasswordRequest,
  ): Promise<ForgotPasswordResponse> => {
    return await fetchApi(`auth/forgot-password/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  },
  resetPassword: async (
    data: ResetPasswordRequest,
  ): Promise<ResetPasswordResponse> => {
    return await fetchApi(`auth/reset-password/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  },
};

export { authentication };
