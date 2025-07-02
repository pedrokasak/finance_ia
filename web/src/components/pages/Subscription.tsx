import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { 
  Check, 
  Crown, 
  Zap, 
  Star,
  CreditCard,
  Calendar,
  TrendingUp,
  Shield,
  Download,
  Bot,
  BarChart3
} from 'lucide-react';
import { cn } from '@/lib/utils';

const plans = [
  {
    id: 'basic',
    name: 'Básico',
    price: 0,
    interval: 'mês',
    description: 'Perfeito para começar',
    icon: TrendingUp,
    color: 'text-gray-600',
    bgColor: 'bg-gray-100',
    features: [
      'Controle básico de receitas e despesas',
      'Categorização manual',
      'Gráficos simples',
      'Até 100 transações/mês',
      'Suporte por email'
    ],
    limitations: [
      'Sem relatórios avançados',
      'Sem exportação de dados',
      'Sem insights de IA'
    ]
  },
  {
    id: 'pro',
    name: 'Pro',
    price: 29.90,
    interval: 'mês',
    description: 'Ideal para controle avançado',
    icon: Crown,
    color: 'text-blue-600',
    bgColor: 'bg-blue-100',
    popular: true,
    features: [
      'Transações ilimitadas',
      'Categorização automática',
      'Gráficos avançados e interativos',
      'Metas financeiras personalizadas',
      'Relatórios detalhados',
      'Exportação em Excel/PDF',
      'Sincronização em nuvem',
      'Suporte prioritário'
    ],
    limitations: []
  },
  {
    id: 'premium',
    name: 'Premium',
    price: 49.90,
    interval: 'mês',
    description: 'Máximo controle financeiro',
    icon: Star,
    color: 'text-purple-600',
    bgColor: 'bg-purple-100',
    features: [
      'Todos os recursos do Pro',
      'Análise de IA personalizada',
      'Projeções e previsões avançadas',
      'Alertas inteligentes',
      'Consultoria financeira automatizada',
      'Dashboard executivo',
      'API para integrações',
      'Suporte 24/7 via chat'
    ],
    limitations: []
  }
];

const currentPlan = 'pro';

