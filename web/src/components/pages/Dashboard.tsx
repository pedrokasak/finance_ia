import { useEffect, useState, useCallback } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  TrendingUp,
  TrendingDown,
  DollarSign,
  PlusCircle,
  Target,
  AlertCircle,
  CheckCircle,
  Bot,
  Flame,
  Loader2,
  AlertTriangle,
  FolderOpen,
  Edit2,
  Trash2,
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
import financeService, {
  DashboardSummary,
  Category,
} from '@/services/financeService';
import aiService, { AIInsight } from '@/services/aiService';
import { toast } from 'sonner';

const levelColors: Record<string, string> = {
  Diamante: 'text-cyan-400',
  Platina: 'text-slate-300',
  Ouro: 'text-yellow-400',
  Prata: 'text-gray-300',
  Bronze: 'text-amber-600',
};

const levelEmojis: Record<string, string> = {
  Diamante: '💎',
  Platina: '⭐',
  Ouro: '🥇',
  Prata: '🥈',
  Bronze: '🥉',
};

function HealthScoreRing({ score, level }: { score: number; level: string }) {
  const radius = 40;
  const circumference = 2 * Math.PI * radius;
  const progress = (score / 1000) * circumference;

  return (
    <div className="flex flex-col items-center gap-1">
      <div className="relative w-24 h-24">
        <svg className="w-24 h-24 -rotate-90" viewBox="0 0 100 100">
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke="currentColor"
            strokeWidth="8"
            className="text-muted/30"
          />
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke="currentColor"
            strokeWidth="8"
            strokeDasharray={`${progress} ${circumference}`}
            strokeLinecap="round"
            className={cn(
              'transition-all duration-1000',
              levelColors[level] || 'text-primary',
            )}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="text-xl font-black">{score}</span>
          <span className="text-xs text-muted-foreground">/1000</span>
        </div>
      </div>
      <Badge variant="secondary" className="text-xs">
        {levelEmojis[level]} {level}
      </Badge>
    </div>
  );
}

function AIInsightCard({ plan }: { plan?: string }) {
  const [insight, setInsight] = useState<AIInsight | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    aiService
      .getInsight()
      .then(setInsight)
      .catch((err) => {
        const msg = err.response?.data?.error || 'Erro ao carregar insight';
        setError(msg);
      })
      .finally(() => setLoading(false));
  }, []);

  const isRateLimited = error?.includes('rate_limited');
  const needsUpgrade = error?.includes('upgrade');

  const typeStyles = {
    warning: 'bg-yellow-500/10 border-yellow-500/30 text-yellow-600',
    tip: 'bg-blue-500/10 border-blue-500/30 text-blue-600',
    achievement: 'bg-green-500/10 border-green-500/30 text-green-600',
    projection: 'bg-purple-500/10 border-purple-500/30 text-purple-600',
  };

  return (
    <div className="flex items-start space-x-3 p-3 bg-muted/30 rounded-lg border">
      <Bot className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <p className="text-sm font-medium">IA Financeira</p>
          {plan && (
            <Badge variant="outline" className="text-xs px-1.5 py-0 capitalize">
              {plan}
            </Badge>
          )}
        </div>
        {loading && (
          <div className="flex items-center gap-2 text-muted-foreground">
            <Loader2 className="h-3 w-3 animate-spin" />
            <span className="text-xs">Analisando seus dados...</span>
          </div>
        )}
        {insight && !loading && (
          <div
            className={cn(
              'rounded-md p-2 border text-xs leading-relaxed',
              typeStyles[insight.type] || typeStyles.tip,
            )}>
            <p className="font-medium mb-0.5">{insight.title}</p>
            <p>{insight.content}</p>
          </div>
        )}
        {!loading && (isRateLimited || needsUpgrade) && (
          <div className="text-xs text-muted-foreground">
            <p>Limite de insights da semana atingido.</p>
            <Button variant="link" className="h-auto p-0 text-xs text-primary">
              Faça upgrade para mais insights →
            </Button>
          </div>
        )}
        {!loading && error && !isRateLimited && !needsUpgrade && (
          <p className="text-xs text-muted-foreground">
            Configure sua chave de IA para ver insights personalizados.
          </p>
        )}
      </div>
    </div>
  );
}

