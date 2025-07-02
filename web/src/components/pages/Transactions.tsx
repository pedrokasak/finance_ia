import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger 
} from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { 
  PlusCircle, 
  Search, 
  Filter,
  TrendingUp,
  TrendingDown,
  Calendar,
  Download
} from 'lucide-react';

const transactions = [
  {
    id: 1,
    type: 'receita',
    description: 'Salário',
    category: 'Trabalho',
    amount: 5500,
    date: '2024-01-15',
    status: 'confirmado'
  },
  {
    id: 2,
    type: 'despesa',
    description: 'Supermercado',
    category: 'Alimentação',
    amount: 350,
    date: '2024-01-14',
    status: 'confirmado'
  },
  {
    id: 3,
    type: 'despesa',
    description: 'Combustível',
    category: 'Transporte',
    amount: 180,
    date: '2024-01-13',
    status: 'pendente'
  },
  {
    id: 4,
    type: 'receita',
    description: 'Freelance',
    category: 'Trabalho',
    amount: 800,
    date: '2024-01-12',
    status: 'confirmado'
  },
  {
    id: 5,
    type: 'despesa',
    description: 'Academia',
    category: 'Saúde',
    amount: 120,
    date: '2024-01-11',
    status: 'confirmado'
  },
];

const categories = {
  receita: ['Trabalho', 'Investimentos', 'Freelance', 'Outros'],
  despesa: ['Alimentação', 'Transporte', 'Saúde', 'Lazer', 'Casa', 'Educação', 'Outros']
};

export function Transactions() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [transactionType, setTransactionType] = useState<'receita' | 'despesa'>('receita');
  const [searchTerm, setSearchTerm] = useState('');

  const filteredTransactions = transactions.filter(transaction =>
    transaction.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
    transaction.category.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Aqui você adicionaria a lógica para salvar a transação
    setDialogOpen(false);
  };

  return (
    <div className="space-y-6">
      {/* Header com estatísticas rápidas */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total este Mês</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">R$ 6.300</div>
            <p className="text-xs text-muted-foreground">
              Receitas do mês atual
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Gastos este Mês</CardTitle>
            <TrendingDown className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">R$ 650</div>
            <p className="text-xs text-muted-foreground">
              Despesas do mês atual
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Transações</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{transactions.length}</div>
            <p className="text-xs text-muted-foreground">
              Total de registros
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Transações</CardTitle>
              <CardDescription>
                Gerencie todas suas receitas e despesas
              </CardDescription>
            </div>
            
            <div className="flex space-x-2">
              <Button variant="outline" size="sm">
                <Download className="h-4 w-4 mr-2" />
                Exportar
              </Button>
              
              <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                <DialogTrigger asChild>
                  <Button>
                    <PlusCircle className="h-4 w-4 mr-2" />
                    Nova Transação
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[425px]">
                  <DialogHeader>
                    <DialogTitle>Nova Transação</DialogTitle>
                    <DialogDescription>
                      Adicione uma nova receita ou despesa ao seu controle financeiro.
                    </DialogDescription>
                  </DialogHeader>
                  
                  <form onSubmit={handleSubmit} className="space-y-4">
                    <Tabs value={transactionType} onValueChange={(value) => setTransactionType(value as 'receita' | 'despesa')}>
                      <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="receita" className="text-green-600">
                          <TrendingUp className="h-4 w-4 mr-2" />
                          Receita
                        </TabsTrigger>
                        <TabsTrigger value="despesa" className="text-red-600">
                          <TrendingDown className="h-4 w-4 mr-2" />
                          Despesa
                        </TabsTrigger>
                      </TabsList>
                    </Tabs>

                    <div className="space-y-2">
                      <Label htmlFor="description">Descrição</Label>
                      <Input 
                        id="description" 
                        placeholder="Ex: Salário, Supermercado..." 
                        required 
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="amount">Valor</Label>
                      <Input 
                        id="amount" 
                        type="number" 
                        placeholder="0,00" 
                        step="0.01"
                        required 
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="category">Categoria</Label>
                      <Select required>
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione uma categoria" />
                        </SelectTrigger>
                        <SelectContent>
                          {categories[transactionType].map((category) => (
                            <SelectItem key={category} value={category}>
                              {category}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="date">Data</Label>
                      <Input 
                        id="date" 
                        type="date" 
                        required 
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="notes">Observações (opcional)</Label>
                      <Textarea 
                        id="notes" 
                        placeholder="Adicione detalhes sobre esta transação..."
                        rows={3}
                      />
                    </div>

                    <div className="flex justify-end space-x-2 pt-4">
                      <Button type="button" variant="outline" onClick={() => setDialogOpen(false)}>
                        Cancelar
                      </Button>
                      <Button type="submit" className={transactionType === 'receita' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'}>
                        Salvar
                      </Button>
                    </div>
                  </form>
                </DialogContent>
              </Dialog>
            </div>
          </div>
        </CardHeader>
        
        <CardContent>
          {/* Filtros */}
          <div className="flex items-center space-x-2 mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Buscar transações..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-8"
              />
            </div>
            <Button variant="outline" size="sm">
              <Filter className="h-4 w-4 mr-2" />
              Filtros
            </Button>
          </div>

          {/* Tabela de Transações */}
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Descrição</TableHead>
                <TableHead>Categoria</TableHead>
                <TableHead>Data</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Valor</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredTransactions.map((transaction) => (
                <TableRow key={transaction.id}>
                  <TableCell className="font-medium">
                    {transaction.description}
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">
                      {transaction.category}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {new Date(transaction.date).toLocaleDateString('pt-BR')}
                  </TableCell>
                  <TableCell>
                    <Badge 
                      variant={transaction.status === 'confirmado' ? 'default' : 'outline'}
                      className={transaction.status === 'confirmado' ? 'bg-green-100 text-green-800' : ''}
                    >
                      {transaction.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <span className={transaction.type === 'receita' ? 'text-green-600 font-medium' : 'text-red-600 font-medium'}>
                      {transaction.type === 'receita' ? '+' : '-'} R$ {transaction.amount.toLocaleString('pt-BR')}
                    </span>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}