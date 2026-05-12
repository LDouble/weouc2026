import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:weouc_mobile/app/weouc_app.dart';

void main() {
  testWidgets('Flutter 壳层首页可渲染', (WidgetTester tester) async {
    await tester.pumpWidget(const ProviderScope(child: WeoucApp()));
    await tester.pumpAndSettle();

    expect(find.text('校园服务'), findsOneWidget);
    expect(find.text('weouc2026 Flutter 壳层'), findsOneWidget);
  });
}
