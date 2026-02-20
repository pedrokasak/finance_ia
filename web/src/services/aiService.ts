import { api } from '@/api/client';

export interface AIInsight {
    id: string;
    type: 'warning' | 'tip' | 'achievement' | 'projection';
    title: string;
    content: string;
    plan: string;
    generated_at: string;
    expires_at: string;
}

export interface HealthScore {
    score: number;
    level: string;
    savings_rate: number;
}

const aiService = {
    getInsight: () =>
        api.get<AIInsight>('/ai/insight').then((r) => r.data),

    getFullAnalysis: () =>
        api.get<{ insights: AIInsight[] }>('/ai/analysis').then((r) => r.data),

    getHealthScore: () =>
        api.get<HealthScore>('/ai/health-score').then((r) => r.data),
};

export default aiService;
