# WeOUC Next API 文档

> Base URL: `http://<host>:8080/api`
>
> 认证方式：大部分接口需要在请求头中携带 `Authorization: Bearer <token>`

---

## 目录

- [1. 认证模块](#1-认证模块)
  - [1.1 用户注册](#11-用户注册)
  - [1.2 用户登录](#12-用户登录)
  - [1.3 微信登录](#13-微信登录)
- [2. 学生模块](#2-学生模块)
  - [2.1 创建学生档案](#21-创建学生档案)
  - [2.2 获取学生档案](#22-获取学生档案)
  - [2.3 更新学生档案](#23-更新学生档案)
  - [2.4 通过OpenID更新通知设置](#24-通过openid更新通知设置)
- [3. 课程模块](#3-课程模块)
  - [3.1 获取课程列表](#31-获取课程列表)
  - [3.2 同步课程数据](#32-同步课程数据)
- [4. 成绩模块](#4-成绩模块)
  - [4.1 获取成绩详情](#41-获取成绩详情)
  - [4.2 同步成绩数据](#42-同步成绩数据)
- [5. 排名模块](#5-排名模块)
  - [5.1 获取排名配置](#51-获取排名配置)
  - [5.2 动态排名计算](#52-动态排名计算)
- [6. 培养方案模块](#6-培养方案模块)
  - [6.1 同步培养方案](#61-同步培养方案)
  - [6.2 获取培养方案详情](#62-获取培养方案详情)
- [7. Feed 动态模块](#7-feed-动态模块)
  - [7.1 获取动态列表](#71-获取动态列表)
  - [7.2 获取审核状态](#72-获取审核状态)
- [8. 跑腿模块](#8-跑腿模块)
  - [8.1 获取跑腿列表](#81-获取跑腿列表)
  - [8.2 获取跑腿详情](#82-获取跑腿详情)
  - [8.3 发布跑腿](#83-发布跑腿)
  - [8.4 接受跑腿](#84-接受跑腿)
- [9. 二手交易模块](#9-二手交易模块)
  - [9.1 获取二手交易列表](#91-获取二手交易列表)
  - [9.2 获取二手交易详情](#92-获取二手交易详情)
  - [9.3 发布二手交易](#93-发布二手交易)
  - [9.4 收藏二手交易](#94-收藏二手交易)
- [10. 拼车模块](#10-拼车模块)
  - [10.1 获取拼车列表](#101-获取拼车列表)
  - [10.2 发布拼车](#102-发布拼车)
- [11. 资料模块](#11-资料模块)
  - [11.1 获取资料列表](#111-获取资料列表)
  - [11.2 获取资料详情](#112-获取资料详情)
  - [11.3 发布资料](#113-发布资料)
  - [11.4 收藏资料](#114-收藏资料)
- [12. 失物招领模块](#12-失物招领模块)
  - [12.1 获取失物招领列表](#121-获取失物招领列表)
  - [12.2 获取失物招领详情](#122-获取失物招领详情)
  - [12.3 发布失物招领](#123-发布失物招领)
- [13. 文件上传模块](#13-文件上传模块)
  - [13.1 上传文件](#131-上传文件)
- [14. 微信回调模块](#14-微信回调模块)
  - [14.1 微信内容审核回调](#141-微信内容审核回调)
- [15. 通知模块](#15-通知模块)
  - [15.1 发送通知](#151-发送通知)
  - [15.2 获取通知任务状态](#152-获取通知任务状态)

---

## 1. 认证模块

### 1.1 用户注册

`POST /auth/register`

**无需认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（唯一） |
| password | string | 是 | 密码（最少6位） |

响应 `201`：

```json
{
  "message": "User registered successfully"
}
```

错误响应 `400`：

```json
{
  "error": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}
```

错误响应 `500`：

```json
{
  "error": "Failed to create user (username might be taken)"
}
```

---

### 1.2 用户登录

`POST /auth/login`

**无需认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

响应 `200`：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

错误响应 `401`：

```json
{
  "error": "Invalid credentials"
}
```

---

### 1.3 微信登录

`POST /auth/wechat/login`

**无需认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 微信登录凭证 code |
| app_id | string | 是 | 微信小程序 AppID |

响应 `200`：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "openid": "oXXXXXXXXXXXX",
    "app_id": "wx1234567890abcdef",
    "isNewUser": true,
    "userInfo": {
      "userId": 1,
      "nickname": "oXXXXXXXXXXXX",
      "avatarUrl": ""
    }
  }
}
```

---

## 2. 学生模块

### 2.1 创建学生档案

`POST /student`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| student_id | string | 是 | 学号 |
| major | string | 是 | 专业 |
| college | string | 否 | 学院 |
| education_type | string | 否 | 培养类型（如：本科生、硕士） |
| grade | string | 是 | 年级（如：2023） |
| notify_grade | bool | 否 | 是否接收成绩通知，默认 true |
| notify_course | bool | 否 | 是否接收课程通知，默认 true |
| notify_task | bool | 否 | 是否接收待办通知，默认 true |
| notify_quota | int | 否 | 通知配额，默认 0 |

响应 `201`：

```json
{
  "student_id": "2023001",
  "major": "计算机科学与技术",
  "college": "计算机学院",
  "education_type": "本科生",
  "grade": "2023",
  "notify_grade": true,
  "notify_course": true,
  "notify_task": true,
  "notify_quota": 0,
  "created_at": "2026-05-01T00:00:00Z",
  "updated_at": "2026-05-01T00:00:00Z"
}
```

---

### 2.2 获取学生档案

`GET /student`

**需要 Bearer Token 认证**

响应 `200`：

```json
{
  "student_id": "2023001",
  "major": "计算机科学与技术",
  "college": "计算机学院",
  "education_type": "本科生",
  "grade": "2023",
  "notify_grade": true,
  "notify_course": true,
  "notify_task": true,
  "notify_quota": 0,
  "created_at": "2026-05-01T00:00:00Z",
  "updated_at": "2026-05-01T00:00:00Z"
}
```

错误响应 `404`：

```json
{
  "error": "Student profile not found"
}
```

---

### 2.3 更新学生档案

`PUT /student`

**需要 Bearer Token 认证**

请求体（所有字段均为可选，仅传需要更新的字段）：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| major | string | 否 | 专业 |
| college | string | 否 | 学院 |
| education_type | string | 否 | 培养类型 |
| grade | string | 否 | 年级 |
| notify_grade | bool | 否 | 是否接收成绩通知 |
| notify_course | bool | 否 | 是否接收课程通知 |
| notify_task | bool | 否 | 是否接收待办通知 |

响应 `200`：返回更新后的完整 Student 对象

---

### 2.4 通过OpenID更新通知设置

`PUT /student/notification/openid`

**无需 Bearer Token 认证**（通过 OpenID 鉴权）

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| open_id | string | 是 | 微信 OpenID |
| notify_grade | bool | 否 | 是否接收成绩通知 |
| notify_course | bool | 否 | 是否接收课程通知 |
| notify_task | bool | 否 | 是否接收待办通知 |

响应 `200`：返回更新后的完整 Student 对象

---

## 3. 课程模块

### 3.1 获取课程列表

`GET /student/courses`

**需要 Bearer Token 认证**

响应 `200`：

```json
[
  {
    "id": 1,
    "student_id": "2023001",
    "kch": "CS101",
    "kc_mc": "数据结构",
    "jg0101mc": "张老师",
    "jsgh": "T001",
    "kt_mc": "计科2301班",
    "xf": 3.0,
    "zxs": 48,
    "xsfl0": 32,
    "sktime": "周一1-2节",
    "skddmc": "教学楼A301",
    "skxqmc": "主校区",
    "kcxz": "专业必修",
    "kclb": "必修",
    "kkyx": "计算机学院",
    "khfs": "考试",
    "pkrs": 60,
    "xkrs": 55,
    "xkh": "202301001",
    "jx0404id": "JX001",
    "zhouxs": "2",
    "bj": "",
    "term": "2025-春季",
    "has_grade": false,
    "created_at": "2026-05-01T00:00:00Z",
    "updated_at": "2026-05-01T00:00:00Z"
  }
]
```

---

### 3.2 同步课程数据

`POST /student/courses`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| term | string | 否 | 学期标识（如：2025-春季） |
| data | UserCourse[] | 是 | 课程数据数组 |

UserCourse 字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| kch | string | 课程号 |
| kc_mc | string | 课程名称 |
| jg0101mc | string | 教师姓名 |
| jsgh | string | 教师工号 |
| kt_mc | string | 课堂名称 |
| xf | float | 学分 |
| zxs | int | 总学时 |
| xsfl0 | int | 讲课学时 |
| sktime | string | 上课时间 |
| skddmc | string | 上课地点 |
| skxqmc | string | 校区 |
| kcxz | string | 课程性质 |
| kclb | string | 课程类别 |
| kkyx | string | 开课院系 |
| khfs | string | 考核方式 |
| pkrs | int | 计划人数 |
| xkrs | int | 选课人数 |
| xkh | string | 选课号（唯一标识） |
| jx0404id | string | 教学班ID |
| zhouxs | string | 周学时 |
| bj | string | 标记 |

响应 `200`：

```json
{
  "message": "Courses synced successfully",
  "count": 5
}
```

---

## 4. 成绩模块

### 4.1 获取成绩详情

`GET /grade-detail`

**需要 Bearer Token 认证**

响应 `200`：

```json
[
  {
    "id": 1,
    "student_id": "2023001",
    "cj0708id": "G001",
    "kch": "CS101",
    "kc_mc": "数据结构",
    "xf": 3.0,
    "zxs": 48,
    "kcsx": "必修",
    "kcxzmc": "专业基础必修课程",
    "kccm": "专业类课程",
    "zcj": 92.5,
    "zcjstr": "92.5",
    "jd": 4.0,
    "ksfs": "闭卷",
    "ksxz": "初修取得",
    "xdlx": "初修",
    "xqmc": "2025秋季学期",
    "xqstr": "2025-2026-1",
    "xnxqid": "2025秋季学期",
    "ksdw": "计算机学院",
    "skjs": "张老师",
    "xkh": "202301001",
    "xs0101id": "XS001",
    "jx0404id": "JX001",
    "kz": 0,
    "cjbs": "",
    "created_at": "2026-05-01T00:00:00Z",
    "updated_at": "2026-05-01T00:00:00Z"
  }
]
```

---

### 4.2 同步成绩数据

`POST /grade-detail`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| data | GradeDetail[] | 是 | 成绩数据数组 |

GradeDetail 主要字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| cj0708id | string | 成绩唯一ID |
| kch | string | 课程号 |
| kc_mc | string | 课程名称 |
| xf | float | 学分 |
| zxs | int | 总学时 |
| kcsx | string | 课程属性 |
| kcxzmc | string | 课程性质名称 |
| kccm | string | 课程类别 |
| zcj | float | 总成绩 |
| zcjstr | string | 总成绩（字符串） |
| jd | float | 绩点 |
| ksfs | string | 考试方式 |
| ksxz | string | 考试性质 |
| xdlx | string | 修读类型 |
| xqmc | string | 学期名称 |
| xqstr | string | 学期代码 |
| xnxqid | string | 学年学期ID |
| ksdw | string | 考试单位 |
| skjs | string | 任课教师 |
| xkh | string | 选课号 |
| cjbs | string | 成绩标识 |

响应 `200`：

```json
{
  "message": "Grade details synced successfully",
  "count": 10
}
```

---

## 5. 排名模块

### 5.1 获取排名配置

`GET /ranking/config`

**无需 Bearer Token 认证**（通过查询参数鉴权）

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| openid | string | 是 | 微信 OpenID |
| xh | string | 是 | 学号 |

响应 `200`：

```json
{
  "data": {
    "categories": ["专业类课程", "公共基础课", "通识教育课", "体育"],
    "retake_modes": [
      { "label": "取最高分", "value": "max" },
      { "label": "取初修分", "value": "first" },
      { "label": "不及格重修按60", "value": "replace_fail_60" }
    ],
    "exemption_modes": [
      { "label": "过滤/忽略", "value": "filter" },
      { "label": "取最高分", "value": "max" },
      { "label": "取平均分", "value": "avg" }
    ],
    "scopes": [
      {
        "value": "计算机学院",
        "label": "计算机学院",
        "children": [
          {
            "value": "计算机科学与技术",
            "label": "计算机科学与技术",
            "children": [
              { "value": "2023", "label": "2023" }
            ]
          }
        ]
      }
    ],
    "defaults": {
      "exclude_categories": [],
      "exclude_keywords": [],
      "general_credit_limit": 0,
      "sport_num_limit": 0,
      "retake_mode": "max",
      "exemption_handling": "ignore",
      "calculator_type": "weighted"
    },
    "user_info": {
      "college": "计算机学院",
      "major": "计算机科学与技术",
      "grade": "2023"
    },
    "recommend_keywords": ["大学英语", "形势与政策"]
  },
  "status": 200
}
```

---

### 5.2 动态排名计算

`POST /ranking/dynamic`

**无需 Bearer Token 认证**（通过请求体鉴权）

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| openid | string | 是 | 微信 OpenID |
| xh | string | 是 | 学号 |
| term_code | string | 否 | 学期代码，为空则计算全量 |
| options | object | 是 | 排名配置选项 |
| options.exclude_categories | string[] | 否 | 排除的课程类别 |
| options.exclude_keywords | string[] | 否 | 排除的课程名称关键字 |
| options.general_credit_limit | float | 否 | 通识课学分上限（0=不限制） |
| options.sport_num_limit | float | 否 | 体育课数量上限（0=不限制） |
| options.retake_mode | string | 否 | 重修策略：max/first/replace_fail_60 |
| options.exemption_handling | string | 否 | 免修策略：max/avg/filter |
| options.calculator_type | string | 否 | 计算器类型：weighted/difficulty/bonus/custom |
| options.category_weight_config | object | 否 | 类别权重配置 |
| compare_college | string | 否 | 筛选学院 |
| compare_major | string | 否 | 筛选专业 |
| compare_grade | string | 否 | 筛选年级 |

响应 `200`：

```json
{
  "data": {
    "score": 88.5,
    "total_weight": 120.0,
    "strategy_id": "weighted_max_filter",
    "description": "标准加权 | 重修取最高 | 免修过滤",
    "rank": 15,
    "total_students": 200,
    "history": [
      {
        "term_code": "2023-2024-1",
        "score": 85.2,
        "total_weight": 22.0,
        "rank": 20,
        "total_students": 200
      },
      {
        "term_code": "ALL",
        "score": 88.5,
        "total_weight": 120.0,
        "rank": 15,
        "total_students": 200
      }
    ]
  },
  "status": 200
}
```

---

## 6. 培养方案模块

### 6.1 同步培养方案

`POST /training-programs/sync`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| college | string | 是 | 学院 |
| education_type | string | 是 | 培养类型 |
| major_name | string | 是 | 专业名称（含年级） |
| requirements | TrainingProgramRequirement[] | 否 | 学分要求列表 |
| courses | TrainingProgramCourse[] | 否 | 方案课程列表 |

响应 `200`：

```json
{
  "message": "Training program synced successfully"
}
```

---

### 6.2 获取培养方案详情

`GET /training-programs/detail`

**需要 Bearer Token 认证**

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| college | string | 是 | 学院 |
| education_type | string | 是 | 培养类型 |
| major_name | string | 是 | 专业名称（含年级） |

响应 `200`：

```json
{
  "college": "计算机学院",
  "education_type": "本科生",
  "major_name": "计算机科学与技术2023",
  "requirements": [
    {
      "id": 1,
      "college": "计算机学院",
      "education_type": "本科生",
      "major_name": "计算机科学与技术2023",
      "category": "专业类课程",
      "course_type": "必修",
      "required_credits": 60.0
    }
  ],
  "courses": [
    {
      "id": 1,
      "college": "计算机学院",
      "education_type": "本科生",
      "major_name": "计算机科学与技术2023",
      "course_code": "CS101",
      "course_name": "数据结构",
      "course_category": "专业类课程",
      "course_type": "必修",
      "credits": 3.0,
      "suggested_term": 3
    }
  ]
}
```

---

## 7. Feed 动态模块

### 7.1 获取动态列表

`GET /feed/list`

**需要 Bearer Token 认证**

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| feed_types | string | 否 | Feed类型，逗号分隔（如：errand,market） |
| keyword | string | 否 | 搜索关键词 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

feed_type 可选值：`errand`(跑腿) / `market`(二手交易) / `carpool`(拼车) / `resource`(资料) / `lostFound`(失物招领)

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "69f3a5debc1e8764471c7112",
        "feed_type": "errand",
        "feed_type_label": "跑腿",
        "title": "帮忙取快递",
        "desc": "北门快递站取件",
        "image": "",
        "publisher": "用户A",
        "publisher_initial": "Y",
        "created_at": "2026-05-01T00:00:00Z",
        "review_status": "approved",
        "review_remark": null,
        "extra": {
          "reward": "5元",
          "urgent": false,
          "route_start": "北门",
          "route_end": "宿舍楼",
          "deadline": "2026-05-02"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 7.2 获取审核状态

`GET /feed/review-status/:id`

**需要 Bearer Token 认证**

路径参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | FeedItem 的 MongoDB ID |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "review_status": "reviewing",
    "review_remark": null,
    "reviewed_at": null
  }
}
```

review_status 可选值：`reviewing`(审核中) / `approved`(已通过) / `rejected`(已拒绝)

---

## 8. 跑腿模块

### 8.1 获取跑腿列表

`GET /errand/list`

**需要 Bearer Token 认证**

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 分类筛选 |
| keyword | string | 否 | 搜索关键词 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "69f3a5debc1e8764471c7112",
        "feed_type": "errand",
        "app_id": "",
        "user_id": 1,
        "title": "帮忙取快递",
        "desc": "北门快递站取件",
        "image": "",
        "publisher": "",
        "publisher_initial": "",
        "created_at": "2026-05-01T00:00:00Z",
        "updated_at": "2026-05-01T00:00:00Z",
        "review_status": "reviewing",
        "reward": "5元",
        "route_start": "北门",
        "route_end": "宿舍楼",
        "deadline": "2026-05-02",
        "category": "快递",
        "is_accepted": false,
        "user_role": "publisher"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

`user_role` 可选值：`publisher`(发布者) / `acceptor`(接受者) / `viewer`(浏览者)

---

### 8.2 获取跑腿详情

`GET /errand/detail/:id`

**需要 Bearer Token 认证**

路径参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | FeedItem 的 MongoDB ID |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "item": { "...FeedItem对象..." },
    "user_role": "publisher"
  }
}
```

---

### 8.3 发布跑腿

`POST /errand/publish`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题 |
| desc | string | 否 | 描述 |
| category | string | 否 | 分类 |
| route_start | string | 否 | 出发地 |
| route_end | string | 否 | 目的地 |
| deadline | string | 否 | 截止时间 |
| reward | string | 否 | 酬劳 |
| contact | string | 否 | 联系方式 |

响应 `200`：

```json
{
  "code": 200,
  "message": "内容正在审核中，仅自己可见",
  "data": {
    "id": "69f3a5debc1e8764471c7112",
    "feed_type": "errand",
    "title": "帮忙取快递",
    "review_status": "reviewing",
    "..."
  }
}
```

---

### 8.4 接受跑腿

`POST /errand/accept`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | string | 是 | FeedItem 的 MongoDB ID |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```

---

## 9. 二手交易模块

### 9.1 获取二手交易列表

`GET /market/list`

**需要 Bearer Token 认证**

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 分类筛选 |
| keyword | string | 否 | 搜索关键词 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

响应 `200`：同 Feed 列表格式，feed_type 为 `market`

---

### 9.2 获取二手交易详情

`GET /market/detail/:id`

**需要 Bearer Token 认证**

路径参数：`id` - FeedItem 的 MongoDB ID

响应 `200`：返回 FeedItem 对象

---

### 9.3 发布二手交易

`POST /market/publish`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题 |
| desc | string | 否 | 描述 |
| price | string | 否 | 售价 |
| original_price | string | 否 | 原价 |
| category | string | 否 | 分类 |
| condition | string | 否 | 成色（如：九成新） |
| trade_mode | string | 否 | 交易方式（如：面交、邮寄） |
| contact | string | 否 | 联系方式 |
| images | string[] | 否 | 图片URL列表 |

响应 `200`：同跑腿发布格式，feed_type 为 `market`

---

### 9.4 收藏二手交易

`POST /market/favorite`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| product_id | string | 是 | FeedItem 的 MongoDB ID |
| action | string | 否 | 操作类型 |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```

---

## 10. 拼车模块

### 10.1 获取拼车列表

`GET /carpool/list`

**需要 Bearer Token 认证**

查询参数：同其他 Feed 列表

响应 `200`：同 Feed 列表格式，feed_type 为 `carpool`

---

### 10.2 发布拼车

`POST /carpool/publish`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 分类（如：顺风车） |
| from | string | 否 | 出发地 |
| to | string | 否 | 目的地 |
| time | string | 否 | 出发时间 |
| type | string | 否 | 类型 |
| seats_text | string | 否 | 座位描述（如：3座） |
| price | string | 否 | 费用 |
| note | string | 否 | 备注 |
| tags | string[] | 否 | 标签 |
| contact | string | 否 | 联系方式 |

响应 `200`：同跑腿发布格式，feed_type 为 `carpool`，title 自动生成为 `{from} → {to}`

---

## 11. 资料模块

### 11.1 获取资料列表

`GET /resource/list`

**需要 Bearer Token 认证**

查询参数：同其他 Feed 列表

响应 `200`：同 Feed 列表格式，feed_type 为 `resource`

---

### 11.2 获取资料详情

`GET /resource/detail/:id`

**需要 Bearer Token 认证**

路径参数：`id` - FeedItem 的 MongoDB ID

响应 `200`：返回 FeedItem 对象

---

### 11.3 发布资料

`POST /resource/publish`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题 |
| desc | string | 否 | 描述 |
| category | string | 否 | 分类 |
| course_name | string | 否 | 课程名称 |
| contact | string | 否 | 联系方式 |
| file_paths | string[] | 否 | 已上传文件的 COS 对象键列表（与上传接口返回的 `path` 一致） |

响应 `200`：同跑腿发布格式，feed_type 为 `resource`

---

### 11.4 收藏资料

`POST /resource/favorite`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| resource_id | string | 是 | FeedItem 的 MongoDB ID |
| action | string | 否 | 操作类型 |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```

---

## 12. 失物招领模块

### 12.1 获取失物招领列表

`GET /lostFound/list`

**需要 Bearer Token 认证**

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 分类筛选 |
| keyword | string | 否 | 搜索关键词 |
| type | string | 否 | 类型筛选（lost/found/all） |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

响应 `200`：同 Feed 列表格式，feed_type 为 `lostFound`

---

### 12.2 获取失物招领详情

`GET /lostFound/detail/:id`

**需要 Bearer Token 认证**

路径参数：`id` - FeedItem 的 MongoDB ID

响应 `200`：返回 FeedItem 对象

---

### 12.3 发布失物招领

`POST /lostFound/publish`

**需要 Bearer Token 认证**

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 类型：lost(寻物) / found(招领) |
| category | string | 否 | 分类（如：证件、钥匙） |
| title | string | 是 | 标题 |
| desc | string | 否 | 描述 |
| location | string | 否 | 地点 |
| event_time | string | 否 | 发生时间 |
| item_feature | string | 否 | 物品特征 |
| contact | string | 否 | 联系方式 |
| reward | string | 否 | 酬谢 |

响应 `200`：同跑腿发布格式，feed_type 为 `lostFound`

---

## 13. 文件上传模块

### 13.1 上传文件

`POST /upload/file`

**需要 Bearer Token 认证**

请求格式：`multipart/form-data`

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 上传的文件 |

响应 `200`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "path": "uploads/app-id/ab/abcdef1234567890abcdef1234567890.pdf",
    "url": "https://bucket.cos.ap-xxx.myqcloud.com/...&q-sign-algorithm=...",
    "name": "笔记.pdf",
    "file_type": "application/pdf",
    "file_size": "2.5 MB",
    "md5": "d41d8cd98f00b204e9800998ecf8427e"
  }
}
```

---

## 14. 微信回调模块

### 14.1 微信内容审核回调

`POST /wechat/callback`

**无需认证**（微信服务器调用）

请求格式：`application/xml`

请求体示例：

```xml
<xml>
  <ToUserName><![CDATA[gh_123456]]></ToUserName>
  <FromUserName><![CDATA[ouser123]]></FromUserName>
  <CreateTime>1234567890</CreateTime>
  <MsgType>event</MsgType>
  <Event>wxa_media_check</Event>
  <appid>wx1234567890abcdef</appid>
  <trace_id>trace_id_here</trace_id>
  <version>2</version>
  <errcode>0</errcode>
  <result>
    <suggest>pass</suggest>
    <label>0</label>
  </result>
</xml>
```

响应 `200`：返回纯文本 `success`

---

## 15. 通知模块

### 15.1 发送通知

`POST /notification/send`

**需要 X-API-Key 认证**（请求头 `X-API-Key: <key>`）

请求体：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| student_ids | string[] | 否 | 目标学号列表（send_all=false 时必填） |
| send_all | bool | 否 | 是否发送给全部学生，默认 false |
| template_data | object | 是 | 模板数据，key 为模板变量名 |
| template_data.*.value | string | 是 | 变量值 |
| template_data.*.color | string | 否 | 文字颜色 |
| page | string | 否 | 点击跳转页面路径 |
| template_id | string | 否 | 指定模板ID（不填使用默认） |
| idempotency_key | string | 否 | 幂等键，防止重复发送 |

响应 `200`：

```json
{
  "task_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "processing"
}
```

---

### 15.2 获取通知任务状态

`GET /notification/task/:task_id`

**需要 X-API-Key 认证**

路径参数：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | string | 是 | 任务ID |

响应 `200`：

```json
{
  "task_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "completed",
  "total": 50,
  "success": 48,
  "failed": 2,
  "not_found_ids": ["2023999"],
  "error": "",
  "created_at": "2026-05-01T10:00:00Z",
  "updated_at": "2026-05-01T10:00:05Z"
}
```

---

## 通用错误码

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 / Token 无效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用（如 API Key 未配置） |

## 认证说明

### Bearer Token 认证

大部分接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

Token 通过 `/auth/login` 或 `/auth/wechat/login` 接口获取。

### X-API-Key 认证

通知模块接口需要在请求头中携带 API Key：

```
X-API-Key: your_api_key_here
```

## FeedItem 通用字段

所有 Feed 类型共享以下基础字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | MongoDB ObjectID |
| feed_type | string | 类型：errand/market/carpool/resource/lostFound |
| app_id | string | 应用ID |
| user_id | uint | 发布者用户ID |
| title | string | 标题 |
| desc | string | 描述 |
| image | string | 封面图URL |
| publisher | string | 发布者昵称 |
| publisher_initial | string | 发布者首字母 |
| created_at | string | 创建时间（ISO 8601） |
| updated_at | string | 更新时间（ISO 8601） |
| review_status | string | 审核状态：reviewing/approved/rejected |
| review_remark | string/null | 审核备注 |