export function Dashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState('month');
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [loading, setLoading] = useState(true);

  // Quick-add transaction modal state
  const [txDialogOpen, setTxDialogOpen] = useState(false);
  const [txType, setTxType] = useState<'income' | 'expense'>('income');
  const [txDesc, setTxDesc] = useState('');
  const [txAmount, setTxAmount] = useState('');
  const [txDate, setTxDate] = useState(new Date().toISOString().split('T')[0]);
  const [txCatId, setTxCatId] = useState('');
  const [categories, setCategories] = useState<Category[]>([]);
  const [txSaving, setTxSaving] = useState(false);

  // Category management modal state
  const [catDialogOpen, setCatDialogOpen] = useState(false);
  const [catSaving, setCatSaving] = useState(false);
  const [catEditId, setCatEditId] = useState<string | null>(null);
  const [catName, setCatName] = useState('');
  const [catType, setCatType] = useState<'income' | 'expense'>('expense');
  const [catColor, setCatColor] = useState('#EF4444');

  const openTxDialog = (type: 'income' | 'expense' = 'expense') => {
    setTxType(type);
    setTxDesc('');
    setTxAmount('');
    setTxCatId('');
    setTxDate(new Date().toISOString().split('T')[0]);
    setTxDialogOpen(true);
  };

  const handleTxSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!txAmount || parseFloat(txAmount) <= 0) {
      toast.error('Valor inválido');
      return;
    }
    setTxSaving(true);
    try {
      await financeService.createTransaction({
        type: txType,
        amount: parseFloat(txAmount),
        description: txDesc,
        date: txDate,
        category_id: txCatId || undefined,
      });
      toast.success(
        `${txType === 'income' ? 'Receita' : 'Despesa'} adicionada!`,
      );
      setTxDialogOpen(false);
      fetchDashboard();
    } catch {
      toast.error('Erro ao salvar');
    } finally {
      setTxSaving(false);
    }
  };

  const fetchDashboard = useCallback(() => {
    setLoading(true);
    fetchCategories();
    financeService
      .getDashboard()
      .then(setSummary)
      .catch(console.error)
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const fetchCategories = useCallback(() => {
    financeService
      .getCategories()
      .then((cats) => setCategories(cats || []))
      .catch(console.error);
  }, []);

  const handleCatSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!catName) return;
    setCatSaving(true);
    try {
      if (catEditId) {
        await financeService.updateCategory(catEditId, {
          name: catName,
          type: catType,
          color: catColor,
          icon: 'plus',
        });
        toast.success('Categoria atualizada!');
      } else {
        await financeService.createCategory({
          name: catName,
          type: catType,
          color: catColor,
          icon: 'plus',
        });
        toast.success('Categoria criada!');
      }
      setCatEditId(null);
      setCatName('');
      fetchCategories();
    } catch {
      toast.error('Erro ao salvar categoria');
    } finally {
      setCatSaving(false);
    }
  };

  const handleDeleteCategory = async (id: string) => {
    if (
      !confirm(
        'Deseja realmente apagar esta categoria? Transações atreladas ficarão sem categoria.',
      )
    )
      return;
    try {
      await financeService.deleteCategory(id);
      toast.success('Categoria removida');
      fetchCategories();
    } catch {
      toast.error(
        'Erro ao remover (pode ser padrão ou haver problemas de rede)',
      );
    }
  };

  useEffect(() => {
    fetchDashboard();
  }, [fetchDashboard]);

  // Fallback to mock data while loading or if no data
  const totalIncome = summary?.total_income ?? 0;
  const totalExpenses = summary?.total_expenses ?? 0;
  const balance = summary?.balance ?? 0;
  const healthScore = summary?.health_score ?? 0;
  const healthLevel = summary?.health_level ?? 'Bronze';
  const savingsRate = summary?.savings_rate ?? 0;
  const monthlyTrend = summary?.monthly_trend ?? [];
  const categoryBreakdown = summary?.category_breakdown ?? [];
  const daysUntilNegative = summary?.days_until_negative;
  const budget = summary?.budget;

  const COLORS = [
    '#3B82F6',
    '#10B981',
    '#F59E0B',
    '#EF4444',
    '#8B5CF6',
    '#EC4899',
    '#06B6D4',
  ];
  const filteredCats = categories.filter(
    (c) => c.type === (txType === 'income' ? 'income' : 'expense'),
  );

  return (
    <div className="space-y-6">
      {/* New Transaction Dialog */}
      <Dialog open={txDialogOpen} onOpenChange={setTxDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              Nova {txType === 'income' ? 'Receita' : 'Despesa'}
            </DialogTitle>
            <DialogDescription>
              Adicione uma movimentação ao seu controle financeiro.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleTxSubmit} className="space-y-4">
            <Tabs
              value={txType}
              onValueChange={(v) => {
                setTxType(v as 'income' | 'expense');
                setTxCatId('');
              }}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="income" className="text-green-600">
                  <TrendingUp className="h-4 w-4 mr-2" /> Receita
                </TabsTrigger>
                <TabsTrigger value="expense" className="text-red-600">
                  <TrendingDown className="h-4 w-4 mr-2" /> Despesa
                </TabsTrigger>
              </TabsList>
            </Tabs>
            <div className="space-y-2">
              <Label>Descrição</Label>
              <Input
                placeholder="Ex: Salário, Supermercado..."
                value={txDesc}
                onChange={(e) => setTxDesc(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label>Valor (R$)</Label>
              <Input
                type="number"
                placeholder="0,00"
                step="0.01"
                min="0.01"
                value={txAmount}
                onChange={(e) => setTxAmount(e.target.value)}
                required
              />
            </div>
            {filteredCats.length > 0 && (
              <div className="space-y-2">
                <Label>Categoria</Label>
                <Select value={txCatId} onValueChange={setTxCatId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Selecionar categoria" />
                  </SelectTrigger>
                  <SelectContent>
                    {filteredCats.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
            <div className="space-y-2">
              <Label>Data</Label>
              <Input
                type="date"
                value={txDate}
                onChange={(e) => setTxDate(e.target.value)}
                required
              />
            </div>
            <div className="flex gap-2 pt-1">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                onClick={() => setTxDialogOpen(false)}>
                Cancelar
              </Button>
              <Button
                type="submit"
                disabled={txSaving}
                className={`flex-1 ${txType === 'income' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'}`}>
                {txSaving ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  'Salvar'
                )}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Category Management Dialog */}
      <Dialog open={catDialogOpen} onOpenChange={setCatDialogOpen}>
        <DialogContent className="sm:max-w-md max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Gerenciar Categorias</DialogTitle>
            <DialogDescription>
              Crie ou edite as suas categorias financeiras.
            </DialogDescription>
          </DialogHeader>

          <form
            onSubmit={handleCatSubmit}
            className="space-y-4 pt-2 border-b pb-6">
            <div className="flex items-center gap-2">
              <h4 className="text-sm font-medium">
                {catEditId ? 'Editar Categoria' : 'Nova Categoria'}
              </h4>
              {catEditId && (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-6 text-xs"
                  onClick={() => {
                    setCatEditId(null);
                    setCatName('');
                  }}>
                  Cancelar Edição
                </Button>
              )}
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2 col-span-2">
                <Label>Nome</Label>
                <Input
                  placeholder="Ex: Assinaturas"
                  value={catName}
                  onChange={(e) => setCatName(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Tipo</Label>
                <Select
                  value={catType}
                  onValueChange={(v: 'income' | 'expense') => setCatType(v)}
                  disabled={!!catEditId}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="expense">Despesa</SelectItem>
                    <SelectItem value="income">Receita</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Cor</Label>
                <div className="flex gap-2">
                  <Input
                    type="color"
                    value={catColor}
                    onChange={(e) => setCatColor(e.target.value)}
                    className="w-12 h-10 p-1"
                  />
                  <Input
                    value={catColor}
                    onChange={(e) => setCatColor(e.target.value)}
                    placeholder="#000000"
                    className="flex-1"
                  />
                </div>
              </div>
            </div>
            <Button type="submit" disabled={catSaving} className="w-full">
              {catSaving ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : catEditId ? (
                'Salvar Alterações'
              ) : (
                'Criar Categoria'
              )}
            </Button>
          </form>

          <div className="space-y-2">
            <h4 className="text-sm font-medium mt-2">Suas Categorias</h4>
            {categories.length === 0 && (
              <p className="text-xs text-muted-foreground">
                Nenhuma categoria encontrada.
              </p>
            )}
            {categories.map((c) => (
              <div
                key={c.id}
                className="flex flex-row items-center justify-between p-2 rounded-md border text-sm">
                <div className="flex items-center gap-2">
                  <div
                    className="w-3 h-3 rounded-full"
                    style={{ backgroundColor: c.color }}
                  />
                  <span>{c.name}</span>
                  <Badge variant="outline" className="text-[10px] ml-1">
                    {c.type === 'income' ? 'Receita' : 'Despesa'}
                  </Badge>
                  {c.is_default && (
                    <Badge variant="secondary" className="text-[10px]">
                      Padrão
                    </Badge>
                  )}
                </div>
                {!c.is_default && (
                  <div className="flex items-center">
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-blue-500"
                      onClick={() => {
                        setCatEditId(c.id);
                        setCatName(c.name);
                        setCatType(c.type);
                        setCatColor(c.color);
                      }}>
                      <Edit2 className="h-3 w-3" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-red-500"
                      onClick={() => handleDeleteCategory(c.id)}>
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                )}
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      {/* Quick action buttons */}
      <div className="flex flex-wrap gap-3">
        <Button
          onClick={() => openTxDialog()}
          className="bg-primary hover:bg-primary/90 text-primary-foreground">
          <PlusCircle className="h-4 w-4 mr-2" /> Nova Transação
        </Button>
        <Button
          onClick={() => setCatDialogOpen(true)}
          variant="outline"
          className="font-medium">
          <FolderOpen className="h-4 w-4 mr-2 text-primary" /> Gerenciar
          Categorias
        </Button>
      </div>

      {daysUntilNegative !== undefined && daysUntilNegative <= 15 && (
        <div className="flex items-center gap-3 p-4 bg-red-500/10 border border-red-500/30 rounded-lg">
          <AlertTriangle className="h-5 w-5 text-red-500 flex-shrink-0" />
          <div>
            <p className="text-sm font-medium text-red-600 dark:text-red-400">
              Risco de ficar no vermelho em {daysUntilNegative} dias
            </p>
            <p className="text-xs text-red-500/80">
              Reduza seus gastos variáveis para evitar saldo negativo.
            </p>
          </div>
        </div>
      )}

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Receitas</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {loading ? (
                <Loader2 className="h-6 w-6 animate-spin" />
              ) : (
                `R$ ${totalIncome.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`
              )}
            </div>
            <p className="text-xs text-muted-foreground">Este mês</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Despesas</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {loading ? (
                <Loader2 className="h-6 w-6 animate-spin" />
              ) : (
                `R$ ${totalExpenses.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`
              )}
            </div>
            <p className="text-xs text-muted-foreground">Este mês</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo</CardTitle>
            <DollarSign
              className={cn(
                'h-4 w-4',
                balance >= 0 ? 'text-blue-500' : 'text-red-500',
              )}
            />
          </CardHeader>
          <CardContent>
            <div
              className={cn(
                'text-2xl font-bold',
                balance >= 0 ? 'text-blue-600' : 'text-red-600',
              )}>
              {loading ? (
                <Loader2 className="h-6 w-6 animate-spin" />
              ) : (
                `R$ ${balance.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Taxa de poupança: {savingsRate.toFixed(1)}%
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Score de Saúde
            </CardTitle>
            <Flame className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Loader2 className="h-6 w-6 animate-spin" />
            ) : (
              <HealthScoreRing score={healthScore} level={healthLevel} />
            )}
          </CardContent>
        </Card>
      </div>

      {/* Budget distribution */}
      {budget && (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Distribuição do Orçamento</CardTitle>
            <CardDescription>
              Renda: R${' '}
              {budget.total_income.toLocaleString('pt-BR', {
                minimumFractionDigits: 2,
              })}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-3">
              {[
                {
                  label: 'Necessidades',
                  amount: budget.needs_amount ?? 0,
                  percent: budget.needs_percent,
                  color: 'text-emerald-500',
                  bg: 'bg-emerald-500/10',
                },
                {
                  label: 'Desejos',
                  amount: budget.wants_amount ?? 0,
                  percent: budget.wants_percent,
                  color: 'text-blue-500',
                  bg: 'bg-blue-500/10',
                },
                {
                  label: 'Investimentos',
                  amount: budget.savings_amount ?? 0,
                  percent: budget.savings_percent,
                  color: 'text-purple-500',
                  bg: 'bg-purple-500/10',
                },
              ].map((item) => (
                <div
                  key={item.label}
                  className={cn('rounded-lg p-3 text-center', item.bg)}>
                  <p className={cn('text-lg font-bold', item.color)}>
                    {item.percent}%
                  </p>
                  <p className="text-xs text-muted-foreground">{item.label}</p>
                  <p className="text-xs font-medium mt-0.5">
                    R${' '}
                    {item.amount.toLocaleString('pt-BR', {
                      minimumFractionDigits: 0,
                    })}
                  </p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Chart */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Fluxo Financeiro</CardTitle>
                <CardDescription>
                  Receitas vs Despesas ao longo do tempo
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
            {monthlyTrend.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={monthlyTrend}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                  <YAxis tick={{ fontSize: 12 }} />
                  <Tooltip content={<CustomTooltip />} />
                  <Area
                    type="monotone"
                    dataKey="income"
                    stackId="1"
                    stroke="#10B981"
                    fill="#10B981"
                    fillOpacity={0.6}
                    name="Receitas"
                  />
                  <Area
                    type="monotone"
                    dataKey="expenses"
                    stackId="2"
                    stroke="#EF4444"
                    fill="#EF4444"
                    fillOpacity={0.6}
                    name="Despesas"
                  />
                </AreaChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                <div className="text-center">
                  <Target className="h-12 w-12 mx-auto mb-3 opacity-30" />
                  <p className="text-sm">
                    Adicione transações para ver o gráfico
                  </p>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Category Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>Categorias de Gastos</CardTitle>
            <CardDescription>
              Distribuição das despesas este mês
            </CardDescription>
          </CardHeader>
          <CardContent>
            {categoryBreakdown.length > 0 ? (
              <>
                <ResponsiveContainer width="100%" height={200}>
                  <PieChart>
                    <Pie
                      data={categoryBreakdown}
                      cx="50%"
                      cy="50%"
                      innerRadius={40}
                      outerRadius={80}
                      paddingAngle={5}
                      dataKey="total">
                      {categoryBreakdown.map((entry, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={entry.color || COLORS[index % COLORS.length]}
                        />
                      ))}
                    </Pie>
                    <Tooltip content={<PieTooltip />} />
                  </PieChart>
                </ResponsiveContainer>
                <div className="space-y-2 mt-4">
                  {categoryBreakdown.slice(0, 5).map((cat, i) => (
                    <div
                      key={i}
                      className="flex items-center justify-between text-sm">
                      <div className="flex items-center space-x-2">
                        <div
                          className="w-3 h-3 rounded-full"
                          style={{
                            backgroundColor:
                              cat.color || COLORS[i % COLORS.length],
                          }}
                        />
                        <span>{cat.category_name}</span>
                      </div>
                      <span className="font-medium">
                        R${' '}
                        {cat.total.toLocaleString('pt-BR', {
                          minimumFractionDigits: 0,
                        })}
                      </span>
                    </div>
                  ))}
                </div>
              </>
            ) : (
              <div className="h-[200px] flex items-center justify-center text-muted-foreground text-sm text-center">
                <div>
                  <PlusCircle className="h-8 w-8 mx-auto mb-2 opacity-30" />
                  Sem despesas este mês
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* AI Insights + Quick Actions */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Insights Financeiros</CardTitle>
            <CardDescription>Análise personalizada da IA</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <AIInsightCard />

            {balance >= 0 && savingsRate >= 10 && (
              <div className="flex items-start space-x-3 p-3 bg-green-50 dark:bg-green-950/30 rounded-lg">
                <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
                <div className="flex-1">
                  <p className="text-sm font-medium text-green-800 dark:text-green-200">
                    Boa taxa de poupança!
                  </p>
                  <p className="text-xs text-green-600 dark:text-green-300">
                    Você está poupando {savingsRate.toFixed(1)}% da sua renda.
                  </p>
                </div>
              </div>
            )}

            {totalExpenses > totalIncome && (
              <div className="flex items-start space-x-3 p-3 bg-red-50 dark:bg-red-950/30 rounded-lg">
                <AlertCircle className="h-5 w-5 text-red-500 mt-0.5" />
                <div className="flex-1">
                  <p className="text-sm font-medium text-red-800 dark:text-red-200">
                    Gastos acima da renda
                  </p>
                  <p className="text-xs text-red-600 dark:text-red-300">
                    Seus gastos superam sua renda em R${' '}
                    {(totalExpenses - totalIncome).toLocaleString('pt-BR', {
                      minimumFractionDigits: 2,
                    })}
                    .
                  </p>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Ações Rápidas</CardTitle>
            <CardDescription>Registre suas movimentações</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Button
              onClick={() => openTxDialog()}
              className="w-full justify-start bg-primary hover:bg-primary/90 text-primary-foreground">
              <PlusCircle className="h-4 w-4 mr-2" />
              Nova Transação
            </Button>
            <Button
              onClick={() => setCatDialogOpen(true)}
              variant="secondary"
              className="w-full justify-start">
              <FolderOpen className="h-4 w-4 mr-2 text-primary" />
              Gerenciar Categorias
            </Button>
            <Button
              variant="secondary"
              className="w-full justify-start"
              onClick={fetchDashboard}>
              <TrendingUp className="h-4 w-4 mr-2" />
              Atualizar Visão
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
