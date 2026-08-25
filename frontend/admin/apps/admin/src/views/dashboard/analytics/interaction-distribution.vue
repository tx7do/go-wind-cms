<script lang="ts" setup>
import { ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { EchartsUI, type EchartsUIType, useEcharts } from '@vben/plugins/echarts';

import { useGetInteractionStats } from '#/api/composables/stats';

const chartRef = ref<EchartsUIType>();
const { renderEcharts } = useEcharts(chartRef);

const { data } = useGetInteractionStats();

watch(
  data,
  (val) => {
    if (!val) return;
    renderEcharts({
      legend: { bottom: 0 },
      series: [
        {
          avoidLabelOverlap: true,
          data: [
            { itemStyle: { color: '#019680' }, name: $t('page.analytics.totalLikes'), value: Number(val.totalLikes ?? 0) },
            { itemStyle: { color: '#5ab1ef' }, name: $t('page.analytics.totalWatches'), value: Number(val.totalWatches ?? 0) },
          ],
          emphasis: {
            label: { fontSize: 16, fontWeight: 'bold', show: true },
          },
          label: { formatter: '{b}: {c}', show: true },
          radius: ['45%', '75%'],
          type: 'pie' as const,
        },
      ],
      tooltip: { trigger: 'item' },
    });
  },
  { immediate: true },
);
</script>

<template>
  <EchartsUI ref="chartRef" />
</template>