export function Subscription() {
  const [billingInterval, setBillingInterval] = useState<'monthly' | 'yearly'>('monthly');

  const getDiscountedPrice = (price: number) => {
    return billingInterval === 'yearly' ? price * 10 : price; // 2 meses grátis no anual
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
          <Label htmlFor="billing-toggle" className={cn(
            "text-sm font-medium",
            billingInterval === 'monthly' ? 'text-foreground' : 'text-muted-foreground'
          )}>
            Mensal
          </Label>
          <Switch
            id="billing-toggle"
            checked={billingInterval === 'yearly'}
            onCheckedChange={(checked) => setBillingInterval(checked ? 'yearly' : 'monthly')}
          />
          <Label htmlFor="billing-toggle" className={cn(
            "text-sm font-medium",
            billingInterval === 'yearly' ? 'text-foreground' : 'text-muted-foreground'
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
      <div className="grid gap-6 lg:grid-cols-3">
        {plans.map((plan) => {
          const Icon = plan.icon;
          const isCurrentPlan = plan.id === currentPlan;
          const discountedPrice = getDiscountedPrice(plan.price);
          
          return (
            <Card 
              key={plan.id} 
              className={cn(
                "relative transition-all duration-200 hover:shadow-lg",
                plan.popular && "ring-2 ring-blue-500 shadow-lg scale-105",
                isCurrentPlan && "bg-blue-50 dark:bg-blue-950 border-blue-200"
              )}
            >
              {plan.popular && (
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
                <div className={cn("w-12 h-12 mx-auto rounded-full flex items-center justify-center", plan.bgColor)}>
                  <Icon className={cn("h-6 w-6", plan.color)} />
                </div>
                
                <div>
                  <CardTitle className="text-xl">{plan.name}</CardTitle>
                  <CardDescription>{plan.description}</CardDescription>
                </div>
                
                <div className="space-y-1">
                  <div className="flex items-baseline justify-center">
                    <span className="text-3xl font-bold">
                      R$ {discountedPrice.toFixed(2).replace('.', ',')}
                    </span>
                    <span className="text-muted-foreground ml-1">
                      /{billingInterval === 'yearly' ? 'ano' : 'mês'}
                    </span>
                  </div>
                  {billingInterval === 'yearly' && plan.price > 0 && (
                    <p className="text-sm text-muted-foreground">
                      ou R$ {plan.price.toFixed(2).replace('.', ',')}/mês cobrado anualmente
                    </p>
                  )}
                </div>
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Features */}
                <div className="space-y-3">
                  {plan.features.map((feature, index) => (
                    <div key={index} className="flex items-start space-x-3">
                      <Check className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                      <span className="text-sm">{feature}</span>
                    </div>
                  ))}
                </div>

                {/* Limitations */}
                {plan.limitations.length > 0 && (
                  <div className="space-y-2 pt-4 border-t">
                    <p className="text-sm font-medium text-muted-foreground">Limitações:</p>
                    {plan.limitations.map((limitation, index) => (
                      <div key={index} className="flex items-start space-x-3">
                        <div className="w-4 h-4 mt-0.5 flex-shrink-0">
                          <div className="w-1 h-1 bg-gray-400 rounded-full mx-auto mt-1.5" />
                        </div>
                        <span className="text-sm text-muted-foreground">{limitation}</span>
                      </div>
                    ))}
                  </div>
                )}

                {/* Action Button */}
                <Button 
                  className={cn(
                    "w-full",
                    isCurrentPlan && "opacity-50 cursor-not-allowed",
                    plan.popular && "bg-blue-600 hover:bg-blue-700"
                  )}
                  variant={plan.id === 'basic' ? 'outline' : 'default'}
                  disabled={isCurrentPlan}
                >
                  {isCurrentPlan ? 'Plano Atual' : 
                   plan.id === 'basic' ? 'Downgrade' : 
                   currentPlan === 'basic' ? 'Upgrade' : 'Alterar Plano'}
                </Button>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Informações de Pagamento */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <CreditCard className="h-5 w-5" />
              <span>Método de Pagamento</span>
            </CardTitle>
            <CardDescription>
              Gerencie suas formas de pagamento
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex items-center space-x-3">
                <div className="w-8 h-8 bg-blue-100 rounded flex items-center justify-center">
                  <CreditCard className="h-4 w-4 text-blue-600" />
                </div>
                <div>
                  <p className="font-medium">•••• •••• •••• 1234</p>
                  <p className="text-sm text-muted-foreground">Cartão principal</p>
                </div>
              </div>
              <Badge>Ativo</Badge>
            </div>
            
            <Button variant="outline" className="w-full">
              Alterar Método de Pagamento
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Calendar className="h-5 w-5" />
              <span>Histórico de Cobranças</span>
            </CardTitle>
            <CardDescription>
              Veja suas últimas faturas
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-3">
              {[
                { date: '15/01/2024', amount: 29.90, status: 'Pago' },
                { date: '15/12/2023', amount: 29.90, status: 'Pago' },
                { date: '15/11/2023', amount: 29.90, status: 'Pago' },
              ].map((invoice, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">{invoice.date}</p>
                    <p className="text-sm text-muted-foreground">
                      R$ {invoice.amount.toFixed(2).replace('.', ',')}
                    </p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge className="bg-green-100 text-green-800">
                      {invoice.status}
                    </Badge>
                    <Button variant="ghost" size="sm">
                      <Download className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
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
                  { name: 'Transações mensais', basic: '100', pro: 'Ilimitadas', premium: 'Ilimitadas' },
                  { name: 'Categorização automática', basic: '❌', pro: '✅', premium: '✅' },
                  { name: 'Relatórios avançados', basic: '❌', pro: '✅', premium: '✅' },
                  { name: 'Exportação de dados', basic: '❌', pro: '✅', premium: '✅' },
                  { name: 'Análise de IA', basic: '❌', pro: '❌', premium: '✅' },
                  { name: 'Consultoria automatizada', basic: '❌', pro: '❌', premium: '✅' },
                  { name: 'Suporte', basic: 'Email', pro: 'Prioritário', premium: '24/7 Chat' },
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