import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';

/// API回应
@immutable
sealed class ApiResponse<T> extends Equatable {
  const ApiResponse();
}

/// 初始状态
final class Initial<T> extends ApiResponse<T> {
  @override
  List<Object> get props => [];
}

/// 请求状态
final class Loading<T> extends ApiResponse<T> {
  @override
  List<Object> get props => [];
}

/// 成功状态
final class Success<T> extends ApiResponse<T> {
  final T data;

  const Success(this.data);

  @override
  List<Object> get props => [
        <T>[data]
      ];
}

/// 失败状态
final class Error<T> extends ApiResponse<T> {
  final String errorMessage;

  const Error(this.errorMessage);

  @override
  List<Object> get props => [errorMessage];
}
