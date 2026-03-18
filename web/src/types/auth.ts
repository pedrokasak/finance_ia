export type SignupRequest = {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  confirmPassword: string;
  acceptTerms: boolean;
};

export type SignupResponse = {
  success: boolean;
};

export type LoginRequest = {
  email: string;
  password: string;
  code?: string;
  keepAlive?: boolean;
};

export type AuthenticationResponse = {
  email: string;
  token?: string;
  refreshToken?: string;
  error?: string;
  success: boolean;
};

export type AuthenticatorResponseError = {
  message: string;
  success: boolean;
};

export type LogoutResponse = {
  success: boolean;
};

export type ForgotPasswordRequest = {
  email: string;
};

export type ForgotPasswordResponse = {
  message: string;
};

export type ResetPasswordRequest = {
  token: string;
  newPassword: string;
  confirmNewPassword: string;
};

export type ResetPasswordResponse = {
  message: string;
};

export type AuthPage = "login" | "signup" | "forgot-password" | "app";
