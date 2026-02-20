import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { Mail, Bell, Shield, CreditCard, Loader2, Sparkles, Zap, Lock } from 'lucide-react';
import { toast } from 'sonner';
import { api } from '@/api/client';
import { QRCodeSVG } from 'qrcode.react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { FinancialMethodModal } from '@/components/finance/FinancialMethodModal';
import { useUser } from '@/contexts/UserContext';

interface UserInfo {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  plan?: string;
  avatar_url?: string;
  financial_method_id?: string;
  notifications_enabled?: boolean;
  two_fa_enabled?: boolean;
}

interface FinancialMethod {
  id: string;
  name: string;
  description: string;
  tagline?: string;
}

// Extract user info from JWT without needing a /me endpoint immediately
function getJWTPayload(): { email?: string; plan?: string; user_id?: string } {
  try {
    const token = localStorage.getItem('authToken');
    if (!token) return {};
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload;
  } catch {
    return {};
  }
}

const planConfig: Record<string, { label: string; color: string; price: string; icon: React.ElementType }> = {
  free: { label: 'Gratuito', color: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300', price: 'R$ 0/mês', icon: Zap },
  premium: { label: 'Premium', color: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300', price: 'R$ 14,90/mês', icon: Sparkles },
  pro: { label: 'Pro', color: 'bg-gradient-to-r from-purple-600 to-blue-600 text-white', price: 'R$ 29,90/mês', icon: Lock },
};

export function Profile() {
  const jwt = getJWTPayload();
  const { updateAvatar, refreshProfile } = useUser();
  const [user, setUser] = useState<UserInfo>({
    id: jwt.user_id || '',
    email: jwt.email || '',
    plan: jwt.plan || 'free',
  });
  const [saving, setSaving] = useState(false);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [avatarUrl, setAvatarUrl] = useState('');
  const [financialMethodId, setFinancialMethodId] = useState<string>('');
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [methods, setMethods] = useState<FinancialMethod[]>([]);
  const [isMethodModalOpen, setIsMethodModalOpen] = useState(false);

  // 2FA Setup States
  const [is2FASetupOpen, setIs2FASetupOpen] = useState(false);
  const [setup2FAData, setSetup2FAData] = useState<{ secret: string; otpauth_url: string } | null>(null);
  const [verifyCode, setVerifyCode] = useState('');
  const [verifying2FA, setVerifying2FA] = useState(false);

  // Try to fetch user details from API
  useEffect(() => {
    if (!jwt.user_id) return;

    // Fetch user details
    api.get(`/user/${jwt.user_id}`).then((res: any) => {
      const u = res.data?.user || res.data;
      if (u) {
        setUser((prev) => ({ ...prev, ...u }));
        setFirstName(u.first_name || '');
        setLastName(u.last_name || '');
        setAvatarUrl(u.avatar_url || '');
        setFinancialMethodId(u.financial_method_id || '');
        if (u.notifications_enabled !== undefined) {
          setNotificationsEnabled(u.notifications_enabled);
        }
      }
    }).catch(() => {
      setFirstName(jwt.email?.split('@')[0] || '');
    });

    // Fetch financial methods
    api.get('/finance/methods').then((res: any) => {
      setMethods(res.data?.data || res.data || []);
    }).catch(console.error);
  }, []);

  const handleSave = async () => {
    if (!jwt.user_id) return;
    setSaving(true);
    try {
      await api.put(`/user/update/${jwt.user_id}`, {
        first_name: firstName,
        last_name: lastName,
        avatar_url: avatarUrl || null,
        financial_method_id: financialMethodId || null,
        notifications_enabled: notificationsEnabled
      });
      toast.success('Perfil atualizado com sucesso!', {
        style: {
          backgroundColor: '#025439ff',
          color: '#fff',
          border: 'none',
        },
      });
      // Atualiza contexto global (Header)
      refreshProfile();
    } catch {
      toast.error('Erro ao atualizar perfil.', {
        style: {
          backgroundColor: '#7b0821ff',
          color: '#fff',
          border: 'none',
        }
      });
    } finally {
      setSaving(false);
    }
  };

  const handleOpen2FASetup = async () => {
    try {
      const res = await api.post('/auth/2fa/setup');
      setSetup2FAData(res.data);
      setIs2FASetupOpen(true);
    } catch (e: any) {
      toast.error(e.response?.data?.error || 'Erro ao preparar ativação do 2FA', {
        style: {
          backgroundColor: '#7b0821ff',
          color: '#fff',
          border: 'none',
        }
      });
    }
  };

  const handleVerify2FA = async () => {
    setVerifying2FA(true);
    try {
      await api.post('/auth/2fa/verify', { code: verifyCode });
      toast.success('2FA ativado com sucesso!', {
        style: {
          backgroundColor: '#025439ff',
          color: '#fff',
          border: 'none',
        },
      });
      setUser((prev) => ({ ...prev, two_fa_enabled: true }));
      setIs2FASetupOpen(false);
      setVerifyCode('');
    } catch (e: any) {
      toast.error(e.response?.data?.error || 'Código inválido', {
        style: {
          backgroundColor: '#7b0821ff',
          color: '#fff',
          border: 'none',
        }
      });
    } finally {
      setVerifying2FA(false);
    }
  };

  const handleDisable2FA = async () => {
    try {
      await api.post('/auth/2fa/disable');
      toast.success('2FA desativado com sucesso!');
      setUser((prev) => ({ ...prev, two_fa_enabled: false }));
    } catch (e: any) {
      toast.error(e.response?.data?.error || 'Erro ao desativar 2FA', {
        style: {
          backgroundColor: '#7b0821ff',
          color: '#fff',
          border: 'none',
        }
      });
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const reader = new FileReader();
      reader.onloadend = () => {
        const base64 = reader.result as string;
        setAvatarUrl(base64);
        // Propaga imediatamente para o Header
        updateAvatar(base64);
      };
      reader.readAsDataURL(file);
    }
  };

  const currentMethod = methods.find(m => m.id === financialMethodId);

  const plan = planConfig[user.plan || 'free'] || planConfig.free;
  const PlanIcon = plan.icon;

  const initials = ((firstName?.[0] || '') + (lastName?.[0] || '')).toUpperCase() || user.email?.[0]?.toUpperCase() || '?';

  return (
    <div className="space-y-6">
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main profile card */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Informações Pessoais</CardTitle>
            <CardDescription>Atualize suas informações de conta</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center gap-4">
              <Avatar className="h-20 w-20 text-2xl">
                <AvatarImage src={avatarUrl} alt="Avatar" className="object-cover" />
                <AvatarFallback className="bg-primary/10 text-primary text-xl font-bold">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 space-y-1">
                <p className="font-semibold">{firstName} {lastName}</p>
                <p className="text-sm text-muted-foreground">{user.email}</p>
                <Badge className={`mt-1 text-xs ${plan.color}`}>
                  <PlanIcon className="h-3 w-3 mr-1" /> Plano {plan.label}
                </Badge>
              </div>
            </div>

            <Separator />

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input placeholder="Seu nome" value={firstName} onChange={(e) => setFirstName(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Sobrenome</Label>
                <Input placeholder="Seu sobrenome" value={lastName} onChange={(e) => setLastName(e.target.value)} />
              </div>
            </div>

            <div className="space-y-2">
              <Label>Email</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input type="email" value={user.email} readOnly className="pl-10 bg-muted cursor-not-allowed" />
              </div>
              <p className="text-xs text-muted-foreground">O email não pode ser alterado.</p>
            </div>

            <Separator />

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="avatar-upload">Foto do Perfil</Label>
                <Input
                  id="avatar-upload"
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                />
              </div>
              <div className="space-y-2">
                <Label>Método de Planejamento</Label>
                {currentMethod ? (
                  <div className="flex items-center justify-between border rounded-md p-2">
                    <div className="flex flex-col">
                      <span className="font-semibold text-sm">{currentMethod.name}</span>
                      <span className="text-xs text-muted-foreground">{currentMethod.tagline}</span>
                    </div>
                    <Button variant="outline" size="sm" onClick={() => setIsMethodModalOpen(true)}>Alterar</Button>
                  </div>
                ) : (
                  <Button variant="outline" className="w-full text-muted-foreground justify-start font-normal h-10 px-3" onClick={() => setIsMethodModalOpen(true)}>
                    Sabe qual a melhor divisão para o seu dinheiro? Selecione.
                  </Button>
                )}
              </div>
            </div>

            <Button onClick={handleSave} disabled={saving} className="w-full">
              {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
              Salvar Alterações
            </Button>
          </CardContent>
        </Card>

        <div className="space-y-6">
          {/* Account Status */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Shield className="h-5 w-5" /> Status da Conta
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm">Email verificado</span>
                <Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300">✓ Sim</Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Autenticação 2FA</span>
                {user.two_fa_enabled ? (
                  <Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300">✓ Ativo</Badge>
                ) : (
                  <Badge variant="outline">Inativo</Badge>
                )}
              </div>
              {user.two_fa_enabled ? (
                <Button variant="destructive" size="sm" className="w-full" onClick={handleDisable2FA}>
                  Desativar 2FA
                </Button>
              ) : (
                <Button variant="outline" size="sm" className="w-full" onClick={handleOpen2FASetup}>
                  Ativar 2FA
                </Button>
              )}
            </CardContent>
          </Card>

          {/* Current Plan — free by default */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <CreditCard className="h-5 w-5" /> Plano Atual
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="text-center">
                <Badge className={`text-sm px-3 py-1 ${plan.color}`}>
                  <PlanIcon className="h-3.5 w-3.5 mr-1.5" /> {plan.label}
                </Badge>
                <p className="text-2xl font-bold mt-2">{plan.price}</p>
                {user.plan === 'free' && (
                  <p className="text-xs text-muted-foreground mt-1">
                    Faça upgrade para desbloquear análises avançadas
                  </p>
                )}
              </div>
              {user.plan !== 'free' ? (
                <Button variant="outline" size="sm" className="w-full" onClick={() => window.location.href = '/?page=subscription'}>
                  Gerenciar Assinatura
                </Button>
              ) : (
                <Button size="sm" className="w-full bg-gradient-to-r from-purple-600 to-blue-600 text-white" onClick={() => window.location.href = '/?page=subscription'}>
                  Fazer Upgrade
                </Button>
              )}
            </CardContent>
          </Card>

          {/* Notifications */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Bell className="h-5 w-5" /> Notificações
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="notificationsEnabled" className="text-sm cursor-pointer">Notificações por Email</Label>
                <Switch
                  id="notificationsEnabled"
                  checked={notificationsEnabled}
                  onCheckedChange={setNotificationsEnabled}
                />
              </div>
              <div className="flex items-center justify-between opacity-50">
                <Label htmlFor="push-disabled" className="text-sm cursor-pointer">Notificações Push (Em Breve)</Label>
                <Switch id="push-disabled" disabled checked={false} />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <FinancialMethodModal
        open={isMethodModalOpen}
        onOpenChange={setIsMethodModalOpen}
        currentMethodId={financialMethodId}
        onSelect={setFinancialMethodId}
      />

      {/* 2FA Setup Dialog */}
      <Dialog open={is2FASetupOpen} onOpenChange={setIs2FASetupOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Configurar Autenticação em 2 Etapas</DialogTitle>
            <DialogDescription>
              Escaneie o QRCode abaixo usando um aplicativo autenticador como Google Authenticator ou Authy.
            </DialogDescription>
          </DialogHeader>
          {setup2FAData && (
            <div className="flex flex-col items-center justify-center space-y-6 py-4">
              <div className="bg-white p-4 rounded-md shadow-sm">
                <QRCodeSVG value={setup2FAData.otpauth_url} size={200} />
              </div>
              <div className="w-full space-y-2">
                <Label>Código de Verificação</Label>
                <Input
                  type="text"
                  placeholder="Digite o código de 6 dígitos"
                  value={verifyCode}
                  onChange={(e) => setVerifyCode(e.target.value)}
                  maxLength={6}
                  className="text-center tracking-widest text-lg"
                />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setIs2FASetupOpen(false)}>Cancelar</Button>
            <Button onClick={handleVerify2FA} disabled={verifying2FA || verifyCode.length < 6}>
              {verifying2FA && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Ativar 2FA
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
