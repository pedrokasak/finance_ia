import { useState, useEffect } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from 'next-themes';
import { Toaster } from '@/components/ui/sonner';
import { Sidebar } from '@/components/layout/Sidebar';
import { Header } from '@/components/layout/Header';
import { Dashboard } from '@/components/pages/Dashboard';
import { Transactions } from '@/components/pages/Transactions';
import { Profile } from '@/components/pages/Profile';
import { Subscription } from '@/components/pages/Subscription';
import { Reports } from '@/components/pages/Reports';
import { AuthContainer } from '@/components/auth/AuthContainer';
import { OnboardingFlow } from '@/components/onboarding/OnboardingFlow';
import { UserProvider } from '@/contexts/UserContext';
import { Goals } from './components/pages/Goals';

export type Page =
  | 'dashboard'
  | 'transactions'
  | 'reports'
  | 'profile'
  | 'subscription'
  | 'goals';

const queryClient = new QueryClient();

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('authToken'));
  const [onboardingCompleted, setOnboardingCompleted] = useState(
    () => localStorage.getItem('onboarding_completed') === 'true'
  );
  const [currentPage, setCurrentPage] = useState<Page>(() => {
    const params = new URLSearchParams(window.location.search);
    const page = params.get('page');
    if (page && ['dashboard', 'transactions', 'reports', 'profile', 'subscription', 'goals'].includes(page)) {
      return page as Page;
    }
    return 'dashboard';
  });
  const [sidebarOpen, setSidebarOpen] = useState(false);

  // Listen for auth:logout events dispatched by the API interceptor on 401
  // Using an event instead of window.location.reload() keeps React state stable
  useEffect(() => {
    const handleAuthLogout = () => {
      setIsAuthenticated(false);
      setOnboardingCompleted(false);
    };
    window.addEventListener('auth:logout', handleAuthLogout);
    return () => window.removeEventListener('auth:logout', handleAuthLogout);
  }, []);

  const handleLogin = () => {
    setIsAuthenticated(true);
    // Re-read onboarding state after login (user may have completed it before)
    setOnboardingCompleted(localStorage.getItem('onboarding_completed') === 'true');
  };

  const handleOnboardingComplete = () => {
    localStorage.setItem('onboarding_completed', 'true');
    setOnboardingCompleted(true);
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <Dashboard />;
      case 'transactions':
        return <Transactions />;
      case 'reports':
        return <Reports />;
      case 'profile':
        return <Profile />;
      case 'subscription':
        return <Subscription />;
      case 'goals':
        return <Goals />;
      default:
        return <Dashboard />;
    }
  };

  if (!isAuthenticated) {
    return (
      <QueryClientProvider client={queryClient}>
        <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
          <AuthContainer onLogin={handleLogin} />
          <Toaster />
        </ThemeProvider>
      </QueryClientProvider>
    );
  }

  if (!onboardingCompleted) {
    return (
      <QueryClientProvider client={queryClient}>
        <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
          <OnboardingFlow onComplete={handleOnboardingComplete} />
          <Toaster />
        </ThemeProvider>
      </QueryClientProvider>
    );
  }

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
        <UserProvider>
          <div className="min-h-screen bg-background">
            <div className="flex h-screen">
              <Sidebar
                currentPage={currentPage}
                setCurrentPage={setCurrentPage}
                isOpen={sidebarOpen}
                setIsOpen={setSidebarOpen}
              />

              <div className="flex-1 flex flex-col min-w-0">
                <Header
                  currentPage={currentPage}
                  setCurrentPage={setCurrentPage}
                  toggleSidebar={() => setSidebarOpen(!sidebarOpen)}
                />

                <main className="flex-1 overflow-auto bg-background">
                  <div className="container mx-auto px-4 sm:px-6 py-8 max-w-7xl">
                    {renderPage()}
                  </div>
                </main>
              </div>
            </div>
          </div>

          <Toaster />
        </UserProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;

