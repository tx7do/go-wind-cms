import 'package:flutter/material.dart';

import 'app_theme.dart';
import 'palette.dart';

/// 晴风（浅色）。全部规则在 buildAppTheme 里，这里只选色板。
ThemeData getLightTheme() =>
    buildAppTheme(p: kDayPalette, brightness: Brightness.light);
