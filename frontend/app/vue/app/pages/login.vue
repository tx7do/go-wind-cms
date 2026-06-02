<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const localePath = useLocalePath()

useHead({ title: t('authentication.login.title') })

const activeTab = ref<string>('account')

const tabs = computed(() => [
  { key: 'account', label: t('authentication.login.tab_account') },
  { key: 'email', label: t('authentication.login.tab_email') },
  { key: 'phone', label: t('authentication.login.tab_phone') },
  { key: 'other', label: t('authentication.login.tab_other') },
])

const textBtn = 'cursor-pointer border-none bg-transparent text-sm text-primary transition-colors hover:text-primary/80 hover:underline'
</script>

<template>
  <AuthLayout
    :title="t('authentication.login.title')"
    :subtitle="t('authentication.login.login_with')"
    :active-tab="activeTab"
    :tabs="tabs"
    @update:active-tab="activeTab = $event"
  >
    <AuthAccountLogin v-if="activeTab === 'account'" />
    <AuthEmailLogin v-else-if="activeTab === 'email'" />
    <AuthPhoneLogin v-else-if="activeTab === 'phone'" />
    <AuthOtherLogin v-else-if="activeTab === 'other'" />

    <template #switchLink>
      {{ t('authentication.login.no_account') }}
      <button :class="textBtn" class="ml-1" @click="navigateTo(localePath('/register'))">
        {{ t('authentication.login.register_now') }}
      </button>
    </template>
  </AuthLayout>
</template>
