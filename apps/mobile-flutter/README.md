# Flutter App

## 技术选型

- `Flutter`
- `flutter_riverpod`
- `go_router`
- `dio`
- `json_serializable`

## 当前状态

- 已初始化可运行 Flutter 工程（`android + ios`）
- 已建立 `feature-first + MVVM + repository` 基线目录与首个 `home` 特性壳层
- 当前首页模块数据由 `service` 层静态提供，后续再切换真实接口

## 目录基线

```text
lib/
├── app/
│   ├── router/
│   └── weouc_app.dart
├── features/
│   └── home/
│       ├── data/
│       │   ├── models/
│       │   ├── repositories/
│       │   └── services/
│       ├── domain/
│       │   └── models/
│       └── presentation/
│           ├── viewmodels/
│           └── views/
├── shared/
│   ├── core/
│   ├── data/
│   └── ui/
└── main.dart
```

## 本地开发

```bash
cd /Users/liangluo/code/weouc2026/apps/mobile-flutter
flutter pub get
flutter run
```

## 本地校验

```bash
cd /Users/liangluo/code/weouc2026/apps/mobile-flutter
flutter analyze
flutter test
```

## 约束

- 页面层不承载业务规则，业务规则统一放在后端或 `repository` 层
- 网络与缓存访问只通过 `data/services`、`data/repositories`
- 联系方式等受限字段只消费后端裁决结果，不在客户端推断
