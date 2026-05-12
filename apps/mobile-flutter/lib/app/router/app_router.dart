import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:weouc_mobile/features/home/presentation/views/home_page.dart';
import 'package:weouc_mobile/features/home/presentation/views/module_placeholder_page.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/',
    routes: [
      GoRoute(
        path: '/',
        name: 'home',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/module/:moduleId',
        name: 'module',
        builder: (context, state) {
          final moduleId = state.pathParameters['moduleId'] ?? 'unknown';
          final moduleTitle = state.uri.queryParameters['title'] ?? moduleId;
          return ModulePlaceholderPage(
            moduleId: moduleId,
            moduleTitle: moduleTitle,
          );
        },
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      appBar: AppBar(title: const Text('页面不存在')),
      body: Center(child: Text('无法打开路由：${state.uri}')),
    ),
  );
});
