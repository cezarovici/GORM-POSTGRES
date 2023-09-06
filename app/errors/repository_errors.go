package apperrors

import "fmt"

type RepositoryError struct {
	Issue      error
	Caller     string
	MethodName string
}

const areaErrRepo = "Repository"

func (e *RepositoryError) Error() string {
	var res [4]string

	res[0] = fmt.Sprintf("\nArea: %s", areaErrRepo)
	res[1] = fmt.Sprintf("Caller: %s", e.Caller)
	res[2] = fmt.Sprintf("Method Name: %s", e.MethodName)
	res[3] = fmt.Sprintf("Issue: %s", e.Issue.Error())

	return res[0] + "\n" + res[1] + "\n" + res[2] + "\n" + res[3]
}
