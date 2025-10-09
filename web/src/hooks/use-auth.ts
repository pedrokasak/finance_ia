import { useState } from 'react';
import { authenticationService } from '../services/Authentication';
import { AuthenticationResponse } from '../types/auth';

const useAuth = () => {
  const [user, setUser] = useState<AuthenticationResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await authenticationService.login({ email, password });
      setUser(response);
      return response;
    } catch (err: { message?: string } | unknown) {
      const msg = (err as { message?: string })?.message || 'Erro no login';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    try {
      await authenticationService.logout();
    } finally {
      setUser(null);
    }
  };

  const forgotPassword = async (email: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await authenticationService.forgotPassword({ email });
      return response;
    } catch (err: { message?: string } | unknown) {
      const msg =
        (err as { message?: string })?.message ||
        'Erro ao solicitar redefinição de senha';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const resetPassword = async (data: {
    token: string;
    newPassword: string;
    confirmNewPassword: string;
  }) => {
    setLoading(true);
    setError(null);
    try {
      const response = await authenticationService.resetPassword(data);
      return response;
    } catch (err: { message?: string } | unknown) {
      const msg =
        (err as { message?: string })?.message ||
        'Erro ao solicitar redefinição de senha';
      setError(msg);
      throw err;
    } finally {
      setLoading(false);
    }
  };
  return { user, loading, error, login, logout, forgotPassword, resetPassword };
};

export default useAuth;
