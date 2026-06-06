import 'dart:convert';

import 'package:protobuf/protobuf.dart';

extension ProtobufEx on GeneratedMessage {
  Map<String, dynamic> writeToProto3JsonMap() {
    return json.decode(json.encode(toProto3Json()));
  }
}
