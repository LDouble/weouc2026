import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:weouc_mobile/features/home/domain/models/home_module.dart';
import 'package:weouc_mobile/features/home/presentation/viewmodels/home_viewmodel.dart';

class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final viewModel = ref.watch(homeViewModelProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('校园服务')),
      body: RefreshIndicator(
        onRefresh: viewModel.retryLoad,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 20),
          children: [
            const _HeaderSection(),
            const SizedBox(height: 16),
            if (viewModel.loadModules.running && viewModel.modules.isEmpty)
              const _LoadingSection()
            else if (viewModel.errorMessage != null &&
                viewModel.modules.isEmpty)
              _ErrorSection(
                message: viewModel.errorMessage!,
                onRetry: viewModel.retryLoad,
              )
            else
              _ModuleGrid(
                modules: viewModel.modules,
                onTap: (module) => _onTapModule(context, module),
              ),
          ],
        ),
      ),
    );
  }

  void _onTapModule(BuildContext context, HomeModule module) {
    if (!module.enabled) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('${module.title} 正在接入中')));
      return;
    }

    final title = Uri.encodeComponent(module.title);
    context.push('${module.route}?title=$title');
  }
}

class _HeaderSection extends StatelessWidget {
  const _HeaderSection();

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'weouc2026 Flutter 壳层',
              style: Theme.of(
                context,
              ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 8),
            Text(
              '当前遵循 feature-first + MVVM + repository，后续会逐步接入真实校园生活和教务能力。',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: Colors.black54,
                height: 1.45,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _LoadingSection extends StatelessWidget {
  const _LoadingSection();

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.symmetric(vertical: 40),
      child: Center(child: CircularProgressIndicator()),
    );
  }
}

class _ErrorSection extends StatelessWidget {
  const _ErrorSection({required this.message, required this.onRetry});

  final String message;
  final Future<void> Function() onRetry;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(message),
            const SizedBox(height: 12),
            FilledButton(onPressed: onRetry, child: const Text('重试')),
          ],
        ),
      ),
    );
  }
}

class _ModuleGrid extends StatelessWidget {
  const _ModuleGrid({required this.modules, required this.onTap});

  final List<HomeModule> modules;
  final ValueChanged<HomeModule> onTap;

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 12,
      runSpacing: 12,
      children: modules
          .map((module) {
            return SizedBox(
              width: (MediaQuery.of(context).size.width - 44) / 2,
              child: InkWell(
                borderRadius: BorderRadius.circular(14),
                onTap: () => onTap(module),
                child: Card(
                  child: Padding(
                    padding: const EdgeInsets.all(14),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        CircleAvatar(
                          radius: 18,
                          backgroundColor: const Color(0xFFE8F0FE),
                          child: Icon(
                            _iconByKey(module.iconKey),
                            color: const Color(0xFF0B57D0),
                          ),
                        ),
                        const SizedBox(height: 12),
                        Text(
                          module.title,
                          style: Theme.of(context).textTheme.titleSmall
                              ?.copyWith(fontWeight: FontWeight.w700),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          module.subtitle,
                          style: Theme.of(context).textTheme.bodySmall
                              ?.copyWith(color: Colors.black54),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          module.enabled ? '已启用' : '规划中',
                          style: Theme.of(context).textTheme.labelSmall
                              ?.copyWith(
                                color: module.enabled
                                    ? const Color(0xFF087443)
                                    : const Color(0xFF9F5F00),
                              ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            );
          })
          .toList(growable: false),
    );
  }

  IconData _iconByKey(String iconKey) {
    const map = {
      'delivery': Icons.local_shipping_outlined,
      'groups': Icons.groups_2_outlined,
      'store': Icons.storefront_outlined,
      'school': Icons.school_outlined,
    };
    return map[iconKey] ?? Icons.apps_outlined;
  }
}
