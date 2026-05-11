# academic

当前已落地：

- 学期列表接口：`GET /api/academic/semesters`
- 课程表接口：`GET /api/academic/schedule`
- 考试安排接口：`GET /api/academic/exams`
- 成绩单接口：`GET /api/academic/grades`
- 基于 `academic_provider` 的 mock 读模型，便于无外部依赖联调
- 课程表、考试、成绩、学期查询成功后统一写入审计日志

约束：

- 只允许已完成教务绑定的登录用户读取教务数据
- 教务查询能力只由后端根据 IAM 绑定资料反查学号，不依赖客户端传学号
- 外部教务系统接入统一经 `providers/academic_provider`，不在模块内散落第三方调用

当前边界：

- 当前为 mock provider 读能力，尚未接入真实教务系统连接器
- 当前未实现课表同步任务、本地缓存快照与更细粒度的成绩统计
