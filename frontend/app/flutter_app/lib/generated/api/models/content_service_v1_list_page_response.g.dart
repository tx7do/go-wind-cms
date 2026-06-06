// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_list_page_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1ListPageResponse _$ContentServiceV1ListPageResponseFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1ListPageResponse(
  items: (json['items'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Page.fromJson(e as Map<String, dynamic>))
      .toList(),
  total: json['total'] as String?,
);

Map<String, dynamic> _$ContentServiceV1ListPageResponseToJson(
  ContentServiceV1ListPageResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
