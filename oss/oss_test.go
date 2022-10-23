package oss

import (
	"context"
	"fmt"
	"testing"
)

var (
	ctx = context.Background()
)

func TestUploadContent(t *testing.T) {
	file := Oss{}
	fmt.Println(file.UploadWithContent(ctx, "testFile.txt", []byte("testFile--First"), nil))
}

func TestUploadWithURL(t *testing.T) {
	file := Oss{}
	fmt.Println(file.UploadWithURL(ctx, "testFile.jpg", "https://t7.baidu.com/it/u=1595072465,3644073269&fm=193&f=GIF", nil))
}

func TestUploadWithPath(t *testing.T) {
	file := Oss{}
	fmt.Println(file.UploadWithPath(ctx, "testFile.go", "./oss_test.go", nil))
}
