import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Menu, Sun, Moon, Bell, Settings, LogOut, Crown } from "lucide-react";
import { useTheme } from "next-themes";
import type { Page } from "@/App";
import useAuth from "../../hooks/use-auth";
import { useUser } from "@/contexts/UserContext";
import { isTrustedAvatarDataUrl } from "@/utils/avatar";

interface HeaderProps {
  currentPage: Page;
  setCurrentPage: (page: Page) => void;
  toggleSidebar: () => void;
}

export function Header({
  currentPage,
  setCurrentPage,
  toggleSidebar,
}: HeaderProps) {
  const { theme, setTheme } = useTheme();
  const auth = useAuth();
  const { profile, refreshProfile } = useUser();

  useEffect(() => {
    refreshProfile();
  }, [refreshProfile]);

  const getPageTitle = () => {
    switch (currentPage) {
      case "dashboard":
        return "Dashboard";
      case "transactions":
        return "Transações";
      case "reports":
        return "Relatórios";
      case "profile":
        return "Perfil";
      case "subscription":
        return "Assinatura";
      case "goals":
        return "Metas Financeiras";
      default:
        return "Dashboard";
    }
  };

  const displayName =
    [profile?.first_name, profile?.last_name].filter(Boolean).join(" ") ||
    profile?.email?.split("@")[0] ||
    "Usuário";

  const initials =
    (profile?.first_name?.[0] || "") + (profile?.last_name?.[0] || "") ||
    displayName[0]?.toUpperCase() ||
    "U";

  const planLabel: Record<string, string> = {
    free: "Gratuito",
    pro: "Pro",
    premium: "Premium",
  };

  const onSubmit = () => {
    auth.logout();
    window.location.href = "/login";
  };

  return (
    <header className="bg-card border-b border-border">
      <div className="flex items-center justify-between h-16 px-4 sm:px-6">
        <div className="flex items-center space-x-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleSidebar}
            className="lg:hidden"
          >
            <Menu className="h-5 w-5" />
          </Button>

          <div>
            <h1 className="text-2xl font-bold text-foreground">
              {getPageTitle()}
            </h1>
            <p className="text-sm text-muted-foreground">
              Gerencie suas finanças de forma inteligente
            </p>
          </div>
        </div>

        <div className="flex items-center space-x-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
          >
            <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          </Button>

          <Button variant="ghost" size="icon">
            <Bell className="h-5 w-5" />
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="relative h-10 w-10 rounded-full"
              >
                <Avatar className="h-10 w-10">
                  {isTrustedAvatarDataUrl(profile?.avatar_url) ? (
                    <AvatarImage src={profile?.avatar_url ?? ""} alt={displayName} />
                  ) : null}
                  <AvatarFallback>{initials}</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium leading-none">
                    {displayName}
                  </p>
                  <p className="text-xs leading-none text-muted-foreground">
                    {profile?.email}
                  </p>
                  <Badge variant="secondary" className="w-fit mt-1">
                    <Crown className="h-3 w-3 mr-1" />
                    {planLabel[profile?.plan || "free"] || "Gratuito"}
                  </Badge>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => setCurrentPage("profile")}>
                <Settings className="mr-2 h-4 w-4" />
                <span>Configurações</span>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setCurrentPage("subscription")}>
                <Crown className="mr-2 h-4 w-4" />
                <span>Assinatura</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <LogOut className="mr-2 h-4 w-4" />
                <Button variant="ghost" className="p-0" onClick={onSubmit}>
                  Sair
                </Button>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
}
