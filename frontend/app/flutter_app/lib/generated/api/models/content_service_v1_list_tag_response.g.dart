// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_list_tag_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1ListTagResponse _$ContentServiceV1ListTagResponseFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1ListTagResponse(
  items: (json['items'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Tag.fromJson(e as Map<String, dynamic>))
      .toList(),
  total: json['total'] as String?,
);

Map<String, dynamic> _$ContentServiceV1ListTagResponseToJson(
  ContentServiceV1ListTagResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
