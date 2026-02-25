import { useState, useTransition } from 'react';
import { api } from '@/api/client';
import { Loader2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { 
  Download, 
  TrendingUp, 
  TrendingDown, 
  Target,
  PieChart,
  Share2
} from 'lucide-react';
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  BarChart,
  Bar,
  LineChart,
  Line,
  Cell,
  PieChart as RechartsPieChart,
  Pie
} from 'recharts';
import { CustomTooltip, PieTooltip, BarTooltip } from '@/components/ui/custom-tooltip';

import { useQuery } from '@tanstack/react-query';

export function Reports() {
  const [selectedPeriod, setSelectedPeriod] = useState('thisMonth');
  const [reportType, setReportType] = useState('overview');

  const { data: reportsData, isPending: isQueryPending } = useQuery({
    queryKey: ['reportsData', selectedPeriod],
    queryFn: async () => {
      const [dashRes, goalsRes] = await Promise.all([
        api.get('/finance/dashboard'),
        api.get('/goals/').catch(() => ({ data: [] }))
      ]);
      return {
        summary: dashRes.data,
        goals: dashRes.data?.goals || goalsRes.data || []
      };
    }
  });

  const loading = isQueryPending;
  const summary = reportsData?.summary || null;
  const goals = reportsData?.goals || [];

  if (loading) {
    return (
      <div className="flex h-64 w-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  // Mapeamentos de dados pro Recharts e Interface
  const monthlyTrends = summary?.monthly_trend?.map((t: { month: string; income: number; expenses: number }) => ({
    month: t.month,
    receitas: t.income,
    despesas: t.expenses,
    economia: t.income - t.expenses,
  })) || [];

  const expenseDistribution = summary?.category_breakdown?.map((c: { category_name: string; percentage: number; color?: string; total: number }) => ({
    name: c.category_name,
    value: c.percentage,
    color: c.color || '#3B82F6',
    total: c.total,
  })) || [];

  const categoryComparison = summary?.category_breakdown?.map((c: { category_name: string; total: number }) => ({
    category: c.category_name,
    atual: c.total,
    anterior: 0, 
    meta: c.total * 1.1, 
  })) || [];

  const totalIncome = summary?.total_income || 0;
  const totalExpenses = summary?.total_expenses || 0;
  const totalSavings = totalIncome - totalExpenses;
  const savingsRate = summary?.savings_rate || 0;

  return (
    <div className="space-y-6">
      {/* Header com controles */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        <div>
          <h1 className="text-2xl font-bold">Relatórios Financeiros</h1>
          <p className="text-muted-foreground">
            Análise detalhada das suas finanças
          </p>
        </div>
        
        <div className="flex space-x-2">
          <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="thisWeek">Esta Semana</SelectItem>
              <SelectItem value="thisMonth">Este Mês</SelectItem>
              <SelectItem value="lastMonth">Mês Passado</SelectItem>
              <SelectItem value="thisYear">Este Ano</SelectItem>
              <SelectItem value="custom">Período Customizado</SelectItem>
            </SelectContent>
          </Select>
          
          <Button variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Exportar
          </Button>
          
          <Button variant="outline">
            <Share2 className="h-4 w-4 mr-2" />
            Compartilhar
          </Button>
        </div>
      </div>

      {/* Métricas Principais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Receitas Totais</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              R$ {totalIncome.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-green-500" />
              <span>Em relação ao período selecionado</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Despesas Totais</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              R$ {totalExpenses.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingDown className="h-3 w-3 text-red-500" />
              <span>Em relação ao período selecionado</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Economia Total</CardTitle>
            <Target className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              R$ {totalSavings.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-blue-500" />
              <span>Saldo final do período</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Taxa de Economia</CardTitle>
            <PieChart className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">{savingsRate.toFixed(1)}%</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-purple-500" />
              <span>Ideal acima de 20%</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={reportType} onValueChange={setReportType}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Visão Geral</TabsTrigger>
          <TabsTrigger value="trends">Tendências</TabsTrigger>
          <TabsTrigger value="categories">Categorias</TabsTrigger>
          <TabsTrigger value="goals">Metas</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Fluxo Financeiro */}
            <Card>
              <CardHeader>
                <CardTitle>Fluxo Financeiro Mensal</CardTitle>
                <CardDescription>
                  Evolução das receitas, despesas e economia
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <AreaChart data={monthlyTrends}>
                    <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                    <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                    <YAxis tick={{ fontSize: 12 }} />
                    <Tooltip content={<CustomTooltip />} />
                    <Area 
                      type="monotone" 
                      dataKey="receitas" 
                      stackId="1"
                      stroke="#10B981" 
                      fill="#10B981" 
                      fillOpacity={0.6}
                    />
                    <Area 
                      type="monotone" 
                      dataKey="despesas" 
                      stackId="2"
                      stroke="#EF4444" 
                      fill="#EF4444" 
                      fillOpacity={0.6}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {/* Distribuição de Gastos */}
            <Card>
              <CardHeader>
                <CardTitle>Distribuição de Gastos</CardTitle>
                <CardDescription>
                  Percentual por categoria de despesa
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <RechartsPieChart>
                    <Pie
                      data={expenseDistribution}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={120}
                      paddingAngle={5}
                      dataKey="value"
                    >
                      {expenseDistribution.map((entry: { name: string; value: number; color?: string; total: number }, index: number) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip content={<PieTooltip />} />
                  </RechartsPieChart>
                </ResponsiveContainer>
                <div className="grid grid-cols-2 gap-2 mt-4">
                  {expenseDistribution.map((item: { name: string; color?: string; value: number; total: number }, index: number) => (
                    <div key={index} className="flex items-center space-x-2 text-sm">
                      <div 
                        className="w-3 h-3 rounded-full" 
                        style={{ backgroundColor: item.color }}
                      />
                      <span>{item.name}: {item.value}%</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="trends" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Tendência de Economia</CardTitle>
              <CardDescription>
                Acompanhe como sua capacidade de economia evoluiu
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={monthlyTrends}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                  <YAxis tick={{ fontSize: 12 }} />
                  <Tooltip content={<CustomTooltip />} />
                  <Line 
                    type="monotone" 
                    dataKey="economia" 
                    stroke="#3B82F6" 
                    strokeWidth={3}
                    dot={{ r: 6, fill: '#3B82F6' }}
                    activeDot={{ r: 8, fill: '#3B82F6' }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="categories" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Comparação por Categorias</CardTitle>
              <CardDescription>
                Compare os gastos atuais com o período anterior e suas metas
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <BarChart data={categoryComparison}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis dataKey="category" tick={{ fontSize: 12 }} />
                  <YAxis tick={{ fontSize: 12 }} />
                  <Tooltip content={<BarTooltip />} />
                  <Bar dataKey="anterior" fill="#94A3B8" name="Período Anterior" radius={[2, 2, 0, 0]} />
                  <Bar dataKey="atual" fill="#3B82F6" name="Período Atual" radius={[2, 2, 0, 0]} />
                  <Bar dataKey="meta" fill="#10B981" name="Meta" radius={[2, 2, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          <div className="grid gap-4">
            {categoryComparison.map((category: { category: string; atual: number; anterior: number; meta: number }, index: number) => (
              <Card key={index}>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="font-medium">{category.category}</h4>
                      <p className="text-sm text-muted-foreground">
                        Meta: R$ {category.meta.toLocaleString('pt-BR')}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold">
                        R$ {category.atual.toLocaleString('pt-BR')}
                      </p>
                      <div className="flex items-center space-x-2">
                        {category.atual > category.anterior ? (
                          <TrendingUp className="h-4 w-4 text-red-500" />
                        ) : (
                          <TrendingDown className="h-4 w-4 text-green-500" />
                        )}
                        <span className={`text-sm ${
                          category.atual > category.anterior ? 'text-red-500' : 'text-green-500'
                        }`}>
                          {((category.atual - category.anterior) / category.anterior * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="goals" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Progresso das Metas</CardTitle>
                <CardDescription>
                  Acompanhe o progresso das suas metas financeiras
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {goals.length === 0 ? (
                  <div className="text-center py-6 text-muted-foreground">
                    <p>Você ainda não tem metas criadas.</p>
                  </div>
                ) : (
                  goals.map((goal: { id: string; name: string; current_amount: number; target_amount: number; deadline: string; color?: string }, index: number) => {
                    const progress = goal.target_amount ? (goal.current_amount / goal.target_amount) * 100 : 0;
                    return (
                      <div key={index} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <span className="font-medium">{goal.deadline || goal.name}</span>
                          <Badge variant="outline">
                            {goal.deadline 
                              ? new Date(goal.deadline).toLocaleDateString('pt-BR') 
                              : 'Sem prazo'}
                          </Badge>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                          <div 
                            className="bg-blue-600 h-2 rounded-full transition-all"
                            style={{ width: `${Math.min(progress, 100)}%` }}
                          />
                        </div>
                        <div className="flex justify-between text-sm text-muted-foreground">
                          <span>R$ {(goal.current_amount || 0).toLocaleString('pt-BR')}</span>
                          <span>{progress.toFixed(1)}%</span>
                          <span>R$ {(goal.target_amount || 0).toLocaleString('pt-BR')}</span>
                        </div>
                      </div>
                    );
                  })
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recomendações</CardTitle>
                <CardDescription>
                  Dicas para atingir suas metas mais rapidamente
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="p-4 bg-blue-50 dark:bg-blue-950 rounded-lg">
                  <h4 className="font-medium text-blue-800 dark:text-blue-200">
                    Aumente sua economia em 10%
                  </h4>
                  <p className="text-sm text-blue-600 dark:text-blue-300 mt-1">
                    Reduza os gastos com alimentação para atingir a meta de reserva mais rapidamente
                  </p>
                </div>
                
                <div className="p-4 bg-green-50 dark:bg-green-950 rounded-lg">
                  <h4 className="font-medium text-green-800 dark:text-green-200">
                    Meta da viagem em dia
                  </h4>
                  <p className="text-sm text-green-600 dark:text-green-300 mt-1">
                    Mantendo o ritmo atual, você atingirá a meta 2 meses antes do prazo
                  </p>
                </div>
                
                <div className="p-4 bg-yellow-50 dark:bg-yellow-950 rounded-lg">
                  <h4 className="font-medium text-yellow-800 dark:text-yellow-200">
                    Ajuste necessário
                  </h4>
                  <p className="text-sm text-yellow-600 dark:text-yellow-300 mt-1">
                    Para o novo carro, considere aumentar a economia mensal em R$ 200
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}