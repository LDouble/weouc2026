import 'package:weouc_mobile/features/home/data/services/home_remote_service.dart';
import 'package:weouc_mobile/features/home/domain/models/home_module.dart';

class HomeRepository {
  HomeRepository(this._remoteService);

  final HomeRemoteService _remoteService;

  Future<List<HomeModule>> loadModules() async {
    final moduleDtos = await _remoteService.fetchHomeModules();
    return moduleDtos
        .map(
          (dto) => HomeModule(
            id: dto.id,
            title: dto.title,
            subtitle: dto.subtitle,
            iconKey: dto.iconKey,
            route: dto.route,
            enabled: dto.enabled,
          ),
        )
        .toList(growable: false);
  }
}
