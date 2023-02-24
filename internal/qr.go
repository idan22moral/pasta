package internal

import "github.com/skip2/go-qrcode"

func PrintQR(content string) error {
	const WHITE_SPACE string = "\033[47m  \033[0m"
	const BLACK_SPACE string = "\033[40m  \033[0m"

	q, err := qrcode.New(content, qrcode.Low)

	if err != nil {
		return err
	}

	for _, bitRow := range q.Bitmap() {
		for _, bit := range bitRow {
			if bit {
				print(WHITE_SPACE)
			} else {
				print(BLACK_SPACE)
			}
		}
		print("\n")
	}

	return nil
}
