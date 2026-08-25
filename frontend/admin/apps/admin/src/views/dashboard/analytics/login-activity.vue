<script lang="ts" setup>
import { ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { EchartsUI, type EchartsUIType, useEcharts } from '@vben/plugins/echarts';

import { useGetLoginActivity } from '#/api/composables/stats';

const chartRef = ref<EchartsUIType>();
const { renderEcharts } = useEcharts(chartRef);

const { data } = useGetLoginActivity();

// 日期截短为 MM-DD,减小 x 轴密度
function shortDate(date: string | undefined) {
  return date?.length === 10 ? date.slice(5) : (date ?? '');
}

watch(
  data,
  (val) => {
    if (!val) return;
    const success = val.success ?? [];
    const dates = success.map((item) => shortDate(item.date));
    const series = [
      {
        areaStyle: {},
        data: success.map((item) => Number(item.value ?? 0)),
        itemStyle: { color: '#019680' },
        name: $t('page.analytics.loginSuccess'),
        smooth: true,
        type: 'line' as const,
      },
      {
        areaStyle: {},
        data: (val.failed ?? []).map((item) => Number(item.value ?? 0)),
        itemStyle: { color: '#d48265' },
        name: $t('page.analytics.loginFailed'),
        smooth: true,
        type: 'line' as const,
      },
    ];
    renderEcharts({
      grid: { bottom: 40, containLabel: true, left: '1%', right: '1%', top: '5%' },
      legend: { bottom: 0, data: series.map((item) => item.name) },
      series,
      tooltip: { axisPointer: { type: 'line' }, trigger: 'axis' },
      xAxis: {
        axisTick: { show: false },
        boundaryGap: false,
        data: dates,
        type: 'category',
      },
      yAxis: { minInterval: 1, type: 'value' },
    });
  },
  { immediate: true },
);
</script>

<template>
  <EchartsUI ref="chartRef" />
</template>
