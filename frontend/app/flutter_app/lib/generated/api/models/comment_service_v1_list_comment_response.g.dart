// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'comment_service_v1_list_comment_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CommentServiceV1ListCommentResponse
_$CommentServiceV1ListCommentResponseFromJson(Map<String, dynamic> json) =>
    CommentServiceV1ListCommentResponse(
      items: (json['items'] as List<dynamic>?)
          ?.map(
            (e) => CommentServiceV1Comment.fromJson(e as Map<String, dynamic>),
          )
          .toList(),
      total: json['total'] as String?,
    );

Map<String, dynamic> _$CommentServiceV1ListCommentResponseToJson(
  CommentServiceV1ListCommentResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
