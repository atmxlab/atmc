package errors

func Combine(err1 error, err2 error) error {
	j := NewJoiner()

	j.Join(err1, err2)

	return j.Err()
}
