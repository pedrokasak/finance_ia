import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import {
  Check,
  Crown,
  Star,
  TrendingUp,
  CreditCard,
  Calendar,
  Download,
  Loader2,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { api } from '@/api/client';
import { toast } from 'sonner';

// Plan data comes from the backend (GET /subscription/plans)
interface Plan {
  id: string; // same as slug: 'free' | 'pro' | 'premium'
  slug: string;
  name: string;
  description: string;
  price_monthly: number;
  price_yearly: number;
  features: Array<string | { description: string }>;
  limitations?: Array<string | { description: string }>;
  popular?: boolean;
  is_active: boolean;
}

// Visual metadata per slug (icon/color cannot come from JSON)
const planVisuals: Record<
  string,
  { icon: React.ElementType; color: string; bgColor: string; popular?: boolean }
> = {
  free: { icon: TrendingUp, color: 'text-gray-600', bgColor: 'bg-gray-100' },
  pro: {
    icon: Crown,
    color: 'text-blue-600',
    bgColor: 'bg-blue-100',
    popular: true,
  },
  premium: { icon: Star, color: 'text-purple-600', bgColor: 'bg-purple-100' },
};

interface SubscriptionInfo {
  plan: string;
  status: string;
  current_period_end?: string;
  cancel_at_period_end?: boolean;
}

interface Invoice {
  id: string;
  date: string;
  amount: number;
  status: string;
  currency: string;
  hosted_invoice_url?: string;
}

export function Subscription() {
  const [billingInterval, setBillingInterval] = useState<'monthly' | 'yearly'>(
    'monthly',
  );
  const [loadingPlan, setLoadingPlan] = useState<string | null>(null);

  // Toast on return from Stripe Checkout
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('success')) {
      toast.success('Assinatura realizada com sucesso!', { duration: 5000 });
      window.history.replaceState(null, '', '/?page=subscription');
    } else if (params.get('canceled')) {
      toast.info('Checkout cancelado. Nenhuma cobrança foi feita.');
      window.history.replaceState(null, '', '/?page=subscription');
    }
  }, []);

  const {
    data: subData,
    isLoading: loadingSub,
    isError: subError,
  } = useQuery({
    queryKey: ['subscription'],
    queryFn: () => api.get('/subscription/').then((r) => r.data),
    staleTime: 1000 * 30, // 30s
    retry: 1,
    refetchOnWindowFocus: false,
  });

  const {
    data: plansData,
    isLoading: loadingPlans,
    isError: plansError,
  } = useQuery({
    queryKey: ['subscription-plans'],
    queryFn: () => api.get('/subscription/plans').then((r) => r.data),
    staleTime: 1000 * 60 * 5, // 5min — plans rarely change
    retry: 1,
    refetchOnWindowFocus: false,
  });

  const { data: invoicesData } = useQuery({
    queryKey: ['subscription-invoices'],
    queryFn: () => api.get('/subscription/invoices').then((r) => r.data),
    staleTime: 1000 * 60, // 1min
  });

  const loadingData = loadingSub || loadingPlans;
  const planRequestFailed = subError || plansError;

  const subInfo: SubscriptionInfo = {
    plan: subData?.subscription?.plan || subData?.plan || 'free',
    status: subData?.subscription?.status || subData?.status || 'active',
    current_period_end: subData?.subscription?.current_period_end,
    cancel_at_period_end: subData?.subscription?.cancel_at_period_end,
  };

  const invoices: Invoice[] = invoicesData?.invoices || invoicesData || [];

  const backendPlans: Plan[] = (plansData?.plans || []).map((p: Plan) => ({
    id: p.slug, // use slug as id so plan.id === currentPlan works
    slug: p.slug,
    name: p.name,
    description: p.description,
    price_monthly: p.price_monthly,
    price_yearly: p.price_yearly,
    features: p.features || [],
    limitations: p.limitations || [],
    popular: planVisuals[p.slug]?.popular ?? false,
    is_active: p.is_active,
  }));

  const fallbackPlans: Plan[] = [
    {
      id: 'free',
      slug: 'free',
      name: 'Básico',
      description: 'Perfeito para começar',
      price_monthly: 0,
      price_yearly: 0,
      features: ['Controle básico de receitas e despesas'],
      is_active: true,
    },
    {
      id: 'pro',
      slug: 'pro',
      name: 'Pro',
      description: 'Ideal para controle avançado',
      price_monthly: 29.9,
      price_yearly: 299,
      features: ['Transações ilimitadas', 'Relatórios avançados'],
      is_active: true,
    },
    {
      id: 'premium',
      slug: 'premium',
      name: 'Premium',
      description: 'Máximo controle financeiro com IA',
      price_monthly: 49.9,
      price_yearly: 499,
      features: ['Todos os recursos do Pro', 'Análise de IA personalizada'],
      is_active: true,
    },
  ];

  const plans: Plan[] =
    backendPlans.length > 0 ? backendPlans : fallbackPlans;

  const currentPlan = subInfo.plan || 'free';
  const isPaidPlan = currentPlan !== 'free';

  const handleSubscribe = async (planId: string) => {
    // 'free' plan = no payment needed (downgrade handled via portal)
    if (planId === 'free') {
      if (isPaidPlan) {
        handleManagePortal();
      }
      return;
    }
    setLoadingPlan(planId);
    try {
      const { data } = await api.post('/subscription/checkout', {
        plan: planId,
        billing_type: billingInterval,
      });
      if (data.url) {
        window.location.href = data.url;
      }
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(
        err.response?.data?.error || 'Erro ao criar sessão de pagamento',
      );
    } finally {
      setLoadingPlan(null);
    }
  };

  const handleManagePortal = async () => {
    try {
      const { data } = await api.post('/subscription/portal');
      if (data.url) window.location.href = data.url;
    } catch {
      toast.error('Erro ao abrir portal de gerenciamento');
    }
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="text-center space-y-4">
        <h1 className="text-3xl font-bold">Escolha seu Plano</h1>
        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
          Desbloqueie todo o potencial do FinanceAI com nossos planos flexíveis
        </p>

        {/* Toggle de cobrança */}
        <div className="flex items-center justify-center space-x-4">
          <Label
            htmlFor="billing-toggle"
            className={cn(
              'text-sm font-medium',
              billingInterval === 'monthly'
                ? 'text-foreground'
                : 'text-muted-foreground',
            )}>
            Mensal
          </Label>
          <Switch
            id="billing-toggle"
            checked={billingInterval === 'yearly'}
            onCheckedChange={(checked) =>
              setBillingInterval(checked ? 'yearly' : 'monthly')
            }
          />
          <Label
            htmlFor="billing-toggle"
            className={cn(
              'text-sm font-medium',
              billingInterval === 'yearly'
                ? 'text-foreground'
                : 'text-muted-foreground',
            )}>
            Anual
          </Label>
          {billingInterval === 'yearly' && (
            <Badge className="bg-green-100 text-green-800">
              2 meses grátis
            </Badge>
          )}
        </div>
      </div>

      {/* Cards dos Planos */}
      {planRequestFailed && (
        <div className="text-center text-sm text-amber-600">
          Não foi possível sincronizar todos os dados de assinatura agora.
          Exibindo planos padrão.
        </div>
      )}
      {loadingData ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <div className="grid gap-6 lg:grid-cols-3">
          {plans.map((plan) => {
            const visuals = planVisuals[plan.slug] ?? {
              icon: TrendingUp,
              color: 'text-gray-600',
              bgColor: 'bg-gray-100',
            };
            const Icon = visuals.icon;
            const isCurrentPlan = plan.id === currentPlan; // plan.id === slug === currentPlan from API
            const price =
              billingInterval === 'yearly'
                ? plan.price_yearly
                : plan.price_monthly;

            return (
              <Card
                key={plan.id}
                className={cn(
                  'relative transition-all duration-200 hover:shadow-lg',
                  visuals.popular && 'ring-2 ring-blue-500 shadow-lg scale-105',
                  isCurrentPlan &&
                    'bg-blue-50 dark:bg-blue-950 border-blue-200',
                )}>
                {visuals.popular && (
                  <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                    <Badge className="bg-blue-600 text-white px-3 py-1">
                      Mais Popular
                    </Badge>
                  </div>
                )}

                {isCurrentPlan && (
                  <div className="absolute -top-3 right-4">
                    <Badge className="bg-green-600 text-white px-3 py-1">
                      Plano Atual
                    </Badge>
                  </div>
                )}

                <CardHeader className="text-center space-y-4">
                  <div
                    className={cn(
                      'w-12 h-12 mx-auto rounded-full flex items-center justify-center',
                      visuals.bgColor,
                    )}>
                    <Icon className={cn('h-6 w-6', visuals.color)} />
                  </div>

                  <div>
                    <CardTitle className="text-xl">{plan.name}</CardTitle>
                    <CardDescription>{plan.description}</CardDescription>
                  </div>

                  <div className="space-y-1">
                    <div className="flex items-baseline justify-center">
                      <span className="text-3xl font-bold">
                        R$ {price.toFixed(2).replace('.', ',')}
                      </span>
                      <span className="text-muted-foreground ml-1">
                        /{billingInterval === 'yearly' ? 'ano' : 'mês'}
                      </span>
                    </div>
                    {billingInterval === 'yearly' && plan.price_monthly > 0 && (
                      <p className="text-sm text-muted-foreground">
                        ou R$ {plan.price_monthly.toFixed(2).replace('.', ',')}
                        /mês cobrado anualmente
                      </p>
                    )}
                  </div>
                </CardHeader>

                <CardContent className="space-y-6">
                  {/* Features */}
                  <div className="space-y-3">
                    {plan.features.map((feature, index: number) => {
                      const text =
                        typeof feature === 'string'
                          ? feature
                          : feature.description;
                      return (
                        <div key={index} className="flex items-start space-x-3">
                          <Check className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                          <span className="text-sm">{text}</span>
                        </div>
                      );
                    })}
                  </div>

                  {/* Limitations */}
                  {(plan.limitations ?? []).length > 0 && (
                    <div className="space-y-2 pt-4 border-t">
                      <p className="text-sm font-medium text-muted-foreground">
                        Limitações:
                      </p>
                      {(plan.limitations ?? []).map(
                        // eslint-disable-next-line @typescript-eslint/no-explicit-any
                        (limitation: any, index: number) => {
                          const text =
                            typeof limitation === 'string'
                              ? limitation
                              : limitation.description;
                          return (
                            <div
                              key={index}
                              className="flex items-start space-x-3">
                              <div className="w-4 h-4 mt-0.5 flex-shrink-0">
                                <div className="w-1 h-1 bg-gray-400 rounded-full mx-auto mt-1.5" />
                              </div>
                              <span className="text-sm text-muted-foreground">
                                {text}
                              </span>
                            </div>
                          );
                        },
                      )}
                    </div>
                  )}

                  {/* Action Button */}
                  <Button
                    className={cn(
                      'w-full',
                      isCurrentPlan && 'opacity-50 cursor-not-allowed',
                      visuals.popular &&
                        !isCurrentPlan &&
                        'bg-blue-600 hover:bg-blue-700',
                    )}
                    variant={plan.id === 'free' ? 'outline' : 'default'}
                    disabled={isCurrentPlan || loadingPlan !== null}
                    onClick={() => handleSubscribe(plan.id)}>
                    {loadingPlan === plan.id ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : isCurrentPlan ? (
                      '✓ Plano Atual'
                    ) : plan.id === 'free' ? (
                      isPaidPlan ? (
                        'Cancelar Assinatura'
                      ) : (
                        'Plano Gratuito'
                      )
                    ) : currentPlan === 'free' ? (
                      'Fazer Upgrade'
                    ) : (
                      'Alterar Plano'
                    )}
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      {/* Informações de Pagamento */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <CreditCard className="h-5 w-5" />
              <span>Método de Pagamento</span>
            </CardTitle>
            <CardDescription>
              Gerencie suas formas de pagamento via portal Stripe
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {currentPlan !== 'free' ? (
              <Button
                variant="outline"
                className="w-full"
                onClick={handleManagePortal}>
                <CreditCard className="h-4 w-4 mr-2" />
                Abrir Portal de Gerenciamento
              </Button>
            ) : (
              <p className="text-sm text-muted-foreground text-center py-4">
                Faça upgrade para gerenciar métodos de pagamento.
              </p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Calendar className="h-5 w-5" />
              <span>Histórico de Cobranças</span>
            </CardTitle>
            <CardDescription>Suas últimas faturas</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {loadingData ? (
              <div className="flex justify-center py-6">
                <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
              </div>
            ) : invoices.length > 0 ? (
              <div className="space-y-3">
                {invoices.map((invoice) => (
                  <div
                    key={invoice.id}
                    className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">{invoice.date}</p>
                      <p className="text-sm text-muted-foreground">
                        R$ {invoice.amount.toFixed(2).replace('.', ',')}
                      </p>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Badge
                        className={
                          invoice.status === 'paid'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-yellow-100 text-yellow-800'
                        }>
                        {invoice.status === 'paid' ? 'Pago' : invoice.status}
                      </Badge>
                      {invoice.hosted_invoice_url && (
                        <Button variant="ghost" size="sm" asChild>
                          <a
                            href={invoice.hosted_invoice_url}
                            target="_blank"
                            rel="noreferrer">
                            <Download className="h-4 w-4" />
                          </a>
                        </Button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground text-center py-4">
                Nenhuma cobrança encontrada.
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recursos Exclusivos */}
      <Card>
        <CardHeader>
          <CardTitle>Recursos por Plano</CardTitle>
          <CardDescription>
            Compare os recursos disponíveis em cada plano
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-3 pr-4">Recurso</th>
                  <th className="text-center py-3 px-4">Básico</th>
                  <th className="text-center py-3 px-4">Pro</th>
                  <th className="text-center py-3 px-4">Premium</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {[
                  {
                    name: 'Transações mensais',
                    basic: '100',
                    pro: 'Ilimitadas',
                    premium: 'Ilimitadas',
                  },
                  {
                    name: 'Categorização automática',
                    basic: '❌',
                    pro: '✅',
                    premium: '✅',
                  },
                  {
                    name: 'Relatórios avançados',
                    basic: '❌',
                    pro: '✅',
                    premium: '✅',
                  },
                  {
                    name: 'Exportação de dados',
                    basic: '❌',
                    pro: '✅',
                    premium: '✅',
                  },
                  {
                    name: 'Análise de IA',
                    basic: '❌',
                    pro: '❌',
                    premium: '✅',
                  },
                  {
                    name: 'Consultoria automatizada',
                    basic: '❌',
                    pro: '❌',
                    premium: '✅',
                  },
                  {
                    name: 'Suporte',
                    basic: 'Email',
                    pro: 'Prioritário',
                    premium: '24/7 Chat',
                  },
                ].map((feature, index) => (
                  <tr key={index}>
                    <td className="py-3 pr-4 font-medium">{feature.name}</td>
                    <td className="text-center py-3 px-4">{feature.basic}</td>
                    <td className="text-center py-3 px-4">{feature.pro}</td>
                    <td className="text-center py-3 px-4">{feature.premium}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
