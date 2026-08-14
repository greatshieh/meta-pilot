package utils

import (
	"server/pkg/errcode"
)

// 客户端错误 2000xx
const (
	// 发送短信验证码失败
	ErrSendSmsCode int = iota + 200001 // 发送短信验证码失败
	// 短信验证码错误
	ErrBadSmsCode
	// 手机号码或学生姓名错误，请联系班主任或管理员
	ErrWrongStudentInfo
	// 学生信息获取失败
	ErrGetStuInfo
	// 获取学生报告失败
	ErrGetReport
	// 验证码已过期
	ErrExpiredCode
	ErrSendRepeat   // 请勿重复发送
	ErrGetWorksList // 查找作品失败
)

// 首页相关错误 2001xx
const (
	ErrCategoryCount           int = iota + 200101 // 获取分类数据错误
	ErrGradeSubjectAverages                        // 获取年级平均分错误
	ErrGradeHistoricalAverages                     // 获取分数趋势错误
	ErrGradeSubjectEntry                           // 获取录入进度失败
)

// 学生管理错误 2002xx
const (
	ErrStudentList            int = iota + 200201 // 获取学生信息列表失败
	ErrStudentDuplicateQuery                      // 获取重名学生信息失败
	ErrStudentDeletedQuery                        // 获取已删除学生信息失败
	ErrStudentAlreadyExists                       // 学生信息已经存在
	ErrStudentCreate                              // 创建学生信息失败
	ErrStudentUpdate                              // 学生信息更新失败
	ErrStudentNotFound                            // 学生信息不存在
	ErrGetStudent                                 // 获取学生信息失败
	ErrStudentDelete                              // 学生信息删除失败
	ErrStudentRestore                             // 学生信息恢复失败
	ErrStudentOperation                           // 学生信息操作失败
	ErrStudentPermanentDelete                     // 永久删除学生信息失败
	ErrStudentFileRead                            // 读取文件失败
	ErrStudentSheetOpen                           // 打开工作表失败
	ErrStudentHeaderInvalid                       // 表头不正确
	ErrStudentClear                               // 清空已删除学生失败
	ErrStudentUnkonw                              // 未知错误
)

// 成绩管理错误 2003xx
const (
	ErrScoreNotFound             int = iota + 200301 // 考试信息不存在
	ErrUpdateScore                                   // 成绩更新失败
	ErrCreateScore                                   // 创建成绩失败
	ErrGetScoreList                                  // 获取成绩列表失败
	ErrScoreDeletedQuery                             // 获取已删除成绩失败
	ErrScoreDelete                                   // 成绩删除失败
	ErrScoreRestore                                  // 成绩恢复失败
	ErrScoreOperation                                // 成绩操作失败
	ErrScoreClearDeleted                             // 成绩清空失败
	ErrScorePermanentDelete                          // 永久删除成绩失败
	ErrGetClassAvgScore                              //按年级获取班级平均分失败
	ErrGetClassScoreDistribution                     // 按年级获取班级分数段人数分布失败
	ErrGetClassScoreRate                             // 按年级获取班级优秀率和及格率失败
	ErrGetClassTopRate                               //按年级获取班级top20和top50失败
	ErrGetClassTopList                               //按年级获取班级top名单失败
	ErrIncorrectEexamInfo                            // 考试信息不争取
)

// 作品管理错误 2004xx
const (
	ErrWorksSave            int = iota + 200401 // 保存作品失败
	ErrWorksGetDeleted                          // 获取已删除作品失败
	ErrWorksList                                // 获取作品列表失败
	ErrStudentDuplicateName                     // 学生姓名重复
	ErrStudentInvalidPhone                      // 学生手机号格式不正确
	ErrWorksAlreadyExists                       // 作品已存在
	ErrWorksClearDeleted                        // 清空已删除作品失败
	ErrWorksPermanentDelete                     // 永久删除作品失败
	ErrWorksDelete                              // 作品删除失败
	ErrWorksRestore                             // 作品恢复失败
	ErrWorksUnkonw                              // 未知错误
)

