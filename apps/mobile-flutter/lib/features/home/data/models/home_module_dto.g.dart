// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'home_module_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

HomeModuleDto _$HomeModuleDtoFromJson(Map<String, dynamic> json) =>
    HomeModuleDto(
      id: json['id'] as String,
      title: json['title'] as String,
      subtitle: json['subtitle'] as String,
      iconKey: json['icon_key'] as String,
      route: json['route'] as String,
      enabled: json['enabled'] as bool,
    );

Map<String, dynamic> _$HomeModuleDtoToJson(HomeModuleDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'title': instance.title,
      'subtitle': instance.subtitle,
      'icon_key': instance.iconKey,
      'route': instance.route,
      'enabled': instance.enabled,
    };
