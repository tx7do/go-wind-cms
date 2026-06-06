import 'package:logger/logger.dart' as logger;

class Logger {
  final String _tag;
  final logger.Logger _logger;

  Logger(String tag)
      : _tag = tag,
        _logger = logger.Logger(
          printer: logger.PrettyPrinter(
            methodCount: 0,
            colors: true,
            printEmojis: true,
            dateTimeFormat: logger.DateTimeFormat.none,
          ),
        );

  void t(String msg) {
    _logger.t("$_tag :: $msg");
  }

  void d(String msg) {
    _logger.d("$_tag :: $msg");
  }

  void i(String msg) {
    _logger.i("$_tag :: $msg");
  }

  void w(String msg) {
    _logger.w("$_tag :: $msg");
  }

  void e(String msg, [Object? err, StackTrace? stackTrace]) =>
      _logger.e('$_tag :: $msg', error: err, stackTrace: stackTrace);

  void f(String msg, [Object? err, StackTrace? stackTrace]) =>
      _logger.f('$_tag :: $msg', error: err, stackTrace: stackTrace);
}

final Logger _logger = Logger('im');

final trace = _logger.t;

final debug = _logger.d;

final info = _logger.i;

final warn = _logger.w;

final error = _logger.e;

final fatal = _logger.f;
