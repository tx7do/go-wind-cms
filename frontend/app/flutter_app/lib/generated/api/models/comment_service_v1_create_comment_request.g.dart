// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'comment_service_v1_create_comment_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CommentServiceV1CreateCommentRequest
_$CommentServiceV1CreateCommentRequestFromJson(Map<String, dynamic> json) =>
    CommentServiceV1CreateCommentRequest(
      data: json['data'] == null
          ? null
          : CommentServiceV1Comment.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$CommentServiceV1CreateCommentRequestToJson(
  CommentServiceV1CreateCommentRequest instance,
) => <String, dynamic>{'data': instance.data};
