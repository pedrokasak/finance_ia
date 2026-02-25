import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useUser } from '@/contexts/UserContext';
import { api } from '@/api/client';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Progress } from '@/components/ui/progress';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Lock, Target, Plus, Plane, Car, Home, TrendingUp, Trash } from 'lucide-react';
import { toast } from 'sonner';

interface Goal {
  id: string;
  name: string;
  target_amount: number;
  current_amount: number;
  target_date: string;
  icon: string;
}

const ICONS: Record<string, React.ReactNode> = {
  flag: <Target className="h-5 w-5 text-blue-500" />,
  plane: <Plane className="h-5 w-5 text-sky-500" />,
  car: <Car className="h-5 w-5 text-orange-500" />,
  home: <Home className="h-5 w-5 text-green-500" />,
  trending: <TrendingUp className="h-5 w-5 text-purple-500" />
};

export function Goals() {
  const { profile } = useUser();
  const queryClient = useQueryClient();
  const [isAddOpen, setIsAddOpen] = useState(false);

  const isFreePlan = profile?.plan === 'free';

  const { data: goals, isLoading } = useQuery<Goal[]>({
    queryKey: ['goals'],
    queryFn: () => api.get('/goals/').then((res) => res.data),
    enabled: !isFreePlan, // Only fetch if not free
  });

  const createGoal = useMutation({
    mutationFn: (newGoal: Partial<Goal>) => api.post('/goals/', newGoal),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      setIsAddOpen(false);
      toast.success('Meta criada com sucesso!');
    },
    onError: () => toast.error('Erro ao criar meta'),
  });

  const deleteGoal = useMutation({
    mutationFn: (id: string) => api.delete(`/goals/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      toast.success('Meta removida!');
    },
    onError: () => toast.error('Erro ao remover meta'),
  });

  const handleCreate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const targetAmountStr = formData.get('target_amount') as string;
    
    createGoal.mutate({
      name: formData.get('name') as string,
      target_amount: parseFloat(targetAmountStr.replace(',', '.')),
      target_date: formData.get('target_date') as string,
      icon: formData.get('icon') as string || 'flag',
    });
  };

  const renderContent = () => {
    if (isLoading && !isFreePlan) return <div className="p-8 text-center text-muted-foreground">Carregando metas...</div>;

    const list = goals || [];

    if (list.length === 0 && !isFreePlan) {
      return (
        <Card className="border-dashed shadow-none bg-transparent">
          <CardContent className="flex flex-col items-center justify-center p-12 text-center">
            <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center mb-4">
              <Target className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-xl font-semibold mb-2">Nenhuma meta definida</h3>
            <p className="text-muted-foreground max-w-sm mb-6">
              Comece a planejar seu futuro criando metas financeiras. Pode ser a viagem dos sonhos, um carro novo ou sua reserva de emergência.
            </p>
            <Button onClick={() => setIsAddOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Criar Primeira Meta
            </Button>
          </CardContent>
        </Card>
      );
    }

    return (
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {list.map((goal) => {
          const progress = Math.min(100, Math.max(0, (goal.current_amount / goal.target_amount) * 100));
          const dateObj = new Date(goal.target_date);
          const formattedDate = dateObj.toLocaleDateString('pt-BR', { month: 'short', year: 'numeric' });
          const IconComponent = ICONS[goal.icon] || ICONS['flag'];

          return (
            <Card key={goal.id} className="relative overflow-hidden group">
              <CardHeader className="pb-4">
                <div className="flex items-start justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 bg-muted rounded-lg">
                      {IconComponent}
                    </div>
                    <div>
                      <CardTitle className="text-lg">{goal.name}</CardTitle>
                      <CardDescription>Para {formattedDate}</CardDescription>
                    </div>
                  </div>
                  <Button 
                    variant="ghost" 
                    size="icon" 
                    className="opacity-0 group-hover:opacity-100 transition-opacity text-destructive"
                    onClick={() => deleteGoal.mutate(goal.id)}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="font-medium text-foreground">
                      {new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(goal.current_amount)}
                    </span>
                    <span className="text-muted-foreground">
                      / {new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(goal.target_amount)}
                    </span>
                  </div>
                  <Progress value={progress} className="h-2" />
                  <p className="text-xs text-right text-muted-foreground font-medium pt-1">
                    {progress.toFixed(0)}% concluído
                  </p>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    );
  };

  return (
    <div className="relative min-h-[500px] space-y-6">
      
      {/* HEADER */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Metas de Vida</h2>
          <p className="text-muted-foreground mt-1">Defina objetivos e acompanhe seu progresso.</p>
        </div>
        <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
          <DialogTrigger asChild>
            <Button disabled={isFreePlan}>
              <Plus className="mr-2 h-4 w-4" /> Nova Meta
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Criar Nova Meta</DialogTitle>
              <DialogDescription>
                Dê um nome, defina o valor total que precisa e até quando deseja alcançar.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleCreate} className="space-y-4 pt-4">
              <div className="space-y-2">
                <Label htmlFor="name">Nome da Meta</Label>
                <Input id="name" name="name" placeholder="Ex: Viagem para o Japão" required />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="target_amount">Valor Alvo (R$)</Label>
                  <Input id="target_amount" name="target_amount" type="number" step="0.01" min="1" placeholder="Ex: 5000" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="target_date">Data Alvo</Label>
                  <Input id="target_date" name="target_date" type="date" required />
                </div>
              </div>
              <div className="space-y-2">
                <Label>Ícone</Label>
                <div className="flex gap-2">
                  {Object.keys(ICONS).map(k => (
                    <label key={k} className="flex-1 cursor-pointer">
                      <input type="radio" name="icon" value={k} className="peer sr-only" defaultChecked={k === 'flag'} />
                      <div className="h-10 border rounded-md flex justify-center items-center peer-checked:bg-primary/10 peer-checked:border-primary transition-all">
                        {ICONS[k]}
                      </div>
                    </label>
                  ))}
                </div>
              </div>
              <Button type="submit" className="w-full mt-4" disabled={createGoal.isPending}>
                {createGoal.isPending ? 'Salvando...' : 'Criar Meta'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* PAYWALL BLUR */}
      {isFreePlan && (
        <div className="absolute inset-0 z-50 flex flex-col items-center justify-center bg-background/60 backdrop-blur-md rounded-xl border border-border/50">
          <div className="bg-card p-8 rounded-2xl shadow-xl max-w-sm text-center border-2 border-primary/20 scale-100 transition-transform">
            <div className="w-16 h-16 bg-blue-100 dark:bg-blue-900/40 rounded-full flex items-center justify-center mx-auto mb-4">
              <Lock className="h-8 w-8 text-blue-600 dark:text-blue-400" />
            </div>
            <h3 className="text-2xl font-bold mb-2">Exclusivo Pro/Premium</h3>
            <p className="text-muted-foreground mb-6">
              O módulo de Metas Inteligentes está disponível apenas para assinantes. Planeje viagens, carros novos e muito mais!
            </p>
            <Button size="lg" className="w-full text-md font-semibold" onClick={() => (window.location.search = '?page=subscription')}>
              Fazer Upgrade Agora
            </Button>
          </div>
        </div>
      )}

      {/* CONTENT DEMO PARA USUÁRIOS FREE, OU CONTEÚDO REAL */}
      <div className={isFreePlan ? 'opacity-30 pointer-events-none select-none blur-sm' : ''}>
        {isFreePlan ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {[
              { name: 'Reserva de Emergência', amount: 5000, target: 10000, icon: 'home' },
              { name: 'Viagem dos Sonhos', amount: 2000, target: 15000, icon: 'plane' }
            ].map((mock, i) => (
               <Card key={i} className="relative overflow-hidden">
                 <CardHeader className="pb-4">
                   <div className="flex items-center space-x-3">
                     <div className="p-2 bg-muted rounded-lg">{ICONS[mock.icon]}</div>
                     <div>
                       <CardTitle className="text-lg">{mock.name}</CardTitle>
                       <CardDescription>Mock para Dez 2026</CardDescription>
                     </div>
                   </div>
                 </CardHeader>
                 <CardContent>
                   <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span className="font-medium text-foreground">R$ {mock.amount}</span>
                        <span className="text-muted-foreground">/ R$ {mock.target}</span>
                      </div>
                      <Progress value={50} className="h-2" />
                      <p className="text-xs text-right text-muted-foreground font-medium pt-1">50% concluído</p>
                   </div>
                 </CardContent>
               </Card>
            ))}
          </div>
        ) : (
          renderContent()
        )}
      </div>

    </div>
  );
}
