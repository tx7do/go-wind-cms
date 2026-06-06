import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

part 'page_state.dart';

/// type definition for API request function
typedef LoadCallback<T> = Future Function(
  int page,
  int limit,
);

/// Page Cubit for maintaining page state
class PageCubit<T extends Object> extends Cubit<PageState<T>> {
  /// Creates a new instance of [PageCubit].
  ///
  /// - [callback] is a required callback function for loading data.
  /// - [page] is a required integer representing the current page number.
  /// - [items] is an optional list of items.
  /// - [limit] is an optional integer representing the maximum number of items per page. Default value is 10.
  PageCubit({
    required this.callback,
    required this.page,
    List<T>? items,
    this.limit = 10,
  }) : super(items == null ? PageInitial<T>() : PageLoaded<T>(items: items));

  /// A callback function for loading data.
  final LoadCallback<T> callback;

  /// The current page number.
  final int page;

  /// The maximum number of items per page.
  final int limit;

  /// Fetch data with [callback] function
  void fetch() {
    // Check if the state is not PageInitial
    // load data only if state is initial
    if (state is! PageInitial<T>) return;

    // Emit a PageLoading event
    emit(PageLoading<T>());

    // Call the callback function and handle the result
    callback(page, limit).then(
      (result) {
        // Handle the success case
        result.fold((success) {
          // Emit a PageLoaded event with the success.data as items
          emit(PageLoaded<T>(items: success.data));
        }, (failure) {
          // Handle the failure case
          // Emit a PageFailure event with the failure value as error
          emit(PageFailure<T>(error: failure));
        });
      },
    ).catchError((Object error) {
      // Handle any uncaught error during execution
      // Emit a PageFailure event with the error argument
      emit(PageFailure<T>(error: error));
    });
  }
}
