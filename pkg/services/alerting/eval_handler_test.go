package alerting

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type conditionStub struct {
	firing   bool
	operator string
	matches  []*EvalMatch
	noData   bool
}

func (c *conditionStub) Eval(context *EvalContext) (*ConditionResult, error) {
	return &ConditionResult{Firing: c.firing, EvalMatches: c.matches, Operator: c.operator, NoDataFound: c.noData}, nil
}

func TestAlertingEvaluationHandler(t *testing.T) {
	Convey("Test alert evaluation handler", t, func() {
		handler := NewEvalHandler()

		Convey("Show return triggered with single passing condition", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{&conditionStub{
					firing: true,
					matches: []*EvalMatch{
						{
							Metric: "abc",
							Firing: true,
						},
					},
				}},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, true)
			So(context.ConditionEvals, ShouldEqual, "abc: true = true")
		})

		Convey("Show return triggered with single passing condition2", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{&conditionStub{
					firing:   true,
					operator: "and",
					matches: []*EvalMatch{
						{
							Metric:   "abc",
							Firing:   true,
							Operator: "and",
						},
					},
				}},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, true)
			So(context.ConditionEvals, ShouldEqual, "abc: true = true")
		})

		Convey("Show return false if any of the condition is failing with AND operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, false)
			So(context.ConditionEvals, ShouldEqual, "abc1: [true AND false] = false")
		})

		Convey("Show return true if any of the condition is passing with OR operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, true)
			So(context.ConditionEvals, ShouldEqual, "abc1: [true OR false] = true")
		})

		Convey("Show return true if one condition is failing with nested OR operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, true)
			So(context.ConditionEvals, ShouldEqual, "abc1: [[true AND true] OR false] = true")
		})

		Convey("Show return false if one condition is passing with nested OR operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, false)
			So(context.ConditionEvals, ShouldEqual, "abc1: [[true AND false] OR false] = false")
		})

		Convey("Show return false if a condition is failing with nested AND operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, false)
			So(context.ConditionEvals, ShouldEqual, "abc1: [[true AND false] AND true] = false")
		})

		Convey("Show return true if a condition is passing with nested OR operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   true,
						operator: "and",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						firing:   true,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   true,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, true)
			So(context.ConditionEvals, ShouldEqual, "abc1: [[true OR false] OR true] = true")
		})

		Convey("Should return false if no condition is firing using OR operator", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, false)
			So(context.ConditionEvals, ShouldEqual, "abc1: [[false OR false] OR false] = false")
		})

		Convey("Should retuasdfrn no data if one condition has nodata", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						noData:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Valid:    true,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						noData:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Valid:    true,
								Operator: "or",
							},
						},
					},
					&conditionStub{
						noData:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Metric:   "abc1",
								Valid:    true,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.NoDataFound, ShouldBeFalse)
		})

		Convey("Should return no data if one condition has nodata", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Valid:    false,
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.Firing, ShouldEqual, false)
			So(context.NoDataFound, ShouldBeTrue)
		})

		Convey("Should return no data if both conditions have no data and using AND", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Valid:    false,
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Valid:    true,
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.NoDataFound, ShouldBeFalse)
		})

		Convey("Should not return no data if both conditions have no data and using OR", func() {
			context := NewEvalContext(context.TODO(), &Rule{
				MatchSerie: 1,
				Conditions: []Condition{
					&conditionStub{
						firing:   false,
						operator: "and",
						matches: []*EvalMatch{
							{
								Valid:    false,
								Metric:   "abc1",
								Firing:   false,
								Operator: "and",
							},
						},
					},
					&conditionStub{
						firing:   false,
						operator: "or",
						matches: []*EvalMatch{
							{
								Valid:    true,
								Metric:   "abc1",
								Firing:   false,
								Operator: "or",
							},
						},
					},
				},
			})

			handler.Eval(context)
			So(context.NoDataFound, ShouldBeTrue)
		})
	})
}
