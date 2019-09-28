package validation

import "github.com/sean-tech/webservice/logging"

func FormValid(form interface{}) error {
	// form valid
	valid := Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return err
	}
	if !check {
		for _, err := range valid.Errors {
			return err
		}
	}
	return nil;
}

func MarkErrors(errors []*Error)  {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
	return
}
