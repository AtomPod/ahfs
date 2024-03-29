package errcode

const (
	FileNotExist          ErrorCode = 400200
	FileDownloadDirError  ErrorCode = 400201 // 不能下载一个文件夹
	FileNotDirError       ErrorCode = 400202 // 不是文件夹
	FileStorageFulled     ErrorCode = 400203 // 文件储存已满
	FileRootOperateError  ErrorCode = 400204 // 不能对根文件进行操作
	FileDirNotExists      ErrorCode = 400205 // 文件夹不存在
	FileParentNotDirError ErrorCode = 400206 // 父结点不是一个文件夹
	FileAlreadyExists     ErrorCode = 400207 // 文件（夹）已经存在
	FileTooLarge          ErrorCode = 400208
	FilenameFormatError   ErrorCode = 400209 // 文件名格式错误
)
