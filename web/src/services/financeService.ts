import { api } from "@/api/client";

export interface Transaction {
  id?: string;
  category_id?: string;
  type: "income" | "expense";
  amount: number;
  description: string;
  tagline?: string;
  date: string;
  is_recurring?: boolean;
  idempotency_key?: string;
  category?: { id: string; name: string; color: string; icon: string };
}

export interface Category {
  id: string;
  name: string;
  type: "income" | "expense";
  color: string;
  icon: string;
  is_default: boolean;
}

export interface FinancialMethodSplit {
  label: string;
  percent: number;
  color: string;
}

export interface FinancialMethod {
  id: string;
  key: string;
  name: string;
  tagline?: string;
  description: string;
  for_who: string;
  icon: string;
  color: string;
  bg: string;
  split: FinancialMethodSplit[];
}

export const DEFAULT_FINANCIAL_METHODS: FinancialMethod[] = [
  {
    id: "default-50-30-20",
    key: "50-30-20",
    name: "Regra 50-30-20",
    tagline: "O método mais popular do mundo",
    description:
      "50% da renda para necessidades básicas (moradia, alimentação, saúde), 30% para desejos (lazer, restaurantes) e 20% para poupança e investimentos.",
    for_who: "Ideal para quem está começando a organizar as finanças.",
    icon: "PieChart",
    color: "text-emerald-400",
    bg: "bg-emerald-500/10 border-emerald-500/30",
    split: [
      { label: "Necessidades", percent: 50, color: "bg-emerald-500" },
      { label: "Desejos", percent: 30, color: "bg-blue-500" },
      { label: "Investimentos", percent: 20, color: "bg-purple-500" },
    ],
  },
  {
    id: "default-60-20-20",
    key: "60-20-20",
    name: "Regra 60-20-20",
    tagline: "Mais moderna e adaptável",
    description:
      "60% para despesas fixas, 20% para investimentos e 20% para lazer e objetivos. Mais flexível que a 50-30-20 para famílias ou quem mora em cidade cara.",
    for_who:
      "Ideal para famílias, quem mora em cidade cara ou precisa de flexibilidade.",
    icon: "Layers",
    color: "text-cyan-400",
    bg: "bg-cyan-500/10 border-cyan-500/30",
    split: [
      { label: "Despesas Fixas", percent: 60, color: "bg-cyan-500" },
      { label: "Investimentos", percent: 20, color: "bg-violet-500" },
      { label: "Lazer e Metas", percent: 20, color: "bg-pink-500" },
    ],
  },
  {
    id: "default-70-20-10",
    key: "70-20-10",
    name: "Regra 70-20-10",
    tagline: "Gastos, poupança e dívidas",
    description:
      "70% para gastos mensais (necessários e desejos), 20% para poupança e investimentos, 10% para quitação de dívidas ou doações.",
    for_who: "Ideal para quem tem dívidas e quer estruturar a saída delas.",
    icon: "Flame",
    color: "text-red-400",
    bg: "bg-red-500/10 border-red-500/30",
    split: [
      { label: "Gastos", percent: 70, color: "bg-red-500" },
      { label: "Poupança", percent: 20, color: "bg-pink-500" },
      { label: "Dívidas/Doação", percent: 10, color: "bg-fuchsia-500" },
    ],
  },
  {
    id: "default-goal-based",
    key: "goal-based",
    name: "Planejamento por Objetivos",
    tagline: "Separe por metas: curto, médio e longo prazo",
    description:
      "Cada objetivo tem seu horizonte temporal e tipo de investimento. Viagem → renda fixa. Carro → renda fixa moderada. Aposentadoria → renda variável.",
    for_who: "Ideal para quem tem múltiplos objetivos financeiros simultâneos.",
    icon: "TrendingUp",
    color: "text-green-400",
    bg: "bg-green-500/10 border-green-500/30",
    split: [
      { label: "Curto Prazo (até 2a)", percent: 20, color: "bg-green-500" },
      { label: "Médio Prazo (2-5a)", percent: 25, color: "bg-teal-500" },
      { label: "Longo Prazo (5a+)", percent: 20, color: "bg-emerald-700" },
      { label: "Gastos", percent: 35, color: "bg-gray-400" },
    ],
  },
  {
    id: "default-envelopes",
    key: "envelopes",
    name: "Método dos Envelopes",
    tagline: "Controle total por categoria",
    description:
      "Cada categoria recebe um 'envelope' com valor fixo. Quando o envelope acaba, parou de gastar naquela categoria. Simples e visual.",
    for_who: "Ideal para quem gasta de forma impulsiva por categorias.",
    icon: "Wallet",
    color: "text-amber-400",
    bg: "bg-amber-500/10 border-amber-500/30",
    split: [
      { label: "Moradia", percent: 30, color: "bg-amber-500" },
      { label: "Alimentação", percent: 20, color: "bg-orange-500" },
      { label: "Transporte", percent: 15, color: "bg-yellow-500" },
      { label: "Outros", percent: 35, color: "bg-red-500" },
    ],
  },
  {
    id: "default-zero-based",
    key: "zero-based",
    name: "Orçamento Base Zero",
    tagline: "Cada real tem uma destinação",
    description:
      "Renda - Despesas = 0. Todo real é alocado intencionalmente. Máximo controle, pois nenhum dinheiro 'some' sem destino definido.",
    for_who: "Ideal para quem quer controle total e obsessivo das finanças.",
    icon: "Target",
    color: "text-blue-400",
    bg: "bg-blue-500/10 border-blue-500/30",
    split: [
      { label: "Fixos", percent: 40, color: "bg-blue-500" },
      { label: "Variáveis", percent: 35, color: "bg-indigo-500" },
      { label: "Reservas", percent: 25, color: "bg-violet-500" },
    ],
  },
];

