import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:flutter_app/src/core/utilities/logger.dart' show fatal;
import 'api_state.dart';

export 'api_state.dart';

class ListApiCubit<T> extends Cubit<ApiResponse<List<T>>> {
  ListApiCubit() : super(Initial());

  Future<bool> listItems(Future<List<T>> Function() loader) async {
    emit(Loading());

    try {
      final items = await loader();
      emit(Success<List<T>>(items));
      return true;
    } on Exception catch (e) {
      fatal("list items from api exception: $e");
      emit(Error<List<T>>(e.toString()));
      return false;
    }
  }
}
