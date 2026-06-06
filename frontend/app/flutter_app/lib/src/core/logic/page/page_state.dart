// ignore_for_file: public_member_api_docs, sort_constructors_first
part of 'page_cubit.dart';

// Sealed class representing the state of a page
sealed class PageState<T extends Object> extends Equatable {
  const PageState({this.limit = 10});
  final int limit;

  @override
  List<Object> get props => [limit];
}

// Represents the initial state of a page
class PageInitial<T extends Object> extends PageState<T> {}

// Represents the loading state of a page
class PageLoading<T extends Object> extends PageState<T> {}

// Represents the loaded state of a page
class PageLoaded<T extends Object> extends PageState<T> {
  final List<T> items;
  const PageLoaded({
    required this.items,
  });
  @override
  List<Object> get props => [items, limit];
}

// Represents the failure state of a page
class PageFailure<T extends Object> extends PageState<T> {
  final Object error;
  const PageFailure({
    required this.error,
  });
  @override
  List<Object> get props => [error];
}
