import 'package:flutter/material.dart';

import 'app_theme.dart';
import 'palette.dart';

/// 夜风（深色）。与晴风共用一套组件规则，差异只在色板。
ThemeData getDarkTheme() =>
    buildAppTheme(p: kNightPalette, brightness: Brightness.dark);
