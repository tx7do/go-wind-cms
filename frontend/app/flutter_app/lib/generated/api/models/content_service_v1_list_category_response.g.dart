// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_list_category_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1ListCategoryResponse
_$ContentServiceV1ListCategoryResponseFromJson(Map<String, dynamic> json) =>
    ContentServiceV1ListCategoryResponse(
      items: (json['items'] as List<dynamic>?)
          ?.map(
            (e) => ContentServiceV1Category.fromJson(e as Map<String, dynamic>),
          )
          .toList(),
      total: json['total'] as String?,
    );

Map<String, dynamic> _$ContentServiceV1ListCategoryResponseToJson(
  ContentServiceV1ListCategoryResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
