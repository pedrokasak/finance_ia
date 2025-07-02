import React from 'react';
import { cn } from '@/lib/utils';

interface CustomTooltipProps {
  active?: boolean;
  payload?: any[];
  label?: string;
  className?: string;
}

export function CustomTooltip({ active, payload, label, className }: CustomTooltipProps) {
  if (active && payload && payload.length) {
    return (
      <div className={cn(
        "bg-card/95 backdrop-blur-sm border border-border/50 rounded-xl shadow-lg p-4 min-w-[200px]",
        "ring-1 ring-black/5 dark:ring-white/10",
        className
      )}>
        {label && (
          <p className="text-sm font-medium text-foreground mb-2 border-b border-border/30 pb-2">
            {label}
          </p>
        )}
        <div className="space-y-2">
          {payload.map((entry, index) => (
            <div key={index} className="flex items-center justify-between space-x-4">
              <div className="flex items-center space-x-2">
                <div 
                  className="w-3 h-3 rounded-full ring-2 ring-white/20"
                  style={{ backgroundColor: entry.color }}
                />
                <span className="text-sm text-muted-foreground capitalize">
                  {entry.dataKey === 'receitas' ? 'Receitas' : 
                   entry.dataKey === 'despesas' ? 'Despesas' : 
                   entry.dataKey === 'economia' ? 'Economia' :
                   entry.name || entry.dataKey}
                </span>
              </div>
              <span className="text-sm font-semibold text-foreground">
                {typeof entry.value === 'number' 
                  ? `R$ ${entry.value.toLocaleString('pt-BR')}` 
                  : entry.value}
              </span>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return null;
}

export function PieTooltip({ active, payload }: CustomTooltipProps) {
  if (active && payload && payload.length) {
    const data = payload[0];
    return (
      <div className="bg-card/95 backdrop-blur-sm border border-border/50 rounded-xl shadow-lg p-4 min-w-[180px] ring-1 ring-black/5 dark:ring-white/10">
        <div className="flex items-center space-x-3">
          <div 
            className="w-4 h-4 rounded-full ring-2 ring-white/20"
            style={{ backgroundColor: data.payload.color }}
          />
          <div className="flex-1">
            <p className="text-sm font-medium text-foreground">
              {data.payload.name}
            </p>
            <div className="flex items-center justify-between mt-1">
              <span className="text-sm font-semibold text-foreground">
                R$ {data.value.toLocaleString('pt-BR')}
              </span>
              <span className="text-xs text-muted-foreground ml-2">
                {((data.value / data.payload.total) * 100).toFixed(1)}%
              </span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return null;
}

export function BarTooltip({ active, payload, label }: CustomTooltipProps) {
  if (active && payload && payload.length) {
    return (
      <div className="bg-card/95 backdrop-blur-sm border border-border/50 rounded-xl shadow-lg p-4 min-w-[220px] ring-1 ring-black/5 dark:ring-white/10">
        <p className="text-sm font-medium text-foreground mb-3 border-b border-border/30 pb-2">
          {label}
        </p>
        <div className="space-y-2">
          {payload.map((entry, index) => (
            <div key={index} className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <div 
                  className="w-3 h-3 rounded-sm ring-1 ring-white/20"
                  style={{ backgroundColor: entry.color }}
                />
                <span className="text-sm text-muted-foreground capitalize">
                  {entry.dataKey === 'atual' ? 'Atual' :
                   entry.dataKey === 'anterior' ? 'Anterior' :
                   entry.dataKey === 'meta' ? 'Meta' :
                   entry.name || entry.dataKey}
                </span>
              </div>
              <span className="text-sm font-semibold text-foreground">
                R$ {entry.value.toLocaleString('pt-BR')}
              </span>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return null;
}