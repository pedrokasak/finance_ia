import React, { useState } from 'react';
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
  Calendar,
  Target,
  PieChart,
  BarChart3,
  FileText,
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

const monthlyTrends = [
  { month: 'Jan', receitas: 4500, despesas: 3200, economia: 1300 },
  { month: 'Fev', receitas: 4800, despesas: 3100, economia: 1700 },
  { month: 'Mar', receitas: 5200, despesas: 3800, economia: 1400 },
  { month: 'Abr', receitas: 4900, despesas: 3600, economia: 1300 },
  { month: 'Mai', receitas: 5500, despesas: 4200, economia: 1300 },
  { month: 'Jun', receitas: 5800, despesas: 4100, economia: 1700 },
];

const categoryComparison = [
  { category: 'Alimentação', atual: 1200, anterior: 1100, meta: 1000 },
  { category: 'Transporte', atual: 800, anterior: 850, meta: 700 },
  { category: 'Saúde', atual: 600, anterior: 400, meta: 500 },
  { category: 'Lazer', atual: 400, anterior: 300, meta: 350 },
  { category: 'Casa', atual: 300, anterior: 280, meta: 250 },
];

const expenseDistribution = [
  { name: 'Alimentação', value: 35, color: '#3B82F6', total: 100 },
  { name: 'Transporte', value: 25, color: '#10B981', total: 100 },
  { name: 'Saúde', value: 18, color: '#F59E0B', total: 100 },
  { name: 'Lazer', value: 12, color: '#EF4444', total: 100 },
  { name: 'Outros', value: 10, color: '#8B5CF6', total: 100 },
];

export function Reports() {
  const [selectedPeriod, setSelectedPeriod] = useState('thisMonth');
  const [reportType, setReportType] = useState('overview');

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
            <div className="text-2xl font-bold text-green-600">R$ 32.400</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-green-500" />
              <span>+12% vs período anterior</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Despesas Totais</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">R$ 22.800</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingDown className="h-3 w-3 text-green-500" />
              <span>-5% vs período anterior</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Economia Total</CardTitle>
            <Target className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">R$ 9.600</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-green-500" />
              <span>+18% vs período anterior</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Taxa de Economia</CardTitle>
            <PieChart className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">29.6%</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              <TrendingUp className="h-3 w-3 text-green-500" />
              <span>+2.1% vs período anterior</span>
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
                      {expenseDistribution.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip content={<PieTooltip />} />
                  </RechartsPieChart>
                </ResponsiveContainer>
                <div className="grid grid-cols-2 gap-2 mt-4">
                  {expenseDistribution.map((item, index) => (
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
            {categoryComparison.map((category, index) => (
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
                {[
                  { name: 'Reserva de Emergência', current: 8500, target: 15000, deadline: '2024-12-31' },
                  { name: 'Viagem Europa', current: 2800, target: 5000, deadline: '2024-07-01' },
                  { name: 'Novo Carro', current: 12000, target: 25000, deadline: '2025-06-01' },
                ].map((goal, index) => {
                  const progress = (goal.current / goal.target) * 100;
                  return (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="font-medium">{goal.name}</span>
                        <Badge variant="outline">
                          {new Date(goal.deadline).toLocaleDateString('pt-BR')}
                        </Badge>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                        <div 
                          className="bg-blue-600 h-2 rounded-full transition-all"
                          style={{ width: `${Math.min(progress, 100)}%` }}
                        />
                      </div>
                      <div className="flex justify-between text-sm text-muted-foreground">
                        <span>R$ {goal.current.toLocaleString('pt-BR')}</span>
                        <span>{progress.toFixed(1)}%</span>
                        <span>R$ {goal.target.toLocaleString('pt-BR')}</span>
                      </div>
                    </div>
                  );
                })}
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