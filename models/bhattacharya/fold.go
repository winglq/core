package bhattacharya

import (
	"coralreef-ci/models/issues"
	"errors"
	"math"
	"fmt"
)

func (model *Model) Fold(issues []issues.Issue) (float64, error) {
	issueCount := len(issues)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED")
	}

	score := 0.00

	for i := 0.10; i < 0.90; i += 0.10 {
		correct := 0

		trainCount := int(Round(i * float64(issueCount)))
		testCount := issueCount - trainCount

		model.Learn(issues[0:trainCount])

		for j := trainCount + 1; j < issueCount; j++ {
			assignees := model.Predict(issues[j])
			if assignees[0] == issues[j].Assignee || assignees[1] == issues[j].Assignee || assignees[2] == issues[j].Assignee{
				correct += 1
			} else {
				continue
			}
		}
		fmt.Println("Fold ", Round(i * 10))
		fmt.Println(" Accuracy ", float64(correct) / float64(testCount))
		score += float64(correct) / float64(testCount)
	}
	return score / 10.00, nil
}

func AppendCopy(slice []issues.Issue, elements ...issues.Issue) []issues.Issue {
    n := len(slice)
    total := len(slice) + len(elements)
		newSlice := make([]issues.Issue, total)
    if total > cap(slice) {
        // Reallocate. Grow to 1.5 times the new size, so we can still grow.
        newSize := total*3/2 + 1
        newSlice = make([]issues.Issue, total, newSize)
    }
		copy(newSlice, slice)
    copy(newSlice[n:], elements)
    return newSlice
}


func Append(slice []issues.Issue, elements ...issues.Issue) []issues.Issue {
    n := len(slice)
    total := len(slice) + len(elements)
    if total > cap(slice) {
        // Reallocate. Grow to 1.5 times the new size, so we can still grow.
        newSize := total*3/2 + 1
        newSlice := make([]issues.Issue, total, newSize)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[:total]
    copy(slice[n:], elements)
    return slice
}

func (model *Model) FoldImpl(train []issues.Issue, test []issues.Issue, tossingGraphLength int) float64 {
	testCount, correct := len(test), 0
	model.Learn(train)
	for j := 0; j < len(test); j++ {
		assignees := model.Predict(test[j])
		if assignees[0] == test[j].Assignee {
			correct += 1
		} else if assignees[1] == test[j].Assignee && tossingGraphLength > 1 {
			correct += 1
		} else if assignees[2] == test[j].Assignee && tossingGraphLength == 3 {
			correct += 1
		} else {
			continue
		}
	}
	return float64(correct) / float64(testCount)
}

func (model *Model) TwoFold(issues []issues.Issue, tossingGraphLength int) (float64, error) {
	length := len(issues)
	trainEndPos := int(0.50 * float64(length))
	trainIssues := AppendCopy(issues[0:trainEndPos])
	testIssues := AppendCopy(issues[trainEndPos+1:length])
	score := 0.00

	score += model.FoldImpl(trainIssues, testIssues, tossingGraphLength)
	score += model.FoldImpl(testIssues, trainIssues, tossingGraphLength)
	return score / 2.00, nil
}

func (model *Model) TenFold(issues []issues.Issue) (float64, error) {
	pStart, testStartPos, testEndPos, testCount := 0, 0, 0, 0
	score := 0.00

	length := len(issues)
	for i := 0.10; i < 0.90; i += 0.10 {
		correct := 0
		testStartPos = int(i * float64(length))
		testEndPos = int((i + 0.10) * float64(length))
		trainIssues := AppendCopy(issues[pStart:testStartPos-1], issues[testEndPos+1:length]...)
		trainIssuesLength := len(trainIssues)
		testIssues := AppendCopy(issues[testStartPos:testEndPos])
		testCount = len(testIssues)

		model.Learn(trainIssues)
		for j := 0; j < len(testIssues); j++ {
			assignees := model.Predict(testIssues[j])
			if assignees[0] == issues[j].Assignee || assignees[1] == issues[j].Assignee || assignees[2] == issues[j].Assignee{
				correct += 1
			} else {
				continue
			}
		}
		fmt.Println("Fold ", Round(i * 10))
		fmt.Println("Accuracy ", float64(correct) / float64(testCount))
		fmt.Println("Correct", float64(correct))
		fmt.Println("Train Count", float64(trainIssuesLength))
		fmt.Println("Test Count", float64(testCount))
		score += float64(correct) / float64(testCount)
	}
	return score / 10.00, nil
}

func Round(input float64) float64 {
	return math.Floor(input + 0.5)
}
