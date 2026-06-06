import 'package:flutter/widgets.dart';

import 'package:flutter_app/generated/l10n.dart';

extension LocalizedBuildContext on BuildContext {
  S get loc => S.of(this);
}
