part of 'connectivity_bloc.dart';

abstract class ConnectivityEvent extends Equatable {
  const ConnectivityEvent();

  @override
  List<Object> get props => [];
}

// Event triggered when connectivity changes
class ConnectivityChanged extends ConnectivityEvent {
  final bool isConnected;

  const ConnectivityChanged(this.isConnected);

  @override
  List<Object> get props => [isConnected];
}
