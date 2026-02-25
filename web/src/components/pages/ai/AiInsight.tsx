import { useTransition } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Bot, Sparkles, TrendingUp, AlertTriangle, ArrowRight, Wallet, CheckCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import aiService from '@/services/aiService';
import { Skeleton } from '@/components/ui/skeleton';

export function AiInsight() {
  const [isPending, startTransition] = useTransition();

  const { data: insight, isPending: isQueryPending, error: queryError } = useQuery({
    queryKey: ['aiInsight'],
    queryFn: () => aiService.getInsight(),
  });

  const loading = isPending || isQueryPending;
  const error = queryError ? (queryError as any).response?.data?.error || 'Erro ao carregar insight' : null;

  const typeConfig = {
    warning: { icon: AlertTriangle, color: 'text-amber-500', bg: 'bg-amber-500/10 border-amber-500/20', gradient: 'from-amber-500/20 to-orange-500/5' },
    tip: { icon: Bot, color: 'text-blue-500', bg: 'bg-blue-500/10 border-blue-500/20', gradient: 'from-blue-500/20 to-cyan-500/5' },
    achievement: { icon: CheckCircle, color: 'text-emerald-500', bg: 'bg-emerald-500/10 border-emerald-500/20', gradient: 'from-emerald-500/20 to-green-500/5' },
    projection: { icon: TrendingUp, color: 'text-purple-500', bg: 'bg-purple-500/10 border-purple-500/20', gradient: 'from-purple-500/20 to-fuchsia-500/5' },
  };

  const insightConfig = insight ? typeConfig[insight.type as keyof typeof typeConfig] || typeConfig.tip : typeConfig.tip;
  const ActiveIcon = insightConfig.icon;

  return (
    <div className="space-y-6 max-w-5xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <Badge variant="secondary" className="bg-primary/10 text-primary border-primary/20 uppercase">
              <Sparkles className="w-3 h-3 mr-1" /> Plano {insight?.plan || loading ? 'Carregando...' : 'Free'}
            </Badge>
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Visão Geral Diária</h1>
          <p className="text-muted-foreground mt-1">Sua dose diária de inteligência artificial para o bolso.</p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Insight Card - Takes up 2 cols on desktop */}
        <div className="md:col-span-2 space-y-6">
          <Card className={cn(
            "relative overflow-hidden border transition-all duration-300 shadow-sm hover:shadow-md",
            insightConfig.bg
          )}>
            {/* Background Gradient decoration */}
            <div className={cn(
              "absolute inset-0 bg-gradient-to-br opacity-50 pointer-events-none",
              insightConfig.gradient
            )} />
            
            <CardHeader className="relative z-10 pb-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className={cn("p-2 rounded-xl bg-background shadow-sm border", insightConfig.color)}>
                    <ActiveIcon className="h-5 w-5" />
                  </div>
                  <CardTitle className="text-xl">Dica do Dia</CardTitle>
                </div>
              </div>
            </CardHeader>
            <CardContent className="relative z-10">
              {loading ? (
                <div className="space-y-3 mt-4">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-5/6" />
                </div>
              ) : error ? (
                <div className="mt-4 p-4 bg-background/50 rounded-lg border border-dashed flex flex-col items-center justify-center text-center">
                  <AlertTriangle className="h-8 w-8 text-muted-foreground mb-2" />
                  <p className="text-sm font-medium">{error}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Verifique sua conexão ou tente novamente mais tarde.
                  </p>
                </div>
              ) : insight ? (
                <div className="mt-2 space-y-4">
                  <h3 className="text-2xl font-semibold tracking-tight leading-tight">
                    {insight.title}
                  </h3>
                  <p className="text-muted-foreground leading-relaxed text-base">
                    {insight.content}
                  </p>
                </div>
              ) : null}
            </CardContent>
          </Card>

          {/* Upgrade Teaser for Pro features */}
          <Card className="bg-gradient-to-br from-slate-900 to-slate-800 text-slate-50 border-slate-700 overflow-hidden relative group">
            <div className="absolute right-0 top-0 w-64 h-64 bg-primary/20 blur-[80px] rounded-full group-hover:bg-primary/30 transition-all duration-500" />
            <CardContent className="p-6 relative z-10 flex flex-col justify-between h-full min-h-[160px]">
              <div className="flex justify-between items-start gap-4">
                <div>
                  <Badge variant="outline" className="text-xs border-slate-600 text-slate-300 mb-3">
                    Desbloqueie o Pro
                  </Badge>
                  <h3 className="text-lg font-semibold mb-1 text-white">Diagnóstico Completo</h3>
                  <p className="text-slate-400 text-sm max-w-[200px]">
                    Descubra para onde seu dinheiro escapa com relatórios automáticos.
                  </p>
                  <Button variant="link" className="px-0 text-primary mt-2">
                    Fazer Upgrade
                  </Button>
                </div>
                <div className="h-10 w-10 rounded-full bg-slate-800 border border-slate-700 flex flex-col items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                  <ArrowRight className="h-4 w-4 text-primary" />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Status / Context */}
        <div className="space-y-6 flex flex-col h-full">
          <Card className="h-full flex flex-col">
            <CardHeader className="pb-4">
              <CardTitle className="text-sm font-medium flex items-center gap-2">
                <Wallet className="h-4 w-4 text-muted-foreground" />
                Contexto Analisado
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6 flex-1">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Status do Mês</span>
                  <span className="font-medium text-emerald-500">Dentro da Meta</span>
                </div>
                <div className="h-2 bg-secondary rounded-full overflow-hidden">
                  <div className="h-full bg-emerald-500 w-[60%]" />
                </div>
              </div>
              
              <div className="pt-4 border-t space-y-4">
                <p className="text-xs text-muted-foreground uppercase tracking-wider font-semibold">Fatores Chave</p>
                <div className="flex items-start gap-3 text-sm">
                  <div className="w-2 h-2 rounded-full bg-blue-500 mt-1.5 flex-shrink-0" />
                  <span className="leading-snug">Seu saldo está positivo e cobrirá as despesas fixas.</span>
                </div>
                <div className="flex items-start gap-3 text-sm">
                  <div className="w-2 h-2 rounded-full bg-amber-500 mt-1.5 flex-shrink-0" />
                  <span className="leading-snug">Atenção com gastos em "Alimentação" este fim de semana.</span>
                </div>
                <div className="flex items-start gap-3 text-sm">
                  <div className="w-2 h-2 rounded-full bg-emerald-500 mt-1.5 flex-shrink-0" />
                  <span className="leading-snug">Você guardou 10% da sua renda. Bom trabalho!</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