func init() {
	errcode.NewErrCode(ErrSendSmsCode, 200, "发送短信验证码失败")
	errcode.NewErrCode(ErrBadSmsCode, 200, "短信验证码错误")
	errcode.NewErrCode(ErrWrongStudentInfo, 200, "手机号码或学生姓名错误，请联系班主任或管理员")
	errcode.NewErrCode(ErrGetStuInfo, 200, "学生信息获取失败")
	errcode.NewErrCode(ErrGetReport, 200, "获取学生报告失败")
	errcode.NewErrCode(ErrExpiredCode, 200, "验证码已过期")
	errcode.NewErrCode(ErrSendRepeat, 200, "请勿重复发送")
	errcode.NewErrCode(ErrGetWorksList, 200, "查找作品失败")
	errcode.NewErrCode(ErrCategoryCount, 200, "获取分类数据错误")
	errcode.NewErrCode(ErrGradeSubjectAverages, 200, "获取年级平均分错误")
	errcode.NewErrCode(ErrGradeHistoricalAverages, 200, "获取分数趋势错误")
	errcode.NewErrCode(ErrGradeSubjectEntry, 200, "获取录入进度失败")
	errcode.NewErrCode(ErrStudentList, 200, "获取学生信息列表失败")
	errcode.NewErrCode(ErrStudentDuplicateQuery, 200, "获取重名学生信息失败")
	errcode.NewErrCode(ErrStudentDeletedQuery, 200, "获取已删除学生信息失败")
	errcode.NewErrCode(ErrStudentAlreadyExists, 200, "学生信息已经存在")
	errcode.NewErrCode(ErrStudentCreate, 200, "创建学生信息失败")
	errcode.NewErrCode(ErrStudentUpdate, 200, "学生信息更新失败")
	errcode.NewErrCode(ErrStudentNotFound, 200, "学生信息不存在")
	errcode.NewErrCode(ErrGetStudent, 200, "获取学生信息失败")
	errcode.NewErrCode(ErrStudentDelete, 200, "学生信息删除失败")
	errcode.NewErrCode(ErrStudentRestore, 200, "学生信息恢复失败")
	errcode.NewErrCode(ErrStudentOperation, 200, "学生信息操作失败")
	errcode.NewErrCode(ErrStudentPermanentDelete, 200, "永久删除学生信息失败")
	errcode.NewErrCode(ErrStudentFileRead, 200, "读取文件失败")
	errcode.NewErrCode(ErrStudentSheetOpen, 200, "打开工作表失败")
	errcode.NewErrCode(ErrStudentHeaderInvalid, 200, "表头不正确")
	errcode.NewErrCode(ErrStudentClear, 200, "清空已删除学生失败")
	errcode.NewErrCode(ErrStudentUnkonw, 200, "未知错误")
	errcode.NewErrCode(ErrScoreNotFound, 200, "考试信息不存在")
	errcode.NewErrCode(ErrUpdateScore, 200, "成绩更新失败")
	errcode.NewErrCode(ErrCreateScore, 200, "创建成绩失败")
	errcode.NewErrCode(ErrGetScoreList, 200, "获取成绩列表失败")
	errcode.NewErrCode(ErrScoreDeletedQuery, 200, "获取已删除成绩失败")
	errcode.NewErrCode(ErrScoreDelete, 200, "成绩删除失败")
	errcode.NewErrCode(ErrScoreRestore, 200, "成绩恢复失败")
	errcode.NewErrCode(ErrScoreOperation, 200, "成绩操作失败")
	errcode.NewErrCode(ErrScoreClearDeleted, 200, "成绩清空失败")
	errcode.NewErrCode(ErrScorePermanentDelete, 200, "永久删除成绩失败")
	errcode.NewErrCode(ErrGetClassAvgScore, 200, "按年级获取班级平均分失败")
	errcode.NewErrCode(ErrGetClassScoreDistribution, 200, "按年级获取班级分数段人数分布失败")
	errcode.NewErrCode(ErrGetClassScoreRate, 200, "按年级获取班级优秀率和及格率失败")
	errcode.NewErrCode(ErrGetClassTopRate, 200, "按年级获取班级top20和top50失败")
	errcode.NewErrCode(ErrGetClassTopList, 200, "按年级获取班级top名单失败")
	errcode.NewErrCode(ErrIncorrectEexamInfo, 200, "考试信息不争取")
	errcode.NewErrCode(ErrWorksSave, 200, "保存作品失败")
	errcode.NewErrCode(ErrWorksGetDeleted, 200, "获取已删除作品失败")
	errcode.NewErrCode(ErrWorksList, 200, "获取作品列表失败")
	errcode.NewErrCode(ErrStudentDuplicateName, 200, "学生姓名重复")
	errcode.NewErrCode(ErrStudentInvalidPhone, 200, "学生手机号格式不正确")
	errcode.NewErrCode(ErrWorksAlreadyExists, 200, "作品已存在")
	errcode.NewErrCode(ErrWorksClearDeleted, 200, "清空已删除作品失败")
	errcode.NewErrCode(ErrWorksPermanentDelete, 200, "永久删除作品失败")
	errcode.NewErrCode(ErrWorksDelete, 200, "作品删除失败")
	errcode.NewErrCode(ErrWorksRestore, 200, "作品恢复失败")
	errcode.NewErrCode(ErrWorksUnkonw, 200, "未知错误")
}
