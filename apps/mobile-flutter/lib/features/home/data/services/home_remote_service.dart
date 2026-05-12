import 'dart:async';

import 'package:weouc_mobile/features/home/data/models/home_module_dto.dart';
import 'package:weouc_mobile/shared/data/network/api_client.dart';

class HomeRemoteService {
  HomeRemoteService(this._apiClient);

  final ApiClient _apiClient;

  Future<List<HomeModuleDto>> fetchHomeModules() async {
    // 当前阶段使用静态数据占位，后续对接真实接口。
    // 这里读取 baseUrl 只是为了保持服务层对网络客户端配置的感知。
    final baseUrl = _apiClient.dio.options.baseUrl;
    if (baseUrl.isEmpty) {
      throw Exception('API baseUrl 未配置');
    }

    await Future<void>.delayed(const Duration(milliseconds: 180));

    const response = [
      {
        'id': 'errand',
        'title': '跑腿',
        'subtitle': '帮取快递、代买代办',
        'icon_key': 'delivery',
        'route': '/module/errand',
        'enabled': true,
      },
      {
        'id': 'meetup',
        'title': '组局',
        'subtitle': '学习自习、运动搭子',
        'icon_key': 'groups',
        'route': '/module/meetup',
        'enabled': true,
      },
      {
        'id': 'market',
        'title': '二手',
        'subtitle': '闲置交易与求购',
        'icon_key': 'store',
        'route': '/module/market',
        'enabled': true,
      },
      {
        'id': 'academic',
        'title': '教务',
        'subtitle': '课表成绩（后续接真实连接器）',
        'icon_key': 'school',
        'route': '/module/academic',
        'enabled': false,
      },
    ];

    return response.map(HomeModuleDto.fromJson).toList(growable: false);
  }
}
