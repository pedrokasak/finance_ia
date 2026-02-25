import { useState } from 'react';
import { Flame, Trophy, Target, Star, Lock, Zap, CheckCircle2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

export function AiMissions() {
  const [isPremium] = useState(true);

  if (!isPremium) {
    return (
      <div className="max-w-4xl mx-auto text-center space-y-6 py-20 animate-in fade-in duration-700">
        <div className="inline-flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 mb-4 ring-8 ring-indigo-500/10">
          <Lock className="h-10 w-10 text-white" />
        </div>
        <h1 className="text-4xl font-extrabold tracking-tight">Recurso Exclusivo Premium</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Faça upgrade para acessar as Missões Inteligentes gamificadas, ganhar recompensas virtuais e proteger seu Score.
        </p>
        <Button size="lg" className="mt-4 px-8 tracking-wide bg-gradient-to-r from-indigo-500 to-purple-500 hover:from-indigo-600 hover:to-purple-600 text-white border-0">
          Assinar Premium
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-8 max-w-6xl mx-auto animate-in fade-in slide-in-from-bottom-6 duration-700">
      
      {/* Header section */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
        <div>
          <Badge className="bg-gradient-to-r from-indigo-500 to-purple-500 text-white border-none mb-3 px-3 py-1">
            <Zap className="w-3 h-3 mr-1.5" /> Missões IA Premium
          </Badge>
          <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-slate-50">
            Missões & Conquistas
          </h1>
          <p className="text-muted-foreground mt-2 text-lg">
            Transformando disciplina financeira no seu esporte favorito.
          </p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-4">
        
        {/* Streak & Level Profile (Left Column) */}
        <div className="md:col-span-1 space-y-6">
          <Card className="border-border/50 shadow-sm bg-gradient-to-b from-slate-50 to-white dark:from-slate-900 dark:to-slate-950 overflow-hidden relative">
            <div className="absolute top-0 right-0 w-32 h-32 bg-orange-500/10 rounded-full blur-3xl -mr-10 -mt-10" />
            <CardContent className="pt-6 relative z-10 flex flex-col items-center text-center">
              <div className="w-20 h-20 rounded-full bg-gradient-to-tr from-orange-400 to-amber-500 p-1 mb-4 shadow-lg shadow-orange-500/20">
                <div className="w-full h-full rounded-full bg-slate-900 flex items-center justify-center border-4 border-slate-900">
                  <Flame className="w-8 h-8 text-orange-500" strokeWidth={2.5} />
                </div>
              </div>
              <h2 className="text-3xl font-black text-slate-900 dark:text-white mb-1">12 Dias</h2>
              <p className="text-sm font-medium text-orange-500 mb-6 uppercase tracking-wider">Ofensiva de Economia</p>
              
              <div className="w-full space-y-2 mb-4">
                <div className="flex justify-between text-xs font-medium text-muted-foreground">
                  <span>Nível 4</span>
                  <span>1450 / 2000 XP</span>
                </div>
                <Progress value={72} className="h-2.5 bg-slate-200 dark:bg-slate-800 [&>div]:bg-indigo-500" />
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50 shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <Trophy className="w-4 h-4 text-yellow-500" /> Sala de Troféus
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-yellow-500/20 flex flex-shrink-0 items-center justify-center">
                  <span className="text-lg">🍕</span>
                </div>
                <div>
                  <p className="text-sm font-medium">Mestre da Cozinha</p>
                  <p className="text-xs text-muted-foreground">7 dias sem delivery</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-slate-200 dark:bg-slate-800 flex flex-shrink-0 items-center justify-center">
                  <Lock className="w-4 h-4 text-slate-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Investidor Anjo</p>
                  <p className="text-xs text-slate-500">Invista R$1000 este mês</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Daily & Weekly Quests (Right Column) */}
        <div className="md:col-span-3 space-y-6">
          <Card className="border-border/50 shadow-sm">
            <CardHeader className="border-b bg-slate-50/50 dark:bg-slate-900/50 pb-4">
              <CardTitle className="text-xl flex items-center gap-2">
                <Target className="w-5 h-5 text-indigo-500" />
                Missões da Semana
              </CardTitle>
              <CardDescription>
                Completar missões aumenta seu score financeiro e libera avatares exclusivos.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <div className="divide-y divide-border">
                
                {/* Mission 1 (Completed) */}
                <div className="p-6 flex items-start gap-4 hover:bg-slate-50/50 dark:hover:bg-slate-900/30 transition-colors">
                  <div className="mt-1">
                    <CheckCircle2 className="w-6 h-6 text-emerald-500" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <h4 className="font-semibold text-slate-900 dark:text-slate-100 line-through opacity-70">Revisar Assinaturas</h4>
                      <Badge variant="outline" className="text-emerald-500 border-emerald-500/30 bg-emerald-500/10">
                        +50 XP
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3 line-through opacity-70">
                      Você acessou a aba de Diagnóstico e encerrou uma assinatura não utilizada.
                    </p>
                  </div>
                </div>

                {/* Mission 2 (Active) */}
                <div className="p-6 flex items-start gap-4 bg-indigo-500/5 hover:bg-indigo-500/10 transition-colors">
                  <div className="mt-1 flex-shrink-0 w-6 h-6 rounded-full border-2 border-indigo-500 flex items-center justify-center">
                    <div className="w-2 h-2 rounded-full bg-indigo-500" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <h4 className="font-semibold text-indigo-700 dark:text-indigo-400">Desafio da Lancheira</h4>
                      <Badge className="bg-gradient-to-r from-indigo-500 to-purple-500 border-none">
                        +150 XP
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground mb-4">
                      Gaste menos de R$ 50 em restaurantes e cafeterias até Sexta-Feira.
                    </p>
                    <div className="flex items-center gap-4">
                      <div className="flex-1 h-2 rounded-full bg-slate-200 dark:bg-slate-800">
                        <div className="h-full w-[40%] rounded-full bg-indigo-500" />
                      </div>
                      <span className="text-xs font-medium text-slate-500">R$ 20 / R$ 50</span>
                    </div>
                  </div>
                </div>

                {/* Mission 3 (Active) */}
                <div className="p-6 flex items-start gap-4 hover:bg-slate-50/50 dark:hover:bg-slate-900/30 transition-colors">
                  <div className="mt-1 flex-shrink-0 w-6 h-6 rounded-full border-2 border-slate-300 dark:border-slate-700" />
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <h4 className="font-semibold text-slate-900 dark:text-slate-100">Visão de Longo Prazo</h4>
                      <Badge variant="outline" className="text-amber-500 border-amber-500/30">
                        <Star className="w-3 h-3 mr-1 fill-amber-500" /> Raro
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground mb-4">
                      Simule 3 cenários diferentes de aposentadoria utilizando o Simulador Interativo.
                    </p>
                    <Button variant="secondary" size="sm" className="text-xs">
                      Ir para o Simulador
                    </Button>
                  </div>
                </div>

              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
