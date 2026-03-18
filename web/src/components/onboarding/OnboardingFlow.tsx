import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  ChevronRight,
  ChevronLeft,
  Check,
  Info,
  Wallet,
  PieChart,
  Target,
  TrendingUp,
  Landmark,
  Layers,
  Coins,
  Percent,
  Flame,
  Loader2,
  Sun,
  Moon,
} from "lucide-react";
import { api } from "@/api/client";
import financeService, { FinancialMethod } from "@/services/financeService";
import { useTheme } from "next-themes";

interface OnboardingFlowProps {
  onComplete: () => void;
}

// Icon mapping from DB string to Lucide component
const iconMap: Record<string, React.ComponentType<{ className?: string }>> = {
  PieChart,
  Layers,
  Coins,
  Landmark,
  Wallet,
  Target,
  TrendingUp,
  Percent,
  Flame,
};

export function OnboardingFlow({ onComplete }: OnboardingFlowProps) {
  const { theme, setTheme } = useTheme();
  const [step, setStep] = useState(1);
  const [methods, setMethods] = useState<FinancialMethod[]>([]);
  const [methodsLoading, setMethodsLoading] = useState(true);
  const [selectedMethod, setSelectedMethod] = useState<FinancialMethod | null>(
    null,
  );
  const [detailMethod, setDetailMethod] = useState<FinancialMethod | null>(
    null,
  );
  const [monthlyIncome, setMonthlyIncome] = useState("");
  const [loading, setLoading] = useState(false);

  const totalSteps = 4;

  useEffect(() => {
    financeService
      .getMethods()
      .then((data) => setMethods(data || []))
      .catch(console.error)
      .finally(() => setMethodsLoading(false));
  }, []);

  const handleComplete = async () => {
    if (!selectedMethod || !monthlyIncome) return;
    setLoading(true);
    try {
      const income = parseFloat(monthlyIncome.replace(",", "."));
      const split = selectedMethod.split;
      const needs = split.find((s) =>
        [
          "Necessidades",
          "Fixos",
          "Gastos",
          "Despesas Fixas",
          "Moradia",
          "Viva com o resto",
        ].includes(s.label),
      );

      const wants = split.find((s) =>
        ["Desejos", "Variáveis", "Alimentação", "Lazer e Metas"].includes(
          s.label,
        ),
      );
      const savings = split.find((s) =>
        [
          "Investimentos",
          "Reservas",
          "Poupança",
          "Investimento Imediato",
          "Investimentos (30%+)",
        ].includes(s.label),
      );

      await api.post("/finance/budget", {
        total_income: income,
        needs_percent: needs?.percent ?? 50,
        wants_percent: wants?.percent ?? 30,
        savings_percent: savings?.percent ?? 20,
      });

      onComplete();
    } catch {
      onComplete(); // Don't block the user if API fails
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      {/* Detail Modal */}
      <Dialog open={!!detailMethod} onOpenChange={() => setDetailMethod(null)}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              {detailMethod && (
                <>
                  {iconMap[detailMethod.icon]
                    ? (() => {
                        const Icon = iconMap[detailMethod.icon];
                        return (
                          <Icon className={`h-5 w-5 ${detailMethod.color}`} />
                        );
                      })()
                    : null}
                  {detailMethod.name}
                </>
              )}
            </DialogTitle>
          </DialogHeader>
          {detailMethod && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground leading-relaxed">
                {detailMethod.description}
              </p>
              <div className="space-y-2">
                <p className="text-sm font-medium">Distribuição:</p>
                {detailMethod.split.map((s) => (
                  <div key={s.label} className="flex items-center gap-2">
                    <div className={`w-3 h-3 rounded-full ${s.color}`} />
                    <span className="text-sm flex-1">{s.label}</span>
                    <span className="text-sm font-bold">{s.percent}%</span>
                  </div>
                ))}
              </div>
              <div className="bg-muted/50 rounded-lg p-3">
                <p className="text-xs text-muted-foreground">
                  {detailMethod.for_who}
                </p>
              </div>
              <Button
                className="w-full"
                onClick={() => {
                  setSelectedMethod(detailMethod);
                  setDetailMethod(null);
                }}
              >
                Selecionar este método
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <div className="w-full max-w-2xl">
        {/* Progress */}
        <div className="flex items-center gap-2 mb-8">
          {Array.from({ length: totalSteps }).map((_, i) => (
            <div key={i} className="flex-1">
              <div
                className={`h-1.5 rounded-full transition-all duration-500 ${i + 1 <= step ? "bg-primary" : "bg-muted"}`}
              />
            </div>
          ))}
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {step}/{totalSteps}
          </span>
        </div>

        {/* Step 1: Welcome */}
        {step === 1 && (
          <div className="text-center space-y-6 animate-in fade-in">
            <div className="space-y-2">
              <div className="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <TrendingUp className="h-8 w-8 text-primary" />
              </div>
              <h1 className="text-3xl font-bold tracking-tight">
                Bem-vindo ao FinanceIA
              </h1>
              <p className="text-muted-foreground text-lg max-w-md mx-auto">
                Vamos configurar sua jornada financeira em 3 passos simples.
              </p>
            </div>
            <div className="grid gap-3 text-left max-w-md mx-auto">
              {[
                "Escolha seu tema preferido",
                "Escolha seu método de planejamento",
                "Defina sua renda mensal",
                "Comece a registrar seus gastos",
              ].map((item, i) => (
                <div
                  key={i}
                  className="flex items-center gap-3 p-3 bg-muted/30 rounded-lg"
                >
                  <div className="w-6 h-6 bg-primary rounded-full flex items-center justify-center flex-shrink-0">
                    <span className="text-xs text-primary-foreground font-bold">
                      {i + 1}
                    </span>
                  </div>
                  <span className="text-sm">{item}</span>
                </div>
              ))}
            </div>
            <Button size="lg" onClick={() => setStep(2)} className="px-8">
              Começar <ChevronRight className="ml-2 h-4 w-4" />
            </Button>
          </div>
        )}

        {/* Step 2: Theme Selection */}
        {step === 2 && (
          <div className="space-y-6 animate-in fade-in text-center">
            <div className="space-y-2">
              <h2 className="text-2xl font-bold">
                Como você prefere usar o FinanceIA?
              </h2>
              <p className="text-muted-foreground">
                Escolha o tema que mais agrada aos seus olhos.
              </p>
            </div>
            <div className="grid grid-cols-2 gap-4 max-w-md mx-auto">
              <Card
                className={`cursor-pointer border-2 transition-all hover:border-primary ${theme === "light" ? "border-primary ring-2 ring-primary/20" : "border-transparent bg-muted/30"}`}
                onClick={() => setTheme("light")}
              >
                <CardContent className="p-6 flex flex-col items-center gap-4">
                  <div className="w-16 h-16 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                    <Sun className="h-8 w-8 text-amber-500" />
                  </div>
                  <span className="font-medium text-lg">Claro (Suave)</span>
                </CardContent>
              </Card>
              <Card
                className={`cursor-pointer border-2 transition-all hover:border-primary ${theme === "dark" ? "border-primary ring-2 ring-primary/20" : "border-transparent bg-muted/30"}`}
                onClick={() => setTheme("dark")}
              >
                <CardContent className="p-6 flex flex-col items-center gap-4">
                  <div className="w-16 h-16 rounded-full bg-slate-800 dark:bg-slate-900 flex items-center justify-center">
                    <Moon className="h-8 w-8 text-blue-400" />
                  </div>
                  <span className="font-medium text-lg">Escuro</span>
                </CardContent>
              </Card>
            </div>
            <div className="flex gap-3 max-w-md mx-auto mt-8">
              <Button
                variant="outline"
                onClick={() => setStep(1)}
                className="w-full"
              >
                <ChevronLeft className="mr-2 h-4 w-4" /> Voltar
              </Button>
              <Button className="w-full" onClick={() => setStep(3)}>
                Continuar <ChevronRight className="ml-2 h-4 w-4" />
              </Button>
            </div>
          </div>
        )}

        {/* Step 3: Method Selection */}
        {step === 3 && (
          <div className="space-y-6 animate-in fade-in">
            <div>
              <h2 className="text-2xl font-bold">
                Escolha seu método financeiro
              </h2>
              <p className="text-muted-foreground mt-1">
                Clique em <Info className="inline h-3 w-3" /> para saber mais
                sobre cada método.
              </p>
            </div>
            {methodsLoading ? (
              <div className="flex flex-col items-center justify-center py-10 space-y-4 text-muted-foreground">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p>Carregando métodos...</p>
              </div>
            ) : (
              <div className="grid gap-3 sm:grid-cols-2 max-h-[420px] overflow-y-auto pr-1">
                {methods.map((method) => {
                  const Icon = iconMap[method.icon] || Info;
                  const isSelected = selectedMethod?.id === method.id;
                  return (
                    <Card
                      key={method.id}
                      className={`cursor-pointer border transition-all duration-200 hover:shadow-md ${isSelected ? `${method.bg} ring-2 ring-primary` : "hover:border-primary/30"}`}
                      onClick={() => setSelectedMethod(method)}
                    >
                      <CardContent className="p-4 space-y-3">
                        <div className="flex items-start justify-between">
                          <div className="flex items-center gap-2">
                            <Icon className={`h-5 w-5 ${method.color}`} />
                            <div>
                              <p className="font-semibold text-sm">
                                {method.name}
                              </p>
                              <p className="text-xs text-muted-foreground">
                                {method.tagline}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-1 flex-shrink-0">
                            {isSelected && (
                              <div className="w-5 h-5 bg-primary rounded-full flex items-center justify-center">
                                <Check className="h-3 w-3 text-primary-foreground" />
                              </div>
                            )}
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-6 w-6"
                              onClick={(e) => {
                                e.stopPropagation();
                                setDetailMethod(method);
                              }}
                            >
                              <Info className="h-3 w-3" />
                            </Button>
                          </div>
                        </div>
                        {/* Distribution bar */}
                        <div className="flex h-2 rounded-full overflow-hidden gap-0.5">
                          {method.split.map((s) => (
                            <div
                              key={s.label}
                              className={`${s.color} transition-all`}
                              style={{ width: `${s.percent}%` }}
                            />
                          ))}
                        </div>
                        <div className="flex flex-wrap gap-1">
                          {method.split.map((s) => (
                            <Badge
                              key={s.label}
                              variant="secondary"
                              className="text-xs px-1.5 py-0"
                            >
                              {s.label} {s.percent}%
                            </Badge>
                          ))}
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            )}
            <div className="flex gap-3">
              <Button variant="outline" onClick={() => setStep(2)}>
                <ChevronLeft className="mr-2 h-4 w-4" /> Voltar
              </Button>
              <Button
                className="flex-1"
                disabled={!selectedMethod}
                onClick={() => setStep(4)}
              >
                Continuar <ChevronRight className="ml-2 h-4 w-4" />
              </Button>
            </div>
          </div>
        )}

        {/* Step 4: Income */}
        {step === 4 && selectedMethod && (
          <div className="space-y-6 animate-in fade-in">
            <div>
              <h2 className="text-2xl font-bold">Qual é sua renda mensal?</h2>
              <p className="text-muted-foreground mt-1">
                Com base na <strong>{selectedMethod.name}</strong>, calcularemos
                a distribuição ideal.
              </p>
            </div>
            <div className="space-y-4">
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground font-medium">
                  R$
                </span>
                <input
                  type="number"
                  placeholder="0,00"
                  value={monthlyIncome}
                  onChange={(e) => setMonthlyIncome(e.target.value)}
                  className="w-full pl-12 pr-4 py-4 text-2xl font-bold border rounded-xl bg-background focus:outline-none focus:ring-2 focus:ring-primary text-foreground"
                />
              </div>
              {monthlyIncome && parseFloat(monthlyIncome) > 0 && (
                <Card className="animate-in fade-in">
                  <CardContent className="p-4 space-y-3">
                    <p className="text-sm font-medium text-muted-foreground">
                      Sua distribuição com {selectedMethod.name}:
                    </p>
                    {selectedMethod.split.map((s) => {
                      const amount =
                        (parseFloat(monthlyIncome) * s.percent) / 100;
                      return (
                        <div
                          key={s.label}
                          className="flex items-center justify-between"
                        >
                          <div className="flex items-center gap-2">
                            <div
                              className={`w-3 h-3 rounded-full ${s.color}`}
                            />
                            <span className="text-sm">{s.label}</span>
                            <span className="text-xs text-muted-foreground">
                              ({s.percent}%)
                            </span>
                          </div>
                          <span className="font-bold text-sm">
                            R${" "}
                            {amount.toLocaleString("pt-BR", {
                              minimumFractionDigits: 2,
                            })}
                          </span>
                        </div>
                      );
                    })}
                  </CardContent>
                </Card>
              )}
            </div>
            <div className="flex gap-3">
              <Button variant="outline" onClick={() => setStep(3)}>
                <ChevronLeft className="mr-2 h-4 w-4" /> Voltar
              </Button>
              <Button
                className="flex-1"
                disabled={
                  !monthlyIncome || parseFloat(monthlyIncome) <= 0 || loading
                }
                onClick={handleComplete}
              >
                {loading ? "Configurando..." : "Começar minha jornada 🚀"}
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
