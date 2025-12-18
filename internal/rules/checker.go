package rules

import "sync"

func New(ctx *Context) *Checker {
	return &Checker{ctx: ctx}
}

func (c *Checker) CheckAll(rules []Rule) []CheckResult {
	results := make([]CheckResult, 0, len(rules))

	// Parallel execution
	resultsCh := make(chan CheckResult, len(rules))
	var wg sync.WaitGroup

	for _, r := range rules {
		wg.Add(1)
		go func(r Rule) {
			defer wg.Done()
			resultsCh <- c.Check(r)
		}(r)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		results = append(results, result)
	}

	return results
}
func (c *Checker) Check(rule Rule) CheckResult {
	if rule.FixSpec.CustomCheck != nil {
		res, err := rule.FixSpec.CustomCheck(c.ctx)
		if err != nil {
			return CheckResult{Passed: false, Message: err.Error()}
		}
		return *res
	}
	switch rule.Type {
	case RuleTypeFile:
		return c.CheckFile(rule)
	case RuleTypeFolder:
		return c.CheckFolder(rule)
	case RuleTypeMulti:
		return c.CheckMulti(rule)
	default:
		return CheckResult{RuleID: rule.ID, Passed: false}
	}
}

func (c *Checker) CheckFile(rule Rule) CheckResult {
	return c.CheckRule(rule, CheckSpec{
		Validator: ValidateExists,
	})
}
func (c *Checker) CheckFolder(rule Rule) CheckResult {
	return c.CheckRule(rule, CheckSpec{
		Validator: ValidateExists,
	})
}
func (c *Checker) CheckMulti(rule Rule) CheckResult {
	fileResult := c.CheckFile(rule)
	if fileResult.Passed {
		return fileResult
	}

	folderResult := c.CheckFolder(rule)
	if folderResult.Passed {
		return folderResult
	}

	return CheckResult{
		RuleID: rule.ID,
		Passed: false,
	}
}
