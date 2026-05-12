import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:weouc_mobile/app/router/app_router.dart';
import 'package:weouc_mobile/shared/ui/app_theme.dart';

class WeoucApp extends ConsumerWidget {
  const WeoucApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final appRouter = ref.watch(appRouterProvider);

    return MaterialApp.router(
      title: 'weouc 校园服务',
      theme: AppTheme.light(),
      routerConfig: appRouter,
    );
  }
}
