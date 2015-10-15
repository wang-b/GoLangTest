package common

/*
 * 全局错误处理函数
 */
func CheckError(err error) {
	if err != nil {
		panic(err)   //抛出异常
	}
}
