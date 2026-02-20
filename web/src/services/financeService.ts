import { api } from '@/api/client';

export interface Transaction {
    id?: string;
    category_id?: string;
    type: 'income' | 'expense';
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
    type: 'income' | 'expense';
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
    type?: 'income' | 'expense';
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
        api.get<DashboardSummary>('/finance/dashboard').then((r) => r.data),

    // Transactions
    createTransaction: (tx: Transaction) =>
        api.post<Transaction>('/finance/transactions', tx, {
            headers: { 'Idempotency-Key': crypto.randomUUID() },
        }).then((r) => r.data),

    updateTransaction: (id: string, tx: Partial<Transaction>) =>
        api.put<Transaction>(`/finance/transactions/${id}`, tx).then((r) => r.data),

    listTransactions: (filter: TransactionFilter = {}) => {
        const params = new URLSearchParams();
        if (filter.page) params.set('page', String(filter.page));
        if (filter.limit) params.set('limit', String(filter.limit));
        if (filter.type) params.set('type', filter.type);
        if (filter.start_date) params.set('start_date', filter.start_date);
        if (filter.end_date) params.set('end_date', filter.end_date);
        if (filter.category_id) params.set('category_id', filter.category_id);
        return api.get<TransactionListResponse>(`/finance/transactions?${params}`).then((r) => r.data);
    },

    deleteTransaction: (id: string) =>
        api.delete(`/finance/transactions/${id}`).then((r) => r.data),

    // Categories
    getCategories: () =>
        api.get<Category[]>('/finance/categories').then((r) => r.data),

    createCategory: (cat: Omit<Category, 'id' | 'is_default'>) =>
        api.post<Category>('/finance/categories', cat).then((r) => r.data),

    // Financial Methods
    getMethods: () =>
        api.get<FinancialMethod[]>('/finance/methods').then((r) => r.data),

    // Budget
    getBudget: () =>
        api.get<Budget>('/finance/budget').then((r) => r.data),

    upsertBudget: (budget: Omit<Budget, 'id' | 'period' | 'needs_amount' | 'wants_amount' | 'savings_amount'>) =>
        api.post<Budget>('/finance/budget', budget).then((r) => r.data),
};

export default financeService;
