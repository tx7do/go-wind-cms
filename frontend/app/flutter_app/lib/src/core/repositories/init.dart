import 'package:get_it/get_it.dart';

import 'index.dart';
import 'package:flutter_app/src/core/preference/preference.dart';

Future<void> init() async {
  final getIt = GetIt.instance;

  getIt.registerLazySingleton<UserAuthCache>(() => UserAuthCache());
  getIt.registerLazySingleton<UserPreferenceCache>(() => UserPreferenceCache());
  getIt.registerLazySingleton<LanguageListRepository>(
      () => LanguageListRepository());

  await GetIt.instance<UserAuthCache>().init();
  await GetIt.instance<UserPreferenceCache>().init();
}
