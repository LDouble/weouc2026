import 'package:json_annotation/json_annotation.dart';
import 'package:weouc_mobile/features/home/domain/models/home_module.dart';

part 'home_module_dto.g.dart';

@JsonSerializable(fieldRename: FieldRename.snake)
class HomeModuleDto {
  const HomeModuleDto({
    required this.id,
    required this.title,
    required this.subtitle,
    required this.iconKey,
    required this.route,
    required this.enabled,
  });

  factory HomeModuleDto.fromJson(Map<String, dynamic> json) =>
      _$HomeModuleDtoFromJson(json);

  final String id;
  final String title;
  final String subtitle;
  final String iconKey;
  final String route;
  final bool enabled;

  Map<String, dynamic> toJson() => _$HomeModuleDtoToJson(this);
}

extension HomeModuleDtoMapper on HomeModuleDto {
  HomeModule toDomain() {
    return HomeModule(
      id: id,
      title: title,
      subtitle: subtitle,
      iconKey: iconKey,
      route: route,
      enabled: enabled,
    );
  }
}
