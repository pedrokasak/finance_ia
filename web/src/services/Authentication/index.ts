import { authentication } from '../../api';
import { AuthenticationInterface } from '../../interfaces/auth';
import {
  AuthenticationResponse,
  LoginRequest,
  LogoutResponse,
  SignupRequest,
  SignupResponse,
} from '../../types/auth';

class AuthenticationService implements AuthenticationInterface {
  async login(data: LoginRequest): Promise<AuthenticationResponse> {
    try {
      const response = await authentication.login(data.email, data.password, data.code) as { token: string; refreshToken?: string; email: string };

      localStorage.setItem('authToken', response.token);

      if (response.refreshToken) {
        localStorage.setItem('refreshToken', response.refreshToken);
      }

      return {
        email: response.email,
        token: response.token,
        refreshToken: response.refreshToken,
        success: true,
      };
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      console.error('Login failed:', error);
      return {
        email: data.email,
        error: error?.error || 'Falha no login',
        success: false,
      };
    }
  }

  async logout(): Promise<LogoutResponse> {
    return { success: true };
  }
  async signup(data: SignupRequest): Promise<SignupResponse> {
    console.log(
      `Signing up with email: ${data.firstName} ${data.lastName} <${data.email}> ${data.password} ${data.confirmPassword} ${data.acceptTerms}`,
    );
    return { success: true };
  }
  async forgotPassword(data: { email: string }): Promise<{ message: string }> {
    console.log(`Forgot password for email: ${data.email}`);
    return {
      message:
        'If that email is registered, you will receive a password reset link.',
    };
  }
  async resetPassword(data: {
    token: string;
    newPassword: string;
    confirmNewPassword: string;
  }): Promise<{ message: string }> {
    console.log(`Resetting password with token: ${data.token}`);
    return { message: 'Your password has been successfully reset.' };
  }
}

export const authenticationService = new AuthenticationService();
