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
  let data;
  try {
    data = await response.json();
  } catch (e) {
    data = {};
  }
  if (!response.ok) {
    throw data;
  }
  return data;
};

const authentication = {
  login: async (email: string, password: string, code?: string) => {
    const response = await fetchApi(`auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, code }),
    });

    console.log('Login successful:', response);

    return response;
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
