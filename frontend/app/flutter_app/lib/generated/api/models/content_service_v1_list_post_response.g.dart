// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_list_post_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1ListPostResponse _$ContentServiceV1ListPostResponseFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1ListPostResponse(
  items: (json['items'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Post.fromJson(e as Map<String, dynamic>))
      .toList(),
  total: json['total'] as String?,
);

Map<String, dynamic> _$ContentServiceV1ListPostResponseToJson(
  ContentServiceV1ListPostResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
