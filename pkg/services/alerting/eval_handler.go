package alerting

import (
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
)

// DefaultEvalHandler is responsible for evaluating the alert rule.
type DefaultEvalHandler struct {
	log             log.Logger
	alertJobTimeout time.Duration
}

// NewEvalHandler is the `DefaultEvalHandler` constructor.
func NewEvalHandler() *DefaultEvalHandler {
	return &DefaultEvalHandler{
		log:             log.New("alerting.evalHandler"),
		alertJobTimeout: time.Second * 5,
	}
}

// Eval evaluated the alert rule.
func (e *DefaultEvalHandler) Eval(context *EvalContext) {
	firing := true
	noDataFound := true
	firings := make([]bool, 0, 10)
	noDataFounds := make([]bool, 0, 10)
	conditionEvals := ""

	evalMatches := make(map[string][]*EvalMatch)
	for i := 0; i < len(context.Rule.Conditions); i++ {
		condition := context.Rule.Conditions[i]
		cr, err := condition.Eval(context)
		if err != nil {
			context.Error = err
		}

		// break if condition could not be evaluated
		if context.Error != nil {
			goto END
		}

		for j := 0; j < len(cr.EvalMatches); j++ {
			evalMatch := cr.EvalMatches[j]
			evalMatches[evalMatch.Metric] = append(evalMatches[evalMatch.Metric], &EvalMatch{
				Metric:   evalMatch.Metric,
				Tags:     evalMatch.Tags,
				Value:    evalMatch.Value,
				Firing:   evalMatch.Firing,
				Valid:    evalMatch.Valid,
				Operator: cr.Operator,
			})
		}
	}

	for metric, serieEvalMatches := range evalMatches {
		metricFiring := true
		metricNoDataFound := true
		metricConditionEvals := ""
		for i := 0; i < len(serieEvalMatches); i++ {
			serieEvalMatch := serieEvalMatches[i]
			if i == 0 {
				metricFiring = serieEvalMatch.Firing
				metricNoDataFound = !serieEvalMatch.Valid
			}

			if serieEvalMatch.Operator == "or" {
				metricFiring = metricFiring || serieEvalMatch.Firing
				metricNoDataFound = metricNoDataFound || !serieEvalMatch.Valid
			} else {
				metricFiring = metricFiring && serieEvalMatch.Firing
				metricNoDataFound = metricNoDataFound && !serieEvalMatch.Valid
			}

			if i > 0 {
				metricConditionEvals = "[" + metricConditionEvals + " " + strings.ToUpper(serieEvalMatch.Operator) + " " + strconv.FormatBool(serieEvalMatch.Firing) + "]"
			} else {
				metricConditionEvals = strconv.FormatBool(metricFiring)
			}
		}

		noDataFounds = append(noDataFounds, metricNoDataFound)
		if metricFiring {
			firings = append(firings, metricFiring)
		}

		metricConditionEvals = metric + ": " + metricConditionEvals
		if conditionEvals == "" {
			conditionEvals = "{" + metricConditionEvals + "}"
		} else {
			conditionEvals = conditionEvals + ", {" + metricConditionEvals + "}"
		}

		// The number of curves that satisfy the condition
		context.EvalMatches = append(context.EvalMatches, serieEvalMatches...)
	}

	if len(firings) < context.Rule.MatchSerie {
		firing = false
	}

	for _, metricNoDataFound := range noDataFounds {
		noDataFound = noDataFound && metricNoDataFound
	}

END:
	context.ConditionEvals = conditionEvals + " = " + strconv.FormatBool(firing) + " and matchSerie=" + strconv.Itoa(context.Rule.MatchSerie)
	context.Firing = firing
	context.NoDataFound = noDataFound
	context.EndTime = time.Now()

	elapsedTime := context.EndTime.Sub(context.StartTime).Nanoseconds() / int64(time.Millisecond)
	metrics.MAlertingExecutionTime.Observe(float64(elapsedTime))
}
