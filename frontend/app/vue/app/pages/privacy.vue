<script setup lang="ts">
const { t } = useI18n()

useHead({ title: t('page.legal.privacy.title') })
</script>

<template>
  <LayoutLegalPage
    icon="carbon:locked"
    :title="t('page.legal.privacy.title')"
    :description="t('page.legal.privacy.description')"
    :items="[t('page.legal.privacy.item_1'), t('page.legal.privacy.item_2'), t('page.legal.privacy.item_3')]"
  />
</template>
