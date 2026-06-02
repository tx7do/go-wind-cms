<script setup lang="ts">
const { t } = useI18n()

useHead({ title: t('page.legal.terms.title') })
</script>

<template>
  <LayoutLegalPage
    icon="carbon:document"
    :title="t('page.legal.terms.title')"
    :description="t('page.legal.terms.description')"
    :items="[t('page.legal.terms.item_1'), t('page.legal.terms.item_2'), t('page.legal.terms.item_3')]"
  />
</template>
