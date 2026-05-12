import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:weouc_mobile/features/home/data/repositories/home_repository.dart';
import 'package:weouc_mobile/features/home/data/services/home_remote_service.dart';
import 'package:weouc_mobile/features/home/domain/models/home_module.dart';
import 'package:weouc_mobile/shared/core/command.dart';
import 'package:weouc_mobile/shared/core/result.dart';
import 'package:weouc_mobile/shared/data/network/api_client.dart';

final homeRemoteServiceProvider = Provider<HomeRemoteService>((ref) {
  return HomeRemoteService(ref.read(apiClientProvider));
});

final homeRepositoryProvider = Provider<HomeRepository>((ref) {
  return HomeRepository(ref.read(homeRemoteServiceProvider));
});

final homeViewModelProvider = ChangeNotifierProvider<HomeViewModel>((ref) {
  return HomeViewModel(ref.read(homeRepositoryProvider));
});

class HomeViewModel extends ChangeNotifier {
  HomeViewModel(this._repository) {
    loadModules = Command0(_loadModules);
    loadModules.addListener(_onCommandChanged);
    loadModules.execute();
  }

  final HomeRepository _repository;
  late final Command0<List<HomeModule>> loadModules;

  List<HomeModule> modules = const [];
  String? errorMessage;

  Future<Result<List<HomeModule>>> _loadModules() async {
    try {
      final data = await _repository.loadModules();
      modules = data;
      errorMessage = null;
      return Result.ok(data);
    } catch (error) {
      errorMessage = '首页模块加载失败，请稍后重试';
      return Result.error(Exception(error.toString()));
    }
  }

  Future<void> retryLoad() async {
    await loadModules.execute();
  }

  void _onCommandChanged() {
    notifyListeners();
  }

  @override
  void dispose() {
    loadModules.removeListener(_onCommandChanged);
    loadModules.dispose();
    super.dispose();
  }
}
