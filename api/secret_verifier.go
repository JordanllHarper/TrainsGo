package main

type secretVerifier interface {
	Verify(ref, secret string) (isValid bool, err error)
}

type inMemSecretVerifier map[string]string

func (sv inMemSecretVerifier) Verify(ref, secret string) (isValid bool, err error) {
	val, ok := sv[ref]
	if !ok || val != secret {
		return false, nil
	}
	return true, nil
}
