<script lang="ts" setup>
import type { AnalysisOverviewItem } from '@vben/common-ui';

import { computed } from 'vue';

import { AnalysisChartCard, AnalysisOverview } from '@vben/common-ui';
import {
  SvgBellIcon,
  SvgCakeIcon,
  SvgCardIcon,
  SvgDownloadIcon,
} from '@vben/icons';
import { $t } from '@vben/locales';

import { useGetDashboardOverview } from '#/api/composables/stats';

import ContentTrend from './content-trend.vue';
import InteractionDistribution from './interaction-distribution.vue';
import InteractionTop from './interaction-top.vue';
import LoginActivity from './login-activity.vue';

const { data: overview } = useGetDashboardOverview();

// proto3 int64 经 JSON 序列化为字符串("2"),生成器类型标注为 number 与运行时
// 不符,展示组件(CountToAnimator)会因 prop 类型校验失败,统一转 number。
function toNum(v: number | undefined) {
  return Number(v ?? 0);
}

// 卡片主数字为总量,footer 为近 7 天新增
const overviewItems = computed<AnalysisOverviewItem[]>(() => [
  {
    icon: SvgCardIcon,
    title: $t('page.analytics.userCountTotal'),
    totalTitle: $t('page.analytics.weeklyNew'),
    totalValue: toNum(overview.value?.newUserCountWeek),
    value: toNum(overview.value?.userCount),
  },
  {
    icon: SvgCakeIcon,
    title: $t('page.analytics.postCountTotal'),
    totalTitle: $t('page.analytics.weeklyNew'),
    totalValue: toNum(overview.value?.newPostCountWeek),
    value: toNum(overview.value?.postCount),
  },
  {
    icon: SvgDownloadIcon,
    title: $t('page.analytics.commentCountTotal'),
    totalTitle: $t('page.analytics.weeklyNew'),
    totalValue: toNum(overview.value?.newCommentCountWeek),
    value: toNum(overview.value?.commentCount),
  },
  {
    icon: SvgBellIcon,
    title: $t('page.analytics.interactionCountTotal'),
    totalTitle: $t('page.analytics.weeklyNew'),
    totalValue: toNum(overview.value?.newLikeCountWeek),
    value: toNum(overview.value?.interactionCount),
  },
]);
</script>

<template>
  <div class="p-5">
    <AnalysisOverview :items="overviewItems" />

    <AnalysisChartCard
      :title="$t('page.analytics.contentTrend')"
      class="mt-5"
    >
      <ContentTrend />
    </AnalysisChartCard>

    <div class="mt-5 w-full md:flex">
      <AnalysisChartCard
        class="md:mr-4 md:w-1/3"
        :title="$t('page.analytics.interactionTop')"
      >
        <InteractionTop />
      </AnalysisChartCard>
      <AnalysisChartCard
        class="mt-5 md:mr-4 md:mt-0 md:w-1/3"
        :title="$t('page.analytics.interactionDistribution')"
      >
        <InteractionDistribution />
      </AnalysisChartCard>
      <AnalysisChartCard
        class="mt-5 md:mt-0 md:w-1/3"
        :title="$t('page.analytics.loginActivity')"
      >
        <LoginActivity />
      </AnalysisChartCard>
    </div>
  </div>
</template>
