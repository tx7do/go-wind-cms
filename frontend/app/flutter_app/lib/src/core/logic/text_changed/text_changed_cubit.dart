import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

part 'text_changed_state.dart';

class TextChangedCubit extends Cubit<TextChangedState> {
  TextChangedCubit() : super(const TextChangedState(''));

  set changeText(String text) {
    emit(TextChangedState(text));
  }
}
