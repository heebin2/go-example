package helper

func GetLRC(message []byte) byte {

	var buf byte = 0

	for _, v := range message {
		buf ^= v
	}

	return buf
}
