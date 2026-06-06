// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'comment_service_v1_update_comment_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CommentServiceV1UpdateCommentRequest
_$CommentServiceV1UpdateCommentRequestFromJson(Map<String, dynamic> json) =>
    CommentServiceV1UpdateCommentRequest(
      id: (json['id'] as num?)?.toInt(),
      data: json['data'] == null
          ? null
          : CommentServiceV1Comment.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
      updateMask: json['updateMask'] as String?,
      allowMissing: json['allowMissing'] as bool?,
    );

Map<String, dynamic> _$CommentServiceV1UpdateCommentRequestToJson(
  CommentServiceV1UpdateCommentRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
