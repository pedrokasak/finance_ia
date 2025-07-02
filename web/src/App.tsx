import { useState } from 'react';
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

export type Page =
  | 'dashboard'
  | 'transactions'
  | 'reports'
  | 'profile'
  | 'subscription';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [currentPage, setCurrentPage] = useState<Page>('dashboard');
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const handleLogin = () => {
    setIsAuthenticated(true);
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
      default:
        return <Dashboard />;
    }
  };

  if (!isAuthenticated) {
    return (
      <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
        <AuthContainer onLogin={handleLogin} />
        <Toaster />
      </ThemeProvider>
    );
  }

  return (
    <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
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
    </ThemeProvider>
  );
}

export default App;
