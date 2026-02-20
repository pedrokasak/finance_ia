import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import {
    Check, Info, Wallet, PieChart, Target, TrendingUp, Landmark, Layers, Coins, Percent, Flame, Loader2
} from 'lucide-react';
import financeService, { FinancialMethod } from '@/services/financeService';

interface FinancialMethodModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    currentMethodId?: string;
    onSelect: (methodId: string) => void;
}

const iconMap: Record<string, React.ComponentType<{ className?: string }>> = {
    PieChart, Layers, Coins, Landmark, Wallet, Target, TrendingUp, Percent, Flame
};

export function FinancialMethodModal({ open, onOpenChange, currentMethodId, onSelect }: FinancialMethodModalProps) {
    const [methods, setMethods] = useState<FinancialMethod[]>([]);
    const [loading, setLoading] = useState(true);
    const [detailMethod, setDetailMethod] = useState<FinancialMethod | null>(null);

    useEffect(() => {
        if (open && methods.length === 0) {
            setLoading(true);
            financeService.getMethods()
                .then((data) => setMethods(data || []))
                .catch(console.error)
                .finally(() => setLoading(false));
        }
    }, [open, methods.length]);

    return (
        <>
            <Dialog open={open} onOpenChange={onOpenChange}>
                <DialogContent className="sm:max-w-2xl max-h-[85vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle>Método de Planejamento Financeiro</DialogTitle>
                        <DialogDescription>
                            Escolha ou altere a estratégia que usaremos para classificar seu orçamento mensal.
                        </DialogDescription>
                    </DialogHeader>

                    {loading ? (
                        <div className="flex flex-col items-center justify-center py-10 space-y-4 text-muted-foreground">
                            <Loader2 className="h-8 w-8 animate-spin text-primary" />
                            <p>Carregando métodos...</p>
                        </div>
                    ) : (
                        <div className="grid gap-3 sm:grid-cols-2 mt-4">
                            {methods.map((method) => {
                                const Icon = iconMap[method.icon] || Info;
                                const isSelected = currentMethodId === method.id;
                                return (
                                    <Card
                                        key={method.id}
                                        className={`cursor-pointer border transition-all duration-200 hover:shadow-md ${isSelected ? `${method.bg} ring-2 ring-primary` : 'hover:border-primary/30'}`}
                                        onClick={() => {
                                            onSelect(method.id);
                                            onOpenChange(false);
                                        }}
                                    >
                                        <CardContent className="p-4 space-y-3">
                                            <div className="flex items-start justify-between">
                                                <div className="flex items-center gap-2">
                                                    <Icon className={`h-5 w-5 ${method.color}`} />
                                                    <div>
                                                        <p className="font-semibold text-sm">{method.name}</p>
                                                        <p className="text-xs text-muted-foreground">{method.tagline}</p>
                                                    </div>
                                                </div>
                                                <div className="flex items-center gap-1 flex-shrink-0">
                                                    {isSelected && (
                                                        <div className="w-5 h-5 bg-primary rounded-full flex items-center justify-center">
                                                            <Check className="h-3 w-3 text-primary-foreground" />
                                                        </div>
                                                    )}
                                                    <Button variant="ghost" size="icon" className="h-6 w-6" onClick={(e) => { e.stopPropagation(); setDetailMethod(method); }}>
                                                        <Info className="h-3 w-3" />
                                                    </Button>
                                                </div>
                                            </div>
                                            <div className="flex h-2 rounded-full overflow-hidden gap-0.5">
                                                {method.split.map((s) => (
                                                    <div key={s.label} className={`${s.color} transition-all`} style={{ width: `${s.percent}%` }} />
                                                ))}
                                            </div>
                                            <div className="flex flex-wrap gap-1">
                                                {method.split.map((s) => (
                                                    <Badge key={s.label} variant="secondary" className="text-xs px-1.5 py-0">
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
                </DialogContent>
            </Dialog>

            {/* Exibir detalhes num Sub-Dialog */}
            {detailMethod && (
                <Dialog open={!!detailMethod} onOpenChange={() => setDetailMethod(null)}>
                    <DialogContent className="max-w-md">
                        <DialogHeader>
                            <DialogTitle className="flex items-center gap-2">
                                {(() => {
                                    const Icon = iconMap[detailMethod.icon] || Info;
                                    return <Icon className={`h-5 w-5 ${detailMethod.color}`} />;
                                })()}
                                {detailMethod.name}
                            </DialogTitle>
                        </DialogHeader>
                        <div className="space-y-4">
                            <p className="text-sm text-muted-foreground leading-relaxed">{detailMethod.description}</p>
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
                                <p className="text-xs text-muted-foreground">{detailMethod.for_who}</p>
                            </div>
                            <Button className="w-full" onClick={() => {
                                onSelect(detailMethod.id);
                                setDetailMethod(null);
                                onOpenChange(false);
                            }}>
                                Utilizar este método
                            </Button>
                        </div>
                    </DialogContent>
                </Dialog>
            )}
        </>
    );
}
