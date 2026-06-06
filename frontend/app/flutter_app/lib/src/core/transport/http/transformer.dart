import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

import 'package:flutter_app/src/core/utilities/logger.dart' show debug;

class DioTransformer extends BackgroundTransformer {
  @override
  Future transformResponse(
    RequestOptions options,
    ResponseBody responseBody,
  ) async {
    final responseType = options.responseType;
    // Do not handled the body for streams.
    if (responseType == ResponseType.stream) {
      return responseBody;
    }

    debug('DioTransformer.transformResponse');

    final chunks = await responseBody.stream.toList();
    final responseBytes = Uint8List.fromList(chunks.expand((c) => c).toList());

    // Return the finalized bytes if the response type is bytes.
    if (responseType == ResponseType.bytes) {
      return responseBytes;
    }

    final isJsonContent = Transformer.isJsonMimeType(
      responseBody.headers[Headers.contentTypeHeader]?.first,
    );

    final String? response;
    if (options.responseDecoder != null) {
      final decodeResponse = options.responseDecoder!(
        responseBytes,
        options,
        responseBody..stream = const Stream.empty(),
      );

      if (decodeResponse is Future) {
        response = await decodeResponse;
      } else {
        response = decodeResponse;
      }
    } else if (!isJsonContent || responseBytes.isNotEmpty) {
      response = utf8.decode(responseBytes, allowMalformed: true);
    } else {
      response = null;
    }

    if (response != null &&
        response.isNotEmpty &&
        responseType == ResponseType.json &&
        isJsonContent) {
      return jsonDecodeCallback(response);
    }
    return response;
  }
}
