import { useState } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  TrendingUp,
  TrendingDown,
  DollarSign,
  PlusCircle,
  Calendar,
  Target,
  AlertCircle,
  CheckCircle,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { CustomTooltip, PieTooltip } from '@/components/ui/custom-tooltip';

const monthlyData = [
  { month: 'Jan', receitas: 4500, despesas: 3200 },
  { month: 'Fev', receitas: 4800, despesas: 3100 },
  { month: 'Mar', receitas: 5200, despesas: 3800 },
  { month: 'Abr', receitas: 4900, despesas: 3600 },
  { month: 'Mai', receitas: 5500, despesas: 4200 },
  { month: 'Jun', receitas: 5800, despesas: 4100 },
];

const expenseCategories = [
  { name: 'Alimentação', value: 1200, color: '#3B82F6', total: 3300 },
  { name: 'Transporte', value: 800, color: '#10B981', total: 3300 },
  { name: 'Saúde', value: 600, color: '#F59E0B', total: 3300 },
  { name: 'Lazer', value: 400, color: '#EF4444', total: 3300 },
  { name: 'Outros', value: 300, color: '#8B5CF6', total: 3300 },
];

const goals = [
  {
    name: 'Reserva de Emergência',
    current: 8500,
    target: 15000,
    color: 'bg-green-500',
  },
  { name: 'Viagem', current: 2800, target: 5000, color: 'bg-blue-500' },
  { name: 'Novo Carro', current: 12000, target: 25000, color: 'bg-purple-500' },
];

export function Dashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState('month');

  const totalReceitas = 5800;
  const totalDespesas = 4100;
  const saldo = totalReceitas - totalDespesas;

  return (
    <div className="space-y-6">
      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Receitas</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              R$ {totalReceitas.toLocaleString('pt-BR')}
            </div>
            <p className="text-xs text-muted-foreground">
              +12% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Despesas</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              R$ {totalDespesas.toLocaleString('pt-BR')}
            </div>
            <p className="text-xs text-muted-foreground">
              -5% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo</CardTitle>
            <DollarSign className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div
              className={cn(
                'text-2xl font-bold',
                saldo >= 0 ? 'text-green-600' : 'text-red-600',
              )}>
              R$ {saldo.toLocaleString('pt-BR')}
            </div>
            <p className="text-xs text-muted-foreground">Economia mensal</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Meta do Mês</CardTitle>
            <Target className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">85%</div>
            <p className="text-xs text-muted-foreground">
              R$ 1.500 de R$ 2.000
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Gráfico Principal */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Fluxo Financeiro</CardTitle>
                <CardDescription>
                  Acompanhe suas receitas e despesas ao longo do tempo
                </CardDescription>
              </div>
              <Tabs value={selectedPeriod} onValueChange={setSelectedPeriod}>
                <TabsList>
                  <TabsTrigger value="week">7D</TabsTrigger>
                  <TabsTrigger value="month">30D</TabsTrigger>
                  <TabsTrigger value="year">1A</TabsTrigger>
                </TabsList>
              </Tabs>
            </div>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={monthlyData}>
                <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                <XAxis
                  dataKey="month"
                  className="text-xs"
                  tick={{ fontSize: 12 }}
                />
                <YAxis className="text-xs" tick={{ fontSize: 12 }} />
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

        {/* Categorias de Gastos */}
        <Card>
          <CardHeader>
            <CardTitle>Categorias de Gastos</CardTitle>
            <CardDescription>
              Distribuição das suas despesas este mês
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  data={expenseCategories}
                  cx="50%"
                  cy="50%"
                  innerRadius={40}
                  outerRadius={80}
                  paddingAngle={5}
                  dataKey="value">
                  {expenseCategories.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip content={<PieTooltip />} />
              </PieChart>
            </ResponsiveContainer>
            <div className="space-y-2 mt-4">
              {expenseCategories.map((category, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between text-sm">
                  <div className="flex items-center space-x-2">
                    <div
                      className="w-3 h-3 rounded-full"
                      style={{ backgroundColor: category.color }}
                    />
                    <span>{category.name}</span>
                  </div>
                  <span className="font-medium">
                    R$ {category.value.toLocaleString('pt-BR')}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Metas Financeiras */}
        <Card>
          <CardHeader>
            <CardTitle>Metas Financeiras</CardTitle>
            <CardDescription>
              Acompanhe o progresso das suas metas
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {goals.map((goal, index) => {
              const progress = (goal.current / goal.target) * 100;
              return (
                <div key={index} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">{goal.name}</span>
                    <span className="text-sm text-muted-foreground">
                      {progress.toFixed(0)}%
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                    <div
                      className={cn(
                        'h-2 rounded-full transition-all',
                        goal.color,
                      )}
                      style={{ width: `${Math.min(progress, 100)}%` }}
                    />
                  </div>
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>R$ {goal.current.toLocaleString('pt-BR')}</span>
                    <span>R$ {goal.target.toLocaleString('pt-BR')}</span>
                  </div>
                </div>
              );
            })}
          </CardContent>
        </Card>

        {/* Insights e Dicas */}
        <Card>
          <CardHeader>
            <CardTitle>Insights Financeiros</CardTitle>
            <CardDescription>
              Dicas personalizadas para melhorar suas finanças
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-start space-x-3 p-3 bg-green-50 dark:bg-green-950 rounded-lg">
              <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-green-800 dark:text-green-200">
                  Parabéns! Você economizou 15% este mês
                </p>
                <p className="text-xs text-green-600 dark:text-green-300">
                  Continue assim para atingir suas metas mais rapidamente
                </p>
              </div>
            </div>

            <div className="flex items-start space-x-3 p-3 bg-yellow-50 dark:bg-yellow-950 rounded-lg">
              <AlertCircle className="h-5 w-5 text-yellow-500 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-yellow-800 dark:text-yellow-200">
                  Gastos com alimentação aumentaram 8%
                </p>
                <p className="text-xs text-yellow-600 dark:text-yellow-300">
                  Considere planejar refeições para reduzir custos
                </p>
              </div>
            </div>

            <div className="flex items-start space-x-3 p-3 bg-blue-50 dark:bg-blue-950 rounded-lg">
              <Target className="h-5 w-5 text-blue-500 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-blue-800 dark:text-blue-200">
                  Meta de reserva de emergência
                </p>
                <p className="text-xs text-blue-600 dark:text-blue-300">
                  Faltam apenas R$ 6.500 para completar sua reserva
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Ações Rápidas */}
      <div className="flex flex-wrap gap-4">
        <Button className="bg-green-600 hover:bg-green-700">
          <PlusCircle className="h-4 w-4 mr-2" />
          Adicionar Receita
        </Button>
        <Button
          variant="outline"
          className="border-red-200 text-red-600 hover:bg-red-50">
          <PlusCircle className="h-4 w-4 mr-2" />
          Adicionar Despesa
        </Button>
        <Button variant="secondary">
          <Calendar className="h-4 w-4 mr-2" />
          Planejar Orçamento
        </Button>
      </div>
    </div>
  );
}
