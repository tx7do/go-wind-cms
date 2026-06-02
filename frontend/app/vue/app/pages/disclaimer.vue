<script setup lang="ts">
const { t } = useI18n()
</script>

<template>
  <LayoutLegalPage
    icon="carbon:warning"
    :title="t('page.legal.disclaimer.title')"
    :description="t('page.legal.disclaimer.description')"
    :items="[t('page.legal.disclaimer.item_1'), t('page.legal.disclaimer.item_2'), t('page.legal.disclaimer.item_3')]"
  />
</template>
