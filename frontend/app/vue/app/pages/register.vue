<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const localePath = useLocalePath()

useHead({ title: t('authentication.register.title') })

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
    :title="t('authentication.register.title')"
    :subtitle="t('authentication.register.register_with')"
    :active-tab="activeTab"
    :tabs="tabs"
    @update:active-tab="activeTab = $event"
  >
    <AuthAccountRegister v-if="activeTab === 'account'" />
    <AuthEmailRegister v-else-if="activeTab === 'email'" />
    <AuthPhoneRegister v-else-if="activeTab === 'phone'" />
    <AuthOtherRegister v-else-if="activeTab === 'other'" />

    <template #switchLink>
      {{ t('authentication.register.already_have_account') }}
      <button :class="textBtn" class="ml-1" @click="navigateTo(localePath('/login'))">
        {{ t('authentication.register.login_now') }}
      </button>
    </template>
  </AuthLayout>
</template>
