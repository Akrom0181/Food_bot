package utils

func HandleError(err error) {
	if err != nil {
		Error(err)
	}
}
