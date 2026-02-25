import { useState, useRef, useEffect, useTransition } from 'react';
import { Send, Bot, Lock, RefreshCcw, Sparkles, Activity } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { api } from '@/api/client';
import { toast } from 'sonner';
import { useQuery } from '@tanstack/react-query';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
}

export function AiCoach() {
  const [isPending, startTransition] = useTransition();
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      role: 'assistant',
      content: 'Olá! Sou seu Coach Financeiro de IA. Analisei suas transações recentes e notei que podemos otimizar seus gastos com Ifood. Como posso ajudar você hoje?'
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const scrollRef = useRef<HTMLDivElement>(null);

  const { data: coachStatus, isPending: isQueryPending, error: queryError } = useQuery({
    queryKey: ['coachStatus'],
    queryFn: () => api.get('/ai/coach').then(res => res.data),
    retry: false,
  });

  const isPremium = !(queryError && ((queryError as any).response?.status === 403 || (queryError as any).response?.status === 401));
  const remainingMsgRaw = coachStatus?.remaining ?? null;
  const maxMsg = coachStatus?.max ?? 50;

  // Local state for tracking decrements before next refetch (optimistic native update)
  const [localRemaining, setLocalRemaining] = useState<number | null>(null);
  
  // Sync localRemaining when query resolves
  useEffect(() => {
    if (remainingMsgRaw !== null) {
      setLocalRemaining(remainingMsgRaw);
    }
  }, [remainingMsgRaw]);

  const remainingMsg = localRemaining !== null ? localRemaining : remainingMsgRaw;
  const isInitializing = isQueryPending;

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages, isPending]);

  if (!isPremium) {
    return (
      <div className="max-w-4xl mx-auto text-center space-y-6 py-20 animate-in fade-in duration-700">
        <div className="inline-flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-tr from-amber-400 to-orange-500 mb-4 ring-8 ring-amber-500/10">
          <Lock className="h-10 w-10 text-white" />
        </div>
        <h1 className="text-4xl font-extrabold tracking-tight">Recurso Exclusivo Premium</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Faça upgrade para a experiência VIP e desbloqueie o Coach Financeiro conversacional, Missões Gamificadas e Automação Bancária total.
        </p>
        <Button size="lg" className="mt-4 px-8 tracking-wide bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 text-white">
          Assinar Premium
        </Button>
      </div>
    );
  }

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || isPending) return;
    if (remainingMsg !== null && remainingMsg <= 0) {
        toast.error('Limite Semanal Atingido', { description: 'Você esgotou os 50 tokens de IA do Coach esta semana. Tente novamente na próxima rotina!' });
        return;
    }

    const newUserMsg: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: inputValue
    };

    startTransition(() => {
      setMessages(prev => [...prev, newUserMsg]);
      setInputValue('');
    });

    try {
      const res = await api.post('/ai/coach', { message: newUserMsg.content });
      
      startTransition(() => {
        setLocalRemaining(res.data.remaining);
        const newBotMsg: Message = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: res.data.message || 'Excelente reflexão! Considerando que você ainda tem R$ 450 do seu orçamento de Lazer para este mês, você pode pedir aquele jantar sem comprometer as suas metas de longo prazo.'
        };
        setMessages(prev => [...prev, newBotMsg]);
      });
    } catch (error: any) {
      startTransition(() => {
        if (error.response?.status === 429) {
          toast.error('Limite de Tokens Excedido', { description: error.response.data.error });
          setLocalRemaining(0);
        } else {
          toast.error('Ocorreu um Erro', { description: 'Não foi possível conectar com o Coach no momento.' });
        }
        setMessages(prev => prev.filter(m => m.id !== newUserMsg.id)); // rollback
      });
    }
  };

  if (isInitializing) return null;

  return (
    <div className="h-[calc(100vh-8rem)] max-w-4xl mx-auto flex flex-col animate-in fade-in duration-500">
      
      {/* Header section */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 mb-4">
        <div>
          <Badge className="bg-gradient-to-r from-amber-500 to-orange-500 text-white border-none mb-3 px-3 py-1">
            <Sparkles className="w-3 h-3 mr-1.5" /> Coach IA Premium
          </Badge>
          <h1 className="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-slate-50">
            Seu Conselheiro Financeiro Pessoal
          </h1>
        </div>
        <div className="flex gap-2">
           {remainingMsg !== null && (
               <Badge variant="secondary" className="px-3 py-1 text-xs">
                 <Activity className="w-3 h-3 mr-2" />
                 {remainingMsg}/{maxMsg} tokens na semana
               </Badge>
           )}
          <Button variant="outline" size="sm" onClick={() => setMessages([{ id: '1', role: 'assistant', content: 'Olá! Como posso ajudar você hoje?' }])} className="text-muted-foreground">
            <RefreshCcw className="w-4 h-4 mr-2" /> Novo Papo
          </Button>
        </div>
      </div>

      {/* Chat Interface */}
      <Card className="flex-1 flex flex-col overflow-hidden border-border/50 shadow-sm">
        <CardHeader className="bg-slate-50 dark:bg-slate-900 border-b py-3 px-4">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
              <Bot className="h-5 w-5 text-primary" />
            </div>
            <div>
              <CardTitle className="text-sm">FinanceAI Coach</CardTitle>
              <CardDescription className="text-xs text-emerald-500 font-medium">Online agora</CardDescription>
            </div>
          </div>
        </CardHeader>
        
        <CardContent className="flex-1 p-0 flex flex-col min-h-0 bg-slate-50/30 dark:bg-slate-950/30 relative">
          
          {/* Messages Area */}
          <div ref={scrollRef} className="flex-1 overflow-y-auto p-4 space-y-6">
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={cn(
                  "flex items-end gap-2 max-w-[80%]",
                  msg.role === 'user' ? "ml-auto flex-row-reverse" : ""
                )}
              >
                {msg.role === 'assistant' && (
                  <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 mb-1">
                    <Bot className="h-4 w-4 text-primary" />
                  </div>
                )}
                <div
                  className={cn(
                    "px-4 py-3 rounded-2xl text-sm leading-relaxed shadow-sm",
                    msg.role === 'user'
                      ? "bg-primary text-primary-foreground rounded-br-sm"
                      : "bg-background border rounded-bl-sm"
                  )}
                >
                  {msg.content}
                </div>
              </div>
            ))}
            
            {isPending && (
              <div className="flex items-end gap-2 max-w-[80%]">
                <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 mb-1">
                  <Bot className="h-4 w-4 text-primary" />
                </div>
                <div className="px-4 py-4 rounded-2xl bg-background border rounded-bl-sm flex gap-1">
                  <span className="w-2 h-2 rounded-full bg-primary/40 animate-bounce" style={{ animationDelay: '0ms' }} />
                  <span className="w-2 h-2 rounded-full bg-primary/40 animate-bounce" style={{ animationDelay: '150ms' }} />
                  <span className="w-2 h-2 rounded-full bg-primary/40 animate-bounce" style={{ animationDelay: '300ms' }} />
                </div>
              </div>
            )}
          </div>

          {/* Input Area */}
          <div className="p-4 bg-background border-t">
            <form onSubmit={handleSend} className="relative flex items-center">
              <Input
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                placeholder="Pergunte algo como: 'Posso pedir pizza hoje?'"
                className="pr-12 py-6 rounded-full bg-slate-100 dark:bg-slate-900 border-none shadow-none focus-visible:ring-1"
                disabled={isPending}
              />
              <Button
                type="submit"
                size="icon"
                disabled={!inputValue.trim() || isPending}
                className="absolute right-1.5 rounded-full h-9 w-9"
              >
                <Send className="h-4 w-4 ml-0.5" />
              </Button>
            </form>
            <div className="mt-3 flex gap-2 overflow-x-auto pb-1 scrollbar-hide">
              <Badge variant="outline" className="cursor-pointer hover:bg-accent whitespace-nowrap" onClick={() => setInputValue('Dá para viajar no fim do ano?')}>
                Dá para viajar no fim do ano?
              </Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-accent whitespace-nowrap" onClick={() => setInputValue('Como corto gastos sem sofrer?')}>
                Como corto gastos sem sofrer?
              </Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-accent whitespace-nowrap" onClick={() => setInputValue('Estou gastando muito com ifood?')}>
                Estou gastando muito com ifood?
              </Badge>
            </div>
          </div>
          
        </CardContent>
      </Card>
      
    </div>
  );
}
