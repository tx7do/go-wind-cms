part of 'text_changed_cubit.dart';

final class TextChangedState extends Equatable {
  final String text;

  const TextChangedState(this.text);

  @override
  List<Object> get props => [text];
}