export interface Budget {
  id?: string;
  total_income: number;
  needs_percent: number;
  wants_percent: number;
  savings_percent: number;
  needs_amount?: number;
  wants_amount?: number;
  savings_amount?: number;
  period?: string;
}

export interface DashboardSummary {
  total_income: number;
  total_expenses: number;
  balance: number;
  savings_rate: number;
  health_score: number;
  health_level: string;
  budget?: Budget;
  category_breakdown: Array<{
    category_id: string;
    category_name: string;
    color: string;
    total: number;
    percentage: number;
  }>;
  monthly_trend: Array<{ month: string; income: number; expenses: number }>;
  days_until_negative?: number;
}

export interface TransactionFilter {
  page?: number;
  limit?: number;
  type?: "income" | "expense";
  start_date?: string;
  end_date?: string;
  category_id?: string;
}

export interface TransactionListResponse {
  data: Transaction[];
  total: number;
  page: number;
  limit: number;
}

const financeService = {
  // Dashboard
  getDashboard: () =>
    api.get<DashboardSummary>("/finance/dashboard").then((r) => r.data),

  // Transactions
  createTransaction: (tx: Transaction) =>
    api
      .post<Transaction>("/finance/transactions", tx, {
        headers: { "Idempotency-Key": crypto.randomUUID() },
      })
      .then((r) => r.data),

  updateTransaction: (id: string, tx: Partial<Transaction>) =>
    api.put<Transaction>(`/finance/transactions/${id}`, tx).then((r) => r.data),

  listTransactions: (filter: TransactionFilter = {}) => {
    const params = new URLSearchParams();
    if (filter.page) params.set("page", String(filter.page));
    if (filter.limit) params.set("limit", String(filter.limit));
    if (filter.type) params.set("type", filter.type);
    if (filter.start_date) params.set("start_date", filter.start_date);
    if (filter.end_date) params.set("end_date", filter.end_date);
    if (filter.category_id) params.set("category_id", filter.category_id);
    return api
      .get<TransactionListResponse>(`/finance/transactions?${params}`)
      .then((r) => r.data);
  },

  deleteTransaction: (id: string) =>
    api.delete(`/finance/transactions/${id}`).then((r) => r.data),

  // Categories
  getCategories: () =>
    api.get<Category[]>("/finance/categories").then((r) => r.data),

  createCategory: (cat: Omit<Category, "id" | "is_default">) =>
    api.post<Category>("/finance/categories", cat).then((r) => r.data),

  updateCategory: (
    id: string,
    cat: Partial<Omit<Category, "id" | "is_default">>,
  ) => api.put<Category>(`/finance/categories/${id}`, cat).then((r) => r.data),

  deleteCategory: (id: string) =>
    api.delete(`/finance/categories/${id}`).then((r) => r.data),

  // Financial Methods
  getMethods: async () => {
    try {
      const data = await api
        .get<FinancialMethod[]>("/finance/methods")
        .then((r) => r.data);
      return Array.isArray(data) && data.length > 0
        ? data
        : DEFAULT_FINANCIAL_METHODS;
    } catch {
      return DEFAULT_FINANCIAL_METHODS;
    }
  },

  // Budget
  getBudget: () => api.get<Budget>("/finance/budget").then((r) => r.data),

  upsertBudget: (
    budget: Omit<
      Budget,
      "id" | "period" | "needs_amount" | "wants_amount" | "savings_amount"
    >,
  ) => api.post<Budget>("/finance/budget", budget).then((r) => r.data),
};

export default financeService;
