import { useCallback, useReducer, useRef } from 'react';
import { authenticationService } from '../services/Authentication';

interface AuthenticationResponse {
  email: string;
  token: string;
  password?: string;
  success: boolean;
}

interface LoginRequest {
  email: string;
  password: string;
  code?: string;
}

interface SignupRequest {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  confirmPassword: string;
  acceptTerms: boolean;
}

interface ResetPasswordRequest {
  token: string;
  newPassword: string;
  confirmNewPassword: string;
}

interface AuthState {
  user: AuthenticationResponse | null;
  loading: boolean;
  error: string | null;
}

type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: AuthenticationResponse }
  | { type: 'AUTH_ERROR'; payload: string }
  | { type: 'AUTH_RESET' }
  | { type: 'LOGOUT' };

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'AUTH_START':
      return { ...state, loading: true, error: null };
    case 'AUTH_SUCCESS':
      return { user: action.payload, loading: false, error: null };
    case 'AUTH_ERROR':
      return { ...state, loading: false, error: action.payload };
    case 'AUTH_RESET':
      return { ...state, loading: false, error: null };
    case 'LOGOUT':
      return { user: null, loading: false, error: null };
    default:
      return state;
  }
};

const initialState: AuthState = {
  user: null,
  loading: false,
  error: null,
};

const useAuth = () => {
  const [state, dispatch] = useReducer(authReducer, initialState);
  const isUnmountedRef = useRef(false);

  const handleError = useCallback(
    (err: unknown, defaultMsg: string): string => {
      const message =
        err instanceof Error
          ? err.message
          : (err as { message?: string })?.message || defaultMsg;
      return message;
    },
    [],
  );

  const login = useCallback(
    async (email: string, password: string, code?: string) => {
      dispatch({ type: 'AUTH_START' });
      try {
        const response = await authenticationService.login({ email, password, code });
        if (!isUnmountedRef.current) {
          dispatch({
            type: 'AUTH_SUCCESS',
            payload: response as AuthenticationResponse,
          });
        }
        return response;
      } catch (err) {
        const message = handleError(err, 'Erro no login');
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_ERROR', payload: message });
        }
        throw err;
      }
    },
    [handleError],
  );

  const logout = useCallback(async () => {
    try {
      await authenticationService.logout();
    } finally {
      if (!isUnmountedRef.current) {
        dispatch({ type: 'LOGOUT' });
      }
    }
  }, []);

  const signup = useCallback(
    async (data: SignupRequest) => {
      dispatch({ type: 'AUTH_START' });
      try {
        const response = await authenticationService.signup(data);
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_RESET' });
        }
        return response;
      } catch (err) {
        const message = handleError(err, 'Erro ao criar conta');
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_ERROR', payload: message });
        }
        throw err;
      }
    },
    [handleError],
  );

  const forgotPassword = useCallback(
    async (email: string) => {
      dispatch({ type: 'AUTH_START' });
      try {
        const response = await authenticationService.forgotPassword({ email });
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_RESET' });
        }
        return response;
      } catch (err) {
        const message = handleError(
          err,
          'Erro ao solicitar redefinição de senha',
        );
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_ERROR', payload: message });
        }
        throw err;
      }
    },
    [handleError],
  );

  const resetPassword = useCallback(
    async (data: ResetPasswordRequest) => {
      dispatch({ type: 'AUTH_START' });
      try {
        const response = await authenticationService.resetPassword(data);
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_RESET' });
        }
        return response;
      } catch (err) {
        const message = handleError(err, 'Erro ao redefinir senha');
        if (!isUnmountedRef.current) {
          dispatch({ type: 'AUTH_ERROR', payload: message });
        }
        throw err;
      }
    },
    [handleError],
  );

  return {
    user: state.user,
    loading: state.loading,
    error: state.error,
    login,
    logout,
    signup,
    forgotPassword,
    resetPassword,
  };
};

export default useAuth;
export type {
  AuthenticationResponse,
  LoginRequest,
  SignupRequest,
  ResetPasswordRequest,
};
