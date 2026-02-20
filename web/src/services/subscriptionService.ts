import { api } from '@/api/client';

export interface Subscription {
    plan: string;
    status: string;
    current_period_end?: string;
    features: PlanFeatures;
}

export interface PlanFeatures {
    plan: string;
    ai_insights_per_week: number;
    full_analysis: boolean;
    projections: boolean;
    smart_alerts: boolean;
    export_data: boolean;
    max_transactions: number;
    features: string[];
}

export interface CheckoutResponse {
    session_id: string;
    url: string;
}

const subscriptionService = {
    getMySubscription: () =>
        api.get<Subscription>('/subscription/').then((r) => r.data),

    createCheckout: (plan: 'premium' | 'pro', billing_type: 'monthly' | 'yearly') =>
        api.post<CheckoutResponse>('/subscription/checkout', { plan, billing_type }).then((r) => r.data),

    createPortal: () =>
        api.post<{ url: string }>('/subscription/portal').then((r) => r.data),
};

export default subscriptionService;
