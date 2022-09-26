package helper

import (
	"os"
	"path/filepath"
)

// CodeSave 保存代码到本地
func CodeSave(code []byte) (string, error) {
	filename := GetUUID()
	dirname := filepath.Join("code", filename)
	filePath := dirname + "/main.go"

	err := os.Mkdir(dirname, os.ModePerm)
	if err != nil {
		return "", err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(code)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

//
//func CodeRun(file string, testCase *models.TestCase) (msg string, err error) {
//	cmd := exec.Command("go", "run", file)
//	var out, stdErr bytes.Buffer
//	cmd.Stdout = &out
//	cmd.Stderr = &stdErr
//
//	pipe, err := cmd.StdinPipe()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	io.WriteString(pipe, testCase.Input)
//
//	var start, end runtime.MemStats // 计算内存
//	runtime.ReadMemStats(&start)
//	// 根据测试的输入案例，进行运行，拿到输出结果和标准的输出结果进行比对
//	if err := cmd.Run(); err != nil {
//		log.Println(err, stdErr.String())
//		if err.Error() == "exit status 2" {
//			msg = stdErr.String()
//			return msg, errors.New("编译错误")
//		}
//	}
//	// 运行超内存
//	runtime.ReadMemStats(&end)
//	if start.Alloc / 1000 - (end.Alloc / 1000) >
//	// 答案错误
//	if testCase.Output != out.String() {
//		return "", errors.New("答案错误")
//	}
//
//}
