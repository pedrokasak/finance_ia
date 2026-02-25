import { useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Area, AreaChart } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Slider } from '@/components/ui/slider';
import { Badge } from '@/components/ui/badge';
import { LineChart as LineChartIcon, Lock, Sparkles, TrendingUp } from 'lucide-react';
import aiService from '@/services/aiService';
import { Skeleton } from '@/components/ui/skeleton';
import { useQuery } from '@tanstack/react-query';

// Mock data generator for projection
const generateProjectionData = (months: number, baseSavings: number, extraInvestment: number) => {
  const data = [];
  const currentBalance = 10000;
  for (let i = 0; i < months; i++) {
    const conservative = currentBalance + (baseSavings + extraInvestment) * i * 1.05;
    const aggressive = currentBalance + (baseSavings + extraInvestment) * i * 1.1;
    data.push({
      month: `Mês ${i + 1}`,
      conservative: Math.round(conservative),
      aggressive: Math.round(aggressive),
    });
  }
  return data;
};

// Mock data for daily balance
const dailyBalanceData = Array.from({ length: 30 }).map((_, i) => {
  const balance = 2500 - (i * 80) + (i === 15 ? 4000 : 0); // salary drop on day 15
  return {
    day: i + 1,
    balance: Math.round(balance),
  };
});

export function AiSimulator() {
  const { data: res, isPending: isQueryPending, error: queryError } = useQuery({
    queryKey: ['aiSimulator'],
    queryFn: () => aiService.getSimulator(),
    retry: false,
  });

  const loading = isQueryPending;
  const err = queryError as { response?: { status?: number; data?: { error?: string } } } | null;
  const isPro = !(err && (err.response?.status === 403 || err.response?.status === 402 || err.response?.data?.error?.includes('upgrade')));
  const insight = res?.projection;
  
  // Simulator states
  const [extraInvestment, setExtraInvestment] = useState([500]);
  const [months, setMonths] = useState([12]);
  
  const projectionData = generateProjectionData(months[0], 1000, extraInvestment[0]);

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
          Faça upgrade para acessar o Simulador de Futuro Interativo e a Previsão Inteligente de Saldo Diário.
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
          <Badge variant="secondary" className="bg-blue-500/10 text-blue-600 border-blue-500/20 mb-3 px-3 py-1">
            <LineChartIcon className="w-3 h-3 mr-1.5" /> Previsão Avançada
          </Badge>
          <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-slate-50">
            Simulador "What-If" Interativo
          </h1>
          <p className="text-muted-foreground mt-2 text-lg">
            Descubra o impacto das suas decisões de hoje no amanhã.
          </p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Interactive Controls */}
        <Card className="md:col-span-1 border-border/50 shadow-sm">
          <CardHeader className="bg-slate-50/50 dark:bg-slate-900/50 border-b">
            <CardTitle className="text-lg flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" /> Cenários
            </CardTitle>
            <CardDescription>
              Ajuste as variáveis e veja a mágica acontecer.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-8 pt-6">
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <label className="text-sm font-medium">Aporte Mensal Extra</label>
                <span className="text-sm font-bold text-emerald-500">R$ {extraInvestment[0]}</span>
              </div>
              <Slider
                value={extraInvestment}
                onValueChange={setExtraInvestment}
                max={5000}
                step={100}
                className="py-4"
              />
              <p className="text-xs text-muted-foreground text-balance">
                "O que acontece se eu investir esse valor a mais todo mês, ao invés de gastar no iFood?"
              </p>
            </div>

            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <label className="text-sm font-medium">Tempo (Meses)</label>
                <span className="text-sm font-bold text-primary">{months[0]}m</span>
              </div>
              <Slider
                value={months}
                onValueChange={setMonths}
                max={60}
                min={6}
                step={6}
                className="py-4"
              />
            </div>
            
            <div className="p-4 rounded-xl bg-primary/5 border border-primary/20 space-y-2 mt-4">
              <div className="flex items-center gap-2 font-medium text-primary">
                <TrendingUp className="h-4 w-4" /> Inteligência Artificial Conclui:
              </div>
              {insight ? (
                <p className="text-sm whitespace-pre-wrap leading-relaxed text-muted-foreground mt-2">
                  {insight.content}
                </p>
              ) : (
                <p className="text-sm text-muted-foreground">Projeção não disponível.</p>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Projection Chart */}
        <Card className="md:col-span-2 border-border/50 shadow-sm overflow-hidden flex flex-col">
          <CardHeader className="pb-2">
            <CardTitle>Crescimento Projetado</CardTitle>
          </CardHeader>
          <CardContent className="flex-1 pt-4 min-h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={projectionData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorAggressive" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorConservative" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <XAxis dataKey="month" fontSize={12} tickLine={false} axisLine={false} />
                <YAxis fontSize={12} tickLine={false} axisLine={false} tickFormatter={(value) => `R$${value/1000}k`} />
                <Tooltip 
                  formatter={(value: number) => [`R$ ${value}`, ""]}
                  contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                />
                <Area type="monotone" dataKey="aggressive" name="Agressivo (10% a.a.)" stroke="#10b981" fillOpacity={1} fill="url(#colorAggressive)" />
                <Area type="monotone" dataKey="conservative" name="Conservador (5% a.a.)" stroke="#3b82f6" fillOpacity={1} fill="url(#colorConservative)" />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Daily Balance Forecast */}
      <Card className="border-border/50 overflow-hidden">
        <CardHeader className="bg-slate-50 dark:bg-slate-900/30 border-b">
          <div className="flex justify-between items-center">
            <div>
              <CardTitle className="text-xl">Radar de Fluxo Diário</CardTitle>
              <CardDescription className="mt-1">
                A IA prevê o saldo da sua conta dia a dia para o próximo mês para evitar cheque especial.
              </CardDescription>
            </div>
            <Badge variant="outline" className="text-rose-500 border-rose-500/50 bg-rose-500/10">
              Alerta Dia 13: Risco de Saldo Negativo
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="h-64 pt-6">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={dailyBalanceData}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.3} />
              <XAxis dataKey="day" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis fontSize={12} tickLine={false} axisLine={false} />
              <Tooltip 
                formatter={(value: number) => [`R$ ${value}`, "Saldo Previsto"]}
                labelFormatter={(label) => `Dia ${label}`}
              />
              <Line 
                type="monotone" 
                dataKey="balance" 
                stroke="#6366f1" 
                strokeWidth={3}
                dot={(props) => {
                  const { cx, cy, payload } = props;
                  if (payload.balance < 500) {
                    return <circle cx={cx} cy={cy} r={4} fill="#ef4444" stroke="none" />;
                  }
                  return <circle cx={cx} cy={cy} r={0} />;
                }}
              />
            </LineChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
      
    </div>
  );
}
