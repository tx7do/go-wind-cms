<script lang="ts" setup>
import { ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { EchartsUI, type EchartsUIType, useEcharts } from '@vben/plugins/echarts';

import { useGetInteractionStats } from '#/api/composables/stats';

const chartRef = ref<EchartsUIType>();
const { renderEcharts } = useEcharts(chartRef);

const { data } = useGetInteractionStats();

// 无标题帖子与超长标题的兜底展示
function displayTitle(title: string | undefined) {
  if (!title) {
    return $t('page.analytics.titleFallback');
  }
  return title.length > 16 ? `${title.slice(0, 16)}…` : title;
}

watch(
  data,
  (val) => {
    if (!val) return;
    const items = val.topLikedPosts ?? [];
    // echarts 的 y 类轴自下而上,reverse 使 TOP1 显示在顶部
    const titles = items
      .map((item) => displayTitle(item.title))
      .reverse();
    const counts = items.map((item) => Number(item.likeCount ?? 0)).reverse();
    const fullTitles = items
      .map((item) => item.title || $t('page.analytics.titleFallback'))
      .reverse();
    renderEcharts({
      grid: { bottom: 0, containLabel: true, left: '1%', right: '5%', top: '2%' },
      series: [
        {
          data: counts,
          itemStyle: { color: '#019680' },
          label: { position: 'right', show: true },
          name: $t('page.analytics.totalLikes'),
          type: 'bar' as const,
        },
      ],
      tooltip: {
        formatter: (params: any) => {
          const idx = Array.isArray(params) ? params[0]!.dataIndex : params.dataIndex;
          return `${fullTitles[idx]}<br/>${$t('page.analytics.totalLikes')}: ${counts[idx]}`;
        },
        trigger: 'axis',
      },
      xAxis: { minInterval: 1, type: 'value' },
      yAxis: {
        axisTick: { show: false },
        data: titles,
        type: 'category',
      },
    });
  },
  { immediate: true },
);
</script>

<template>
  <EchartsUI ref="chartRef" />
</template>
