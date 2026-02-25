import { ReactElement, useTransition } from 'react';
import { 
  Activity, ShieldAlert, Target, TrendingDown, 
  TrendingUp, Search, Lock, UserCircle, PieChart 
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ResponsiveContainer, RadarChart, PolarGrid, PolarAngleAxis, Radar } from 'recharts';
import aiService from '@/services/aiService';
import { Skeleton } from '@/components/ui/skeleton';
import { useQuery } from '@tanstack/react-query';

// Mock data to demonstrate the Pro functionality
const profileData = [
  { subject: 'Poupança', A: 80, fullMark: 100 },
  { subject: 'Controle', A: 65, fullMark: 100 },
  { subject: 'Investimento', A: 40, fullMark: 100 },
  { subject: 'Lazer', A: 90, fullMark: 100 },
  { subject: 'Essenciais', A: 85, fullMark: 100 },
];

export function AiDiagnostic(): ReactElement {
  const [isPending, startTransition] = useTransition();

  const { data: res, isPending: isQueryPending, error: queryError } = useQuery({
    queryKey: ['aiDiagnostic'],
    queryFn: () => aiService.getDiagnostic(),
    retry: false, // Don't retry on 402/403
  });

  const loading = isPending || isQueryPending;
  const isPro = !(queryError && ((queryError as any).response?.status === 403 || (queryError as any).response?.status === 402 || (queryError as any).response?.data?.error?.includes('upgrade')));
  const insight = res?.diagnostic;

  if (loading) {
    return (
      <div className="space-y-6 max-w-6xl mx-auto py-10">
        <Skeleton className="h-[400px] w-full rounded-xl" />
      </div>
    );
  }

  if (!isPro) {
    return (
      <div className="max-w-4xl mx-auto text-center space-y-6 py-20 animate-in fade-in duration-700">
        <div className="inline-flex h-20 w-20 items-center justify-center rounded-full bg-primary/10 mb-4 ring-8 ring-primary/5">
          <Lock className="h-10 w-10 text-primary" />
        </div>
        <h1 className="text-4xl font-extrabold tracking-tight">Recurso Exclusivo Pro</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Faça upgrade para acessar o Detetive de Vazamentos, Relatórios Mensais aprofundados e descobrir o seu Perfil de Gasto.
        </p>
        <Button size="lg" className="mt-4 px-8 tracking-wide">
          Desbloquear Tudo
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-8 max-w-6xl mx-auto animate-in fade-in slide-in-from-bottom-6 duration-700">
      
      {/* Header section */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
        <div>
          <Badge variant="secondary" className="bg-fuchsia-500/10 text-fuchsia-600 border-fuchsia-500/20 mb-3 px-3 py-1">
            <Activity className="w-3 h-3 mr-1.5" /> Diagnóstico Pro
          </Badge>
          <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-slate-50">
            Saúde & Perfil de Gasto
          </h1>
          <p className="text-muted-foreground mt-2 text-lg">
            Análise comportamental das suas finanças baseada em Inteligência Artificial.
          </p>
        </div>
      </div>

      {/* Profile & Leak Detective Row */}
      <div className="grid gap-6 md:grid-cols-2">
        
        {/* User Profile Analysis */}
        <Card className="border-border/50 shadow-sm overflow-hidden flex flex-col">
          <CardHeader className="bg-slate-50/50 dark:bg-slate-900/50 border-b pb-6">
            <div className="flex justify-between items-start">
              <div>
                <CardTitle className="text-xl flex items-center gap-2">
                  <UserCircle className="w-5 h-5 text-primary" />
                  Seu Perfil Principal
                </CardTitle>
                <CardDescription className="mt-1.5">Análise dos seus hábitos de consumo</CardDescription>
              </div>
              <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/20 hover:bg-emerald-500/20">
                O Planejador Cauteloso
              </Badge>
            </div>
          </CardHeader>
          <CardContent className="pt-6 flex-1 flex flex-col md:flex-row items-center gap-6">
            <div className="w-full md:w-1/2 h-48">
              <ResponsiveContainer width="100%" height="100%">
                <RadarChart cx="50%" cy="50%" outerRadius="70%" data={profileData}>
                  <PolarGrid strokeOpacity={0.2} />
                  <PolarAngleAxis dataKey="subject" tick={{ fill: 'currentColor', fontSize: 11 }} />
                  <Radar name="Você" dataKey="A" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.3} />
                </RadarChart>
              </ResponsiveContainer>
            </div>
            <div className="w-full md:w-1/2 space-y-4">
              <p className="text-sm leading-relaxed text-muted-foreground">
                <span className="font-semibold text-foreground">Sua tendência:</span> Você demonstra alto controle nas despesas essenciais e no lazer, mas pode otimizar seus investimentos para fazer o dinheiro render mais.
              </p>
              <div className="space-y-2">
                <div className="flex items-center text-sm gap-2">
                  <Target className="w-4 h-4 text-primary" /> Foco: Construção de Reserva
                </div>
                <div className="flex items-center text-sm gap-2">
                  <TrendingUp className="w-4 h-4 text-emerald-500" /> Ponto Forte: Orçamento Fixo
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Leak Detective */}
        <Card className="border-rose-500/20 shadow-sm shadow-rose-500/5 bg-gradient-to-bl from-rose-500/5 to-transparent relative overflow-hidden flex flex-col">
          <div className="absolute -right-12 -top-12 w-40 h-40 bg-rose-500/10 rounded-full blur-3xl pointer-events-none" />
          <CardHeader className="pb-4">
            <CardTitle className="text-xl flex items-center gap-2 text-rose-600 dark:text-rose-400">
              <Search className="w-5 h-5" />
              Detetive de Vazamentos
            </CardTitle>
            <CardDescription>
              Compras ocultas ou pequenos hábitos que estão drenando a conta.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-5 flex-1">
            <div className="p-4 rounded-xl bg-background shadow-sm border border-rose-100 dark:border-rose-900/30 flex gap-4">
              <div className="w-10 h-10 rounded-full bg-rose-100 dark:bg-rose-900/30 flex items-center justify-center flex-shrink-0">
                <ShieldAlert className="w-5 h-5 text-rose-600" />
              </div>
              <div>
                <h4 className="font-semibold text-sm">Assinaturas Esquecidas</h4>
                <p className="text-xs text-muted-foreground mt-1 leading-relaxed">
                  Detectamos 2 cobranças recorrentes que não são utilizadas há mais de 3 meses (Total: R$ 89,90/mês).
                </p>
                <Button variant="link" className="text-rose-600 h-auto p-0 mt-2 text-xs">
                  Revisar assinaturas
                </Button>
              </div>
            </div>
            
            <div className="p-4 rounded-xl bg-background shadow-sm border border-amber-100 dark:border-amber-900/30 flex gap-4">
              <div className="w-10 h-10 rounded-full bg-amber-100 dark:bg-amber-900/30 flex items-center justify-center flex-shrink-0">
                <PieChart className="w-5 h-5 text-amber-600" />
              </div>
              <div>
                <h4 className="font-semibold text-sm">Pequenos Luxos Acumulados</h4>
                <p className="text-xs text-muted-foreground mt-1 leading-relaxed">
                  Gastos com "Delivery" subiram 24% em relação à sua média mensal. Isso compromete sua poupança.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Monthly Report Deep Dive */}
      <Card className="overflow-hidden border-border/50">
        <CardHeader className="bg-slate-50 dark:bg-slate-900/30 border-b">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <CardTitle className="text-xl">Relatório Mensal Profundo</CardTitle>
              <CardDescription className="mt-1">
                A IA analisa o contexto de cada real gasto no mês
              </CardDescription>
            </div>
            <Button variant="outline" size="sm">Baixar em PDF</Button>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          <div className="grid md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-border">
            <div className="p-6 space-y-2 lg:p-8">
              <div className="flex items-center text-sm font-medium text-muted-foreground mb-4">
                <TrendingDown className="w-4 h-4 mr-2 text-emerald-500" /> Redução de Custos
              </div>
              <h3 className="text-2xl font-bold">- 12.4%</h3>
              <p className="text-sm text-muted-foreground">Você economizou com transporte este mês comparado à media histórica.</p>
            </div>
            <div className="p-6 space-y-2 lg:p-8">
              <div className="flex items-center text-sm font-medium text-muted-foreground mb-4">
                <TrendingUp className="w-4 h-4 mr-2 text-rose-500" /> Inflação Pessoal
              </div>
              <h3 className="text-2xl font-bold">+ 5.2%</h3>
              <p className="text-sm text-muted-foreground">Seus custos e mercado tiveram alta. Verifique a categoria 'Supermercado'.</p>
            </div>
            <div className="p-6 space-y-2 lg:p-8 bg-slate-50/50 dark:bg-slate-900/20 col-span-1 md:col-span-3">
              <h4 className="font-semibold text-lg mb-4 text-primary">Conclusão da Inteligência Artificial</h4>
              {insight ? (
                <div className="text-sm text-muted-foreground whitespace-pre-wrap leading-relaxed">
                  {insight.content}
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">Diagnóstico não disponível.</p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
      
    </div>
  );
}
