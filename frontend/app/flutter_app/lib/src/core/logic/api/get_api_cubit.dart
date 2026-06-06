import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:flutter_app/src/core/utilities/logger.dart' show fatal;
import 'api_state.dart';

export 'api_state.dart';

class GetApiCubit<T> extends Cubit<ApiResponse<T>> {
  GetApiCubit() : super(Initial());

  Future<bool> listItems(Future<T> Function() loader) async {
    emit(Loading());

    try {
      final item = await loader();
      emit(Success<T>(item));
      return true;
    } on Exception catch (e) {
      fatal("get item from api exception: $e");
      emit(Error<T>(e.toString()));
      return false;
    }
  }
}
