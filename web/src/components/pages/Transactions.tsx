import React, { useState, useTransition } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { PlusCircle, Search, TrendingUp, TrendingDown, Calendar, Loader2, Trash2, Pencil } from 'lucide-react';
import { toast } from 'sonner';
import financeService, { Transaction } from '@/services/financeService';

type TxType = 'income' | 'expense';

interface FormState {
  type: TxType;
  description: string;
  amount: string;
  category_id: string;
  date: string;
}

const emptyForm = (): FormState => ({
  type: 'income',
  description: '',
  amount: '',
  category_id: '',
  date: new Date().toISOString().split('T')[0],
});

export function Transactions() {
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm());
  const [isPending, startTransition] = useTransition();
  const [searchTerm, setSearchTerm] = useState('');

  const { data: pageData, isPending: isQueryPending } = useQuery({
    queryKey: ['transactionsPage'],
    queryFn: async () => {
      const [txRes, catRes, dashRes] = await Promise.all([
        financeService.listTransactions({ limit: 50 }),
        financeService.getCategories(),
        financeService.getDashboard(),
      ]);
      return {
        transactions: txRes.data || [],
        categories: catRes || [],
        totalIncome: dashRes.total_income || 0,
        totalExpenses: dashRes.total_expenses || 0,
      };
    }
  });

  const transactions = pageData?.transactions || [];
  const categories = pageData?.categories || [];
  const totalIncome = pageData?.totalIncome || 0;
  const totalExpenses = pageData?.totalExpenses || 0;
  const loading = isPending || isQueryPending;

  const filteredCategories = categories.filter(
    (c) => c.type === (form.type === 'income' ? 'income' : 'expense')
  );

  const filteredTransactions = transactions.filter((tx) =>
    tx.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleOpenDialog = (type: TxType) => {
    setEditingId(null);
    setForm({ ...emptyForm(), type });
    setDialogOpen(true);
  };

  const handleEdit = (tx: Transaction) => {
    setEditingId(tx.id || null);
    setForm({
      type: tx.type,
      description: tx.description,
      amount: String(tx.amount),
      category_id: tx.category_id || '',
      date: tx.date.split('T')[0],
    });
    setDialogOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.amount || parseFloat(form.amount) <= 0) {
      toast.error('Valor inválido');
      return;
    }
    
    try {
      const payload = {
        type: form.type,
        amount: parseFloat(form.amount),
        description: form.description,
        category_id: form.category_id || undefined,
        date: form.date,
      };

      if (editingId) {
        await financeService.updateTransaction(editingId, payload);
      } else {
        await financeService.createTransaction(payload);
      }

      startTransition(() => {
        toast.success(`${form.type === 'income' ? 'Receita' : 'Despesa'} ${editingId ? 'atualizada' : 'adicionada'}!`);
        setDialogOpen(false);
        setEditingId(null);
        setForm(emptyForm());
        queryClient.invalidateQueries({ queryKey: ['transactionsPage'] });
        queryClient.invalidateQueries({ queryKey: ['dashboardSummary'] });
      });
    } catch {
      startTransition(() => {
        toast.error(editingId ? 'Erro ao atualizar transação' : 'Erro ao salvar transação');
      });
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await financeService.deleteTransaction(id);
      startTransition(() => {
        toast.success('Transação removida');
        queryClient.invalidateQueries({ queryKey: ['transactionsPage'] });
        queryClient.invalidateQueries({ queryKey: ['dashboardSummary'] });
      });
    } catch {
      startTransition(() => {
        toast.error('Erro ao remover transação');
      });
    }
  };

  return (
    <div className="space-y-6">
      {/* Summary cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Receitas do Mês</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {loading ? <Loader2 className="h-5 w-5 animate-spin" /> : `R$ ${totalIncome.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Despesas do Mês</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {loading ? <Loader2 className="h-5 w-5 animate-spin" /> : `R$ ${totalExpenses.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Total de Transações</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{transactions.length}</div>
          </CardContent>
        </Card>
      </div>

      {/* Quick action buttons */}
      <div className="flex gap-3">
        <Button onClick={() => handleOpenDialog('expense')} className="bg-primary hover:bg-primary/90 text-primary-foreground">
          <PlusCircle className="h-4 w-4 mr-2" /> Nova Transação
        </Button>
      </div>

      {/* Transaction list */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Transações</CardTitle>
              <CardDescription>Gerencie todas suas receitas e despesas</CardDescription>
            </div>
            <Button onClick={() => handleOpenDialog('income')}>
              <PlusCircle className="h-4 w-4 mr-2" /> Nova Transação
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="relative mb-4">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Buscar transações..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-9"
            />
          </div>

          {loading ? (
            <div className="flex justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : filteredTransactions.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              <TrendingUp className="h-12 w-12 mx-auto mb-3 opacity-30" />
              <p>Nenhuma transação encontrada.</p>
              <Button className="mt-4" onClick={() => handleOpenDialog('expense')}>
                <PlusCircle className="h-4 w-4 mr-2" /> Adicionar primeira transação
              </Button>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Descrição</TableHead>
                  <TableHead>Categoria</TableHead>
                  <TableHead>Data</TableHead>
                  <TableHead className="text-right">Valor</TableHead>
                  <TableHead className="w-20 text-right">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredTransactions.map((tx) => (
                  <TableRow key={tx.id}>
                    <TableCell className="font-medium">{tx.description || '—'}</TableCell>
                    <TableCell>
                      {tx.category ? (
                        <Badge variant="secondary" style={{ borderColor: tx.category.color }}>
                          {tx.category.name}
                        </Badge>
                      ) : '—'}
                    </TableCell>
                    <TableCell>{tx.date ? new Date(tx.date).toLocaleDateString('pt-BR') : '—'}</TableCell>
                    <TableCell className="text-right">
                      <span className={tx.type === 'income' ? 'text-green-600 font-medium' : 'text-red-600 font-medium'}>
                        {tx.type === 'income' ? '+' : '-'} R$ {tx.amount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-1">
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-primary" onClick={() => handleEdit(tx)}>
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-red-500" onClick={() => tx.id && handleDelete(tx.id)}>
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* New/Edit Transaction Dialog */}
      <Dialog open={dialogOpen} onOpenChange={(open) => { setDialogOpen(open); if (!open) setEditingId(null); }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{editingId ? 'Editar' : 'Nova'} {form.type === 'income' ? 'Receita' : 'Despesa'}</DialogTitle>
            <DialogDescription>{editingId ? 'Modifique os detalhes desta movimentação.' : 'Adicione uma nova movimentação ao seu controle financeiro.'}</DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Tabs value={form.type} onValueChange={(v) => setForm({ ...form, type: v as TxType, category_id: '' })}>
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
              <Input placeholder="Ex: Salário, Supermercado..." value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })} required />
            </div>

            <div className="space-y-2">
              <Label>Valor (R$)</Label>
              <Input type="number" placeholder="0,00" step="0.01" min="0.01" value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })} required />
            </div>

            {filteredCategories.length > 0 && (
              <div className="space-y-2">
                <Label>Categoria</Label>
                <Select value={form.category_id} onValueChange={(v) => setForm({ ...form, category_id: v })}>
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione uma categoria" />
                  </SelectTrigger>
                  <SelectContent>
                    {filteredCategories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>{cat.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            <div className="space-y-2">
              <Label>Data</Label>
              <Input type="date" value={form.date}
                onChange={(e) => setForm({ ...form, date: e.target.value })} required />
            </div>

            <div className="flex gap-2 pt-2">
              <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? 'Salvando...' : 'Salvar Transação'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}