import { useState } from 'react';
import { LoginPage } from './LoginPage';
import { SignupPage } from './SignupPage';
import { ForgotPasswordPage } from './ForgotPasswordPage';

interface AuthContainerProps {
  onLogin: () => void;
}

export function AuthContainer({ onLogin }: AuthContainerProps) {
  const [currentPage, setCurrentPage] = useState<
    'login' | 'signup' | 'forgot-password'
  >('login');

  const handleNavigate = (
    page: 'login' | 'signup' | 'forgot-password' | 'app',
  ) => {
    if (page === 'app') {
      onLogin();
    } else {
      setCurrentPage(page);
    }
  };

  switch (currentPage) {
    case 'login':
      return <LoginPage onNavigate={handleNavigate} />;
    case 'signup':
      return <SignupPage onNavigate={handleNavigate} />;
    case 'forgot-password':
      return <ForgotPasswordPage onNavigate={handleNavigate} />;
    default:
      return <LoginPage onNavigate={handleNavigate} />;
  }
}
