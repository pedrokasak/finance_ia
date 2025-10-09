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
    authentication
      .login(data.email, data.password)
      .then((response) => {
        console.log('Login successful:', response);
        localStorage.setItem('authToken', response.token);
        return {
          email: response.email,
          password: response.password,
          token: response.token,
          success: true,
        };
      })
      .catch((error) => {
        throw new Error(`Login failed: ${error.message}`);
      });
    return {
      email: data.email,
      password: data.password,
      success: false,
    };
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
