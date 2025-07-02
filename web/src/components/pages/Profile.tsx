import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import {
  Mail,
  Phone,
  MapPin,
  Bell,
  Shield,
  CreditCard,
  Download,
  Trash2,
} from 'lucide-react';

export function Profile() {
  return (
    <div className="space-y-6">
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Perfil Principal */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Informações Pessoais</CardTitle>
            <CardDescription>
              Atualize suas informações pessoais e de contato
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Avatar */}
            <div className="flex items-center space-x-4">
              <Avatar className="h-20 w-20">
                <AvatarImage
                  src="https://avatars.githubusercontent.com/u/14133553?s=400&u=7e45f209b9d7876ee65773357a0d6ef0fe64f012&v=4"
                  alt="Perfil"
                />
                <AvatarFallback>JD</AvatarFallback>
              </Avatar>
              <div className="space-y-2">
                <Button variant="outline" size="sm">
                  Alterar Foto
                </Button>
                <p className="text-xs text-muted-foreground">
                  JPG, PNG ou GIF (máx. 2MB)
                </p>
              </div>
            </div>

            <Separator />

            {/* Formulário */}
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="firstName">Nome</Label>
                <Input id="firstName" placeholder="João" defaultValue="João" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Sobrenome</Label>
                <Input
                  id="lastName"
                  placeholder="Silva"
                  defaultValue="da Silva"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="email"
                  type="email"
                  placeholder="joao@exemplo.com"
                  defaultValue="joao@exemplo.com"
                  className="pl-10"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="phone">Telefone</Label>
              <div className="relative">
                <Phone className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="phone"
                  placeholder="(11) 99999-9999"
                  defaultValue="(11) 99999-9999"
                  className="pl-10"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="location">Localização</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="location"
                  placeholder="São Paulo, SP"
                  defaultValue="São Paulo, SP"
                  className="pl-10"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="currency">Moeda Padrão</Label>
              <Select defaultValue="BRL">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="BRL">Real Brasileiro (R$)</SelectItem>
                  <SelectItem value="USD">Dólar Americano ($)</SelectItem>
                  <SelectItem value="EUR">Euro (€)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <Button className="w-full">Salvar Alterações</Button>
          </CardContent>
        </Card>

        {/* Sidebar com informações da conta */}
        <div className="space-y-6">
          {/* Status da Conta */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Shield className="h-5 w-5" />
                <span>Status da Conta</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm">Verificação de Email</span>
                <Badge
                  variant="default"
                  className="bg-green-100 text-green-800">
                  Verificado
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Autenticação 2FA</span>
                <Badge variant="outline">Inativo</Badge>
              </div>
              <Button variant="outline" size="sm" className="w-full">
                Ativar 2FA
              </Button>
            </CardContent>
          </Card>

          {/* Plano Atual */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <CreditCard className="h-5 w-5" />
                <span>Plano Atual</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="text-center">
                <Badge className="bg-gradient-to-r from-purple-500 to-blue-500 text-white">
                  Plano Pro
                </Badge>
                <p className="text-2xl font-bold mt-2">R$ 29,90/mês</p>
                <p className="text-sm text-muted-foreground">
                  Próxima cobrança: 15/02/2024
                </p>
              </div>
              <Button variant="outline" size="sm" className="w-full">
                Gerenciar Assinatura
              </Button>
            </CardContent>
          </Card>

          {/* Configurações Rápidas */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Bell className="h-5 w-5" />
                <span>Notificações</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="emailNotifications" className="text-sm">
                  Notificações por Email
                </Label>
                <Switch id="emailNotifications" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="pushNotifications" className="text-sm">
                  Notificações Push
                </Label>
                <Switch id="pushNotifications" />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="weeklyReports" className="text-sm">
                  Relatórios Semanais
                </Label>
                <Switch id="weeklyReports" defaultChecked />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Seção de Dados e Privacidade */}
      <Card>
        <CardHeader>
          <CardTitle className="text-red-600">Zona de Perigo</CardTitle>
          <CardDescription>Ações irreversíveis para sua conta</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between p-4 border border-yellow-200 rounded-lg bg-yellow-50 dark:bg-yellow-950 dark:border-yellow-800">
            <div>
              <h4 className="font-medium text-yellow-800 dark:text-yellow-200">
                Exportar Dados
              </h4>
              <p className="text-sm text-yellow-600 dark:text-yellow-300">
                Baixe uma cópia de todos os seus dados
              </p>
            </div>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Exportar
            </Button>
          </div>

          <div className="flex items-center justify-between p-4 border border-red-200 rounded-lg bg-red-50 dark:bg-red-950 dark:border-red-800">
            <div>
              <h4 className="font-medium text-red-800 dark:text-red-200">
                Excluir Conta
              </h4>
              <p className="text-sm text-red-600 dark:text-red-300">
                Exclua permanentemente sua conta e todos os dados
              </p>
            </div>
            <Button variant="destructive" size="sm">
              <Trash2 className="h-4 w-4 mr-2" />
              Excluir
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
