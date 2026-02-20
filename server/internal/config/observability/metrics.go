package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestsTotal counts all HTTP requests
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	// RequestDuration measures HTTP request latency
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "finance_ia_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	// TransactionsCreated counts new transactions
	TransactionsCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_transactions_created_total",
		Help: "Total number of financial transactions created",
	}, []string{"type"}) // income / expense

	// AIInsightsGenerated counts AI insight generation
	AIInsightsGenerated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_ai_insights_generated_total",
		Help: "Total number of AI insights generated",
	}, []string{"plan", "insight_type"})

	// AIResponseDuration measures time to generate AI insight
	AIResponseDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "finance_ia_ai_response_duration_seconds",
		Help:    "Time to generate AI insight",
		Buckets: []float64{0.5, 1, 2, 5, 10, 20, 30},
	}, []string{"plan"})

	// PaymentEvents counts Stripe webhook events
	PaymentEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_payment_events_total",
		Help: "Total number of payment events processed",
	}, []string{"event_type", "status"})

	// QueuePublished counts messages published to queues
	QueuePublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_queue_published_total",
		Help: "Messages published to queues",
	}, []string{"queue"})

	// QueueConsumed counts messages consumed from queues
	QueueConsumed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_queue_consumed_total",
		Help: "Messages consumed from queues",
	}, []string{"queue", "status"}) // status: success / error

	// CacheHits counts Redis cache hits
	CacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_cache_hits_total",
		Help: "Redis cache hits",
	}, []string{"key_prefix"})

	// CacheMisses counts Redis cache misses
	CacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finance_ia_cache_misses_total",
		Help: "Redis cache misses",
	}, []string{"key_prefix"})

	// HealthScoreDistribution tracks user health score levels
	HealthScoreDistribution = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "finance_ia_health_score",
		Help:    "Distribution of user financial health scores",
		Buckets: []float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
	})
)
