import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:flutter_app/src/core/widgets/adaptive_shell.dart';
import 'package:flutter_app/src/features/cms/services/nav_destinations.dart';

ShellTab _default(int branchIndex, String route, String label) => ShellTab(
      branchIndex: branchIndex,
      route: route,
      label: label,
      icon: Icons.home_outlined,
      activeIcon: Icons.home,
    );

List<ShellTab> _defaults() => [
      _default(0, '/', '首页'),
      _default(1, '/explore', '发现'),
      _default(2, '/bookmarks', '收藏'),
      _default(3, '/profile', '我的'),
    ];

NavSlot _slot(String route, {String label = '服务端'}) => NavSlot(
      label: label,
      route: route,
      icon: Icons.explore_outlined,
      activeIcon: Icons.explore,
    );

void main() {
  group('底栏 tab 合并规则（服务端导航 × 内置分支）', () {
    test('服务端为空：原样返回内置四支', () {
      final merged = AdaptiveShell.mergeTabs(_defaults(), const []);
      expect(merged, hasLength(4));
      expect(merged.map((t) => t.branchIndex), [0, 1, 2, 3]);
      expect(merged.map((t) => t.label), ['首页', '发现', '收藏', '我的']);
    });

    test('服务端可重排顺序，branchIndex 跟着 tab 走不跟着位置走', () {
      final merged = AdaptiveShell.mergeTabs(
        _defaults(),
        [_slot('/profile', label: '我'), _slot('/', label: '主页')],
      );
      expect(merged.map((t) => t.route).toList(), ['/profile', '/', '/explore', '/bookmarks']);
      expect(merged.first.branchIndex, 3);
      expect(merged[1].branchIndex, 0);
      // 服务端没提到的支补在尾部
      expect(merged[2].branchIndex, 1);
      expect(merged[3].branchIndex, 2);
    });

    test('服务端只覆盖标签与图标；空标签回退内置', () {
      final merged = AdaptiveShell.mergeTabs(
        _defaults(),
        [_slot('/explore', label: '逛逛'), _slot('/bookmarks', label: '')],
      );
      expect(merged[0].label, '逛逛');
      expect(merged[0].icon, Icons.explore_outlined);
      expect(merged[1].label, '收藏');
    });

    test('未知路由不进底栏，重复项去重', () {
      final merged = AdaptiveShell.mergeTabs(
        _defaults(),
        [_slot('/post/12'), _slot('/explore/'), _slot('/explore', label: '重复')],
      );
      // 未命中的支按内置顺序补在尾部，四支恒在
      expect(merged.map((t) => t.route), ['/explore', '/', '/bookmarks', '/profile']);
      expect(merged.map((t) => t.branchIndex), [1, 0, 2, 3]);
    });

    test('尾部斜杠视为同一路径', () {
      final merged = AdaptiveShell.mergeTabs(
        _defaults(),
        [_slot('/bookmarks/', label: '书架')],
      );
      expect(merged.first.route, '/bookmarks');
      expect(merged.first.label, '书架');
    });
  });
}
