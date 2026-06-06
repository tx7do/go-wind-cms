import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';

/// UI状态
@immutable
sealed class UIState extends Equatable {
  const UIState();
}

/// 初始状态
final class InitialState extends UIState {
  const InitialState();

  @override
  List<Object> get props => [];
}

/// 加载状态
final class LoadingState extends UIState {
  const LoadingState();

  @override
  List<Object> get props => [];
}

/// 成功状态
final class SuccessState<T> extends UIState {
  final T data;

  const SuccessState(this.data);

  @override
  List<Object> get props => [
        <T>[data]
      ];
}

/// 失败状态
final class ErrorState extends UIState {
  final String message;

  const ErrorState(this.message);

  @override
  List<Object> get props => [message];
}
