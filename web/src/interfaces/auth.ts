import {
  AuthenticationResponse,
  ForgotPasswordRequest,
  ForgotPasswordResponse,
  LoginRequest,
  LogoutResponse,
  SignupRequest,
  SignupResponse,
} from '../types/auth';

export interface AuthenticationInterface {
  login: (data: LoginRequest) => Promise<AuthenticationResponse>;
  signup: (data: SignupRequest) => Promise<SignupResponse>;
  logout: () => Promise<LogoutResponse>;
  forgotPassword: (
    data: ForgotPasswordRequest,
  ) => Promise<ForgotPasswordResponse>;
  resetPassword: (data: {
    token: string;
    newPassword: string;
    confirmNewPassword: string;
  }) => Promise<ForgotPasswordResponse>;
}
