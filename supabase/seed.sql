-- Seed data for onboarding financial methods
-- Idempotent: updates existing rows by key.

DO $$
BEGIN
  IF to_regclass('public.financial_methods') IS NULL THEN
    RAISE NOTICE 'public.financial_methods does not exist yet; skipping seed.';
    RETURN;
  END IF;

  INSERT INTO public.financial_methods
    (key, name, tagline, description, for_who, icon, color, bg, split_raw, is_active, created_at, updated_at)
  VALUES
    (
      '5-3-2',
      'Regra 5-3-2',
      'Foco em equilíbrio simples',
      '50% para necessidades essenciais, 30% para estilo de vida e 20% para investimentos e reserva. Leitura rápida da 50-30-20 em formato compacto.',
      'Ideal para onboarding rápido com uma regra fácil de memorizar.',
      'CircleEllipsis',
      'text-lime-400',
      'bg-lime-500/10 border-lime-500/30',
      '[{"label":"Necessidades","percent":50,"color":"bg-lime-500"},{"label":"Estilo de Vida","percent":30,"color":"bg-teal-500"},{"label":"Investimentos","percent":20,"color":"bg-cyan-500"}]',
      true,
      now(),
      now()
    ),
    (
      '60-10-10-10-10',
      'Regra 60-10-10-10-10',
      'Mais granular para metas simultâneas',
      '60% para custos fixos, 10% para investimentos de longo prazo, 10% para reserva, 10% para lazer e 10% para objetivos pessoais.',
      'Ideal para quem quer dividir melhor os 40% restantes em metas específicas.',
      'ListTree',
      'text-indigo-400',
      'bg-indigo-500/10 border-indigo-500/30',
      '[{"label":"Custos Fixos","percent":60,"color":"bg-indigo-500"},{"label":"Invest. Longo Prazo","percent":10,"color":"bg-violet-500"},{"label":"Reserva","percent":10,"color":"bg-sky-500"},{"label":"Lazer","percent":10,"color":"bg-pink-500"},{"label":"Objetivos","percent":10,"color":"bg-amber-500"}]',
      true,
      now(),
      now()
    ),
    (
      '50-30-20',
      'Regra 50-30-20',
      'O método mais popular do mundo',
      '50% da renda para necessidades básicas (moradia, alimentação, saúde), 30% para desejos (lazer, restaurantes) e 20% para poupança e investimentos.',
      'Ideal para quem está começando a organizar as finanças.',
      'PieChart',
      'text-emerald-400',
      'bg-emerald-500/10 border-emerald-500/30',
      '[{"label":"Necessidades","percent":50,"color":"bg-emerald-500"},{"label":"Desejos","percent":30,"color":"bg-blue-500"},{"label":"Investimentos","percent":20,"color":"bg-purple-500"}]',
      true,
      now(),
      now()
    ),
    (
      '60-20-20',
      'Regra 60-20-20',
      'Mais moderna e adaptável',
      '60% para despesas fixas, 20% para investimentos e 20% para lazer e objetivos. Mais flexível que a 50-30-20 para famílias ou quem mora em cidade cara.',
      'Ideal para famílias, quem mora em cidade cara ou precisa de flexibilidade.',
      'Layers',
      'text-cyan-400',
      'bg-cyan-500/10 border-cyan-500/30',
      '[{"label":"Despesas Fixas","percent":60,"color":"bg-cyan-500"},{"label":"Investimentos","percent":20,"color":"bg-violet-500"},{"label":"Lazer e Metas","percent":20,"color":"bg-pink-500"}]',
      true,
      now(),
      now()
    ),
    (
      'pay-yourself-first',
      'Pague-se Primeiro',
      'Invista antes de gastar',
      'Assim que receber o salário, invista imediatamente antes de qualquer gasto. O percentual de investimento é definido por você — o foco está na ordem, não na proporção fixa.',
      'Excelente para construção de patrimônio e quem já tem alguma disciplina financeira.',
      'Coins',
      'text-yellow-400',
      'bg-yellow-500/10 border-yellow-500/30',
      '[{"label":"Investimento Imediato","percent":20,"color":"bg-yellow-500"},{"label":"Viva com o resto","percent":80,"color":"bg-orange-400"}]',
      true,
      now(),
      now()
    ),
    (
      'emergency-reserve',
      'Reserva de Emergência',
      'Sua primeira meta financeira',
      'Foco em construir uma reserva de 3 a 6 meses do custo de vida (ou 12 meses para autônomos). Guardar em Tesouro Selic ou CDB com liquidez diária.',
      'Essencial para quem não possui reserva de emergência ainda.',
      'Landmark',
      'text-sky-400',
      'bg-sky-500/10 border-sky-500/30',
      '[{"label":"Reserva (3-6 meses)","percent":30,"color":"bg-sky-500"},{"label":"Gastos Mensais","percent":70,"color":"bg-slate-400"}]',
      true,
      now(),
      now()
    ),
    (
      'envelopes',
      'Método dos Envelopes',
      'Controle total por categoria',
      'Cada categoria recebe um envelope com valor fixo. Quando o envelope acaba, parou de gastar naquela categoria. Simples e visual.',
      'Ideal para quem gasta de forma impulsiva por categorias.',
      'Wallet',
      'text-amber-400',
      'bg-amber-500/10 border-amber-500/30',
      '[{"label":"Moradia","percent":30,"color":"bg-amber-500"},{"label":"Alimentação","percent":20,"color":"bg-orange-500"},{"label":"Transporte","percent":15,"color":"bg-yellow-500"},{"label":"Outros","percent":35,"color":"bg-red-500"}]',
      true,
      now(),
      now()
    ),
    (
      'zero-based',
      'Orçamento Base Zero',
      'Cada real tem uma destinação',
      'Renda - Despesas = 0. Todo real é alocado intencionalmente. Máximo controle, pois nenhum dinheiro some sem destino definido.',
      'Ideal para quem quer controle total e obsessivo das finanças.',
      'Target',
      'text-blue-400',
      'bg-blue-500/10 border-blue-500/30',
      '[{"label":"Fixos","percent":40,"color":"bg-blue-500"},{"label":"Variáveis","percent":35,"color":"bg-indigo-500"},{"label":"Reservas","percent":25,"color":"bg-violet-500"}]',
      true,
      now(),
      now()
    ),
    (
      'goal-based',
      'Planejamento por Objetivos',
      'Separe por metas: curto, médio e longo prazo',
      'Cada objetivo tem seu horizonte temporal e tipo de investimento. Viagem → renda fixa. Carro → renda fixa moderada. Aposentadoria → renda variável.',
      'Ideal para quem tem múltiplos objetivos financeiros simultâneos.',
      'TrendingUp',
      'text-green-400',
      'bg-green-500/10 border-green-500/30',
      '[{"label":"Curto Prazo (até 2a)","percent":20,"color":"bg-green-500"},{"label":"Médio Prazo (2-5a)","percent":25,"color":"bg-teal-500"},{"label":"Longo Prazo (5a+)","percent":20,"color":"bg-emerald-700"},{"label":"Gastos","percent":35,"color":"bg-gray-400"}]',
      true,
      now(),
      now()
    ),
    (
      'savings-rate',
      'Taxa de Poupança',
      'Maximize o quanto você poupa',
      'Taxa de poupança = valor investido ÷ renda. 10% é básico, 20% é bom, 30%+ acelera a riqueza.',
      'Para quem quer acelerar a construção de patrimônio.',
      'Percent',
      'text-rose-400',
      'bg-rose-500/10 border-rose-500/30',
      '[{"label":"Investimentos (30%+)","percent":30,"color":"bg-rose-500"},{"label":"Gastos (70% ou menos)","percent":70,"color":"bg-pink-300"}]',
      true,
      now(),
      now()
    ),
    (
      '70-20-10',
      'Regra 70-20-10',
      'Gastos, poupança e dívidas',
      '70% para gastos mensais (necessários e desejos), 20% para poupança e investimentos, 10% para quitação de dívidas ou doações.',
      'Ideal para quem tem dívidas e quer estruturar a saída delas.',
      'Flame',
      'text-red-400',
      'bg-red-500/10 border-red-500/30',
      '[{"label":"Gastos","percent":70,"color":"bg-red-500"},{"label":"Poupança","percent":20,"color":"bg-pink-500"},{"label":"Dívidas/Doação","percent":10,"color":"bg-fuchsia-500"}]',
      true,
      now(),
      now()
    )
  ON CONFLICT (key)
  DO UPDATE SET
    name = EXCLUDED.name,
    tagline = EXCLUDED.tagline,
    description = EXCLUDED.description,
    for_who = EXCLUDED.for_who,
    icon = EXCLUDED.icon,
    color = EXCLUDED.color,
    bg = EXCLUDED.bg,
    split_raw = EXCLUDED.split_raw,
    is_active = EXCLUDED.is_active,
    updated_at = now();
END $$;
