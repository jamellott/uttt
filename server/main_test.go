package main

import "testing"

func testHandler(in, out string, shouldFail bool) func(*testing.T) {
	return func(t *testing.T) {
		str, err := ValidatePath(in)
		if shouldFail && err == nil {
			t.Errorf("ValidatePath should have errored but did not")
		} else if !shouldFail && err != nil {
			t.Errorf("ValidatePath returned unexpected error: %v", err.Error())
		}

		if str != out {
			t.Errorf("ValidatePath returned %#v, expected %#v", str, out)
		}
	}
}

func TestValidatePath(t *testing.T) {
	t.Run("PassesOnValidFilename", testHandler("/index.html", "./ui/index.html", false))
	t.Run("PassesOnValidNestedFile", testHandler("/img/img1.png", "./ui/img/img1.png", false))
	t.Run("PassesOnFileWithoutLeadingSlash", testHandler("img/img1.png", "./ui/img/img1.png", false))
	t.Run("NoFailOnEmptyString", testHandler("", "./ui", false))
	t.Run("FailOnRelativePathBreakout", testHandler("../../../etc/passwd", "", true))
}
