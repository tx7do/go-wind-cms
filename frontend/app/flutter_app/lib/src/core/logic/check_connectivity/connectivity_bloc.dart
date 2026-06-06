import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:equatable/equatable.dart';

part 'connectivity_event.dart';

part 'connectivity_state.dart';

class ConnectivityBloc extends Bloc<ConnectivityEvent, ConnectivityState> {
  late StreamSubscription<List<ConnectivityResult>> _connectivitySubscription;

  ConnectivityBloc() : super(ConnectivityInitial()) {
    _registerListener();

    initialize();
  }

  @override
  Future<void> close() {
    _connectivitySubscription.cancel();
    return super.close();
  }

  void initialize() async {
    final connectivityResult = await Connectivity().checkConnectivity();
    bool isConnected = !connectivityResult.contains(ConnectivityResult.none);
    add(ConnectivityChanged(isConnected));
  }

  void _registerListener() {
    // Listen to connectivity changes
    _connectivitySubscription = Connectivity()
        .onConnectivityChanged
        .listen((List<ConnectivityResult> result) {
      // Determine if there is an internet connection
      bool isConnected = !result.contains(ConnectivityResult.none);

      // debug('ConnectivityBloc, $result');
      // Add the ConnectivityChanged event
      add(ConnectivityChanged(isConnected));
    });

    // Handle the ConnectivityChanged event
    on<ConnectivityChanged>((event, emit) {
      if (event.isConnected) {
        emit(ConnectivitySuccess(event.isConnected));
      } else {
        emit(ConnectivityFailure());
      }
    });
  }
}
