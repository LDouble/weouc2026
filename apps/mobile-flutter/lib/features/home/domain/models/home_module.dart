class HomeModule {
  const HomeModule({
    required this.id,
    required this.title,
    required this.subtitle,
    required this.iconKey,
    required this.route,
    required this.enabled,
  });

  final String id;
  final String title;
  final String subtitle;
  final String iconKey;
  final String route;
  final bool enabled;
}
